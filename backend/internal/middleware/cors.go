package middleware

import (
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
	if len(security.AllowedOrigins) == 0 {
		return &CORSMiddleware{enabled: false}
	}

	config := cors.Config{
		AllowOrigins:     security.AllowedOrigins,
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
