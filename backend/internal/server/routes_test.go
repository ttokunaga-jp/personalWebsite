package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service"
	"github.com/takumi/personal-website/internal/service/auth"
)

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	engine := gin.New()

	profileSvc := service.NewProfileService(inmemory.NewProfileRepository())
	projectSvc := service.NewProjectService(inmemory.NewProjectRepository())
	researchSvc := service.NewResearchService(inmemory.NewResearchRepository())
	contactSvc := service.NewContactService(inmemory.NewContactRepository())

	jwtVerifier := auth.NewJWTVerifier(config.AuthConfig{
		JWTSecret:        "test-secret",
		Issuer:           "personal-website",
		Audience:         []string{"personal-website-admin"},
		ClockSkewSeconds: 30,
		Disabled:         true,
	})
	jwtMiddleware := middleware.NewJWTMiddleware(jwtVerifier)

	registerRoutes(
		engine,
		handler.NewHealthHandler(),
		handler.NewProfileHandler(profileSvc),
		handler.NewProjectHandler(projectSvc),
		handler.NewResearchHandler(researchSvc),
		handler.NewContactHandler(contactSvc),
		handler.NewAuthHandler(&stubAuthService{}),
		jwtMiddleware,
	)

	t.Run("health route ok", func(t *testing.T) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"status":"ok"`)
	})

	t.Run("profile route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/profile", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("projects route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/projects", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("research route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/research", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("contact route accepts payload", func(t *testing.T) {
		t.Helper()
		body, err := json.Marshal(model.ContactRequest{
			Name:    "Ada Lovelace",
			Email:   "ada@example.com",
			Message: "I'd like to learn more about your research.",
		})
		require.NoError(t, err)

		rec := performRequest(engine, http.MethodPost, "/api/contact", body)
		require.Equal(t, http.StatusAccepted, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("auth login route responds", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/auth/login", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"state"`)
	})

	t.Run("auth callback route responds", func(t *testing.T) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/callback?state=stub&code=ok", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"token"`)
	})
}

func performRequest(engine *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

type stubAuthService struct{}

func (s *stubAuthService) StartLogin(context.Context, string) (*auth.LoginResult, error) {
	return &auth.LoginResult{
		AuthURL: "",
		State:   "stub-state",
	}, nil
}

func (s *stubAuthService) HandleCallback(context.Context, string, string) (*auth.CallbackResult, error) {
	return &auth.CallbackResult{
		Token:       "stub-token",
		ExpiresAt:   123,
		RedirectURI: "/admin",
	}, nil
}
