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

type projectDocumentRepository struct {
	db *sqlx.DB
}

// NewProjectDocumentRepository returns a repository backed by MySQL for project aggregates.
func NewProjectDocumentRepository(db *sqlx.DB) repository.ProjectDocumentRepository {
	return &projectDocumentRepository{db: db}
}

const listProjectDocumentsQuery = `
SELECT
    p.id,
    p.slug,
    p.title_ja,
    p.title_en,
    p.summary_ja,
    p.summary_en,
    p.description_ja,
    p.description_en,
    p.cover_image_url,
    p.primary_link_url,
    p.period_start,
    p.period_end,
    p.created_at,
    p.updated_at,
    p.published,
    p.highlight,
    p.sort_order
FROM projects p
%s
ORDER BY p.sort_order, p.created_at DESC, p.id DESC`

type projectDocumentRow struct {
	ID            uint64         `db:"id"`
	Slug          string         `db:"slug"`
	TitleJA       sql.NullString `db:"title_ja"`
	TitleEN       sql.NullString `db:"title_en"`
	SummaryJA     sql.NullString `db:"summary_ja"`
	SummaryEN     sql.NullString `db:"summary_en"`
	DescriptionJA sql.NullString `db:"description_ja"`
	DescriptionEN sql.NullString `db:"description_en"`
	CoverImageURL sql.NullString `db:"cover_image_url"`
	PrimaryLink   sql.NullString `db:"primary_link_url"`
	PeriodStart   sql.NullTime   `db:"period_start"`
	PeriodEnd     sql.NullTime   `db:"period_end"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	Published     bool           `db:"published"`
	Highlight     bool           `db:"highlight"`
	SortOrder     sql.NullInt64  `db:"sort_order"`
}

func (r *projectDocumentRepository) ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error) {
	where := ""
	if !includeDrafts {
		where = "WHERE p.published = 1"
	}

	query := fmt.Sprintf(listProjectDocumentsQuery, where)

	var rows []projectDocumentRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("select projects: %w", err)
	}

	if len(rows) == 0 {
		return []model.ProjectDocument{}, nil
	}

	documents := make([]model.ProjectDocument, 0, len(rows))
	index := make(map[uint64]*model.ProjectDocument, len(rows))
	ids := make([]uint64, 0, len(rows))

	for _, row := range rows {
		ids = append(ids, row.ID)
		document := model.ProjectDocument{
			ID:          row.ID,
			Slug:        strings.TrimSpace(row.Slug),
			Title:       toLocalizedText(row.TitleJA, row.TitleEN),
			Summary:     toLocalizedText(row.SummaryJA, row.SummaryEN),
			Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
			Period: model.ProjectPeriod{
				Start: nullableTime(row.PeriodStart),
				End:   nullableTime(row.PeriodEnd),
			},
			CoverImageURL: strings.TrimSpace(row.CoverImageURL.String),
			PrimaryLink:   strings.TrimSpace(row.PrimaryLink.String),
			Highlight:     row.Highlight,
			Published:     row.Published,
			Links:         []model.ProjectLink{},
			Tech:          []model.TechMembership{},
		}
		if row.SortOrder.Valid {
			document.SortOrder = int(row.SortOrder.Int64)
		}
		if row.CreatedAt.Valid {
			document.CreatedAt = row.CreatedAt.Time.UTC()
		}
		if row.UpdatedAt.Valid {
			document.UpdatedAt = row.UpdatedAt.Time.UTC()
		}
		documents = append(documents, document)
		index[row.ID] = &documents[len(documents)-1]
	}

	if err := r.attachProjectLinks(ctx, ids, index); err != nil {
		return nil, err
	}
	if err := r.attachProjectTech(ctx, ids, index); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *projectDocumentRepository) attachProjectLinks(ctx context.Context, ids []uint64, projects map[uint64]*model.ProjectDocument) error {
	query, args, err := sqlx.In(`
SELECT
    id,
    project_id,
    link_type,
    label_ja,
    label_en,
    url,
    sort_order
FROM project_links
WHERE project_id IN (?)
ORDER BY project_id, sort_order, id`, ids)
	if err != nil {
		return fmt.Errorf("project links IN: %w", err)
	}

	query = r.db.Rebind(query)

	type linkRow struct {
		ID        uint64         `db:"id"`
		ProjectID uint64         `db:"project_id"`
		Type      string         `db:"link_type"`
		LabelJA   sql.NullString `db:"label_ja"`
		LabelEN   sql.NullString `db:"label_en"`
		URL       sql.NullString `db:"url"`
		SortOrder int            `db:"sort_order"`
	}

	var rows []linkRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("select project links: %w", err)
	}

	for _, row := range rows {
		project := projects[row.ProjectID]
		if project == nil {
			continue
		}
		project.Links = append(project.Links, model.ProjectLink{
			ID:        row.ID,
			ProjectID: row.ProjectID,
			Type:      model.ProjectLinkType(strings.TrimSpace(row.Type)),
			Label:     toLocalizedText(row.LabelJA, row.LabelEN),
			URL:       strings.TrimSpace(row.URL.String),
			SortOrder: row.SortOrder,
		})
	}

	return nil
}

func (r *projectDocumentRepository) attachProjectTech(ctx context.Context, ids []uint64, projects map[uint64]*model.ProjectDocument) error {
	query, args, err := sqlx.In(`
SELECT
    tr.id              AS membership_id,
    tr.entity_id       AS entity_id,
    tr.context         AS context,
    tr.note            AS note,
    tr.sort_order      AS membership_sort_order,
    tc.id              AS tech_id,
    tc.slug            AS tech_slug,
    tc.display_name    AS tech_display_name,
    tc.category        AS tech_category,
    tc.level           AS tech_level,
    tc.icon            AS tech_icon,
    tc.sort_order      AS tech_sort_order,
    tc.is_active       AS tech_is_active,
    tc.created_at      AS tech_created_at,
    tc.updated_at      AS tech_updated_at
FROM tech_relationships tr
JOIN tech_catalog tc ON tc.id = tr.tech_id
WHERE tr.entity_type = 'project' AND tr.entity_id IN (?)
ORDER BY tr.entity_id, tr.sort_order, tr.id`, ids)
	if err != nil {
		return fmt.Errorf("project tech IN: %w", err)
	}
	query = r.db.Rebind(query)

	type membershipRow struct {
		MembershipID    uint64         `db:"membership_id"`
		EntityID        uint64         `db:"entity_id"`
		Context         string         `db:"context"`
		Note            sql.NullString `db:"note"`
		SortOrder       int            `db:"membership_sort_order"`
		TechID          uint64         `db:"tech_id"`
		TechSlug        string         `db:"tech_slug"`
		TechDisplayName string         `db:"tech_display_name"`
		TechCategory    sql.NullString `db:"tech_category"`
		TechLevel       string         `db:"tech_level"`
		TechIcon        sql.NullString `db:"tech_icon"`
		TechSortOrder   int            `db:"tech_sort_order"`
		TechActive      bool           `db:"tech_is_active"`
		TechCreatedAt   sql.NullTime   `db:"tech_created_at"`
		TechUpdatedAt   sql.NullTime   `db:"tech_updated_at"`
	}

	var rows []membershipRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("select project tech: %w", err)
	}

	for _, row := range rows {
		project := projects[row.EntityID]
		if project == nil {
			continue
		}

		entry := model.TechCatalogEntry{
			ID:          row.TechID,
			Slug:        strings.TrimSpace(row.TechSlug),
			DisplayName: strings.TrimSpace(row.TechDisplayName),
			Category:    strings.TrimSpace(row.TechCategory.String),
			Level:       model.TechLevel(strings.TrimSpace(row.TechLevel)),
			Icon:        strings.TrimSpace(row.TechIcon.String),
			SortOrder:   row.TechSortOrder,
			Active:      row.TechActive,
		}
		if row.TechCreatedAt.Valid {
			entry.CreatedAt = row.TechCreatedAt.Time.UTC()
		}
		if row.TechUpdatedAt.Valid {
			entry.UpdatedAt = row.TechUpdatedAt.Time.UTC()
		}

		project.Tech = append(project.Tech, model.TechMembership{
			MembershipID: row.MembershipID,
			EntityType:   "project",
			EntityID:     row.EntityID,
			Tech:         entry,
			Context:      model.TechContext(strings.TrimSpace(row.Context)),
			Note:         strings.TrimSpace(row.Note.String),
			SortOrder:    row.SortOrder,
		})
	}

	return nil
}
