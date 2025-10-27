package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
)

// SecurityHandler exposes endpoints for security utilities (e.g. CSRF token issuance).
type SecurityHandler struct {
	manager *csrfmgr.Manager
	cfg     *config.AppConfig
}

// NewSecurityHandler creates the handler.
func NewSecurityHandler(manager *csrfmgr.Manager, cfg *config.AppConfig) *SecurityHandler {
	return &SecurityHandler{manager: manager, cfg: cfg}
}

// IssueCSRFToken sets the CSRF cookie and returns the token value for client headers.
func (h *SecurityHandler) IssueCSRFToken(c *gin.Context) {
	if h == nil || h.manager == nil || h.cfg == nil {
		respondError(c, errs.New(errs.CodeInternal, http.StatusInternalServerError, "csrf disabled", nil))
		return
	}

	token, err := h.manager.Issue()
	if err != nil {
		respondError(c, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to issue csrf token", err))
		return
	}

	cookie := &http.Cookie{
		Name:     h.cfg.Security.CSRFCookieName,
		Value:    token.Cookie,
		Path:     "/",
		Expires:  token.ExpiresAt,
		HttpOnly: h.cfg.Security.CSRFCookieHTTPOnly,
		Secure:   h.cfg.Security.CSRFCookieSecure,
		SameSite: mapSameSite(h.cfg.Security.CSRFCookieSameSite),
	}
	if domain := strings.TrimSpace(h.cfg.Security.CSRFCookieDomain); domain != "" {
		cookie.Domain = domain
	}

	http.SetCookie(c.Writer, cookie)

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":      token.Value,
			"expires_at": token.ExpiresAt.Format(time.RFC3339),
		},
	})
}

func mapSameSite(mode string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteDefaultMode
	}
}
