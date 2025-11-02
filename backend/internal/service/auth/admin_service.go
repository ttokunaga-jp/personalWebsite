package auth

import (
	"context"
	"log"
	"strings"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
)

// AdminService coordinates the administrator SSO flow via Google OAuth.
type AdminService interface {
	StartLogin(ctx context.Context, redirectURI string) (*AdminLoginResult, error)
	HandleCallback(ctx context.Context, state, code string) (*AdminCallbackResult, error)
}

// AdminLoginResult describes the OAuth initiation result for administrators.
type AdminLoginResult struct {
	AuthURL string
	State   string
}

// AdminCallbackResult carries the outcome of a successful administrator login.
type AdminCallbackResult struct {
	Token        string
	ExpiresAt    int64
	RedirectPath string
}

type adminService struct {
	provider       OAuthProvider
	issuer         TokenIssuer
	stateManager   *StateManager
	authCfg        config.AuthConfig
	adminCfg       config.AdminAuthConfig
	allowedDomains map[string]struct{}
	allowedEmails  map[string]struct{}
	tokenSaver     TokenSaver
}

// NewAdminService assembles the admin authentication service with whitelist enforcement.
func NewAdminService(
	cfg *config.AppConfig,
	provider OAuthProvider,
	issuer TokenIssuer,
	stateManager *StateManager,
	tokenSaver TokenSaver,
) (AdminService, error) {
	if cfg == nil {
		return nil, errs.New(errs.CodeInternal, 500, "admin auth service: missing config", nil)
	}
	if provider == nil {
		return nil, errs.New(errs.CodeInternal, 500, "admin auth service: missing oauth provider", nil)
	}
	if issuer == nil {
		return nil, errs.New(errs.CodeInternal, 500, "admin auth service: missing token issuer", nil)
	}
	if stateManager == nil {
		return nil, errs.New(errs.CodeInternal, 500, "admin auth service: missing state manager", nil)
	}

	adminCfg := cfg.Auth.Admin
	if strings.TrimSpace(adminCfg.DefaultRedirectURI) == "" {
		adminCfg.DefaultRedirectURI = "/admin"
	}

	domainSet := make(map[string]struct{})
	for _, domain := range adminCfg.AllowedDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain != "" {
			domainSet[domain] = struct{}{}
		}
	}

	emailSet := make(map[string]struct{})
	for _, email := range adminCfg.AllowedEmails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			emailSet[email] = struct{}{}
		}
	}

	return &adminService{
		provider:       provider,
		issuer:         issuer,
		stateManager:   stateManager,
		authCfg:        cfg.Auth,
		adminCfg:       adminCfg,
		allowedDomains: domainSet,
		allowedEmails:  emailSet,
		tokenSaver:     tokenSaver,
	}, nil
}

func (s *adminService) StartLogin(ctx context.Context, redirectURI string) (*AdminLoginResult, error) {
	target := sanitizeAdminRedirect(redirectURI, s.adminCfg.DefaultRedirectURI)
	state, err := s.stateManager.Issue(target)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, 500, "failed to issue oauth state", err)
	}

	if s.authCfg.Disabled {
		return &AdminLoginResult{
			AuthURL: "",
			State:   state,
		}, nil
	}

	authURL := s.provider.AuthCodeURL(state)
	if strings.TrimSpace(authURL) == "" {
		return nil, errs.New(errs.CodeInternal, 500, "oauth provider not configured", nil)
	}

	return &AdminLoginResult{
		AuthURL: authURL,
		State:   state,
	}, nil
}

func (s *adminService) HandleCallback(ctx context.Context, state, code string) (*AdminCallbackResult, error) {
	if strings.TrimSpace(state) == "" || strings.TrimSpace(code) == "" {
		return nil, errs.New(errs.CodeInvalidInput, 400, "missing state or code", nil)
	}

	statePayload, err := s.stateManager.Validate(state)
	if err != nil {
		return nil, errs.New(errs.CodeUnauthorized, 401, "invalid oauth state", err)
	}

	redirectPath := sanitizeAdminRedirect(statePayload.RedirectURI, s.adminCfg.DefaultRedirectURI)

	if s.authCfg.Disabled {
		token, expiresAt, err := s.issuer.Issue(ctx, "disabled-auth", "", "admin")
		if err != nil {
			return nil, err
		}
		return &AdminCallbackResult{
			Token:        token,
			ExpiresAt:    expiresAt.Unix(),
			RedirectPath: redirectPath,
		}, nil
	}

	token, err := s.provider.Exchange(ctx, code)
	if err != nil {
		log.Printf("admin auth: oauth exchange failed: %v", err)
		return nil, errs.New(errs.CodeUnauthorized, 401, "oauth exchange failed", err)
	}

	userInfo, err := s.provider.FetchUserInfo(ctx, token)
	if err != nil {
		log.Printf("admin auth: oauth userinfo fetch failed: %v", err)
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
			log.Printf("admin auth service: failed to persist oauth token: %v", err)
			return nil, errs.New(errs.CodeInternal, 500, "failed to persist Google authorization; please retry after granting offline access", err)
		}
	} else {
		log.Printf("admin auth service: token saver not configured; skipping token persistence")
	}

	appToken, expiresAt, err := s.issuer.Issue(ctx, userInfo.Subject, userInfo.Email, "admin")
	if err != nil {
		return nil, err
	}

	return &AdminCallbackResult{
		Token:        appToken,
		ExpiresAt:    expiresAt.Unix(),
		RedirectPath: redirectPath,
	}, nil
}

func (s *adminService) ensureAllowed(info *UserInfo) error {
	email := strings.ToLower(strings.TrimSpace(info.Email))
	if email == "" {
		return errs.New(errs.CodeUnauthorized, 401, "google account missing email", nil)
	}

	if len(s.allowedEmails) == 0 && len(s.allowedDomains) == 0 {
		return errs.New(errs.CodeForbidden, 403, "admin access denied", nil)
	}

	if _, ok := s.allowedEmails[email]; ok {
		return nil
	}

	if domain := domainPart(email); domain != "" {
		if _, ok := s.allowedDomains[domain]; ok {
			return nil
		}
	}

	return errs.New(errs.CodeForbidden, 403, "admin access denied", nil)
}

func sanitizeAdminRedirect(requested, fallback string) string {
	requested = strings.TrimSpace(requested)
	fallback = strings.TrimSpace(fallback)
	if fallback == "" {
		fallback = "/admin"
	}
	if !strings.HasPrefix(fallback, "/") {
		fallback = "/admin"
	}
	if requested == "" {
		return fallback
	}
	if strings.HasPrefix(requested, "http://") || strings.HasPrefix(requested, "https://") {
		return fallback
	}
	if !strings.HasPrefix(requested, "/") {
		return fallback
	}
	return requested
}

var _ AdminService = (*adminService)(nil)
