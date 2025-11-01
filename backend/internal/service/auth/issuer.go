package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
)

// TokenIssuer generates signed JWT tokens for authenticated principals.
type TokenIssuer interface {
	Issue(ctx context.Context, subject, email string, roles ...string) (token string, expiresAt time.Time, err error)
}

type jwtIssuer struct {
	secret   []byte
	issuer   string
	audience []string
	ttl      time.Duration
	disabled bool
	now      func() time.Time
}

// NewJWTIssuer constructs a HS256 token issuer based on application configuration.
func NewJWTIssuer(cfg *config.AppConfig) (TokenIssuer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("jwt issuer: missing config")
	}

	authCfg := cfg.Auth
	if len(authCfg.JWTSecret) == 0 {
		return nil, fmt.Errorf("jwt issuer: jwt_secret must be configured")
	}

	ttl := time.Duration(authCfg.AccessTokenTTLMinutes) * time.Minute
	if ttl <= 0 {
		ttl = 60 * time.Minute
	}

	return &jwtIssuer{
		secret:   []byte(authCfg.JWTSecret),
		issuer:   authCfg.Issuer,
		audience: authCfg.Audience,
		ttl:      ttl,
		disabled: authCfg.Disabled,
		now:      time.Now,
	}, nil
}

func (i *jwtIssuer) Issue(_ context.Context, subject, email string, roles ...string) (string, time.Time, error) {
	if i.disabled {
		expires := i.now().Add(24 * time.Hour)
		return "development-token", expires, nil
	}

	if subject == "" {
		return "", time.Time{}, errs.New(errs.CodeInternal, 500, "token issuer: subject is required", nil)
	}

	now := i.now()
	expiresAt := now.Add(i.ttl)

	headerBytes, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", time.Time{}, errs.New(errs.CodeInternal, 500, "token issuer: marshal header", err)
	}
	header := base64.RawURLEncoding.EncodeToString(headerBytes)

	payload := map[string]any{
		"iss": i.issuer,
		"sub": subject,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	if email != "" {
		payload["email"] = email
	}

	if len(i.audience) > 0 {
		payload["aud"] = i.audience
	}

	if normalized := normalizeRoles(roles); len(normalized) > 0 {
		payload["roles"] = normalized
		payload["role"] = normalized[0]
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, errs.New(errs.CodeInternal, 500, "token issuer: marshal payload", err)
	}
	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signature, err := i.sign(header, payloadSegment)
	if err != nil {
		return "", time.Time{}, err
	}

	return fmt.Sprintf("%s.%s.%s", header, payloadSegment, signature), expiresAt, nil
}

func (i *jwtIssuer) sign(header, payload string) (string, error) {
	mac := hmac.New(sha256.New, i.secret)
	if _, err := mac.Write([]byte(fmt.Sprintf("%s.%s", header, payload))); err != nil {
		return "", errs.New(errs.CodeInternal, 500, "token issuer: unable to sign token", err)
	}

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func normalizeRoles(roles []string) []string {
	if len(roles) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(roles))
	normalized := make([]string, 0, len(roles))

	for _, role := range roles {
		trimmed := strings.TrimSpace(role)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

var _ TokenIssuer = (*jwtIssuer)(nil)
