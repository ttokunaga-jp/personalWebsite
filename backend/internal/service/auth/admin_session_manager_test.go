package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

func TestAdminSessionManagerCreateAndValidate(t *testing.T) {
	t.Parallel()

	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			Admin: config.AdminAuthConfig{
				SessionTTL:           30 * time.Minute,
				SessionIdleTimeout:   15 * time.Minute,
				SessionRefreshWindow: 5 * time.Minute,
			},
		},
	}
	store := newStubSessionRepository()
	manager, err := NewAdminSessionManager(cfg, store)
	require.NoError(t, err)

	session, err := manager.Create(context.Background(), AdminPrincipal{
		Subject: "admin-123",
		Email:   "admin@example.com",
		Roles:   []string{"admin"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, session.ID)
	require.NotEmpty(t, session.TokenHash)
	require.Equal(t, "admin@example.com", session.Email)
	require.Contains(t, session.Roles, "admin")

	validated, err := manager.Validate(context.Background(), session.ID)
	require.NoError(t, err)
	require.Equal(t, session.ID, validated.ID)
	require.True(t, validated.LastAccessedAt.After(session.LastAccessedAt) || validated.LastAccessedAt.Equal(session.LastAccessedAt))
}

func TestAdminSessionManagerRefreshExtendsExpiry(t *testing.T) {
	t.Parallel()

	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			Admin: config.AdminAuthConfig{
				SessionTTL:           45 * time.Minute,
				SessionIdleTimeout:   15 * time.Minute,
				SessionRefreshWindow: 10 * time.Minute,
			},
		},
	}
	store := newStubSessionRepository()
	manager, err := NewAdminSessionManager(cfg, store)
	require.NoError(t, err)

	session, err := manager.Create(context.Background(), AdminPrincipal{
		Subject: "admin-123",
		Email:   "admin@example.com",
		Roles:   []string{"admin"},
	})
	require.NoError(t, err)

	refreshed, err := manager.Refresh(context.Background(), session.ID)
	require.NoError(t, err)
	require.True(t, refreshed.ExpiresAt.After(session.ExpiresAt))
}

func TestAdminSessionManagerValidateHonoursIdleTimeout(t *testing.T) {
	t.Parallel()

	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			Admin: config.AdminAuthConfig{
				SessionTTL:           time.Hour,
				SessionIdleTimeout:   5 * time.Minute,
				SessionRefreshWindow: 1 * time.Minute,
			},
		},
	}
	store := newStubSessionRepository()
	manager, err := NewAdminSessionManager(cfg, store)
	require.NoError(t, err)

	session, err := manager.Create(context.Background(), AdminPrincipal{
		Subject: "admin-123",
		Email:   "admin@example.com",
		Roles:   []string{"admin"},
	})
	require.NoError(t, err)

	// Simulate idle timeout by setting last accessed time far in the past.
	store.withSession(session.TokenHash, func(s *model.AdminSession) {
		s.LastAccessedAt = time.Now().Add(-10 * time.Minute)
	})

	_, err = manager.Validate(context.Background(), session.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, repository.ErrNotFound)
}

type stubSessionRepository struct {
	sessions map[string]*model.AdminSession
}

func newStubSessionRepository() *stubSessionRepository {
	return &stubSessionRepository{
		sessions: make(map[string]*model.AdminSession),
	}
}

func (s *stubSessionRepository) CreateSession(_ context.Context, session *model.AdminSession) (*model.AdminSession, error) {
	cloned := cloneSession(session)
	s.sessions[session.TokenHash] = cloned
	return cloneSession(cloned), nil
}

func (s *stubSessionRepository) FindSessionByHash(_ context.Context, hash string) (*model.AdminSession, error) {
	session, ok := s.sessions[hash]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneSession(session), nil
}

func (s *stubSessionRepository) UpdateSessionActivity(_ context.Context, hash string, lastAccessed time.Time, expiresAt time.Time) (*model.AdminSession, error) {
	session, ok := s.sessions[hash]
	if !ok {
		return nil, repository.ErrNotFound
	}
	session.LastAccessedAt = lastAccessed
	session.ExpiresAt = expiresAt
	session.UpdatedAt = time.Now()
	return cloneSession(session), nil
}

func (s *stubSessionRepository) RevokeSession(_ context.Context, hash string) error {
	if _, ok := s.sessions[hash]; !ok {
		return repository.ErrNotFound
	}
	delete(s.sessions, hash)
	return nil
}

func (s *stubSessionRepository) withSession(hash string, fn func(*model.AdminSession)) {
	if session, ok := s.sessions[hash]; ok {
		fn(session)
	}
}

func cloneSession(session *model.AdminSession) *model.AdminSession {
	if session == nil {
		return nil
	}
	cloned := *session
	if session.Roles != nil {
		cloned.Roles = append([]string(nil), session.Roles...)
	}
	return &cloned
}
