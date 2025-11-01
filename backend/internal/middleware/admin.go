package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/service/auth"
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
		claimsAny, exists := c.Get(ContextClaimsKey)
		if !exists {
			appErr := errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "authentication required", nil)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		claims, ok := claimsAny.(*auth.Claims)
		if !ok || claims == nil {
			appErr := errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid authentication context", nil)
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		if !claims.HasRole("admin") {
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
