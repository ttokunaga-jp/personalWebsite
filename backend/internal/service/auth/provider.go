package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/config"
)

// OAuthProvider abstracts Google OAuth operations for easier testing.
type OAuthProvider interface {
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*OAuthToken, error)
	FetchUserInfo(ctx context.Context, token *OAuthToken) (*UserInfo, error)
}

// UserInfo captures relevant Google profile fields for authorization decisions.
type UserInfo struct {
	Subject       string
	Email         string
	EmailVerified bool
	HostedDomain  string
}

var ErrProviderDisabled = errors.New("oauth provider disabled")

// OAuthToken is the subset of Google token fields required downstream.
type OAuthToken struct {
	AccessToken  string
	IDToken      string
	ExpiresIn    int64
	RefreshToken string
}

type googleProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
	authURL      string
	tokenURL     string
	userInfoURL  string
	httpClient   *http.Client
}

// NewGoogleOAuthProvider builds an OAuth provider configured for Google endpoints.
func NewGoogleOAuthProvider(cfg *config.AppConfig) (OAuthProvider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("google provider: missing config")
	}

	gCfg := cfg.Google
	if gCfg.ClientID == "" || gCfg.ClientSecret == "" || gCfg.RedirectURL == "" {
		return &noopProvider{}, nil
	}

	return &googleProvider{
		clientID:     gCfg.ClientID,
		clientSecret: gCfg.ClientSecret,
		redirectURL:  gCfg.RedirectURL,
		scopes:       gCfg.Scopes,
		authURL:      pickFirstNonEmpty(gCfg.AuthURL, "https://accounts.google.com/o/oauth2/v2/auth"),
		tokenURL:     pickFirstNonEmpty(gCfg.TokenURL, "https://oauth2.googleapis.com/token"),
		userInfoURL:  pickFirstNonEmpty(gCfg.UserInfoURL, "https://openidconnect.googleapis.com/v1/userinfo"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (p *googleProvider) AuthCodeURL(state string) string {
	if p.clientID == "" || p.redirectURL == "" || p.authURL == "" {
		return ""
	}
	values := url.Values{}
	values.Set("client_id", p.clientID)
	values.Set("redirect_uri", p.redirectURL)
	values.Set("response_type", "code")
	values.Set("scope", strings.Join(p.scopes, " "))
	values.Set("state", state)
	values.Set("access_type", "offline")
	values.Set("prompt", "consent")

	return fmt.Sprintf("%s?%s", strings.TrimRight(p.authURL, "?"), values.Encode())
}

func (p *googleProvider) Exchange(ctx context.Context, code string) (*OAuthToken, error) {
	if p.clientID == "" || p.tokenURL == "" {
		return nil, ErrProviderDisabled
	}

	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("redirect_uri", p.redirectURL)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("google provider: build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google provider: exchange code: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("google provider: token endpoint status %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("google provider: decode token response: %w", err)
	}

	return &OAuthToken{
		AccessToken:  payload.AccessToken,
		ExpiresIn:    payload.ExpiresIn,
		IDToken:      payload.IDToken,
		RefreshToken: payload.RefreshToken,
	}, nil
}

func (p *googleProvider) FetchUserInfo(ctx context.Context, token *OAuthToken) (*UserInfo, error) {
	if p.clientID == "" || p.userInfoURL == "" {
		return nil, ErrProviderDisabled
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("google provider: build userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google provider: fetch userinfo: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("google provider: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var payload struct {
		Subject       string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		HostedDomain  string `json:"hd"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("google provider: decode userinfo: %w", err)
	}

	return &UserInfo{
		Subject:       payload.Subject,
		Email:         strings.ToLower(payload.Email),
		EmailVerified: payload.EmailVerified,
		HostedDomain:  strings.ToLower(payload.HostedDomain),
	}, nil
}

type noopProvider struct{}

func (n *noopProvider) AuthCodeURL(string) string {
	return ""
}

func (n *noopProvider) Exchange(context.Context, string) (*OAuthToken, error) {
	return nil, ErrProviderDisabled
}

func (n *noopProvider) FetchUserInfo(context.Context, *OAuthToken) (*UserInfo, error) {
	return nil, ErrProviderDisabled
}

func pickFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
