package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
)

// HTTPSRedirect ensures clients always communicate over TLS in production.
type HTTPSRedirect struct {
	enabled bool
}

// NewHTTPSRedirect returns the middleware instance.
func NewHTTPSRedirect(cfg *config.AppConfig) *HTTPSRedirect {
	if cfg == nil {
		return &HTTPSRedirect{}
	}
	return &HTTPSRedirect{enabled: cfg.Security.HTTPSRedirect}
}

// Handler returns the gin middleware handler.
func (h *HTTPSRedirect) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h == nil || !h.enabled {
			c.Next()
			return
		}

		if isSecure(c.Request) {
			c.Next()
			return
		}

		url := c.Request.URL
		url.Scheme = "https"
		url.Host = c.Request.Host

		c.Redirect(http.StatusPermanentRedirect, url.String())
		c.Abort()
	}
}

func isSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	return strings.EqualFold(proto, "https")
}
