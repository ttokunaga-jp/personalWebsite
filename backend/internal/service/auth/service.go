package auth

import (
	"context"
	"log"
	"strings"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
)

// Service orchestrates the Google OAuth login and JWT issuance flow.
type Service interface {
	StartLogin(ctx context.Context, redirectURI string) (*LoginResult, error)
	HandleCallback(ctx context.Context, state, code string) (*CallbackResult, error)
}

// LoginResult represents the first step of the OAuth flow.
type LoginResult struct {
	AuthURL string
	State   string
}

// CallbackResult returns the application JWT and redirect target.
type CallbackResult struct {
	Token       string
	ExpiresAt   int64
	RedirectURI string
}

type service struct {
	provider       OAuthProvider
	issuer         TokenIssuer
	stateManager   *StateManager
	authCfg        config.AuthConfig
	googleCfg      config.GoogleOAuthConfig
	allowedDomains map[string]struct{}
	allowedEmails  map[string]struct{}
	tokenSaver     TokenSaver
}

// TokenSaver persists OAuth tokens for downstream services such as Gmail API integrations.
type TokenSaver interface {
	Save(ctx context.Context, token *OAuthToken) error
}

// NewService assembles the authentication service from its collaborators.
func NewService(
	cfg *config.AppConfig,
	provider OAuthProvider,
	issuer TokenIssuer,
	stateManager *StateManager,
	tokenSaver TokenSaver,
) (Service, error) {
	if cfg == nil {
		return nil, errs.New(errs.CodeInternal, 500, "auth service: missing config", nil)
	}
	if provider == nil {
		return nil, errs.New(errs.CodeInternal, 500, "auth service: missing oauth provider", nil)
	}
	if issuer == nil {
		return nil, errs.New(errs.CodeInternal, 500, "auth service: missing token issuer", nil)
	}
	if stateManager == nil {
		return nil, errs.New(errs.CodeInternal, 500, "auth service: missing state manager", nil)
	}

	domainSet := make(map[string]struct{})
	for _, domain := range cfg.Google.AllowedDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain != "" {
			domainSet[domain] = struct{}{}
		}
	}

	emailSet := make(map[string]struct{})
	for _, email := range cfg.Google.AllowedEmails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			emailSet[email] = struct{}{}
		}
	}

	return &service{
		provider:       provider,
		issuer:         issuer,
		stateManager:   stateManager,
		authCfg:        cfg.Auth,
		googleCfg:      cfg.Google,
		allowedDomains: domainSet,
		allowedEmails:  emailSet,
		tokenSaver:     tokenSaver,
	}, nil
}

func (s *service) StartLogin(ctx context.Context, redirectURI string) (*LoginResult, error) {
	redirectURI = sanitizeRedirectURI(redirectURI)

	state, err := s.stateManager.Issue(redirectURI)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, 500, "failed to issue oauth state", err)
	}

	if s.authCfg.Disabled {
		return &LoginResult{
			AuthURL: "",
			State:   state,
		}, nil
	}

	authURL := s.provider.AuthCodeURL(state)
	if strings.TrimSpace(authURL) == "" {
		return nil, errs.New(errs.CodeInternal, 500, "oauth provider not configured", nil)
	}

	return &LoginResult{
		AuthURL: authURL,
		State:   state,
	}, nil
}

func (s *service) HandleCallback(ctx context.Context, state, code string) (*CallbackResult, error) {
	if strings.TrimSpace(state) == "" || strings.TrimSpace(code) == "" {
		return nil, errs.New(errs.CodeInvalidInput, 400, "missing state or code", nil)
	}

	statePayload, err := s.stateManager.Validate(state)
	if err != nil {
		return nil, errs.New(errs.CodeUnauthorized, 401, "invalid oauth state", err)
	}

	if s.authCfg.Disabled {
		token, expiresAt, err := s.issuer.Issue(ctx, "disabled-auth", "")
		if err != nil {
			return nil, err
		}
		return &CallbackResult{
			Token:       token,
			ExpiresAt:   expiresAt.Unix(),
			RedirectURI: sanitizeRedirectURI(statePayload.RedirectURI),
		}, nil
	}

	token, err := s.provider.Exchange(ctx, code)
	if err != nil {
		return nil, errs.New(errs.CodeUnauthorized, 401, "oauth exchange failed", err)
	}

	userInfo, err := s.provider.FetchUserInfo(ctx, token)
	if err != nil {
		return nil, errs.New(errs.CodeUnauthorized, 401, "oauth userinfo fetch failed", err)
	}

	if !userInfo.EmailVerified {
		return nil, errs.New(errs.CodeUnauthorized, 401, "email not verified for google account", nil)
	}

	if err := s.ensureAllowed(userInfo); err != nil {
		return nil, err
	}

	if s.tokenSaver != nil {
		if err := s.tokenSaver.Save(ctx, token); err != nil {
			log.Printf("auth service: failed to persist oauth token: %v", err)
			return nil, errs.New(errs.CodeInternal, 500, "failed to persist Google authorization; please retry after granting offline access", err)
		}
	} else {
		log.Printf("auth service: token saver not configured; skipping token persistence")
	}

	appToken, expiresAt, err := s.issuer.Issue(ctx, userInfo.Subject, userInfo.Email)
	if err != nil {
		return nil, err
	}

	return &CallbackResult{
		Token:       appToken,
		ExpiresAt:   expiresAt.Unix(),
		RedirectURI: sanitizeRedirectURI(statePayload.RedirectURI),
	}, nil
}

func (s *service) ensureAllowed(info *UserInfo) error {
	email := strings.ToLower(strings.TrimSpace(info.Email))
	if email == "" {
		return errs.New(errs.CodeUnauthorized, 401, "google account missing email", nil)
	}

	if len(s.allowedEmails) == 0 && len(s.allowedDomains) == 0 {
		return nil
	}

	if _, ok := s.allowedEmails[email]; ok {
		return nil
	}

	if domain := domainPart(email); domain != "" {
		if _, ok := s.allowedDomains[domain]; ok {
			return nil
		}
	}

	return errs.New(errs.CodeUnauthorized, 401, "google account not permitted", nil)
}

func domainPart(email string) string {
	if idx := strings.LastIndex(email, "@"); idx != -1 && idx+1 < len(email) {
		return email[idx+1:]
	}
	return ""
}

func sanitizeRedirectURI(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return "/admin"
	}
	if !strings.HasPrefix(raw, "/") {
		return "/admin"
	}
	return raw
}

var _ Service = (*service)(nil)
