package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
)

// CORSMiddleware configures allowed origins for API access.
type CORSMiddleware struct {
	handler gin.HandlerFunc
	enabled bool
}

// NewCORSMiddleware builds the middleware using application configuration.
func NewCORSMiddleware(cfg *config.AppConfig) *CORSMiddleware {
	if cfg == nil {
		return &CORSMiddleware{}
	}
	security := cfg.Security
	primaryOrigin := normalizeOrigin(security.PrimaryOrigin)
	if primaryOrigin == "" {
		for _, candidate := range security.AllowedOrigins {
			if origin := normalizeOrigin(candidate); origin != "" {
				primaryOrigin = origin
				break
			}
		}
	}
	if primaryOrigin == "" {
		primaryOrigin = normalizeOrigin(cfg.Auth.Admin.DefaultRedirectURI)
	}
	if primaryOrigin == "" {
		return &CORSMiddleware{enabled: false}
	}

	config := cors.Config{
		AllowOrigins:     []string{primaryOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", security.CSRFHeaderName, "X-Request-ID", "X-Requested-With"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: security.AllowCredentials,
		MaxAge:           10 * time.Minute,
	}

	return &CORSMiddleware{
		handler: cors.New(config),
		enabled: true,
	}
}

// Handler returns the gin middleware handler.
func (m *CORSMiddleware) Handler() gin.HandlerFunc {
	if m == nil || !m.enabled || m.handler == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return m.handler
}

func extractOrigin(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func normalizeOrigin(raw string) string {
	origin := strings.TrimSpace(raw)
	if origin == "" {
		return ""
	}
	// Allow values with explicit scheme and host.
	if parsed := extractOrigin(origin); parsed != "" {
		return parsed
	}
	// If the value already looks like scheme://host we would have returned above.
	// Fall back to treating a bare host as https.
	if strings.HasPrefix(origin, "http://") || strings.HasPrefix(origin, "https://") {
		return origin
	}
	u := &url.URL{
		Scheme: "https",
		Host:   origin,
	}
	return u.Scheme + "://" + u.Host
}
