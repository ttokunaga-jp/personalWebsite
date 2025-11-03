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

type homePageConfigRepository struct {
	db *sqlx.DB
}

// NewHomePageConfigRepository returns a HomePageConfigRepository backed by MySQL.
func NewHomePageConfigRepository(db *sqlx.DB) repository.HomePageConfigRepository {
	return &homePageConfigRepository{db: db}
}

const selectHomePageConfigQuery = `
SELECT
    id,
    profile_id,
    hero_subtitle_ja,
    hero_subtitle_en,
    updated_at
FROM home_page_config
ORDER BY id
LIMIT 1`

type homeConfigRow struct {
	ID             uint64         `db:"id"`
	ProfileID      uint64         `db:"profile_id"`
	HeroSubtitleJA sql.NullString `db:"hero_subtitle_ja"`
	HeroSubtitleEN sql.NullString `db:"hero_subtitle_en"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}

func (r *homePageConfigRepository) GetHomePageConfig(ctx context.Context) (*model.HomePageConfigDocument, error) {
	var row homeConfigRow
	if err := r.db.GetContext(ctx, &row, selectHomePageConfigQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select home_page_config: %w", err)
	}

	config := &model.HomePageConfigDocument{
		ID:           row.ID,
		ProfileID:    row.ProfileID,
		HeroSubtitle: toLocalizedText(row.HeroSubtitleJA, row.HeroSubtitleEN),
		QuickLinks:   []model.HomeQuickLink{},
		ChipSources:  []model.HomeChipSource{},
	}
	if row.UpdatedAt.Valid {
		config.UpdatedAt = row.UpdatedAt.Time.UTC()
	}

	if err := r.attachHomeQuickLinks(ctx, row.ID, config); err != nil {
		return nil, err
	}
	if err := r.attachHomeChipSources(ctx, row.ID, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (r *homePageConfigRepository) attachHomeQuickLinks(ctx context.Context, configID uint64, config *model.HomePageConfigDocument) error {
	query := `
SELECT
    id,
    config_id,
    section,
    label_ja,
    label_en,
    description_ja,
    description_en,
    cta_ja,
    cta_en,
    target_url,
    sort_order
FROM home_quick_links
WHERE config_id = ?
ORDER BY sort_order, id`

	type quickLinkRow struct {
		ID            uint64         `db:"id"`
		ConfigID      uint64         `db:"config_id"`
		Section       string         `db:"section"`
		LabelJA       sql.NullString `db:"label_ja"`
		LabelEN       sql.NullString `db:"label_en"`
		DescriptionJA sql.NullString `db:"description_ja"`
		DescriptionEN sql.NullString `db:"description_en"`
		CTAJA         sql.NullString `db:"cta_ja"`
		CTAEN         sql.NullString `db:"cta_en"`
		TargetURL     sql.NullString `db:"target_url"`
		SortOrder     int            `db:"sort_order"`
	}

	var rows []quickLinkRow
	if err := r.db.SelectContext(ctx, &rows, query, configID); err != nil {
		return fmt.Errorf("select home_quick_links: %w", err)
	}

	for _, row := range rows {
		config.QuickLinks = append(config.QuickLinks, model.HomeQuickLink{
			ID:          row.ID,
			ConfigID:    row.ConfigID,
			Section:     strings.TrimSpace(row.Section),
			Label:       toLocalizedText(row.LabelJA, row.LabelEN),
			Description: toLocalizedText(row.DescriptionJA, row.DescriptionEN),
			CTA:         toLocalizedText(row.CTAJA, row.CTAEN),
			TargetURL:   strings.TrimSpace(row.TargetURL.String),
			SortOrder:   row.SortOrder,
		})
	}

	return nil
}

func (r *homePageConfigRepository) attachHomeChipSources(ctx context.Context, configID uint64, config *model.HomePageConfigDocument) error {
	query := `
SELECT
    id,
    config_id,
    source_type,
    limit_count,
    label_ja,
    label_en,
    sort_order
FROM home_chip_sources
WHERE config_id = ?
ORDER BY sort_order, id`

	type chipSourceRow struct {
		ID        uint64         `db:"id"`
		ConfigID  uint64         `db:"config_id"`
		Source    string         `db:"source_type"`
		Limit     int            `db:"limit_count"`
		LabelJA   sql.NullString `db:"label_ja"`
		LabelEN   sql.NullString `db:"label_en"`
		SortOrder int            `db:"sort_order"`
	}

	var rows []chipSourceRow
	if err := r.db.SelectContext(ctx, &rows, query, configID); err != nil {
		return fmt.Errorf("select home_chip_sources: %w", err)
	}

	for _, row := range rows {
		config.ChipSources = append(config.ChipSources, model.HomeChipSource{
			ID:        row.ID,
			ConfigID:  row.ConfigID,
			Source:    strings.TrimSpace(row.Source),
			Label:     toLocalizedText(row.LabelJA, row.LabelEN),
			Limit:     row.Limit,
			SortOrder: row.SortOrder,
		})
	}

	return nil
}
