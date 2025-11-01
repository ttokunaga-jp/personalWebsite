package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/service/auth"
)

func TestJWTMiddlewareUsesAuthorizationHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	verifier := &stubTokenVerifier{
		claims: &auth.Claims{Roles: []string{"admin"}},
	}
	mw := NewJWTMiddleware(verifier)

	engine := gin.New()
	engine.Use(mw.Handler())
	engine.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "header-token", verifier.lastToken)
}

func TestJWTMiddlewareFallsBackToCookie(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	verifier := &stubTokenVerifier{
		claims: &auth.Claims{Roles: []string{"admin"}},
	}
	mw := NewJWTMiddleware(verifier)

	engine := gin.New()
	engine.Use(mw.Handler())
	engine.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.AddCookie(&http.Cookie{Name: "ps_admin_jwt", Value: "cookie-token"})
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "cookie-token", verifier.lastToken)
}

func TestAdminGuardRequiresAdminRole(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	guard := NewAdminGuard()
	engine := gin.New()
	engine.GET("/admin", func(c *gin.Context) {
		c.Set(ContextClaimsKey, &auth.Claims{Roles: []string{"viewer"}})
	}, guard.RequireAdmin(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "admin role required")
}

func TestAdminGuardAllowsAdminRole(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	guard := NewAdminGuard()
	engine := gin.New()
	engine.GET("/admin", func(c *gin.Context) {
		c.Set(ContextClaimsKey, &auth.Claims{Roles: []string{"admin"}})
	}, guard.RequireAdmin(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

type stubTokenVerifier struct {
	claims    *auth.Claims
	err       error
	lastToken string
}

func (s *stubTokenVerifier) Verify(_ context.Context, token string) (*auth.Claims, error) {
	s.lastToken = token
	if s.err != nil {
		return nil, s.err
	}
	if s.claims == nil {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "missing claims", nil)
	}
	return s.claims, nil
}
