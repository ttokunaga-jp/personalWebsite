package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type blacklistRepository struct {
	db *sqlx.DB
}

// NewBlacklistRepository returns a MySQL-backed blacklist repository.
func NewBlacklistRepository(db *sqlx.DB) repository.BlacklistRepository {
	return &blacklistRepository{db: db}
}

const listBlacklistQuery = `
SELECT
	b.id,
	b.email,
	b.reason,
	b.created_at
FROM blacklist b
ORDER BY b.created_at DESC, b.id DESC`

const getBlacklistByEmailQuery = `
SELECT
	b.id,
	b.email,
	b.reason,
	b.created_at
FROM blacklist b
WHERE b.email = ?`

const getBlacklistByIDQuery = `
SELECT
	b.id,
	b.email,
	b.reason,
	b.created_at
FROM blacklist b
WHERE b.id = ?`

const insertBlacklistQuery = `
INSERT INTO blacklist (
	email,
	reason,
	created_at
)
VALUES (?, ?, NOW())`

const updateBlacklistQuery = `
UPDATE blacklist SET
	email = ?,
	reason = ?
WHERE id = ?`

const deleteBlacklistQuery = `DELETE FROM blacklist WHERE id = ?`

type blacklistRow struct {
	ID        int64          `db:"id"`
	Email     sql.NullString `db:"email"`
	Reason    sql.NullString `db:"reason"`
	CreatedAt sql.NullTime   `db:"created_at"`
}

func (r *blacklistRepository) ListBlacklistEntries(ctx context.Context) ([]model.BlacklistEntry, error) {
	var rows []blacklistRow
	if err := r.db.SelectContext(ctx, &rows, listBlacklistQuery); err != nil {
		return nil, fmt.Errorf("select blacklist: %w", err)
	}

	entries := make([]model.BlacklistEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, mapBlacklistRow(row))
	}
	return entries, nil
}

func (r *blacklistRepository) AddBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}

	email := strings.ToLower(strings.TrimSpace(entry.Email))
	if email == "" {
		return nil, repository.ErrInvalidInput
	}

	reason := strings.TrimSpace(entry.Reason)

	if _, err := r.FindBlacklistEntryByEmail(ctx, email); err == nil {
		return nil, repository.ErrDuplicate
	} else if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx, insertBlacklistQuery, email, reason)
	if err != nil {
		return nil, fmt.Errorf("insert blacklist entry: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("blacklist last insert id: %w", err)
	}

	created, err := r.FindBlacklistEntryByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	created.ID = id
	return created, nil
}

func (r *blacklistRepository) UpdateBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}
	if entry.ID == 0 {
		return nil, repository.ErrInvalidInput
	}

	email := strings.ToLower(strings.TrimSpace(entry.Email))
	if email == "" {
		return nil, repository.ErrInvalidInput
	}
	reason := strings.TrimSpace(entry.Reason)

	if existing, err := r.FindBlacklistEntryByEmail(ctx, email); err == nil && existing.ID != entry.ID {
		return nil, repository.ErrDuplicate
	} else if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx, updateBlacklistQuery, email, reason, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("update blacklist entry %d: %w", entry.ID, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected update blacklist %d: %w", entry.ID, err)
	}
	if affected == 0 {
		return nil, repository.ErrNotFound
	}

	return r.getBlacklistEntryByID(ctx, entry.ID)
}

func (r *blacklistRepository) RemoveBlacklistEntry(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, deleteBlacklistQuery, id)
	if err != nil {
		return fmt.Errorf("delete blacklist entry %d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected delete blacklist %d: %w", id, err)
	}
	if affected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *blacklistRepository) FindBlacklistEntryByEmail(ctx context.Context, email string) (*model.BlacklistEntry, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, repository.ErrInvalidInput
	}

	var row blacklistRow
	if err := r.db.GetContext(ctx, &row, getBlacklistByEmailQuery, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get blacklist by email %s: %w", email, err)
	}

	entry := mapBlacklistRow(row)
	return &entry, nil
}

func (r *blacklistRepository) getBlacklistEntryByID(ctx context.Context, id int64) (*model.BlacklistEntry, error) {
	var row blacklistRow
	if err := r.db.GetContext(ctx, &row, getBlacklistByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get blacklist by id %d: %w", id, err)
	}
	entry := mapBlacklistRow(row)
	return &entry, nil
}

func mapBlacklistRow(row blacklistRow) model.BlacklistEntry {
	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = timeNowUTC()
	}
	return model.BlacklistEntry{
		ID:        row.ID,
		Email:     row.Email.String,
		Reason:    row.Reason.String,
		CreatedAt: createdAt,
	}
}

var _ repository.BlacklistRepository = (*blacklistRepository)(nil)
