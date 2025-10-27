package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

// RequestID assigns a stable identifier to each request for traceability.
type RequestID struct{}

// NewRequestID creates a new RequestID middleware.
func NewRequestID() *RequestID {
	return &RequestID{}
}

// Handler returns the gin middleware function.
func (r *RequestID) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestIDFromContext retrieves the request ID if present.
func RequestIDFromContext(c *gin.Context) string {
	if value, exists := c.Get(requestIDKey); exists {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}
