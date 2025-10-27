package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
)

// CSRFMiddleware validates double-submit tokens for state changing requests.
type CSRFMiddleware struct {
	manager      *csrfmgr.Manager
	cfg          config.SecurityConfig
	exemptPrefix []string
}

// NewCSRFMiddleware constructs a middleware backed by the CSRF manager.
func NewCSRFMiddleware(cfg *config.AppConfig, manager *csrfmgr.Manager) *CSRFMiddleware {
	if cfg == nil {
		return &CSRFMiddleware{}
	}
	return &CSRFMiddleware{
		manager:      manager,
		cfg:          cfg.Security,
		exemptPrefix: cfg.Security.CSRFExemptPaths,
	}
}

// Handler returns the gin middleware handler.
func (m *CSRFMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || !m.cfg.EnableCSRF || m.manager == nil {
			c.Next()
			return
		}

		if !isStateChangingMethod(c.Request.Method) {
			c.Next()
			return
		}

		if m.isExemptPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		cookie, err := c.Cookie(m.cfg.CSRFCookieName)
		if err != nil {
			m.reject(c, "missing csrf cookie")
			return
		}

		header := c.GetHeader(m.cfg.CSRFHeaderName)
		if err := m.manager.Validate(cookie, header); err != nil {
			reason := "invalid csrf token"
			if err == csrfmgr.ErrExpiredToken {
				reason = "expired csrf token"
			}
			m.reject(c, reason)
			return
		}

		c.Next()
	}
}

func (m *CSRFMiddleware) isExemptPath(path string) bool {
	for _, prefix := range m.exemptPrefix {
		if prefix == "" {
			continue
		}
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func (m *CSRFMiddleware) reject(c *gin.Context, message string) {
	err := errs.New(errs.CodeUnauthorized, http.StatusForbidden, message, nil)
	respondJSONError(c, err)
	c.Abort()
}

func isStateChangingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func respondJSONError(c *gin.Context, appErr *errs.AppError) {
	payload := gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	}
	if requestID := c.Writer.Header().Get("X-Request-ID"); requestID != "" {
		payload["request_id"] = requestID
	}
	c.JSON(appErr.Status, payload)
}
