package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
)

func TestRateLimiterBlocksExcessRequests(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	cfg := &config.AppConfig{}
	cfg.Security.RateLimitRequestsPerMinute = 1
	cfg.Security.RateLimitBurst = 1

	rl := newRateLimiter(nil, cfg)
	engine.Use(rl.Handler())
	engine.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/", nil)
	engine.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected first request to be allowed, got %d", rec1.Code)
	}

	rec2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/", nil)
	engine.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request to be rate limited, got %d", rec2.Code)
	}
}
