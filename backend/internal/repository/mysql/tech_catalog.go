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

type techCatalogRepository struct {
	db *sqlx.DB
}

// NewTechCatalogRepository returns a MySQL-backed implementation of TechCatalogRepository.
func NewTechCatalogRepository(db *sqlx.DB) repository.TechCatalogRepository {
	return &techCatalogRepository{db: db}
}

const techCatalogSelectBase = `
SELECT
    id,
    slug,
    display_name,
    category,
    level,
    icon,
    sort_order,
    is_active,
    created_at,
    updated_at
FROM tech_catalog`

const listTechCatalogOrderClause = `
ORDER BY sort_order, id`

const getTechCatalogByIDQuery = techCatalogSelectBase + `
WHERE id = ?
LIMIT 1`

const getTechCatalogBySlugQuery = techCatalogSelectBase + `
WHERE slug = ?
LIMIT 1`

const insertTechCatalogQuery = `
INSERT INTO tech_catalog (
    slug,
    display_name,
    category,
    level,
    icon,
    sort_order,
    is_active,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3))`

const updateTechCatalogQuery = `
UPDATE tech_catalog
SET
    slug = ?,
    display_name = ?,
    category = ?,
    level = ?,
    icon = ?,
    sort_order = ?,
    is_active = ?,
    updated_at = NOW(3)
WHERE id = ?`

type techCatalogRow struct {
	ID          uint64         `db:"id"`
	Slug        string         `db:"slug"`
	DisplayName string         `db:"display_name"`
	Category    sql.NullString `db:"category"`
	Level       string         `db:"level"`
	Icon        sql.NullString `db:"icon"`
	SortOrder   int            `db:"sort_order"`
	Active      bool           `db:"is_active"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func (r *techCatalogRepository) ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(techCatalogSelectBase)
	if !includeInactive {
		queryBuilder.WriteString("\nWHERE is_active = 1")
	}
	queryBuilder.WriteString(listTechCatalogOrderClause)
	query := queryBuilder.String()

	var rows []techCatalogRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("select tech catalog: %w", err)
	}

	results := make([]model.TechCatalogEntry, 0, len(rows))
	for _, row := range rows {
		results = append(results, mapTechCatalogRow(row))
	}

	return results, nil
}

func (r *techCatalogRepository) GetTechCatalogEntry(ctx context.Context, id uint64) (*model.TechCatalogEntry, error) {
	if id == 0 {
		return nil, repository.ErrInvalidInput
	}

	var row techCatalogRow
	if err := r.db.GetContext(ctx, &row, getTechCatalogByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get tech catalog by id %d: %w", id, err)
	}
	entry := mapTechCatalogRow(row)
	return &entry, nil
}

func (r *techCatalogRepository) CreateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}

	slug := strings.TrimSpace(entry.Slug)
	displayName := strings.TrimSpace(entry.DisplayName)
	category := strings.TrimSpace(entry.Category)
	icon := strings.TrimSpace(entry.Icon)

	if slug == "" || displayName == "" {
		return nil, repository.ErrInvalidInput
	}
	if !isValidTechLevel(entry.Level) {
		return nil, repository.ErrInvalidInput
	}

	if existing, err := r.getTechCatalogBySlug(ctx, slug); err == nil && existing != nil {
		return nil, repository.ErrDuplicate
	} else if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx, insertTechCatalogQuery,
		slug,
		displayName,
		nullString(category),
		strings.TrimSpace(string(entry.Level)),
		nullString(icon),
		entry.SortOrder,
		entry.Active,
	)
	if err != nil {
		return nil, fmt.Errorf("insert tech catalog entry: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("tech catalog last insert id: %w", err)
	}

	return r.GetTechCatalogEntry(ctx, uint64(id))
}

func (r *techCatalogRepository) UpdateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}
	if entry.ID == 0 {
		return nil, repository.ErrInvalidInput
	}

	if _, err := r.GetTechCatalogEntry(ctx, entry.ID); err != nil {
		return nil, err
	}

	slug := strings.TrimSpace(entry.Slug)
	displayName := strings.TrimSpace(entry.DisplayName)
	category := strings.TrimSpace(entry.Category)
	icon := strings.TrimSpace(entry.Icon)

	if slug == "" || displayName == "" {
		return nil, repository.ErrInvalidInput
	}
	if !isValidTechLevel(entry.Level) {
		return nil, repository.ErrInvalidInput
	}

	if existing, err := r.getTechCatalogBySlug(ctx, slug); err == nil && existing.ID != entry.ID {
		return nil, repository.ErrDuplicate
	} else if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	_, err := r.db.ExecContext(ctx, updateTechCatalogQuery,
		slug,
		displayName,
		nullString(category),
		strings.TrimSpace(string(entry.Level)),
		nullString(icon),
		entry.SortOrder,
		entry.Active,
		entry.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update tech catalog entry %d: %w", entry.ID, err)
	}

	return r.GetTechCatalogEntry(ctx, entry.ID)
}

func (r *techCatalogRepository) getTechCatalogBySlug(ctx context.Context, slug string) (*model.TechCatalogEntry, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, repository.ErrInvalidInput
	}

	var row techCatalogRow
	if err := r.db.GetContext(ctx, &row, getTechCatalogBySlugQuery, slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get tech catalog by slug %s: %w", slug, err)
	}
	entry := mapTechCatalogRow(row)
	return &entry, nil
}

func mapTechCatalogRow(row techCatalogRow) model.TechCatalogEntry {
	entry := model.TechCatalogEntry{
		ID:          row.ID,
		Slug:        strings.TrimSpace(row.Slug),
		DisplayName: strings.TrimSpace(row.DisplayName),
		Category:    strings.TrimSpace(row.Category.String),
		Level:       model.TechLevel(strings.TrimSpace(row.Level)),
		Icon:        strings.TrimSpace(row.Icon.String),
		SortOrder:   row.SortOrder,
		Active:      row.Active,
	}
	if row.CreatedAt.Valid {
		entry.CreatedAt = row.CreatedAt.Time.UTC()
	}
	if row.UpdatedAt.Valid {
		entry.UpdatedAt = row.UpdatedAt.Time.UTC()
	}
	return entry
}

func isValidTechLevel(level model.TechLevel) bool {
	switch strings.TrimSpace(string(level)) {
	case string(model.TechLevelBeginner), string(model.TechLevelIntermediate), string(model.TechLevelAdvanced):
		return true
	default:
		return false
	}
}
