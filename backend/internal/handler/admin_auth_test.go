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
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

func TestAdminAuthCallbackSetsCookieAndRedirects(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	service := &stubAdminAuthService{
		callbackResult: &authsvc.AdminCallbackResult{
			Token:        "jwt-token",
			ExpiresAt:    time.Now().Add(30 * time.Minute).Unix(),
			RedirectPath: "/admin/dashboard",
		},
	}
	handler := NewAdminAuthHandler(service, nil, nil, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/callback?state=abc&code=xyz", nil)
	c.Request = req

	handler.Callback(c)

	require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	require.Equal(t, "/admin/dashboard#token=jwt-token", rec.Header().Get("Location"))
	require.Equal(t, "no-store", rec.Header().Get("Cache-Control"))

	res := rec.Result()
	var jwtCookie *http.Cookie
	for _, cookie := range res.Cookies() {
		if cookie.Name == "ps_admin_jwt" {
			jwtCookie = cookie
			break
		}
	}
	require.NotNil(t, jwtCookie, "ps_admin_jwt cookie should be set")
	require.Equal(t, "jwt-token", jwtCookie.Value)
	require.True(t, jwtCookie.HttpOnly)
	require.True(t, jwtCookie.Secure)
}

func TestAdminAuthCallbackValidatesQueryParameters(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	handler := NewAdminAuthHandler(&stubAdminAuthService{}, nil, nil, testAuthConfig())

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
	handler := NewAdminAuthHandler(service, nil, nil, testAuthConfig())

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

	verifier := &stubTokenVerifier{
		expectedToken: "cookie-token",
		claims: &authsvc.Claims{
			Subject:   "admin-subject",
			Email:     "admin@example.com",
			Roles:     []string{"admin"},
			ExpiresAt: time.Now().Add(5 * time.Minute),
		},
	}
	issuer := &stubTokenIssuer{
		issuedToken: "refreshed-token",
		expiresAt:   time.Now().Add(45 * time.Minute),
	}

	handler := NewAdminAuthHandler(&stubAdminAuthService{}, issuer, verifier, testAuthConfig())

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "ps_admin_jwt", Value: "cookie-token"})
	c.Request = req

	handler.Session(c)

	require.Equal(t, http.StatusOK, rec.Code)

	var payload struct {
		Data struct {
			Active    bool   `json:"active"`
			Token     string `json:"token"`
			ExpiresAt int64  `json:"expiresAt"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	require.True(t, payload.Data.Active)
	require.Equal(t, "refreshed-token", payload.Data.Token)
	require.NotZero(t, payload.Data.ExpiresAt)

	res := rec.Result()
	var sessionCookie *http.Cookie
	for _, cookie := range res.Cookies() {
		if cookie.Name == "ps_admin_jwt" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie)
	require.Equal(t, "refreshed-token", sessionCookie.Value)
}

func testAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		Admin: config.AdminAuthConfig{
			SessionCookieName:     "ps_admin_jwt",
			SessionCookiePath:     "/",
			SessionCookieSecure:   true,
			SessionCookieHTTPOnly: true,
			SessionCookieSameSite: "lax",
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
		Token:        "default-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
		RedirectPath: "/admin/",
	}, nil
}

type stubTokenIssuer struct {
	issuedToken string
	expiresAt   time.Time
	err         error
}

func (s *stubTokenIssuer) Issue(_ context.Context, subject, email string, roles ...string) (string, time.Time, error) {
	_ = subject
	_ = email
	_ = roles
	if s.err != nil {
		return "", time.Time{}, s.err
	}
	return s.issuedToken, s.expiresAt, nil
}

type stubTokenVerifier struct {
	claims        *authsvc.Claims
	err           error
	expectedToken string
}

func (s *stubTokenVerifier) Verify(_ context.Context, token string) (*authsvc.Claims, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.expectedToken != "" && token != s.expectedToken {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "unexpected token", nil)
	}
	return s.claims, nil
}
