package google

import (
	"context"
	"errors"
	"fmt"
)

type fallbackTokenProvider struct {
	providers []TokenProvider
}

// NewFallbackTokenProvider tries the supplied providers in order until one succeeds.
func NewFallbackTokenProvider(providers ...TokenProvider) TokenProvider {
	filtered := make([]TokenProvider, 0, len(providers))
	for _, p := range providers {
		if p != nil {
			filtered = append(filtered, p)
		}
	}
	return &fallbackTokenProvider{providers: filtered}
}

func (f *fallbackTokenProvider) AccessToken(ctx context.Context) (string, error) {
	var lastErr error
	for _, provider := range f.providers {
		token, err := provider.AccessToken(ctx)
		if err == nil {
			return token, nil
		}
		if errors.Is(err, ErrTokenNotFound) {
			lastErr = err
			continue
		}
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("fallback token provider: no providers configured")
}
