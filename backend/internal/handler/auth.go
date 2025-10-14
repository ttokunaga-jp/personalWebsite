package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

// AuthHandler exposes endpoints for Google OAuth + JWT login flow.
type AuthHandler struct {
	service authsvc.Service
}

func NewAuthHandler(service authsvc.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login initiates the OAuth flow by redirecting to Google's consent screen.
func (h *AuthHandler) Login(c *gin.Context) {
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

// Callback finalizes the OAuth flow, exchanging the code for a JWT.
func (h *AuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if state == "" || code == "" {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "missing state or code", nil))
		return
	}

	result, err := h.service.HandleCallback(c.Request.Context(), state, code)
	if err != nil {
		respondError(c, err)
		return
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":        result.Token,
			"expires_at":   result.ExpiresAt,
			"redirect_uri": result.RedirectURI,
		},
	})
}
