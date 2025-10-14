package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/handler"
)

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	engine := gin.New()
	registerRoutes(engine, handler.NewHealthHandler())

	req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"status":"ok"`)
}
