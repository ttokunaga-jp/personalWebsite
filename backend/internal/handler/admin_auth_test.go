package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

func TestAdminAuthCallbackSetsCookieAndRedirects(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	service := &stubAdminAuthService{
		callbackResult: &authsvc.AdminCallbackResult{
			SessionID:    "session-token",
			ExpiresAt:    time.Now().Add(30 * time.Minute).Unix(),
			RedirectPath: "/admin/dashboard",
			Email:        "admin@example.com",
			Roles:        []string{"admin"},
		},
	}
	handler := NewAdminAuthHandler(service, &stubSessionManager{}, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/callback?state=abc&code=xyz", nil)
	c.Request = req

	handler.Callback(c)

	require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	require.Equal(t, "/admin/dashboard", rec.Header().Get("Location"))
	require.Equal(t, "no-store", rec.Header().Get("Cache-Control"))

	res := rec.Result()
	var sessionCookie *http.Cookie
	for _, cookie := range res.Cookies() {
		if cookie.Name == "ps_admin_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "session cookie should be set")
	require.Equal(t, "session-token", sessionCookie.Value)
	require.True(t, sessionCookie.HttpOnly)
	require.True(t, sessionCookie.Secure)
}

func TestAdminAuthCallbackValidatesQueryParameters(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	handler := NewAdminAuthHandler(&stubAdminAuthService{}, &stubSessionManager{}, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/callback?state=&code=", nil)
	c.Request = req

	handler.Callback(c)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "missing state or code")
}

func TestAdminAuthCallbackPropagatesServiceErrors(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	service := &stubAdminAuthService{
		callbackErr: errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid oauth state", nil),
	}
	handler := NewAdminAuthHandler(service, &stubSessionManager{}, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/callback?state=abc&code=invalid", nil)
	c.Request = req

	handler.Callback(c)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid oauth state")
}

func TestAdminSessionRehydratesFromCookie(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	now := time.Now()
	validateSession := &model.AdminSession{
		ID:             "cookie-token",
		TokenHash:      "hash",
		Email:          "admin@example.com",
		Roles:          []string{"admin"},
		ExpiresAt:      now.Add(30 * time.Second),
		LastAccessedAt: now.Add(-time.Minute),
	}
	refreshedSession := &model.AdminSession{
		ID:             "cookie-token",
		TokenHash:      "hash",
		Email:          "admin@example.com",
		Roles:          []string{"admin"},
		ExpiresAt:      now.Add(2 * time.Hour),
		LastAccessedAt: now,
	}

	sessions := &stubSessionManager{
		validateFn: func(string) (*model.AdminSession, error) {
			return cloneSession(validateSession), nil
		},
		refreshFn: func(string) (*model.AdminSession, error) {
			return cloneSession(refreshedSession), nil
		},
	}

	handler := NewAdminAuthHandler(&stubAdminAuthService{}, sessions, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "ps_admin_session", Value: "cookie-token"})
	c.Request = req

	handler.Session(c)

	require.Equal(t, http.StatusOK, rec.Code)

	var payload struct {
		Active    bool     `json:"active"`
		ExpiresAt int64    `json:"expiresAt"`
		Roles     []string `json:"roles"`
		Email     string   `json:"email"`
		Refreshed bool     `json:"refreshed"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.True(t, payload.Active)
	require.NotZero(t, payload.ExpiresAt)
	require.Equal(t, refreshedSession.ExpiresAt.Unix(), payload.ExpiresAt)
	require.Equal(t, []string{"admin"}, payload.Roles)
	require.Equal(t, "admin@example.com", payload.Email)
	require.True(t, payload.Refreshed)

	res := rec.Result()
	var sessionCookie *http.Cookie
	for _, cookie := range res.Cookies() {
		if cookie.Name == "ps_admin_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie)
	require.Equal(t, "cookie-token", sessionCookie.Value)
	require.True(t, sessionCookie.HttpOnly)
	require.True(t, sessionCookie.Secure)
}

func testAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		Admin: config.AdminAuthConfig{
			SessionCookieName:     "ps_admin_session",
			SessionCookiePath:     "/",
			SessionCookieSecure:   true,
			SessionCookieHTTPOnly: true,
			SessionCookieSameSite: "strict",
			SessionTTL:            time.Hour,
			SessionRefreshWindow:  10 * time.Minute,
		},
	}
}

type stubAdminAuthService struct {
	loginResult    *authsvc.AdminLoginResult
	loginErr       error
	callbackResult *authsvc.AdminCallbackResult
	callbackErr    error
}

func (s *stubAdminAuthService) StartLogin(context.Context, string) (*authsvc.AdminLoginResult, error) {
	if s.loginResult != nil || s.loginErr != nil {
		return s.loginResult, s.loginErr
	}
	return &authsvc.AdminLoginResult{
		AuthURL: "https://example.com/oauth?state=seed",
		State:   "seed",
	}, nil
}

func (s *stubAdminAuthService) HandleCallback(context.Context, string, string) (*authsvc.AdminCallbackResult, error) {
	if s.callbackErr != nil || s.callbackResult != nil {
		return s.callbackResult, s.callbackErr
	}
	return &authsvc.AdminCallbackResult{
		SessionID:    "default-session",
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
		RedirectPath: "/admin/",
		Email:        "admin@example.com",
		Roles:        []string{"admin"},
	}, nil
}

type stubSessionManager struct {
	createFn   func(authsvc.AdminPrincipal) (*model.AdminSession, error)
	validateFn func(string) (*model.AdminSession, error)
	refreshFn  func(string) (*model.AdminSession, error)
	revokeFn   func(string) error
}

func (s *stubSessionManager) Create(_ context.Context, principal authsvc.AdminPrincipal) (*model.AdminSession, error) {
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

func (s *stubSessionManager) Validate(_ context.Context, sessionID string) (*model.AdminSession, error) {
	if s.validateFn != nil {
		return s.validateFn(sessionID)
	}
	return &model.AdminSession{
		ID:        sessionID,
		TokenHash: "hash",
		Email:     "admin@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}, nil
}

func (s *stubSessionManager) Refresh(_ context.Context, sessionID string) (*model.AdminSession, error) {
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

func (s *stubSessionManager) Revoke(_ context.Context, sessionID string) error {
	if s.revokeFn != nil {
		return s.revokeFn(sessionID)
	}
	return nil
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
