package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidToken   = errors.New("invalid csrf token")
	errExpiredToken   = errors.New("expired csrf token")
	errMalformedToken = errors.New("malformed csrf token")

	// Exported aliases for external package checks.
	ErrInvalidToken   = errInvalidToken
	ErrExpiredToken   = errExpiredToken
	ErrMalformedToken = errMalformedToken
)

// Manager handles issuance and verification of CSRF tokens using an HMAC signature.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// Token encapsulates the values required to set client cookies and headers.
type Token struct {
	Value     string
	Cookie    string
	ExpiresAt time.Time
}

// NewManager returns a new CSRF token manager.
func NewManager(signingKey string, ttl time.Duration) *Manager {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &Manager{secret: []byte(signingKey), ttl: ttl}
}

// Issue generates a token value alongside the cookie payload with embedded expiry.
func (m *Manager) Issue() (*Token, error) {
	if len(m.secret) == 0 {
		return nil, fmt.Errorf("csrf manager: signing key is empty")
	}

	raw := make([]byte, 32)
	if _, err := randRead(raw); err != nil {
		return nil, fmt.Errorf("csrf manager: generate random: %w", err)
	}

	value := base64.RawURLEncoding.EncodeToString(raw)
	expires := time.Now().UTC().Add(m.ttl).Unix()
	message := fmt.Sprintf("%s:%d", value, expires)
	signature := sign(message, m.secret)
	cookie := fmt.Sprintf("%s:%d:%s", value, expires, signature)

	return &Token{
		Value:     value,
		Cookie:    cookie,
		ExpiresAt: time.Unix(expires, 0).UTC(),
	}, nil
}

// Validate checks that the provided header token matches the signed cookie payload.
func (m *Manager) Validate(cookieValue, headerValue string) error {
	if len(m.secret) == 0 {
		return fmt.Errorf("csrf manager: signing key is empty")
	}
	if cookieValue == "" || headerValue == "" {
		return errInvalidToken
	}

	parts := strings.Split(cookieValue, ":")
	if len(parts) != 3 {
		return errMalformedToken
	}

	token := parts[0]
	expires, err := parseExpires(parts[1])
	if err != nil {
		return errMalformedToken
	}
	signature := parts[2]
	message := fmt.Sprintf("%s:%d", token, expires.Unix())

	if !hmac.Equal([]byte(signature), []byte(sign(message, m.secret))) {
		return errInvalidToken
	}

	if expires.Before(time.Now().UTC()) {
		return errExpiredToken
	}

	if subtleConstantTimeCompare(token, headerValue) == false {
		return errInvalidToken
	}

	return nil
}

func parseExpires(raw string) (time.Time, error) {
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(value, 0).UTC(), nil
}

func sign(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// subtleConstantTimeCompare compares two strings without leaking timing information.
func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// randRead wraps crypto/rand.Read for testability.
var randRead = func(b []byte) (int, error) {
	return rand.Read(b)
}
