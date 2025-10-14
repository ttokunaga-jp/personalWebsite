package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
)

func TestServiceSuccessFlow(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "secret",
			Issuer:                "personal-website",
			Audience:              []string{"admin"},
			AccessTokenTTLMinutes: 30,
			StateSecret:           "state-secret",
			StateTTLSeconds:       300,
		},
		Google: config.GoogleOAuthConfig{
			AllowedDomains: []string{"example.com"},
		},
	}

	stateMgr, err := NewStateManager(cfg)
	require.NoError(t, err)
	now := time.Now().Truncate(time.Second)
	stateMgr.now = func() time.Time { return now }

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)
	iss := issuer.(*jwtIssuer)
	iss.now = func() time.Time { return now }

	provider := &stubProvider{
		authBase: "https://accounts.google.com/o/oauth2/v2/auth?state=",
		token: &OAuthToken{
			AccessToken: "access",
		},
		user: &UserInfo{
			Subject:       "google-123",
			Email:         "admin@example.com",
			EmailVerified: true,
		},
	}

	service, err := NewService(cfg, provider, issuer, stateMgr)
	require.NoError(t, err)

	login, err := service.StartLogin(context.Background(), "/admin/dashboard")
	require.NoError(t, err)
	require.NotEmpty(t, login.AuthURL)
	require.Contains(t, login.AuthURL, "state=")

	callback, err := service.HandleCallback(context.Background(), login.State, "auth-code")
	require.NoError(t, err)
	require.NotEmpty(t, callback.Token)
	require.Equal(t, "/admin/dashboard", callback.RedirectURI)

	verifier := NewJWTVerifier(cfg.Auth)
	claims, err := verifier.Verify(context.Background(), callback.Token)
	require.NoError(t, err)
	require.Equal(t, "google-123", claims.Subject)
	require.Equal(t, "admin@example.com", claims.Email)
}

func TestServiceRejectsUnauthorizedDomain(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "secret",
			AccessTokenTTLMinutes: 30,
			StateSecret:           "state-secret",
			StateTTLSeconds:       300,
		},
		Google: config.GoogleOAuthConfig{
			AllowedDomains: []string{"example.com"},
		},
	}

	stateMgr, err := NewStateManager(cfg)
	require.NoError(t, err)

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)

	provider := &stubProvider{
		token: &OAuthToken{AccessToken: "access"},
		user: &UserInfo{
			Subject:       "id",
			Email:         "user@other.com",
			EmailVerified: true,
		},
	}

	service, err := NewService(cfg, provider, issuer, stateMgr)
	require.NoError(t, err)

	login, err := service.StartLogin(context.Background(), "/admin")
	require.NoError(t, err)

	_, err = service.HandleCallback(context.Background(), login.State, "code")
	require.Error(t, err)
}

func TestServiceDisabledBypassesProvider(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			JWTSecret:             "secret",
			AccessTokenTTLMinutes: 30,
			StateSecret:           "state-secret",
			StateTTLSeconds:       300,
			Disabled:              true,
		},
	}

	stateMgr, err := NewStateManager(cfg)
	require.NoError(t, err)
	now := time.Now().Truncate(time.Second)
	stateMgr.now = func() time.Time { return now }

	issuer, err := NewJWTIssuer(cfg)
	require.NoError(t, err)

	provider := &stubProvider{
		exchangeErr: errors.New("should not be called"),
	}

	service, err := NewService(cfg, provider, issuer, stateMgr)
	require.NoError(t, err)

	login, err := service.StartLogin(context.Background(), "")
	require.NoError(t, err)
	require.Empty(t, login.AuthURL)

	res, err := service.HandleCallback(context.Background(), login.State, "ignored")
	require.NoError(t, err)
	require.Equal(t, "/admin", res.RedirectURI)
	require.NotEmpty(t, res.Token)
}

type stubProvider struct {
	authBase    string
	token       *OAuthToken
	user        *UserInfo
	exchangeErr error
	userErr     error
}

func (s *stubProvider) AuthCodeURL(state string) string {
	return s.authBase + state
}

func (s *stubProvider) Exchange(context.Context, string) (*OAuthToken, error) {
	if s.exchangeErr != nil {
		return nil, s.exchangeErr
	}
	if s.token == nil {
		return nil, errors.New("no token configured")
	}
	return s.token, nil
}

func (s *stubProvider) FetchUserInfo(context.Context, *OAuthToken) (*UserInfo, error) {
	if s.userErr != nil {
		return nil, s.userErr
	}
	if s.user == nil {
		return nil, errors.New("no user configured")
	}
	return s.user, nil
}
