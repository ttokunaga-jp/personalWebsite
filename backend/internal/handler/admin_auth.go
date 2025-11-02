package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

// AdminAuthHandler exposes admin-specific OAuth login endpoints.
type AdminAuthHandler struct {
	service  authsvc.AdminService
	issuer   authsvc.TokenIssuer
	verifier authsvc.TokenVerifier
}

// NewAdminAuthHandler constructs the handler with the provided service.
func NewAdminAuthHandler(service authsvc.AdminService, issuer authsvc.TokenIssuer, verifier authsvc.TokenVerifier) *AdminAuthHandler {
	return &AdminAuthHandler{
		service:  service,
		issuer:   issuer,
		verifier: verifier,
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
		setAdminSessionCookie(c.Writer, result.Token, expires)
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
		if cookie, err := c.Cookie("ps_admin_jwt"); err == nil {
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
		setAdminSessionCookie(c.Writer, issuedToken, issuedExpiry)
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
	path = strings.TrimSpace(path)
	if path == "" {
		path = "/admin"
	}
	if strings.Contains(path, "#") {
		parts := strings.SplitN(path, "#", 2)
		path = parts[0]
	}
	fragment := "token=" + url.QueryEscape(token)
	return path + "#" + fragment
}

func setAdminSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	if strings.TrimSpace(token) == "" {
		return
	}

	cookie := &http.Cookie{
		Name:     "ps_admin_jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
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
