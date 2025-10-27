package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
)

// SecurityHeaders injects defensive HTTP headers.
type SecurityHeaders struct {
	cfg config.SecurityConfig
}

// NewSecurityHeaders initialises the middleware from configuration.
func NewSecurityHeaders(cfg *config.AppConfig) *SecurityHeaders {
	if cfg == nil {
		return &SecurityHeaders{}
	}
	return &SecurityHeaders{cfg: cfg.Security}
}

// Handler returns the middleware function.
func (s *SecurityHeaders) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s == nil {
			c.Next()
			return
		}

		if s.cfg.ContentSecurityPolicy != "" {
			c.Writer.Header().Set("Content-Security-Policy", s.cfg.ContentSecurityPolicy)
		}
		if s.cfg.ReferrerPolicy != "" {
			c.Writer.Header().Set("Referrer-Policy", s.cfg.ReferrerPolicy)
		}
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), fullscreen=(self)")

		if s.cfg.HSTSMaxAgeSeconds > 0 && isHTTPSRequest(c.Request) {
			directives := []string{"max-age=" + strconv.Itoa(s.cfg.HSTSMaxAgeSeconds), "includeSubDomains", "preload"}
			c.Writer.Header().Set("Strict-Transport-Security", strings.Join(directives, "; "))
		}

		c.Next()
	}
}

func isHTTPSRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	return strings.EqualFold(proto, "https")
}
