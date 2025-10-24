package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/service/auth"
)

type refreshingTokenProvider struct {
	store        TokenStore
	clientID     string
	clientSecret string
	tokenURL     string
	httpClient   *http.Client
	clock        func() time.Time

	mu     sync.Mutex
	cached *TokenRecord
}

// NewRefreshingTokenProvider creates a TokenProvider that keeps access tokens fresh using refresh tokens persisted in the token store.
func NewRefreshingTokenProvider(cfg config.GoogleOAuthConfig, store TokenStore, client *http.Client) (TokenProvider, error) {
	if store == nil {
		return nil, errors.New("refreshing token provider: token store is nil")
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("refreshing token provider: missing client credentials")
	}
	tokenURL := cfg.TokenURL
	if strings.TrimSpace(tokenURL) == "" {
		tokenURL = "https://oauth2.googleapis.com/token"
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	return &refreshingTokenProvider{
		store:        store,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		tokenURL:     tokenURL,
		httpClient:   client,
		clock:        time.Now,
	}, nil
}

func (p *refreshingTokenProvider) AccessToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	record := p.cached
	p.mu.Unlock()

	if record == nil {
		var err error
		record, err = p.store.Load(ctx, GmailProvider)
		if err != nil {
			return "", err
		}
		p.mu.Lock()
		p.cached = record
		p.mu.Unlock()
	}

	if record.RefreshToken == "" {
		return "", errors.New("refreshing token provider: refresh token unavailable")
	}

	if p.clock().Before(record.Expiry.Add(-1 * time.Minute)) {
		return record.AccessToken, nil
	}

	updated, err := p.refresh(ctx, record)
	if err != nil {
		return "", err
	}

	p.mu.Lock()
	p.cached = updated
	p.mu.Unlock()

	return updated.AccessToken, nil
}

func (p *refreshingTokenProvider) refresh(ctx context.Context, record *TokenRecord) (*TokenRecord, error) {
	form := url.Values{}
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("refresh_token", record.RefreshToken)
	form.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refreshing token provider: build refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refreshing token provider: call token endpoint: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("refreshing token provider: token endpoint status %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("refreshing token provider: decode response: %w", err)
	}
	if payload.AccessToken == "" {
		return nil, errors.New("refreshing token provider: access token missing in response")
	}

	newRecord := &TokenRecord{
		AccessToken:  payload.AccessToken,
		RefreshToken: record.RefreshToken,
		Expiry:       p.clock().Add(time.Duration(payload.ExpiresIn) * time.Second),
	}

	if strings.TrimSpace(payload.RefreshToken) != "" {
		newRecord.RefreshToken = payload.RefreshToken
	}

	if err := p.store.Save(ctx, GmailProvider, newRecord); err != nil {
		return nil, err
	}

	return newRecord, nil
}

type gmailTokenManager struct {
	store TokenStore
}

// NewGmailTokenManager returns an auth.TokenSaver that records Gmail OAuth tokens in the provided store.
func NewGmailTokenManager(store TokenStore) (auth.TokenSaver, error) {
	if store == nil {
		return nil, nil
	}
	return &gmailTokenManager{store: store}, nil
}

func (m *gmailTokenManager) Save(ctx context.Context, token *auth.OAuthToken) error {
	if token == nil {
		return errors.New("gmail token manager: token is nil")
	}

	now := time.Now()
	expiry := now.Add(55 * time.Minute)
	if token.ExpiresIn > 0 {
		expiry = now.Add(time.Duration(token.ExpiresIn) * time.Second)
	}

	record := &TokenRecord{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       expiry,
	}

	// Preserve existing refresh token if Google did not return one.
	if record.RefreshToken == "" {
		existing, err := m.store.Load(ctx, GmailProvider)
		if err == nil && existing.RefreshToken != "" {
			record.RefreshToken = existing.RefreshToken
		} else if err != nil && !errors.Is(err, ErrTokenNotFound) {
			return fmt.Errorf("gmail token manager: load existing token: %w", err)
		}
	}

	if record.RefreshToken == "" {
		return errors.New("gmail token manager: refresh token missing")
	}

	return m.store.Save(ctx, GmailProvider, record)
}
