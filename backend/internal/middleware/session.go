package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	authsvc "github.com/takumi/personal-website/internal/service/auth"
)

// ContextSessionKey is the request context key for authenticated administrator sessions.
const ContextSessionKey = "auth.session"

// AdminSessionMiddleware validates administrator sessions backed by the session manager.
type AdminSessionMiddleware struct {
	sessions      authsvc.AdminSessionManager
	cookie        authsvc.CookieOptions
	refreshWindow time.Duration
}

// NewAdminSessionMiddleware constructs the middleware.
func NewAdminSessionMiddleware(sessions authsvc.AdminSessionManager, cfg config.AuthConfig) *AdminSessionMiddleware {
	refreshWindow := cfg.Admin.SessionRefreshWindow
	if refreshWindow <= 0 {
		refreshWindow = 20 * time.Minute
	}
	return &AdminSessionMiddleware{
		sessions:      sessions,
		cookie:        authsvc.NewCookieOptions(cfg.Admin),
		refreshWindow: refreshWindow,
	}
}

// Handler returns the Gin middleware function.
func (m *AdminSessionMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, source := m.extractSessionID(c)
		if sessionID == "" {
			m.cookie.Clear(c.Writer)
			abortUnauthorized(c, "missing session token")
			return
		}

		session, err := m.sessions.Validate(c.Request.Context(), sessionID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				m.cookie.Clear(c.Writer)
				abortUnauthorized(c, "session invalid or expired")
				return
			}
			abortInternal(c, err)
			return
		}

		if time.Until(session.ExpiresAt) <= m.refreshWindow {
			refreshed, refreshErr := m.sessions.Refresh(c.Request.Context(), sessionID)
			if refreshErr != nil {
				if errors.Is(refreshErr, repository.ErrNotFound) {
					m.cookie.Clear(c.Writer)
					abortUnauthorized(c, "session expired")
					return
				}
				abortInternal(c, refreshErr)
				return
			}
			session = refreshed
		}

		// Ensure cookie expiry reflects latest session state.
		m.cookie.Write(c.Writer, session.ID, session.ExpiresAt.UTC())

		// Record session and the source for downstream handlers (e.g. auditing/logging).
		c.Set(ContextSessionKey, session)
		c.Set("auth.session.source", source)

		c.Next()
	}
}

func (m *AdminSessionMiddleware) extractSessionID(c *gin.Context) (string, string) {
	authHeader := c.GetHeader("Authorization")
	if token := extractBearerToken(authHeader); token != "" {
		return token, "header"
	}

	cookieName := strings.TrimSpace(m.cookie.Name)
	if cookieName == "" {
		cookieName = "ps_admin_session"
	}

	if cookie, err := c.Cookie(cookieName); err == nil {
		return strings.TrimSpace(cookie), "cookie"
	}

	// Legacy fallback to maintain compatibility with previous cookie names.
	if cookie, err := c.Cookie("ps_admin_jwt"); err == nil {
		return strings.TrimSpace(cookie), "cookie-legacy"
	}

	return "", ""
}

func abortUnauthorized(c *gin.Context, message string) {
	appErr := errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, message, nil)
	c.AbortWithStatusJSON(appErr.Status, gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	})
}

func abortInternal(c *gin.Context, err error) {
	appErr := errs.From(err)
	if appErr.Status < http.StatusInternalServerError {
		c.AbortWithStatusJSON(appErr.Status, gin.H{
			"error":   appErr.Code,
			"message": appErr.Message,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error":   "internal_error",
		"message": "unexpected authentication error",
	})
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// GetSessionFromContext extracts the authenticated session if present.
func GetSessionFromContext(c *gin.Context) (*model.AdminSession, bool) {
	sessionAny, exists := c.Get(ContextSessionKey)
	if !exists {
		return nil, false
	}
	session, ok := sessionAny.(*model.AdminSession)
	return session, ok && session != nil
}
