package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/repository"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

// AdminAuthHandler exposes admin-specific OAuth login endpoints.
type AdminAuthHandler struct {
	service       authsvc.AdminService
	sessions      authsvc.AdminSessionManager
	cookieCfg     authsvc.CookieOptions
	refreshWindow time.Duration
}

// NewAdminAuthHandler constructs the handler with the provided service.
func NewAdminAuthHandler(service authsvc.AdminService, sessions authsvc.AdminSessionManager, authCfg config.AuthConfig) *AdminAuthHandler {
	refreshWindow := authCfg.Admin.SessionRefreshWindow
	if refreshWindow <= 0 {
		refreshWindow = 20 * time.Minute
	}
	return &AdminAuthHandler{
		service:       service,
		sessions:      sessions,
		cookieCfg:     authsvc.NewCookieOptions(authCfg.Admin),
		refreshWindow: refreshWindow,
	}
}

// Login redirects administrators to Google's OAuth consent screen.
func (h *AdminAuthHandler) Login(c *gin.Context) {
	result, err := h.service.StartLogin(c.Request.Context(), c.Query("redirect_uri"))
	if err != nil {
		respondError(c, err)
		return
	}

	if result.AuthURL == "" {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"state":   result.State,
				"message": "authentication disabled",
			},
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, result.AuthURL)
}

// Callback finalizes the admin OAuth flow and redirects to the admin SPA.
func (h *AdminAuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if strings.TrimSpace(state) == "" || strings.TrimSpace(code) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "missing state or code",
		})
		return
	}

	result, err := h.service.HandleCallback(c.Request.Context(), state, code)
	if err != nil {
		respondError(c, err)
		return
	}

	if strings.TrimSpace(result.SessionID) != "" {
		expires := time.Unix(result.ExpiresAt, 0).UTC()
		h.cookieCfg.Write(c.Writer, result.SessionID, expires)
	}

	c.Header("Cache-Control", "no-store")
	target := buildAdminRedirect(result.RedirectPath)
	c.Redirect(http.StatusTemporaryRedirect, target)
}

// Session validates the current administrator session and refreshes the cookie when appropriate.
func (h *AdminAuthHandler) Session(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	sessionID := extractBearerToken(authHeader)
	tokenSource := "header"
	if sessionID == "" {
		if cookieName := strings.TrimSpace(h.cookieCfg.Name); cookieName != "" {
			if cookie, err := c.Cookie(cookieName); err == nil {
				sessionID = strings.TrimSpace(cookie)
				tokenSource = "cookie"
			}
		} else if cookie, err := c.Cookie("ps_admin_session"); err == nil {
			sessionID = strings.TrimSpace(cookie)
			tokenSource = "cookie"
		}
	}

	if strings.TrimSpace(sessionID) == "" {
		h.cookieCfg.Clear(c.Writer)
		respondError(c, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "missing session", nil))
		return
	}

	if h.sessions == nil {
		respondError(c, errs.New(errs.CodeInternal, http.StatusInternalServerError, "session manager not configured", nil))
		return
	}

	session, err := h.sessions.Validate(ctx, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.cookieCfg.Clear(c.Writer)
			respondError(c, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid session", nil))
			return
		}
		respondError(c, err)
		return
	}

	shouldRefresh := time.Until(session.ExpiresAt) <= h.refreshWindow
	if shouldRefresh {
		refreshed, refreshErr := h.sessions.Refresh(ctx, sessionID)
		if refreshErr != nil {
			if errors.Is(refreshErr, repository.ErrNotFound) {
				h.cookieCfg.Clear(c.Writer)
				respondError(c, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "session expired", nil))
				return
			}
			respondError(c, refreshErr)
			return
		}
		session = refreshed
	}

	h.cookieCfg.Write(c.Writer, session.ID, session.ExpiresAt.UTC())
	c.Header("Cache-Control", "no-store")

	c.JSON(http.StatusOK, gin.H{
		"active":    true,
		"expiresAt": session.ExpiresAt.Unix(),
		"email":     session.Email,
		"roles":     session.Roles,
		"source":    tokenSource,
		"refreshed": shouldRefresh,
	})
}

func buildAdminRedirect(path string) string {
	stripped := strings.TrimSpace(path)
	if strings.Contains(stripped, "#") {
		parts := strings.SplitN(stripped, "#", 2)
		stripped = parts[0]
	}
	return normalizeAdminRedirectLocation(stripped)
}

func normalizeAdminRedirectLocation(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/admin/"
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		u, err := url.Parse(trimmed)
		if err != nil {
			return "/admin/"
		}
		u.Path = ensureAdminRedirectPath(u.Path)
		return u.String()
	}

	if !strings.HasPrefix(trimmed, "/") {
		return "/admin/"
	}

	return ensureAdminRedirectPath(trimmed)
}

func ensureAdminRedirectPath(path string) string {
	cleaned := strings.TrimSpace(path)
	if cleaned == "" || cleaned == "/" {
		return "/admin/"
	}
	if cleaned == "/admin" {
		return "/admin/"
	}
	return cleaned
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
