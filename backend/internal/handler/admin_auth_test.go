package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

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
	handler := NewAdminAuthHandler(service)

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

	handler := NewAdminAuthHandler(&stubAdminAuthService{})

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
	handler := NewAdminAuthHandler(service)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/callback?state=abc&code=invalid", nil)
	c.Request = req

	handler.Callback(c)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid oauth state")
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
		RedirectPath: "/admin",
	}, nil
}
