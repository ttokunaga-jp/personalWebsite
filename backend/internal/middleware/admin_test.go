package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/service/auth"
)

func TestAdminGuardRequireAdmin(t *testing.T) {
	t.Parallel()

	guard := NewAdminGuard()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextClaimsKey, &auth.Claims{Roles: []string{"admin"}})
	})
	router.Use(guard.RequireAdmin())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/secure", nil)
	require.NoError(t, err)

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminGuardRequireAdminForbidden(t *testing.T) {
	t.Parallel()

	guard := NewAdminGuard()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextClaimsKey, &auth.Claims{Roles: []string{"viewer"}})
	})
	router.Use(guard.RequireAdmin())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/secure", nil)
	require.NoError(t, err)

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminGuardRequireAdminMissingClaims(t *testing.T) {
	t.Parallel()

	guard := NewAdminGuard()
	router := gin.New()
	router.Use(guard.RequireAdmin())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/secure", nil)
	require.NoError(t, err)

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
