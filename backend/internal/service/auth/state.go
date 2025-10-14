package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/takumi/personal-website/internal/config"
)

// StateManager issues and validates HMAC-signed OAuth state parameters.
type StateManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// StatePayload carries metadata embedded in the OAuth state parameter.
type StatePayload struct {
	Nonce       string
	RedirectURI string
	IssuedAt    time.Time
}

// NewStateManager constructs a state manager with the provided configuration.
func NewStateManager(cfg *config.AppConfig) (*StateManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("state manager: missing config")
	}
	if len(cfg.Auth.StateSecret) == 0 {
		return nil, fmt.Errorf("state manager: state_secret must be configured")
	}

	ttl := time.Duration(cfg.Auth.StateTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &StateManager{
		secret: []byte(cfg.Auth.StateSecret),
		ttl:    ttl,
		now:    time.Now,
	}, nil
}

// Issue builds a signed state string containing a random nonce and optional redirect URI.
func (m *StateManager) Issue(redirectURI string) (string, error) {
	payload := struct {
		Nonce       string `json:"nonce"`
		RedirectURI string `json:"redirect_uri,omitempty"`
		IssuedAt    int64  `json:"issued_at"`
	}{
		Nonce:       randomNonce(),
		RedirectURI: redirectURI,
		IssuedAt:    m.now().Unix(),
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("state manager: marshal payload: %w", err)
	}

	payloadSegment := base64.RawURLEncoding.EncodeToString(encodedPayload)
	signature, err := m.sign(payloadSegment)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", payloadSegment, signature), nil
}

// Validate ensures the state originates from this service and is within the configured TTL.
func (m *StateManager) Validate(state string) (*StatePayload, error) {
	payloadSegment, signatureSegment, found := splitState(state)
	if !found {
		return nil, errors.New("invalid state: malformed value")
	}

	expected, err := m.sign(payloadSegment)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal([]byte(signatureSegment), []byte(expected)) {
		return nil, errors.New("invalid state: signature mismatch")
	}

	rawPayload, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, fmt.Errorf("invalid state: decode payload: %w", err)
	}

	var payload struct {
		Nonce       string `json:"nonce"`
		RedirectURI string `json:"redirect_uri"`
		IssuedAt    int64  `json:"issued_at"`
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return nil, fmt.Errorf("invalid state: decode json: %w", err)
	}
	issuedAt := time.Unix(payload.IssuedAt, 0)
	if m.now().After(issuedAt.Add(m.ttl)) {
		return nil, errors.New("invalid state: expired")
	}

	return &StatePayload{
		Nonce:       payload.Nonce,
		RedirectURI: payload.RedirectURI,
		IssuedAt:    issuedAt,
	}, nil
}

func (m *StateManager) sign(payload string) (string, error) {
	mac := hmac.New(sha256.New, m.secret)
	if _, err := mac.Write([]byte(payload)); err != nil {
		return "", fmt.Errorf("state manager: unable to sign payload: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func splitState(state string) (payload, signature string, ok bool) {
	for i := 0; i < len(state); i++ {
		if state[i] == '.' {
			return state[:i], state[i+1:], true
		}
	}
	return "", "", false
}

func randomNonce() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// Fallback to timestamp-based nonce when randomness fails.
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}
