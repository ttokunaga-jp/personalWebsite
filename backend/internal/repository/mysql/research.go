package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchRepository struct {
	db *sqlx.DB
}

// NewResearchRepository returns a MySQL-backed research repository.
func NewResearchRepository(db *sqlx.DB) repository.ResearchRepository {
	return &researchRepository{db: db}
}

const listResearchQuery = `
SELECT
	r.id,
	r.year,
	r.title_ja,
	r.title_en,
	r.summary_ja,
	r.summary_en,
	r.content_md_ja,
	r.content_md_en
FROM research r
WHERE r.published = TRUE
ORDER BY COALESCE(r.sort_order, r.year * 1000), r.year DESC, r.id`

const listAdminResearchQuery = `
SELECT
	r.id,
	r.year,
	r.title_ja,
	r.title_en,
	r.summary_ja,
	r.summary_en,
	r.content_md_ja,
	r.content_md_en,
	r.published,
	r.sort_order,
	r.created_at,
	r.updated_at
FROM research r
ORDER BY COALESCE(r.sort_order, r.year * 1000), r.year DESC, r.id`

const getAdminResearchQuery = `
SELECT
	r.id,
	r.year,
	r.title_ja,
	r.title_en,
	r.summary_ja,
	r.summary_en,
	r.content_md_ja,
	r.content_md_en,
	r.published,
	r.sort_order,
	r.created_at,
	r.updated_at
FROM research r
WHERE r.id = ?`

const insertResearchQuery = `
INSERT INTO research (
	title_ja,
	title_en,
	summary_ja,
	summary_en,
	content_md_ja,
	content_md_en,
	year,
	published,
	sort_order,
	created_at,
	updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

const updateResearchQuery = `
UPDATE research
SET
	title_ja = ?,
	title_en = ?,
	summary_ja = ?,
	summary_en = ?,
	content_md_ja = ?,
	content_md_en = ?,
	year = ?,
	published = ?,
	sort_order = ?,
	updated_at = NOW()
WHERE id = ?`

const deleteResearchQuery = `DELETE FROM research WHERE id = ?`

type researchRow struct {
	ID        int64          `db:"id"`
	Year      int            `db:"year"`
	TitleJA   sql.NullString `db:"title_ja"`
	TitleEN   sql.NullString `db:"title_en"`
	SummaryJA sql.NullString `db:"summary_ja"`
	SummaryEN sql.NullString `db:"summary_en"`
	ContentJA sql.NullString `db:"content_md_ja"`
	ContentEN sql.NullString `db:"content_md_en"`
	Published sql.NullBool   `db:"published"`
	SortOrder sql.NullInt64  `db:"sort_order"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	var rows []researchRow
	if err := r.db.SelectContext(ctx, &rows, listResearchQuery); err != nil {
		return nil, fmt.Errorf("select research: %w", err)
	}

	research := make([]model.Research, 0, len(rows))
	for _, row := range rows {
		research = append(research, model.Research{
			ID:        row.ID,
			Year:      row.Year,
			Title:     toLocalizedText(row.TitleJA, row.TitleEN),
			Summary:   toLocalizedText(row.SummaryJA, row.SummaryEN),
			ContentMD: toLocalizedText(row.ContentJA, row.ContentEN),
		})
	}

	return research, nil
}

func (r *researchRepository) ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error) {
	var rows []researchRow
	if err := r.db.SelectContext(ctx, &rows, listAdminResearchQuery); err != nil {
		return nil, fmt.Errorf("select admin research: %w", err)
	}

	items := make([]model.AdminResearch, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapResearchRow(row))
	}
	return items, nil
}

func (r *researchRepository) GetAdminResearch(ctx context.Context, id int64) (*model.AdminResearch, error) {
	var row researchRow
	if err := r.db.GetContext(ctx, &row, getAdminResearchQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get research %d: %w", id, err)
	}

	mapped := mapResearchRow(row)
	return &mapped, nil
}

func (r *researchRepository) CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	res, err := r.db.ExecContext(ctx, insertResearchQuery,
		item.Title.Ja,
		item.Title.En,
		item.Summary.Ja,
		item.Summary.En,
		item.ContentMD.Ja,
		item.ContentMD.En,
		item.Year,
		item.Published,
		nullInt(itemSortOrderPtr(item)),
	)
	if err != nil {
		return nil, fmt.Errorf("insert research: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("research last insert id: %w", err)
	}

	created, err := r.GetAdminResearch(ctx, id)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	res, err := r.db.ExecContext(ctx, updateResearchQuery,
		item.Title.Ja,
		item.Title.En,
		item.Summary.Ja,
		item.Summary.En,
		item.ContentMD.Ja,
		item.ContentMD.En,
		item.Year,
		item.Published,
		nullInt(itemSortOrderPtr(item)),
		item.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update research %d: %w", item.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected research %d: %w", item.ID, err)
	}
	if affected == 0 {
		return nil, repository.ErrNotFound
	}

	updated, err := r.GetAdminResearch(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, deleteResearchQuery, id)
	if err != nil {
		return fmt.Errorf("delete research %d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected delete research %d: %w", id, err)
	}
	if affected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func mapResearchRow(row researchRow) model.AdminResearch {
	createdAt := row.CreatedAt.Time
	if !row.CreatedAt.Valid {
		createdAt = time.Time{}
	}
	updatedAt := row.UpdatedAt.Time
	if !row.UpdatedAt.Valid {
		updatedAt = time.Time{}
	}

	return model.AdminResearch{
		ID:        row.ID,
		Title:     toLocalizedText(row.TitleJA, row.TitleEN),
		Summary:   toLocalizedText(row.SummaryJA, row.SummaryEN),
		ContentMD: toLocalizedText(row.ContentJA, row.ContentEN),
		Year:      row.Year,
		Published: row.Published.Bool,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func itemSortOrderPtr(item *model.AdminResearch) *int {
	if item == nil {
		return nil
	}
	// Future enhancement: support explicit sort order on research items.
	return nil
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
