package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
)

// AdminGuard ensures routes are only accessible by users with the admin role.
type AdminGuard struct{}

// NewAdminGuard constructs an authorization helper for admin-only routes.
func NewAdminGuard() *AdminGuard {
	return &AdminGuard{}
}

// RequireAdmin validates JWT claims contain the admin role.
func (g *AdminGuard) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := GetSessionFromContext(c)
		if !ok {
			appErr := errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "authentication required", nil)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		if !hasAdminRole(session.Roles) {
			appErr := errs.New(errs.CodeForbidden, http.StatusForbidden, "admin role required", nil)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		c.Next()
	}
}

func hasAdminRole(roles []string) bool {
	for _, role := range roles {
		if strings.EqualFold(strings.TrimSpace(role), "admin") {
			return true
		}
	}
	return false
}
