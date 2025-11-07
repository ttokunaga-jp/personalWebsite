package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

const researchEntityType = "research_blog"

const listPublicResearchQuery = `
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
	r.published_at,
	r.is_draft
FROM research_blog_entries r
WHERE r.is_draft = 0
ORDER BY r.published_at DESC, r.id DESC`

const baseAdminResearchQuery = `
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
	r.highlight_image_url,
	r.image_alt_ja,
	r.image_alt_en,
	r.published_at,
	r.is_draft,
	r.created_at,
	r.updated_at
FROM research_blog_entries r
%s
ORDER BY r.published_at DESC, r.id DESC`

const insertResearchEntryQuery = `
INSERT INTO research_blog_entries (
	slug,
	kind,
	title_ja,
	title_en,
	overview_ja,
	overview_en,
	outcome_ja,
	outcome_en,
	outlook_ja,
	outlook_en,
	external_url,
	highlight_image_url,
	image_alt_ja,
	image_alt_en,
	published_at,
	is_draft,
	created_at,
	updated_at
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(3), NOW(3)
)`

const updateResearchEntryQuery = `
UPDATE research_blog_entries
SET
	slug = ?,
	kind = ?,
	title_ja = ?,
	title_en = ?,
	overview_ja = ?,
	overview_en = ?,
	outcome_ja = ?,
	outcome_en = ?,
	outlook_ja = ?,
	outlook_en = ?,
	external_url = ?,
	highlight_image_url = ?,
	image_alt_ja = ?,
	image_alt_en = ?,
	published_at = ?,
	is_draft = ?,
	updated_at = NOW(3)
WHERE id = ?`

const deleteResearchEntryQuery = `DELETE FROM research_blog_entries WHERE id = ?`

const deleteResearchTagsQuery = `DELETE FROM research_blog_tags WHERE entry_id = ?`
const deleteResearchLinksQuery = `DELETE FROM research_blog_links WHERE entry_id = ?`
const deleteResearchAssetsQuery = `DELETE FROM research_blog_assets WHERE entry_id = ?`
const deleteResearchTechQuery = `DELETE FROM tech_relationships WHERE entity_type = ? AND entity_id = ?`

const insertResearchTagQuery = `
INSERT INTO research_blog_tags (entry_id, tag, sort_order)
VALUES (?, ?, ?)`

const insertResearchLinkQuery = `
INSERT INTO research_blog_links (
	entry_id,
	link_type,
	label_ja,
	label_en,
	url,
	sort_order
) VALUES (?, ?, ?, ?, ?, ?)`

const insertResearchAssetQuery = `
INSERT INTO research_blog_assets (
	entry_id,
	asset_url,
	caption_ja,
	caption_en,
	sort_order
) VALUES (?, ?, ?, ?, ?)`

const insertResearchTechQuery = `
INSERT INTO tech_relationships (
	entity_type,
	entity_id,
	tech_id,
	context,
	note,
	sort_order,
	created_at
) VALUES (?, ?, ?, ?, ?, ?, NOW(3))`

type researchRepository struct {
	db *sqlx.DB
}

// NewResearchRepository returns a MySQL-backed research repository.
func NewResearchRepository(db *sqlx.DB) repository.ResearchRepository {
	return &researchRepository{db: db}
}

type researchEntryRow struct {
	ID                uint64         `db:"id"`
	Slug              string         `db:"slug"`
	Kind              string         `db:"kind"`
	TitleJA           string         `db:"title_ja"`
	TitleEN           string         `db:"title_en"`
	OverviewJA        sql.NullString `db:"overview_ja"`
	OverviewEN        sql.NullString `db:"overview_en"`
	OutcomeJA         sql.NullString `db:"outcome_ja"`
	OutcomeEN         sql.NullString `db:"outcome_en"`
	OutlookJA         sql.NullString `db:"outlook_ja"`
	OutlookEN         sql.NullString `db:"outlook_en"`
	ExternalURL       string         `db:"external_url"`
	HighlightImageURL sql.NullString `db:"highlight_image_url"`
	ImageAltJA        sql.NullString `db:"image_alt_ja"`
	ImageAltEN        sql.NullString `db:"image_alt_en"`
	PublishedAt       time.Time      `db:"published_at"`
	IsDraft           bool           `db:"is_draft"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	var rows []struct {
		ID          uint64         `db:"id"`
		TitleJA     string         `db:"title_ja"`
		TitleEN     string         `db:"title_en"`
		OverviewJA  sql.NullString `db:"overview_ja"`
		OverviewEN  sql.NullString `db:"overview_en"`
		OutcomeJA   sql.NullString `db:"outcome_ja"`
		OutcomeEN   sql.NullString `db:"outcome_en"`
		PublishedAt time.Time      `db:"published_at"`
	}

	if err := r.db.SelectContext(ctx, &rows, listPublicResearchQuery); err != nil {
		return nil, fmt.Errorf("select research_blog_entries: %w", err)
	}

	if len(rows) == 0 {
		return []model.Research{}, nil
	}

	research := make([]model.Research, 0, len(rows))
	for _, row := range rows {
		year := row.PublishedAt.Year()
		if year <= 0 {
			year = time.Now().Year()
		}

		id, err := safeUintToInt64(row.ID)
		if err != nil {
			return nil, err
		}

		research = append(research, model.Research{
			ID:        id,
			Year:      year,
			Title:     model.NewLocalizedText(row.TitleJA, row.TitleEN),
			Summary:   toLocalizedText(row.OverviewJA, row.OverviewEN),
			ContentMD: toLocalizedText(row.OutcomeJA, row.OutcomeEN),
		})
	}
	return research, nil
}

func (r *researchRepository) ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error) {
	return r.loadAdminResearch(ctx, nil)
}

func (r *researchRepository) GetAdminResearch(ctx context.Context, id uint64) (*model.AdminResearch, error) {
	results, err := r.loadAdminResearch(ctx, []uint64{id})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, repository.ErrNotFound
	}
	return &results[0], nil
}

func (r *researchRepository) CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("research create: begin tx: %w", err)
	}

	var entryID uint64
	if err := func() error {
		defer rollbackOnError(tx, &err)

		id, err := r.insertResearchEntry(ctx, tx, item)
		if err != nil {
			return err
		}
		entryID = id

		if err := r.replaceResearchRelations(ctx, tx, entryID, item); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("research create: commit: %w", err)
		}
		return nil
	}(); err != nil {
		return nil, err
	}

	return r.GetAdminResearch(ctx, entryID)
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}
	if item.ID == 0 {
		return nil, repository.ErrInvalidInput
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("research update: begin tx: %w", err)
	}

	if err := func() error {
		defer rollbackOnError(tx, &err)

		res, err := tx.ExecContext(ctx, updateResearchEntryQuery,
			strings.TrimSpace(item.Slug),
			item.Kind,
			strings.TrimSpace(item.Title.Ja),
			strings.TrimSpace(item.Title.En),
			nullString(item.Overview.Ja),
			nullString(item.Overview.En),
			nullString(item.Outcome.Ja),
			nullString(item.Outcome.En),
			nullString(item.Outlook.Ja),
			nullString(item.Outlook.En),
			strings.TrimSpace(item.ExternalURL),
			nullString(item.HighlightImageURL),
			nullString(item.ImageAlt.Ja),
			nullString(item.ImageAlt.En),
			item.PublishedAt.UTC(),
			item.IsDraft,
			item.ID,
		)
		if err != nil {
			return fmt.Errorf("research update: update entry %d: %w", item.ID, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("research update: rows affected %d: %w", item.ID, err)
		}
		if affected == 0 {
			return repository.ErrNotFound
		}

		if err := r.replaceResearchRelations(ctx, tx, item.ID, item); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("research update: commit %d: %w", item.ID, err)
		}
		return nil
	}(); err != nil {
		return nil, err
	}

	return r.GetAdminResearch(ctx, item.ID)
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id uint64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("research delete: begin tx: %w", err)
	}

	return func() error {
		defer rollbackOnError(tx, &err)

		if _, err := tx.ExecContext(ctx, deleteResearchTechQuery, researchEntityType, id); err != nil {
			return fmt.Errorf("research delete: delete tech %d: %w", id, err)
		}

		res, err := tx.ExecContext(ctx, deleteResearchEntryQuery, id)
		if err != nil {
			return fmt.Errorf("research delete: delete entry %d: %w", id, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("research delete: rows affected %d: %w", id, err)
		}
		if affected == 0 {
			return repository.ErrNotFound
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("research delete: commit %d: %w", id, err)
		}
		return nil
	}()
}

func (r *researchRepository) loadAdminResearch(ctx context.Context, ids []uint64) ([]model.AdminResearch, error) {
	clause := ""
	query := fmt.Sprintf(baseAdminResearchQuery, clause)
	var args []interface{}

	if len(ids) > 0 {
		clause = "WHERE r.id IN (?)"
		var err error
		query, args, err = sqlx.In(fmt.Sprintf(baseAdminResearchQuery, clause), ids)
		if err != nil {
			return nil, fmt.Errorf("research load IN: %w", err)
		}
		query = r.db.Rebind(query)
	}

	return r.queryAdminResearch(ctx, query, args...)
}

func (r *researchRepository) queryAdminResearch(ctx context.Context, query string, args ...interface{}) ([]model.AdminResearch, error) {
	var rows []researchEntryRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("select admin research_blog_entries: %w", err)
	}
	if len(rows) == 0 {
		return []model.AdminResearch{}, nil
	}

	result := make([]model.AdminResearch, len(rows))
	documents := make(map[uint64]*model.ResearchDocument, len(rows))

	for i, row := range rows {
		admin := &result[i]
		*admin = model.AdminResearch{
			ID:                row.ID,
			Slug:              strings.TrimSpace(row.Slug),
			Kind:              model.ResearchKind(strings.TrimSpace(row.Kind)),
			Title:             model.NewLocalizedText(row.TitleJA, row.TitleEN),
			Overview:          toLocalizedText(row.OverviewJA, row.OverviewEN),
			Outcome:           toLocalizedText(row.OutcomeJA, row.OutcomeEN),
			Outlook:           toLocalizedText(row.OutlookJA, row.OutlookEN),
			ExternalURL:       strings.TrimSpace(row.ExternalURL),
			HighlightImageURL: nullableString(row.HighlightImageURL),
			ImageAlt:          toLocalizedText(row.ImageAltJA, row.ImageAltEN),
			PublishedAt:       row.PublishedAt.UTC(),
			IsDraft:           row.IsDraft,
			CreatedAt:         row.CreatedAt.UTC(),
			UpdatedAt:         row.UpdatedAt.UTC(),
		}

		document := &model.ResearchDocument{
			ID:                row.ID,
			Slug:              admin.Slug,
			Kind:              admin.Kind,
			Title:             admin.Title,
			Overview:          admin.Overview,
			Outcome:           admin.Outcome,
			Outlook:           admin.Outlook,
			ExternalURL:       admin.ExternalURL,
			PublishedAt:       admin.PublishedAt,
			UpdatedAt:         admin.UpdatedAt,
			HighlightImageURL: admin.HighlightImageURL,
			ImageAlt:          admin.ImageAlt,
			IsDraft:           admin.IsDraft,
			Tags:              []model.ResearchTag{},
			Links:             []model.ResearchLink{},
			Assets:            []model.ResearchAsset{},
			Tech:              []model.TechMembership{},
		}

		documents[row.ID] = document
	}

	ids := make([]uint64, 0, len(documents))
	for id := range documents {
		ids = append(ids, id)
	}

	docRepo := &researchDocumentRepository{db: r.db}
	if err := docRepo.attachResearchTags(ctx, ids, documents); err != nil {
		return nil, err
	}
	if err := docRepo.attachResearchLinks(ctx, ids, documents); err != nil {
		return nil, err
	}
	if err := docRepo.attachResearchAssets(ctx, ids, documents); err != nil {
		return nil, err
	}
	if err := docRepo.attachResearchTech(ctx, ids, documents); err != nil {
		return nil, err
	}

	for i := range result {
		doc := documents[result[i].ID]
		if doc == nil {
			continue
		}
		result[i].Tags = append([]model.ResearchTag(nil), doc.Tags...)
		result[i].Links = append([]model.ResearchLink(nil), doc.Links...)
		result[i].Assets = append([]model.ResearchAsset(nil), doc.Assets...)
		result[i].Tech = append([]model.TechMembership(nil), doc.Tech...)
	}

	return result, nil
}

func (r *researchRepository) insertResearchEntry(ctx context.Context, tx *sqlx.Tx, item *model.AdminResearch) (uint64, error) {
	res, err := tx.ExecContext(ctx, insertResearchEntryQuery,
		strings.TrimSpace(item.Slug),
		item.Kind,
		strings.TrimSpace(item.Title.Ja),
		strings.TrimSpace(item.Title.En),
		nullString(item.Overview.Ja),
		nullString(item.Overview.En),
		nullString(item.Outcome.Ja),
		nullString(item.Outcome.En),
		nullString(item.Outlook.Ja),
		nullString(item.Outlook.En),
		strings.TrimSpace(item.ExternalURL),
		nullString(item.HighlightImageURL),
		nullString(item.ImageAlt.Ja),
		nullString(item.ImageAlt.En),
		item.PublishedAt.UTC(),
		item.IsDraft,
	)
	if err != nil {
		return 0, fmt.Errorf("insert research_blog_entries: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("research entry last insert id: %w", err)
	}
	return uint64(id), nil
}

func (r *researchRepository) replaceResearchRelations(ctx context.Context, tx *sqlx.Tx, entryID uint64, item *model.AdminResearch) error {
	if _, err := tx.ExecContext(ctx, deleteResearchTagsQuery, entryID); err != nil {
		return fmt.Errorf("research relations: delete tags %d: %w", entryID, err)
	}
	if _, err := tx.ExecContext(ctx, deleteResearchLinksQuery, entryID); err != nil {
		return fmt.Errorf("research relations: delete links %d: %w", entryID, err)
	}
	if _, err := tx.ExecContext(ctx, deleteResearchAssetsQuery, entryID); err != nil {
		return fmt.Errorf("research relations: delete assets %d: %w", entryID, err)
	}
	if _, err := tx.ExecContext(ctx, deleteResearchTechQuery, researchEntityType, entryID); err != nil {
		return fmt.Errorf("research relations: delete tech %d: %w", entryID, err)
	}

	if err := r.insertResearchTags(ctx, tx, entryID, item.Tags); err != nil {
		return err
	}
	if err := r.insertResearchLinks(ctx, tx, entryID, item.Links); err != nil {
		return err
	}
	if err := r.insertResearchAssets(ctx, tx, entryID, item.Assets); err != nil {
		return err
	}
	if err := r.insertResearchTech(ctx, tx, entryID, item.Tech); err != nil {
		return err
	}
	return nil
}

func (r *researchRepository) insertResearchTags(ctx context.Context, tx *sqlx.Tx, entryID uint64, tags []model.ResearchTag) error {
	for _, tag := range tags {
		if strings.TrimSpace(tag.Value) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, insertResearchTagQuery, entryID, strings.TrimSpace(tag.Value), tag.SortOrder); err != nil {
			return fmt.Errorf("insert research tag %d: %w", entryID, err)
		}
	}
	return nil
}

func (r *researchRepository) insertResearchLinks(ctx context.Context, tx *sqlx.Tx, entryID uint64, links []model.ResearchLink) error {
	for _, link := range links {
		if strings.TrimSpace(link.URL) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, insertResearchLinkQuery,
			entryID,
			link.Type,
			nullString(link.Label.Ja),
			nullString(link.Label.En),
			strings.TrimSpace(link.URL),
			link.SortOrder,
		); err != nil {
			return fmt.Errorf("insert research link %d: %w", entryID, err)
		}
	}
	return nil
}

func (r *researchRepository) insertResearchAssets(ctx context.Context, tx *sqlx.Tx, entryID uint64, assets []model.ResearchAsset) error {
	for _, asset := range assets {
		if strings.TrimSpace(asset.URL) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, insertResearchAssetQuery,
			entryID,
			strings.TrimSpace(asset.URL),
			nullString(asset.Caption.Ja),
			nullString(asset.Caption.En),
			asset.SortOrder,
		); err != nil {
			return fmt.Errorf("insert research asset %d: %w", entryID, err)
		}
	}
	return nil
}

func (r *researchRepository) insertResearchTech(ctx context.Context, tx *sqlx.Tx, entryID uint64, tech []model.TechMembership) error {
	for _, membership := range tech {
		if membership.Tech.ID == 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, insertResearchTechQuery,
			researchEntityType,
			entryID,
			membership.Tech.ID,
			membership.Context,
			nullString(membership.Note),
			membership.SortOrder,
		); err != nil {
			return fmt.Errorf("insert research tech %d: %w", entryID, err)
		}
	}
	return nil
}

func safeUintToInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("value %d exceeds int64 range", value)
	}
	return int64(value), nil
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
