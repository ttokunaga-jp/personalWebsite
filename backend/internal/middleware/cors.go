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
	allowedOrigins := append([]string{}, security.AllowedOrigins...)
	if adminOrigin := extractOrigin(cfg.Auth.Admin.DefaultRedirectURI); adminOrigin != "" {
		if !containsString(allowedOrigins, adminOrigin) {
			allowedOrigins = append(allowedOrigins, adminOrigin)
		}
	}
	if len(allowedOrigins) == 0 {
		return &CORSMiddleware{enabled: false}
	}

	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", security.CSRFHeaderName, "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: security.AllowCredentials,
		MaxAge:           12 * time.Hour,
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

func containsString(list []string, target string) bool {
	if target == "" {
		return false
	}
	for _, item := range list {
		if strings.EqualFold(strings.TrimSpace(item), target) {
			return true
		}
	}
	return false
}
