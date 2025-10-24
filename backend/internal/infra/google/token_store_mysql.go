package google

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"
)

type mysqlTokenStore struct {
	db  *sqlx.DB
	key []byte
}

// NewMySQLTokenStore returns a TokenStore backed by the supplied MySQL connection.
// The secret is used to derive an AES-256 key for encrypting token payloads at rest.
func NewMySQLTokenStore(db *sqlx.DB, secret string) (TokenStore, error) {
	if db == nil {
		return nil, fmt.Errorf("mysql token store: db is nil")
	}
	if secret == "" {
		return nil, fmt.Errorf("mysql token store: secret is empty")
	}
	sum := sha256.Sum256([]byte(secret))
	return &mysqlTokenStore{
		db:  db,
		key: sum[:],
	}, nil
}

func (s *mysqlTokenStore) Load(ctx context.Context, provider string) (*TokenRecord, error) {
	const query = `
SELECT
	access_token,
	refresh_token,
	expiry
FROM
	google_oauth_tokens
WHERE
	provider = ?
LIMIT 1`

	var row struct {
		AccessToken  string    `db:"access_token"`
		RefreshToken string    `db:"refresh_token"`
		Expiry       time.Time `db:"expiry"`
	}

	if err := s.db.GetContext(ctx, &row, query, provider); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("mysql token store: select token: %w", err)
	}

	access, err := s.decrypt(row.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("mysql token store: decrypt access token: %w", err)
	}

	refresh := ""
	if row.RefreshToken != "" {
		refresh, err = s.decrypt(row.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("mysql token store: decrypt refresh token: %w", err)
		}
	}

	return &TokenRecord{
		AccessToken:  access,
		RefreshToken: refresh,
		Expiry:       row.Expiry.UTC(),
	}, nil
}

func (s *mysqlTokenStore) Save(ctx context.Context, provider string, record *TokenRecord) error {
	if record == nil {
		return fmt.Errorf("mysql token store: record is nil")
	}

	encAccess, err := s.encrypt(record.AccessToken)
	if err != nil {
		return fmt.Errorf("mysql token store: encrypt access token: %w", err)
	}

	encRefresh := ""
	if record.RefreshToken != "" {
		encRefresh, err = s.encrypt(record.RefreshToken)
		if err != nil {
			return fmt.Errorf("mysql token store: encrypt refresh token: %w", err)
		}
	}

	const stmt = `
INSERT INTO google_oauth_tokens (
	provider,
	access_token,
	refresh_token,
	expiry,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE
	access_token = VALUES(access_token),
	refresh_token = CASE
		WHEN VALUES(refresh_token) = '' THEN refresh_token
		ELSE VALUES(refresh_token)
	END,
	expiry = VALUES(expiry),
	updated_at = NOW()`

	_, err = s.db.ExecContext(ctx, stmt,
		provider,
		encAccess,
		encRefresh,
		record.Expiry.UTC(),
	)
	if err != nil {
		return fmt.Errorf("mysql token store: upsert token: %w", err)
	}
	return nil
}

func (s *mysqlTokenStore) encrypt(value string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

func (s *mysqlTokenStore) decrypt(value string) (string, error) {
	data, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
