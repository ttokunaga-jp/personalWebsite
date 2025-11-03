package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
)

func TestStateManagerIssueValidate(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			StateSecret:     "state-secret",
			StateTTLSeconds: 120,
		},
	}

	manager, err := NewStateManager(cfg)
	require.NoError(t, err)

	sm := manager
	sm.now = func() time.Time { return time.Unix(1700000000, 0) }

	state, err := sm.Issue("/admin/")
	require.NoError(t, err)
	require.NotEmpty(t, state)

	payload, err := sm.Validate(state)
	require.NoError(t, err)
	require.Equal(t, "/admin/", payload.RedirectURI)
	require.Equal(t, time.Unix(1700000000, 0), payload.IssuedAt)
}

func TestStateManagerExpires(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			StateSecret:     "state-secret",
			StateTTLSeconds: 1,
		},
	}

	manager, err := NewStateManager(cfg)
	require.NoError(t, err)

	sm := manager
	base := time.Unix(1700000000, 0)
	sm.now = func() time.Time { return base }

	state, err := sm.Issue("/admin/")
	require.NoError(t, err)

	sm.now = func() time.Time { return base.Add(2 * time.Second) }

	_, err = sm.Validate(state)
	require.Error(t, err)
}

func TestStateManagerTamper(t *testing.T) {
	cfg := &config.AppConfig{
		Auth: config.AuthConfig{
			StateSecret:     "state-secret",
			StateTTLSeconds: 60,
		},
	}

	manager, err := NewStateManager(cfg)
	require.NoError(t, err)

	state, err := manager.Issue("/admin/")
	require.NoError(t, err)

	// Flip the last byte to break the signature.
	runes := []rune(state)
	runes[len(runes)-1] = 'a'

	_, err = manager.Validate(string(runes))
	require.Error(t, err)
}
