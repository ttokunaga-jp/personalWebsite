package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

// AdminAuthHandler exposes admin-specific OAuth login endpoints.
type AdminAuthHandler struct {
	service authsvc.AdminService
}

// NewAdminAuthHandler constructs the handler with the provided service.
func NewAdminAuthHandler(service authsvc.AdminService) *AdminAuthHandler {
	return &AdminAuthHandler{service: service}
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
		cookie := &http.Cookie{
			Name:     "ps_admin_jwt",
			Value:    result.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		if result.ExpiresAt > 0 {
			expires := time.Unix(result.ExpiresAt, 0)
			if expires.After(time.Now()) {
				cookie.Expires = expires
				cookie.MaxAge = int(time.Until(expires).Seconds())
			}
		}
		http.SetCookie(c.Writer, cookie)
	}

	c.Header("Cache-Control", "no-store")
	target := buildAdminRedirect(result.RedirectPath, result.Token)
	c.Redirect(http.StatusTemporaryRedirect, target)
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
