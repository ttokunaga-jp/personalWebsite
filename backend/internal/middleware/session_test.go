package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/model"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

func TestAdminSessionMiddlewareAllowsValidSession(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	session := &model.AdminSession{
		ID:        "session-token",
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	sessions := &middlewareSessionStub{
		validateFn: func(string) (*model.AdminSession, error) {
			return cloneSession(session), nil
		},
	}

	mw := NewAdminSessionMiddleware(sessions, testMiddlewareAuthConfig())

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/secure", func(c *gin.Context) {
		value, ok := GetSessionFromContext(c)
		require.True(t, ok)
		require.Equal(t, "admin@example.com", value.Email)
		c.Status(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.AddCookie(&http.Cookie{Name: "ps_admin_session", Value: "session-token"})

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)

	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
}

func TestAdminSessionMiddlewareRefreshesNearExpiry(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	now := time.Now()
	validateSession := &model.AdminSession{
		ID:        "session-token",
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: now.Add(15 * time.Second),
	}
	refreshed := &model.AdminSession{
		ID:        "session-token",
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: now.Add(time.Hour),
	}

	sessions := &middlewareSessionStub{
		validateFn: func(string) (*model.AdminSession, error) {
			return cloneSession(validateSession), nil
		},
		refreshFn: func(string) (*model.AdminSession, error) {
			return cloneSession(refreshed), nil
		},
	}

	cfg := testMiddlewareAuthConfig()
	cfg.Admin.SessionRefreshWindow = 30 * time.Second
	mw := NewAdminSessionMiddleware(sessions, cfg)

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.AddCookie(&http.Cookie{Name: "ps_admin_session", Value: "session-token"})

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotEmpty(t, rec.Result().Cookies())
}

func TestAdminSessionMiddlewareRejectsMissingSession(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	mw := NewAdminSessionMiddleware(&middlewareSessionStub{}, testMiddlewareAuthConfig())

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

type middlewareSessionStub struct {
	createFn   func(authsvc.AdminPrincipal) (*model.AdminSession, error)
	validateFn func(string) (*model.AdminSession, error)
	refreshFn  func(string) (*model.AdminSession, error)
	revokeFn   func(string) error
}

func (s *middlewareSessionStub) Create(_ context.Context, principal authsvc.AdminPrincipal) (*model.AdminSession, error) {
	if s.createFn != nil {
		return s.createFn(principal)
	}
	return &model.AdminSession{
		ID:        "stub-session",
		TokenHash: "hash",
		Email:     principal.Email,
		Roles:     principal.Roles,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (s *middlewareSessionStub) Validate(_ context.Context, sessionID string) (*model.AdminSession, error) {
	if s.validateFn != nil {
		return s.validateFn(sessionID)
	}
	return &model.AdminSession{
		ID:        sessionID,
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (s *middlewareSessionStub) Refresh(_ context.Context, sessionID string) (*model.AdminSession, error) {
	if s.refreshFn != nil {
		return s.refreshFn(sessionID)
	}
	return &model.AdminSession{
		ID:        sessionID,
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}, nil
}

func (s *middlewareSessionStub) Revoke(_ context.Context, sessionID string) error {
	if s.revokeFn != nil {
		return s.revokeFn(sessionID)
	}
	return nil
}

func testMiddlewareAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		Admin: config.AdminAuthConfig{
			SessionCookieName:     "ps_admin_session",
			SessionCookiePath:     "/",
			SessionCookieSecure:   true,
			SessionCookieHTTPOnly: true,
			SessionCookieSameSite: "strict",
			SessionTTL:            time.Hour,
			SessionRefreshWindow:  20 * time.Second,
		},
	}
}

func cloneSession(src *model.AdminSession) *model.AdminSession {
	if src == nil {
		return nil
	}
	dest := *src
	if src.Roles != nil {
		dest.Roles = append([]string(nil), src.Roles...)
	}
	return &dest
}
