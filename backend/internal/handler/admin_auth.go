package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

type adminSessionCookieOptions struct {
	name     string
	domain   string
	path     string
	secure   bool
	httpOnly bool
	sameSite http.SameSite
}

func newAdminSessionCookieOptions(cfg config.AdminAuthConfig) adminSessionCookieOptions {
	name := strings.TrimSpace(cfg.SessionCookieName)
	if name == "" {
		name = "ps_admin_jwt"
	}

	path := strings.TrimSpace(cfg.SessionCookiePath)
	if path == "" {
		path = "/"
	}

	return adminSessionCookieOptions{
		name:     name,
		domain:   strings.TrimSpace(cfg.SessionCookieDomain),
		path:     path,
		secure:   cfg.SessionCookieSecure,
		httpOnly: cfg.SessionCookieHTTPOnly,
		sameSite: parseSameSite(cfg.SessionCookieSameSite),
	}
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteLaxMode
	}
}

// AdminAuthHandler exposes admin-specific OAuth login endpoints.
type AdminAuthHandler struct {
	service   authsvc.AdminService
	issuer    authsvc.TokenIssuer
	verifier  authsvc.TokenVerifier
	cookieCfg adminSessionCookieOptions
}

// NewAdminAuthHandler constructs the handler with the provided service.
func NewAdminAuthHandler(service authsvc.AdminService, issuer authsvc.TokenIssuer, verifier authsvc.TokenVerifier, authCfg config.AuthConfig) *AdminAuthHandler {
	return &AdminAuthHandler{
		service:   service,
		issuer:    issuer,
		verifier:  verifier,
		cookieCfg: newAdminSessionCookieOptions(authCfg.Admin),
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

	if strings.TrimSpace(result.Token) != "" {
		expires := time.Unix(result.ExpiresAt, 0)
		setAdminSessionCookie(c.Writer, result.Token, expires, h.cookieCfg)
	}

	c.Header("Cache-Control", "no-store")
	target := buildAdminRedirect(result.RedirectPath, result.Token)
	c.Redirect(http.StatusTemporaryRedirect, target)
}

// Session validates the current administrator session and refreshes JWT/Cookie when appropriate.
func (h *AdminAuthHandler) Session(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	token := extractBearerToken(authHeader)
	tokenSource := "header"
	if token == "" {
		if cookieName := strings.TrimSpace(h.cookieCfg.name); cookieName != "" {
			if cookie, err := c.Cookie(cookieName); err == nil {
				token = strings.TrimSpace(cookie)
				tokenSource = "cookie"
			}
		} else if cookie, err := c.Cookie("ps_admin_jwt"); err == nil {
			token = strings.TrimSpace(cookie)
			tokenSource = "cookie"
		}
	}

	if strings.TrimSpace(token) == "" {
		respondError(c, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "missing token", nil))
		return
	}

	if h.verifier == nil {
		respondError(c, errs.New(errs.CodeInternal, http.StatusInternalServerError, "token verifier not configured", nil))
		return
	}

	claims, err := h.verifier.Verify(ctx, token)
	if err != nil {
		respondError(c, err)
		return
	}
	if claims == nil {
		respondError(c, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid token", nil))
		return
	}

	refreshedToken := ""
	expiresAt := claims.ExpiresAt

	shouldRefresh := tokenSource != "header"
	if !shouldRefresh && !claims.ExpiresAt.IsZero() {
		const refreshWindow = 10 * time.Minute
		if time.Until(claims.ExpiresAt) <= refreshWindow {
			shouldRefresh = true
		}
	}

	if shouldRefresh {
		if h.issuer == nil {
			respondError(c, errs.New(errs.CodeInternal, http.StatusInternalServerError, "token issuer not configured", nil))
			return
		}
		issuedToken, issuedExpiry, issueErr := h.issuer.Issue(ctx, claims.Subject, claims.Email, claims.Roles...)
		if issueErr != nil {
			respondError(c, issueErr)
			return
		}
		refreshedToken = issuedToken
		expiresAt = issuedExpiry
		setAdminSessionCookie(c.Writer, issuedToken, issuedExpiry, h.cookieCfg)
	} else if tokenSource == "cookie" {
		// Provide the existing cookie token to the SPA when rehydrating without Authorization headers.
		refreshedToken = token
	}

	c.Header("Cache-Control", "no-store")

	response := gin.H{
		"data": gin.H{
			"active":    true,
			"expiresAt": expiresAt.Unix(),
		},
	}
	if refreshedToken != "" {
		response["data"].(gin.H)["token"] = refreshedToken
	}

	c.JSON(http.StatusOK, response)
}

func buildAdminRedirect(path, token string) string {
	stripped := strings.TrimSpace(path)
	if strings.Contains(stripped, "#") {
		parts := strings.SplitN(stripped, "#", 2)
		stripped = parts[0]
	}
	normalized := normalizeAdminRedirectLocation(stripped)
	fragment := "token=" + url.QueryEscape(token)
	return normalized + "#" + fragment
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

func setAdminSessionCookie(w http.ResponseWriter, token string, expires time.Time, opts adminSessionCookieOptions) {
	if strings.TrimSpace(token) == "" {
		return
	}

	name := strings.TrimSpace(opts.name)
	if name == "" {
		name = "ps_admin_jwt"
	}

	path := strings.TrimSpace(opts.path)
	if path == "" {
		path = "/"
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     path,
		HttpOnly: opts.httpOnly,
		Secure:   opts.secure,
		SameSite: opts.sameSite,
	}

	if domain := strings.TrimSpace(opts.domain); domain != "" {
		cookie.Domain = domain
	}

	if !expires.IsZero() && expires.After(time.Now()) {
		cookie.Expires = expires
		cookie.MaxAge = int(time.Until(expires).Seconds())
	}

	http.SetCookie(w, cookie)
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
