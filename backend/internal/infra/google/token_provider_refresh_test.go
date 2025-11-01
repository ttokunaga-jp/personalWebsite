package google

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/service/auth"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type stubTokenStore struct {
	record *TokenRecord
	saved  *TokenRecord
	err    error
}

func (s *stubTokenStore) Load(ctx context.Context, provider string) (*TokenRecord, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.record == nil {
		return nil, ErrTokenNotFound
	}
	copy := *s.record
	return &copy, nil
}

func (s *stubTokenStore) Save(ctx context.Context, provider string, record *TokenRecord) error {
	copy := *record
	s.saved = &copy
	s.record = &copy
	return nil
}

func TestRefreshingTokenProviderRefreshesWhenExpired(t *testing.T) {
	store := &stubTokenStore{
		record: &TokenRecord{
			AccessToken:  "old-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(-5 * time.Minute),
		},
	}

	transport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, r.Body.Close())

		values, err := url.ParseQuery(string(body))
		require.NoError(t, err)
		require.Equal(t, "refresh-token", values.Get("refresh_token"))

		payload, err := json.Marshal(map[string]any{
			"access_token":  "new-token",
			"expires_in":    3600,
			"refresh_token": "refresh-token",
		})
		require.NoError(t, err)

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(string(payload))),
			Request:    r,
		}, nil
	})

	cfg := config.GoogleOAuthConfig{
		ClientID:     "client",
		ClientSecret: "secret",
		TokenURL:     "https://example.com/token",
	}

	provider, err := NewRefreshingTokenProvider(cfg, store, &http.Client{Transport: transport})
	require.NoError(t, err)

	token, err := provider.AccessToken(context.Background())
	require.NoError(t, err)
	require.Equal(t, "new-token", token)
	require.NotNil(t, store.saved)
	require.Equal(t, "new-token", store.saved.AccessToken)
}

type trackingStore struct {
	loadResp *TokenRecord
	loadErr  error
	saved    *TokenRecord
}

func (s *trackingStore) Load(ctx context.Context, provider string) (*TokenRecord, error) {
	if s.loadErr != nil {
		return nil, s.loadErr
	}
	if s.loadResp == nil {
		return nil, ErrTokenNotFound
	}
	copy := *s.loadResp
	return &copy, nil
}

func (s *trackingStore) Save(ctx context.Context, provider string, record *TokenRecord) error {
	copy := *record
	s.saved = &copy
	return nil
}

func TestGmailTokenManagerPersistsToken(t *testing.T) {
	store := &trackingStore{
		loadResp: &TokenRecord{RefreshToken: "existing-refresh"},
	}

	saver, err := NewGmailTokenManager(store)
	require.NoError(t, err)
	require.NotNil(t, saver)

	token := &auth.OAuthToken{
		AccessToken:  "access",
		RefreshToken: "",
		ExpiresIn:    1800,
	}

	err = saver.Save(context.Background(), token)
	require.NoError(t, err)
	require.NotNil(t, store.saved)
	require.Equal(t, "access", store.saved.AccessToken)
	require.Equal(t, "existing-refresh", store.saved.RefreshToken)
	require.True(t, store.saved.Expiry.After(time.Now()))
}
