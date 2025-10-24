package google

import (
	"context"
	"errors"
	"time"
)

const GmailProvider = "gmail"

var (
	// ErrTokenNotFound is returned when no token record exists for the given provider.
	ErrTokenNotFound = errors.New("google token store: token not found")
)

// TokenRecord represents the OAuth access and refresh token pair persisted for reuse.
type TokenRecord struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// TokenStore abstracts persistence operations for OAuth tokens.
type TokenStore interface {
	Load(ctx context.Context, provider string) (*TokenRecord, error)
	Save(ctx context.Context, provider string, record *TokenRecord) error
}
