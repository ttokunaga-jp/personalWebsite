package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactFormSettingsRepository struct {
	db *sqlx.DB
}

// NewContactFormSettingsRepository returns a ContactFormSettingsRepository backed by MySQL.
func NewContactFormSettingsRepository(db *sqlx.DB) repository.ContactFormSettingsRepository {
	return &contactFormSettingsRepository{db: db}
}

const selectContactSettingsQuery = `
SELECT
    id,
    hero_title_ja,
    hero_title_en,
    hero_description_ja,
    hero_description_en,
    topics,
    consent_text_ja,
    consent_text_en,
    minimum_lead_hours,
    recaptcha_public_key,
    support_email,
    calendar_timezone,
    google_calendar_id,
    booking_window_days,
    created_at,
    updated_at
FROM contact_form_settings
ORDER BY id
LIMIT 1`

type contactSettingsRow struct {
	ID                uint64         `db:"id"`
	HeroTitleJA       sql.NullString `db:"hero_title_ja"`
	HeroTitleEN       sql.NullString `db:"hero_title_en"`
	HeroDescriptionJA sql.NullString `db:"hero_description_ja"`
	HeroDescriptionEN sql.NullString `db:"hero_description_en"`
	TopicsJSON        []byte         `db:"topics"`
	ConsentJA         sql.NullString `db:"consent_text_ja"`
	ConsentEN         sql.NullString `db:"consent_text_en"`
	MinimumLeadHours  int            `db:"minimum_lead_hours"`
	RecaptchaKey      sql.NullString `db:"recaptcha_public_key"`
	SupportEmail      sql.NullString `db:"support_email"`
	CalendarTimezone  sql.NullString `db:"calendar_timezone"`
	CalendarID        sql.NullString `db:"google_calendar_id"`
	BookingWindowDays int            `db:"booking_window_days"`
	CreatedAt         sql.NullTime   `db:"created_at"`
	UpdatedAt         sql.NullTime   `db:"updated_at"`
}

type contactTopicRow struct {
	ID          string        `json:"id"`
	Label       localizedJSON `json:"label"`
	Description localizedJSON `json:"description"`
}

type localizedJSON struct {
	Ja string `json:"ja"`
	En string `json:"en"`
}

func (r *contactFormSettingsRepository) GetContactFormSettings(ctx context.Context) (*model.ContactFormSettingsV2, error) {
	var row contactSettingsRow
	if err := r.db.GetContext(ctx, &row, selectContactSettingsQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("select contact_form_settings: %w", err)
	}

	topics, err := decodeContactTopics(row.TopicsJSON)
	if err != nil {
		return nil, fmt.Errorf("decode contact topics: %w", err)
	}

	settings := &model.ContactFormSettingsV2{
		ID:                row.ID,
		HeroTitle:         toLocalizedText(row.HeroTitleJA, row.HeroTitleEN),
		HeroDescription:   toLocalizedText(row.HeroDescriptionJA, row.HeroDescriptionEN),
		Topics:            topics,
		ConsentText:       toLocalizedText(row.ConsentJA, row.ConsentEN),
		MinimumLeadHours:  row.MinimumLeadHours,
		RecaptchaSiteKey:  strings.TrimSpace(row.RecaptchaKey.String),
		SupportEmail:      strings.TrimSpace(row.SupportEmail.String),
		CalendarTimezone:  strings.TrimSpace(row.CalendarTimezone.String),
		GoogleCalendarID:  strings.TrimSpace(row.CalendarID.String),
		BookingWindowDays: row.BookingWindowDays,
	}

	if row.CreatedAt.Valid {
		settings.CreatedAt = row.CreatedAt.Time.UTC()
	}
	if row.UpdatedAt.Valid {
		settings.UpdatedAt = row.UpdatedAt.Time.UTC()
	}

	return settings, nil
}

func decodeContactTopics(payload []byte) ([]model.ContactTopicV2, error) {
	if len(payload) == 0 {
		return []model.ContactTopicV2{}, nil
	}

	var rows []contactTopicRow
	if err := json.Unmarshal(payload, &rows); err != nil {
		return nil, err
	}

	topics := make([]model.ContactTopicV2, 0, len(rows))
	for _, row := range rows {
		topics = append(topics, model.ContactTopicV2{
			ID: row.ID,
			Label: model.LocalizedText{
				Ja: strings.TrimSpace(row.Label.Ja),
				En: strings.TrimSpace(row.Label.En),
			},
			Description: model.LocalizedText{
				Ja: strings.TrimSpace(row.Description.Ja),
				En: strings.TrimSpace(row.Description.En),
			},
		})
	}

	return topics, nil
}
