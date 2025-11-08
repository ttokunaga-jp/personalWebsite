package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// AdminSessionManager coordinates issuing and validating administrator sessions.
type AdminSessionManager interface {
	Create(ctx context.Context, principal AdminPrincipal) (*model.AdminSession, error)
	Validate(ctx context.Context, sessionID string) (*model.AdminSession, error)
	Refresh(ctx context.Context, sessionID string) (*model.AdminSession, error)
	Revoke(ctx context.Context, sessionID string) error
}

// AdminPrincipal captures the identity metadata to persist for a session.
type AdminPrincipal struct {
	Subject   string
	Email     string
	Roles     []string
	UserAgent string
	IPAddress string
}

type adminSessionManager struct {
	repo          repository.AdminSessionRepository
	ttl           time.Duration
	idleTimeout   time.Duration
	refreshWindow time.Duration
	now           func() time.Time
}

// NewAdminSessionManager constructs a session manager using the configured repository.
func NewAdminSessionManager(cfg *config.AppConfig, repo repository.AdminSessionRepository) (AdminSessionManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("admin session manager: missing config")
	}
	// When authentication is disabled (e.g. smoke tests or local runs),
	// provide a no-op session manager that always behaves as unauthenticated.
	if repo == nil {
		if cfg.Auth.Disabled {
			return noopAdminSessionManager{}, nil
		}
		return nil, fmt.Errorf("admin session manager: missing repository")
	}

	adminCfg := cfg.Auth.Admin
	manager := &adminSessionManager{
		repo:          repo,
		ttl:           normalizeDuration(adminCfg.SessionTTL, 24*time.Hour),
		idleTimeout:   normalizeDuration(adminCfg.SessionIdleTimeout, 2*time.Hour),
		refreshWindow: normalizeDuration(adminCfg.SessionRefreshWindow, 20*time.Minute),
		now:           time.Now,
	}
	return manager, nil
}

func normalizeDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value <= 0 {
		return fallback
	}
	return value
}

func (m *adminSessionManager) Create(ctx context.Context, principal AdminPrincipal) (*model.AdminSession, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("issue session id: %w", err)
	}

	hash := hashSessionID(sessionID)
	now := m.now().UTC()

	session := &model.AdminSession{
		ID:             sessionID,
		TokenHash:      hash,
		Subject:        strings.TrimSpace(principal.Subject),
		Email:          strings.ToLower(strings.TrimSpace(principal.Email)),
		Roles:          normalizeRoles(principal.Roles),
		UserAgent:      strings.TrimSpace(principal.UserAgent),
		IPAddress:      strings.TrimSpace(principal.IPAddress),
		CreatedAt:      now,
		UpdatedAt:      now,
		LastAccessedAt: now,
		ExpiresAt:      now.Add(m.ttl),
	}

	stored, err := m.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}
	stored.ID = sessionID
	return stored, nil
}

func (m *adminSessionManager) Validate(ctx context.Context, sessionID string) (*model.AdminSession, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, repository.ErrNotFound
	}

	hash := hashSessionID(sessionID)
	session, err := m.repo.FindSessionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	now := m.now().UTC()
	if session.RevokedAt != nil {
		return nil, repository.ErrNotFound
	}
	if !session.ExpiresAt.After(now) {
		_ = m.repo.RevokeSession(ctx, hash)
		return nil, repository.ErrNotFound
	}
	if m.idleTimeout > 0 && now.Sub(session.LastAccessedAt) > m.idleTimeout {
		_ = m.repo.RevokeSession(ctx, hash)
		return nil, repository.ErrNotFound
	}

	// Touch the session to extend idle timeout tracking.
	updated, err := m.repo.UpdateSessionActivity(ctx, hash, now, session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	updated.ID = sessionID
	return updated, nil
}

func (m *adminSessionManager) Refresh(ctx context.Context, sessionID string) (*model.AdminSession, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, repository.ErrNotFound
	}

	hash := hashSessionID(sessionID)
	session, err := m.repo.FindSessionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	now := m.now().UTC()
	if session.RevokedAt != nil {
		return nil, repository.ErrNotFound
	}
	if !session.ExpiresAt.After(now) {
		_ = m.repo.RevokeSession(ctx, hash)
		return nil, repository.ErrNotFound
	}

	nextExpiry := now.Add(m.ttl)
	updated, err := m.repo.UpdateSessionActivity(ctx, hash, now, nextExpiry)
	if err != nil {
		return nil, err
	}
	updated.ID = sessionID
	return updated, nil
}

func (m *adminSessionManager) Revoke(ctx context.Context, sessionID string) error {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return repository.ErrNotFound
	}
	hash := hashSessionID(sessionID)
	return m.repo.RevokeSession(ctx, hash)
}

// noopAdminSessionManager is a minimal implementation used when auth is disabled.
// It consistently returns ErrNotFound for lookups and performs no persistence.
type noopAdminSessionManager struct{}

func (noopAdminSessionManager) Create(context.Context, AdminPrincipal) (*model.AdminSession, error) {
	return nil, repository.ErrInvalidInput
}

func (noopAdminSessionManager) Validate(context.Context, string) (*model.AdminSession, error) {
	return nil, repository.ErrNotFound
}

func (noopAdminSessionManager) Refresh(context.Context, string) (*model.AdminSession, error) {
	return nil, repository.ErrNotFound
}

func (noopAdminSessionManager) Revoke(context.Context, string) error { return repository.ErrNotFound }

func generateSessionID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashSessionID(id string) string {
	sum := sha256.Sum256([]byte(id))
	return hex.EncodeToString(sum[:])
}
