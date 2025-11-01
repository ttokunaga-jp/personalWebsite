package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
)

func TestClaimsHasRole(t *testing.T) {
	t.Parallel()

	appCfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "secret",
			Issuer:                "issuer",
			Audience:              []string{"admin"},
			AccessTokenTTLMinutes: 60,
		},
	}

	issuer, err := NewJWTIssuer(appCfg)
	require.NoError(t, err)

	token, _, err := issuer.Issue(context.Background(), "user-1", "user@example.com", "admin")
	require.NoError(t, err)

	verifier := NewJWTVerifier(appCfg.Auth)
	claims, err := verifier.Verify(context.Background(), token)
	require.NoError(t, err)

	require.True(t, claims.HasRole("ADMIN"))
	require.False(t, claims.HasRole("viewer"))
}
