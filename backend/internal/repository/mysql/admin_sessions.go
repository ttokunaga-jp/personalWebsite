package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type adminSessionRepository struct {
	db *sqlx.DB
}

// NewAdminSessionRepository returns a MySQL-backed admin session repository.
func NewAdminSessionRepository(db *sqlx.DB) repository.AdminSessionRepository {
	repo := &adminSessionRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		log.Printf("admin session repository: schema ensure failed: %v", err)
	}
	return repo
}

const insertAdminSessionQuery = `
INSERT INTO admin_sessions (
	session_id_hash,
	subject,
	email,
	roles,
	user_agent,
	ip_address,
	expires_at,
	last_accessed_at,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`

const selectAdminSessionQuery = `
SELECT
	id,
	session_id_hash,
	subject,
	email,
	roles,
	user_agent,
	ip_address,
	expires_at,
	last_accessed_at,
	created_at,
	updated_at,
	revoked_at
FROM admin_sessions
WHERE session_id_hash = ?`

const updateAdminSessionActivityQuery = `
UPDATE admin_sessions
SET
	last_accessed_at = ?,
	expires_at = ?,
	updated_at = NOW(3)
WHERE session_id_hash = ? AND revoked_at IS NULL`

const revokeAdminSessionQuery = `
UPDATE admin_sessions
SET revoked_at = NOW(3), updated_at = NOW(3)
WHERE session_id_hash = ?`

type adminSessionRow struct {
	ID             uint64         `db:"id"`
	SessionIDHash  string         `db:"session_id_hash"`
	Subject        string         `db:"subject"`
	Email          string         `db:"email"`
	Roles          []byte         `db:"roles"`
	UserAgent      sql.NullString `db:"user_agent"`
	IPAddress      sql.NullString `db:"ip_address"`
	ExpiresAt      time.Time      `db:"expires_at"`
	LastAccessedAt time.Time      `db:"last_accessed_at"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	RevokedAt      sql.NullTime   `db:"revoked_at"`
}

func (r *adminSessionRepository) CreateSession(ctx context.Context, session *model.AdminSession) (*model.AdminSession, error) {
	if session == nil {
		return nil, fmt.Errorf("create session: %w", repository.ErrInvalidInput)
	}
	roles, err := json.Marshal(session.Roles)
	if err != nil {
		return nil, fmt.Errorf("marshal roles: %w", err)
	}

	userAgent := sql.NullString{Valid: session.UserAgent != "", String: session.UserAgent}
	ipAddress := sql.NullString{Valid: session.IPAddress != "", String: session.IPAddress}

	if _, err := r.db.ExecContext(
		ctx,
		insertAdminSessionQuery,
		session.TokenHash,
		session.Subject,
		session.Email,
		roles,
		userAgent,
		ipAddress,
		session.ExpiresAt.UTC(),
		session.LastAccessedAt.UTC(),
	); err != nil {
		return nil, fmt.Errorf("insert admin session: %w", err)
	}

	return r.FindSessionByHash(ctx, session.TokenHash)
}

func (r *adminSessionRepository) FindSessionByHash(ctx context.Context, hash string) (*model.AdminSession, error) {
	var row adminSessionRow
	if err := r.db.GetContext(ctx, &row, selectAdminSessionQuery, hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select admin session: %w", err)
	}
	return mapAdminSessionRow(row)
}

func (r *adminSessionRepository) UpdateSessionActivity(ctx context.Context, hash string, lastAccessed time.Time, expiresAt time.Time) (*model.AdminSession, error) {
	res, err := r.db.ExecContext(ctx, updateAdminSessionActivityQuery, lastAccessed.UTC(), expiresAt.UTC(), hash)
	if err != nil {
		return nil, fmt.Errorf("update admin session: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update admin session rows: %w", err)
	}
	if rows == 0 {
		return nil, repository.ErrNotFound
	}
	return r.FindSessionByHash(ctx, hash)
}

func (r *adminSessionRepository) RevokeSession(ctx context.Context, hash string) error {
	res, err := r.db.ExecContext(ctx, revokeAdminSessionQuery, hash)
	if err != nil {
		return fmt.Errorf("revoke admin session: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke admin session rows: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *adminSessionRepository) ensureSchema(ctx context.Context) error {
	if r.db == nil {
		return fmt.Errorf("admin session repo: nil db")
	}
	const ddl = `
CREATE TABLE IF NOT EXISTS admin_sessions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  session_id_hash CHAR(64) NOT NULL UNIQUE,
  subject VARCHAR(255) NOT NULL,
  email VARCHAR(320) NOT NULL,
  roles JSON NOT NULL,
  user_agent VARCHAR(512) NULL,
  ip_address VARCHAR(45) NULL,
  expires_at DATETIME(3) NOT NULL,
  last_accessed_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  revoked_at DATETIME(3) NULL,
  INDEX idx_admin_sessions_expires (expires_at),
  INDEX idx_admin_sessions_last_accessed (last_accessed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := r.db.ExecContext(ctx, ddl)
	return err
}

func mapAdminSessionRow(row adminSessionRow) (*model.AdminSession, error) {
	var roles []string
	if len(row.Roles) > 0 {
		if err := json.Unmarshal(row.Roles, &roles); err != nil {
			return nil, fmt.Errorf("decode roles: %w", err)
		}
	}
	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		val := row.RevokedAt.Time
		revokedAt = &val
	}
	return &model.AdminSession{
		TokenHash:      row.SessionIDHash,
		Subject:        row.Subject,
		Email:          row.Email,
		Roles:          roles,
		UserAgent:      row.UserAgent.String,
		IPAddress:      row.IPAddress.String,
		ExpiresAt:      row.ExpiresAt,
		LastAccessedAt: row.LastAccessedAt,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		RevokedAt:      revokedAt,
	}, nil
}
