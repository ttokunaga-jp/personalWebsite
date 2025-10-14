package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/errs"
)

// Claims captures the subset of JWT fields required by the application.
type Claims struct {
	Subject   string
	Email     string
	Issuer    string
	Audience  []string
	ExpiresAt time.Time
	Raw       map[string]any
}

// TokenVerifier validates JWT tokens and yields the Claims.
type TokenVerifier interface {
	Verify(ctx context.Context, token string) (*Claims, error)
}

type jwtVerifier struct {
	secret    []byte
	issuer    string
	audience  []string
	clockSkew time.Duration
	disabled  bool
}

// NewJWTVerifier constructs a JWT HS256 verifier based on configuration.
func NewJWTVerifier(cfg config.AuthConfig) TokenVerifier {
	return &jwtVerifier{
		secret:    []byte(cfg.JWTSecret),
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
		clockSkew: time.Duration(cfg.ClockSkewSeconds) * time.Second,
		disabled:  cfg.Disabled,
	}
}

func (v *jwtVerifier) Verify(ctx context.Context, token string) (*Claims, error) {
	if v.disabled {
		return &Claims{
			Subject:   "anonymous",
			Email:     "",
			Issuer:    v.issuer,
			Audience:  v.audience,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Raw:       map[string]any{"disabled": true},
		}, nil
	}

	if strings.TrimSpace(token) == "" {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "missing token", nil)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "malformed token", nil)
	}

	headerBytes, err := decodeSegment(parts[0])
	if err != nil {
		return nil, wrapUnauthorized(err, "invalid token header encoding")
	}

	payloadBytes, err := decodeSegment(parts[1])
	if err != nil {
		return nil, wrapUnauthorized(err, "invalid token payload encoding")
	}

	if err := v.verifySignature(parts[0], parts[1], parts[2]); err != nil {
		return nil, err
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, wrapUnauthorized(err, "invalid token header")
	}
	if header.Alg != "HS256" {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "unsupported token algorithm", nil)
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, wrapUnauthorized(err, "invalid token payload")
	}

	claims, err := v.extractClaims(payload)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (v *jwtVerifier) extractClaims(payload map[string]any) (*Claims, error) {
	claims := &Claims{
		Raw: payload,
	}

	if sub, ok := payload["sub"].(string); ok {
		claims.Subject = sub
	}
	if email, ok := payload["email"].(string); ok {
		claims.Email = email
	}
	if iss, ok := payload["iss"].(string); ok {
		claims.Issuer = iss
	}
	rawAud, hasAud := payload["aud"]
	switch aud := rawAud.(type) {
	case string:
		if aud != "" {
			claims.Audience = []string{aud}
		}
	case []any:
		for _, entry := range aud {
			if s, ok := entry.(string); ok {
				claims.Audience = append(claims.Audience, s)
			}
		}
	case nil:
		// no-op
	default:
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid audience in token", nil)
	}

	exp, hasExp := payload["exp"]
	if hasExp {
		switch val := exp.(type) {
		case float64:
			claims.ExpiresAt = time.Unix(int64(val), 0)
		case int64:
			claims.ExpiresAt = time.Unix(val, 0)
		case json.Number:
			parsed, err := val.Int64()
			if err != nil {
				return nil, wrapUnauthorized(err, "invalid expiration claim")
			}
			claims.ExpiresAt = time.Unix(parsed, 0)
		default:
			return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid expiration claim", nil)
		}
	}

	if hasExp && time.Now().Add(-v.clockSkew).After(claims.ExpiresAt) {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "token expired", nil)
	}

	if !hasAud && len(v.audience) > 0 {
		return nil, errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "missing audience", nil)
	}

	return claims, nil
}

func (v *jwtVerifier) verifySignature(header, payload, signature string) error {
	mac := hmac.New(sha256.New, v.secret)
	_, err := mac.Write([]byte(fmt.Sprintf("%s.%s", header, payload)))
	if err != nil {
		return wrapUnauthorized(err, "unable to compute expected signature")
	}
	expected := mac.Sum(nil)

	actual, err := decodeSegment(signature)
	if err != nil {
		return wrapUnauthorized(err, "invalid token signature encoding")
	}

	if !hmac.Equal(expected, actual) {
		return errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, "invalid token signature", nil)
	}
	return nil
}

func decodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}

func wrapUnauthorized(err error, message string) *errs.AppError {
	if err == nil {
		return errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, message, nil)
	}
	return errs.New(errs.CodeUnauthorized, http.StatusUnauthorized, message, err)
}

// Ensure implementation satisfies the interface at compile time.
var _ TokenVerifier = (*jwtVerifier)(nil)
