package csrf

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestIssueAndValidate(t *testing.T) {
	t.Parallel()

	manager := NewManager("test-secret", time.Minute)

	token, err := manager.Issue()
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	if err := manager.Validate(token.Cookie, token.Value); err != nil {
		t.Fatalf("validate token: %v", err)
	}
}

func TestValidateRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	manager := NewManager("another-secret", time.Millisecond)
	token, err := manager.Issue()
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	time.Sleep(2 * time.Millisecond)
	if err := manager.Validate(token.Cookie, token.Value); err == nil {
		t.Fatalf("expected expired token error")
	}
}

func TestValidateRejectsTamperedToken(t *testing.T) {
	t.Parallel()

	manager := NewManager("secret", time.Minute)
	token, err := manager.Issue()
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(token.Value)
	if err != nil {
		t.Fatalf("decode token: %v", err)
	}
	decoded[0] ^= 0xFF
	altered := base64.RawURLEncoding.EncodeToString(decoded)

	if err := manager.Validate(token.Cookie, altered); err == nil {
		t.Fatalf("expected tampering to be detected")
	}
}
