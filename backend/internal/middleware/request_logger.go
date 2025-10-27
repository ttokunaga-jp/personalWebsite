package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger emits structured logs for each HTTP request.
type RequestLogger struct {
	logger *slog.Logger
}

// NewRequestLogger constructs the middleware using the shared slog logger.
func NewRequestLogger(logger *slog.Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// Handler returns the gin middleware implementation.
func (m *RequestLogger) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || m.logger == nil {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		status := c.Writer.Status()
		reqID := RequestIDFromContext(c)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.String("latency", duration.String()),
		}

		if reqID != "" {
			attrs = append(attrs, slog.String("request_id", reqID))
		}

		level := slog.LevelInfo
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
			level = slog.LevelError
		}

		m.logger.LogAttrs(c.Request.Context(), level, "http_request", attrs...)
	}
}
