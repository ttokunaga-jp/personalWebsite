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

type researchDocumentRepository struct {
	db *sqlx.DB
}

// NewResearchDocumentRepository returns a MySQL-backed implementation for research/blog aggregates.
func NewResearchDocumentRepository(db *sqlx.DB) repository.ResearchDocumentRepository {
	return &researchDocumentRepository{db: db}
}

const listResearchDocumentsQuery = `
SELECT
    r.id,
    r.slug,
    r.kind,
    r.title_ja,
    r.title_en,
    r.overview_ja,
    r.overview_en,
    r.outcome_ja,
    r.outcome_en,
    r.outlook_ja,
    r.outlook_en,
    r.external_url,
    r.published_at,
    r.updated_at,
    r.highlight_image_url,
    r.image_alt_ja,
    r.image_alt_en,
    r.is_draft
FROM research_blog_entries r
%s
ORDER BY r.published_at DESC, r.id DESC`

type researchDocumentRow struct {
	ID                uint64         `db:"id"`
	Slug              string         `db:"slug"`
	Kind              string         `db:"kind"`
	TitleJA           sql.NullString `db:"title_ja"`
	TitleEN           sql.NullString `db:"title_en"`
	OverviewJA        sql.NullString `db:"overview_ja"`
	OverviewEN        sql.NullString `db:"overview_en"`
	OutcomeJA         sql.NullString `db:"outcome_ja"`
	OutcomeEN         sql.NullString `db:"outcome_en"`
	OutlookJA         sql.NullString `db:"outlook_ja"`
	OutlookEN         sql.NullString `db:"outlook_en"`
	ExternalURL       sql.NullString `db:"external_url"`
	PublishedAt       sql.NullTime   `db:"published_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
	HighlightImageURL sql.NullString `db:"highlight_image_url"`
	ImageAltJA        sql.NullString `db:"image_alt_ja"`
	ImageAltEN        sql.NullString `db:"image_alt_en"`
	IsDraft           bool           `db:"is_draft"`
}

func (r *researchDocumentRepository) ListResearchDocuments(ctx context.Context, includeDrafts bool) ([]model.ResearchDocument, error) {
	where := ""
	if !includeDrafts {
		where = "WHERE r.is_draft = 0"
	}

	query := fmt.Sprintf(listResearchDocumentsQuery, where)

	var rows []researchDocumentRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("select research_blog_entries: %w", err)
	}

	if len(rows) == 0 {
		return []model.ResearchDocument{}, nil
	}

	documents := make([]model.ResearchDocument, 0, len(rows))
	index := make(map[uint64]*model.ResearchDocument, len(rows))
	ids := make([]uint64, 0, len(rows))

	for _, row := range rows {
		ids = append(ids, row.ID)

		document := model.ResearchDocument{
			ID:                row.ID,
			Slug:              strings.TrimSpace(row.Slug),
			Kind:              model.ResearchKind(strings.TrimSpace(row.Kind)),
			Title:             toLocalizedText(row.TitleJA, row.TitleEN),
			Overview:          toLocalizedText(row.OverviewJA, row.OverviewEN),
			Outcome:           toLocalizedText(row.OutcomeJA, row.OutcomeEN),
			Outlook:           toLocalizedText(row.OutlookJA, row.OutlookEN),
			ExternalURL:       strings.TrimSpace(row.ExternalURL.String),
			HighlightImageURL: strings.TrimSpace(row.HighlightImageURL.String),
			ImageAlt:          toLocalizedText(row.ImageAltJA, row.ImageAltEN),
			IsDraft:           row.IsDraft,
			Tags:              []model.ResearchTag{},
			Links:             []model.ResearchLink{},
			Assets:            []model.ResearchAsset{},
			Tech:              []model.TechMembership{},
		}

		if row.PublishedAt.Valid {
			document.PublishedAt = row.PublishedAt.Time.UTC()
		}
		if row.UpdatedAt.Valid {
			document.UpdatedAt = row.UpdatedAt.Time.UTC()
		}

		documents = append(documents, document)
		index[row.ID] = &documents[len(documents)-1]
	}

	if err := r.attachResearchTags(ctx, ids, index); err != nil {
		return nil, err
	}
	if err := r.attachResearchLinks(ctx, ids, index); err != nil {
		return nil, err
	}
	if err := r.attachResearchAssets(ctx, ids, index); err != nil {
		return nil, err
	}
	if err := r.attachResearchTech(ctx, ids, index); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *researchDocumentRepository) attachResearchTags(ctx context.Context, ids []uint64, documents map[uint64]*model.ResearchDocument) error {
	query, args, err := sqlx.In(`
SELECT
    id,
    entry_id,
    tag,
    sort_order
FROM research_blog_tags
WHERE entry_id IN (?)
ORDER BY entry_id, sort_order, id`, ids)
	if err != nil {
		return fmt.Errorf("research tags IN: %w", err)
	}
	query = r.db.Rebind(query)

	type tagRow struct {
		ID        uint64 `db:"id"`
		EntryID   uint64 `db:"entry_id"`
		Value     string `db:"tag"`
		SortOrder int    `db:"sort_order"`
	}

	var rows []tagRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("select research_blog_tags: %w", err)
	}

	for _, row := range rows {
		doc := documents[row.EntryID]
		if doc == nil {
			continue
		}
		doc.Tags = append(doc.Tags, model.ResearchTag{
			ID:        row.ID,
			EntryID:   row.EntryID,
			Value:     strings.TrimSpace(row.Value),
			SortOrder: row.SortOrder,
		})
	}

	return nil
}

func (r *researchDocumentRepository) attachResearchLinks(ctx context.Context, ids []uint64, documents map[uint64]*model.ResearchDocument) error {
	query, args, err := sqlx.In(`
SELECT
    id,
    entry_id,
    link_type,
    label_ja,
    label_en,
    url,
    sort_order
FROM research_blog_links
WHERE entry_id IN (?)
ORDER BY entry_id, sort_order, id`, ids)
	if err != nil {
		return fmt.Errorf("research links IN: %w", err)
	}
	query = r.db.Rebind(query)

	type linkRow struct {
		ID        uint64         `db:"id"`
		EntryID   uint64         `db:"entry_id"`
		Type      string         `db:"link_type"`
		LabelJA   sql.NullString `db:"label_ja"`
		LabelEN   sql.NullString `db:"label_en"`
		URL       sql.NullString `db:"url"`
		SortOrder int            `db:"sort_order"`
	}

	var rows []linkRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("select research_blog_links: %w", err)
	}

	for _, row := range rows {
		doc := documents[row.EntryID]
		if doc == nil {
			continue
		}
		doc.Links = append(doc.Links, model.ResearchLink{
			ID:        row.ID,
			EntryID:   row.EntryID,
			Type:      model.ResearchLinkType(strings.TrimSpace(row.Type)),
			Label:     toLocalizedText(row.LabelJA, row.LabelEN),
			URL:       strings.TrimSpace(row.URL.String),
			SortOrder: row.SortOrder,
		})
	}

	return nil
}

func (r *researchDocumentRepository) attachResearchAssets(ctx context.Context, ids []uint64, documents map[uint64]*model.ResearchDocument) error {
	query, args, err := sqlx.In(`
SELECT
    id,
    entry_id,
    asset_url,
    caption_ja,
    caption_en,
    sort_order
FROM research_blog_assets
WHERE entry_id IN (?)
ORDER BY entry_id, sort_order, id`, ids)
	if err != nil {
		return fmt.Errorf("research assets IN: %w", err)
	}
	query = r.db.Rebind(query)

	type assetRow struct {
		ID        uint64         `db:"id"`
		EntryID   uint64         `db:"entry_id"`
		URL       sql.NullString `db:"asset_url"`
		CaptionJA sql.NullString `db:"caption_ja"`
		CaptionEN sql.NullString `db:"caption_en"`
		SortOrder int            `db:"sort_order"`
	}

	var rows []assetRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return fmt.Errorf("select research_blog_assets: %w", err)
	}

	for _, row := range rows {
		doc := documents[row.EntryID]
		if doc == nil {
			continue
		}
		doc.Assets = append(doc.Assets, model.ResearchAsset{
			ID:        row.ID,
			EntryID:   row.EntryID,
			URL:       strings.TrimSpace(row.URL.String),
			Caption:   toLocalizedText(row.CaptionJA, row.CaptionEN),
			SortOrder: row.SortOrder,
		})
	}

	return nil
}

func (r *researchDocumentRepository) attachResearchTech(ctx context.Context, ids []uint64, documents map[uint64]*model.ResearchDocument) error {
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
WHERE tr.entity_type = 'research_blog' AND tr.entity_id IN (?)
ORDER BY tr.entity_id, tr.sort_order, tr.id`, ids)
	if err != nil {
		return fmt.Errorf("research tech IN: %w", err)
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
		return fmt.Errorf("select research tech: %w", err)
	}

	for _, row := range rows {
		doc := documents[row.EntityID]
		if doc == nil {
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

		doc.Tech = append(doc.Tech, model.TechMembership{
			MembershipID: row.MembershipID,
			EntityType:   "research_blog",
			EntityID:     row.EntityID,
			Tech:         entry,
			Context:      model.TechContext(strings.TrimSpace(row.Context)),
			Note:         strings.TrimSpace(row.Note.String),
			SortOrder:    row.SortOrder,
		})
	}

	return nil
}
