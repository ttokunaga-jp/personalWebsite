package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
)

func TestCSRFMiddleware(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	appCfg := &config.AppConfig{}
	appCfg.Security.EnableCSRF = true
	appCfg.Security.CSRFSigningKey = "test-secret"
	appCfg.Security.CSRFTokenTTL = time.Minute
	appCfg.Security.CSRFCookieName = "csrf"
	appCfg.Security.CSRFHeaderName = "X-CSRF-Token"

	manager := csrfmgr.NewManager(appCfg.Security.CSRFSigningKey, appCfg.Security.CSRFTokenTTL)
	mw := NewCSRFMiddleware(appCfg, manager)
	engine.Use(mw.Handler())
	engine.POST("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Missing token should be rejected.
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/secure", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden when token missing, got %d", rec.Code)
	}

	// Valid token should succeed.
	token, err := manager.Issue()
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	rec2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/secure", nil)
	req2.Header.Set("X-Requested-With", "XMLHttpRequest")
	req2.Header.Set(appCfg.Security.CSRFHeaderName, token.Value)
	req2.AddCookie(&http.Cookie{Name: appCfg.Security.CSRFCookieName, Value: token.Cookie})
	engine.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected success with valid token, got %d", rec2.Code)
	}
}
