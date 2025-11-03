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

const listTechCatalogQuery = `
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
FROM tech_catalog
%s
ORDER BY sort_order, id`

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
	where := ""
	if !includeInactive {
		where = "WHERE is_active = 1"
	}
	query := fmt.Sprintf(listTechCatalogQuery, where)

	var rows []techCatalogRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("select tech catalog: %w", err)
	}

	results := make([]model.TechCatalogEntry, 0, len(rows))
	for _, row := range rows {
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
		results = append(results, entry)
	}

	return results, nil
}
