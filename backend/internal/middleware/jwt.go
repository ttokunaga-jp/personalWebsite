package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/service/auth"
)

const ContextClaimsKey = "auth.claims"

type JWTMiddleware struct {
	verifier auth.TokenVerifier
}

func NewJWTMiddleware(verifier auth.TokenVerifier) *JWTMiddleware {
	return &JWTMiddleware{
		verifier: verifier,
	}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := extractBearerToken(authHeader)
		if token == "" {
			if cookie, err := c.Cookie("ps_admin_jwt"); err == nil {
				token = strings.TrimSpace(cookie)
			}
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
