package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/config"
)

// CookieOptions encapsulates the configuration for the administrator session cookie.
type CookieOptions struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// NewCookieOptions derives cookie options from the admin auth configuration with safe defaults.
func NewCookieOptions(cfg config.AdminAuthConfig) CookieOptions {
	name := strings.TrimSpace(cfg.SessionCookieName)
	if name == "" {
		name = "ps_admin_session"
	}
	path := strings.TrimSpace(cfg.SessionCookiePath)
	if path == "" {
		path = "/"
	}
	return CookieOptions{
		Name:     name,
		Domain:   strings.TrimSpace(cfg.SessionCookieDomain),
		Path:     path,
		Secure:   cfg.SessionCookieSecure,
		HTTPOnly: cfg.SessionCookieHTTPOnly,
		SameSite: parseCookieSameSite(cfg.SessionCookieSameSite),
	}
}

func parseCookieSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteStrictMode
	}
}

// Write issues the session cookie with the provided value and expiry.
func (o CookieOptions) Write(w http.ResponseWriter, sessionID string, expires time.Time) {
	if strings.TrimSpace(sessionID) == "" {
		return
	}

	cookie := &http.Cookie{
		Name:     valueOrFallback(o.Name, "ps_admin_session"),
		Value:    sessionID,
		Path:     valueOrFallback(o.Path, "/"),
		HttpOnly: o.HTTPOnly,
		Secure:   o.Secure,
		SameSite: o.SameSite,
	}

	if strings.TrimSpace(o.Domain) != "" {
		cookie.Domain = o.Domain
	}
	if !expires.IsZero() && expires.After(time.Now()) {
		expiry := expires
		cookie.Expires = expiry
		cookie.MaxAge = int(time.Until(expiry).Seconds())
	}

	http.SetCookie(w, cookie)
}

// Clear invalidates the session cookie on the client.
func (o CookieOptions) Clear(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     valueOrFallback(o.Name, "ps_admin_session"),
		Value:    "",
		Path:     valueOrFallback(o.Path, "/"),
		HttpOnly: o.HTTPOnly,
		Secure:   o.Secure,
		SameSite: o.SameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	if strings.TrimSpace(o.Domain) != "" {
		cookie.Domain = o.Domain
	}

	http.SetCookie(w, cookie)
}

func valueOrFallback(value string, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}
