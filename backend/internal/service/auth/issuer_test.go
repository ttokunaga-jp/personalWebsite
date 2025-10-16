package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
)

func TestJWTIssuerIssueSuccess(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "test-secret",
			Issuer:                "personal-website",
			Audience:              []string{"admin"},
			AccessTokenTTLMinutes: 30,
		},
	}

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)

	jwtIss := issuer.(*jwtIssuer)
	fixed := time.Now().Truncate(time.Second)
	jwtIss.now = func() time.Time { return fixed }

	token, expiresAt, err := issuer.Issue(context.Background(), "google-123", "user@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, fixed.Add(30*time.Minute), expiresAt)

	verifier := NewJWTVerifier(cfg.Auth)
	claims, err := verifier.Verify(context.Background(), token)
	require.NoError(t, err)
	require.Equal(t, "google-123", claims.Subject)
	require.Equal(t, "user@example.com", claims.Email)
	require.True(t, claims.HasRole("admin"))
}

func TestJWTIssuerRequiresSubject(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "test-secret",
			AccessTokenTTLMinutes: 15,
		},
	}

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)

	_, _, err = issuer.Issue(context.Background(), "", "")
	require.Error(t, err)
}

func TestJWTIssuerDisabled(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "test-secret",
			AccessTokenTTLMinutes: 60,
			Disabled:              true,
		},
	}

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)

	token, expiresAt, err := issuer.Issue(context.Background(), "ignored", "")
	require.NoError(t, err)
	require.Equal(t, "development-token", token)
	require.True(t, expiresAt.After(time.Now()))
}
