package google

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
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

	access, err := decryptToken(row.AccessToken, s.key)
	if err != nil {
		return nil, fmt.Errorf("mysql token store: decrypt access token: %w", err)
	}

	refresh := ""
	if row.RefreshToken != "" {
		refresh, err = decryptToken(row.RefreshToken, s.key)
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

	encAccess, err := encryptToken(record.AccessToken, s.key)
	if err != nil {
		return fmt.Errorf("mysql token store: encrypt access token: %w", err)
	}

	encRefresh := ""
	if record.RefreshToken != "" {
		encRefresh, err = encryptToken(record.RefreshToken, s.key)
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
