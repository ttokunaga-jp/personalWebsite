package google

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"
)

// TokenProvider exposes a minimal interface to retrieve OAuth2 access tokens for Google APIs.
type TokenProvider interface {
	AccessToken(ctx context.Context) (string, error)
}

// EnvTokenProvider fetches an OAuth2 access token from an environment variable with basic caching.
type EnvTokenProvider struct {
	EnvVar string

	mu         sync.Mutex
	token      string
	expires    time.Time
	expiryHint time.Duration
}

// AccessToken returns the cached token or attempts to read a fresh value from the configured environment variable.
// The provider expects tokens to be refreshed externally and supports optional cache invalidation through
// detecting "exp=<unix>" suffixes encoded alongside the token.
func (p *EnvTokenProvider) AccessToken(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.EnvVar == "" {
		return "", ErrTokenNotFound
	}

	now := time.Now()
	if p.token != "" && (p.expires.IsZero() || now.Before(p.expires)) {
		return p.token, nil
	}

	raw := strings.TrimSpace(os.Getenv(p.EnvVar))
	if raw == "" {
		return "", ErrTokenNotFound
	}

	token, expires := splitTokenAndExpiry(raw)
	p.token = token
	p.expires = expires
	return token, nil
}

func splitTokenAndExpiry(raw string) (string, time.Time) {
	parts := strings.Split(raw, "|")
	if len(parts) == 1 {
		return raw, time.Time{}
	}

	token := strings.TrimSpace(parts[0])
	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "exp=") {
			if ts, err := time.Parse(time.RFC3339, strings.TrimPrefix(part, "exp=")); err == nil {
				return token, ts
			}
		}
	}
	return token, time.Time{}
}
