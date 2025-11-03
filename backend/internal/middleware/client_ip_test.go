package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResolveClientIPPrefersForwardedForWhenTrusted(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.1.2.3:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.10, 10.1.2.3")
	c.Request = req

	ip := resolveClientIP(c, nil)
	if ip != "198.51.100.10" {
		t.Fatalf("expected forwarded client IP, got %q", ip)
	}
}

func TestResolveClientIPFallsBackToRemoteWhenUntrusted(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	req.Header.Set("X-Forwarded-For", "198.51.100.10")
	c.Request = req

	ip := resolveClientIP(c, nil)
	if ip != "1.2.3.4" {
		t.Fatalf("expected remote IP fallback, got %q", ip)
	}
}

func TestResolveClientIPHonorsConfiguredTrustedProxies(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.5:8080"
	req.Header.Set("X-Forwarded-For", "198.51.100.99")
	c.Request = req

	trusted := newTrustedProxyList([]string{"203.0.113.0/24"})
	ip := resolveClientIP(c, trusted)
	if ip != "198.51.100.99" {
		t.Fatalf("expected forwarded IP for configured proxy, got %q", ip)
	}
}

func TestResolveClientIPRejectsInvalidForwardHeader(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.1.2.3:1234"
	req.Header.Set("X-Forwarded-For", "not-an-ip")
	c.Request = req

	ip := resolveClientIP(c, nil)
	if ip != "10.1.2.3" {
		t.Fatalf("expected remote IP when header invalid, got %q", ip)
	}
}
