package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	headerContentType = "Content-Type"
	utf8Charset       = "charset=UTF-8"
)

// Charset injects the UTF-8 charset into the Content-Type header.
type Charset struct{}

// NewCharsetMiddleware creates a new Charset middleware.
func NewCharsetMiddleware() *Charset {
	return &Charset{}
}

// Handler returns the middleware function.
func (m *Charset) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only modify headers for successful responses with a content type
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 || c.Writer.Header().Get(headerContentType) == "" {
			return
		}

		contentType := c.Writer.Header().Get(headerContentType)
		if !strings.Contains(strings.ToLower(contentType), "charset=") {
			c.Writer.Header().Set(headerContentType, fmt.Sprintf("%s; %s", contentType, utf8Charset))
		}
	}
}
