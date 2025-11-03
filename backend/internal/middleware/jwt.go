package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/service/auth"
)

const ContextClaimsKey = "auth.claims"

type JWTMiddleware struct {
	verifier   auth.TokenVerifier
	cookieName string
}

func NewJWTMiddleware(verifier auth.TokenVerifier, cfg config.AuthConfig) *JWTMiddleware {
	name := strings.TrimSpace(cfg.Admin.SessionCookieName)
	if name == "" {
		name = "ps_admin_jwt"
	}
	return &JWTMiddleware{
		verifier:   verifier,
		cookieName: name,
	}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := extractBearerToken(authHeader)
		if token == "" {
			if cookieName := strings.TrimSpace(m.cookieName); cookieName != "" {
				if cookie, err := c.Cookie(cookieName); err == nil {
					token = strings.TrimSpace(cookie)
				}
			}
			if token == "" {
				if cookie, err := c.Cookie("ps_admin_jwt"); err == nil {
					token = strings.TrimSpace(cookie)
				}
			}
		}
		if strings.TrimSpace(token) == "" {
			appErr := errs.New(errs.CodeUnauthorized, 401, "missing token", nil)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}
		claims, err := m.verifier.Verify(c.Request.Context(), token)
		if err != nil {
			appErr := errs.From(err)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		c.Set(ContextClaimsKey, claims)
		c.Next()
	}
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
