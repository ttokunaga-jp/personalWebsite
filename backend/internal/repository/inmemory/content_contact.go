package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactFormSettingsRepository struct {
	mu       sync.RWMutex
	settings *model.ContactFormSettingsV2
}

// NewContactFormSettingsRepository returns an in-memory contact form settings repository.
func NewContactFormSettingsRepository() repository.ContactFormSettingsRepository {
	now := time.Now().UTC()
	return &contactFormSettingsRepository{
		settings: &model.ContactFormSettingsV2{
			ID: 1,
			HeroTitle: model.NewLocalizedText(
				"お問い合わせ・予約",
				"Contact & Scheduling",
			),
			HeroDescription: model.NewLocalizedText(
				"研究・開発に関する相談や講演依頼を受け付けています。",
				"Reach out for research collaborations or speaking engagements.",
			),
			Topics: []model.ContactTopicV2{
				{
					ID: "consulting",
					Label: model.NewLocalizedText(
						"開発相談",
						"Consulting",
					),
					Description: model.NewLocalizedText(
						"技術選定やアーキテクチャ設計の支援。",
						"Support for technology choices and architecture design.",
					),
				},
				{
					ID:    "research",
					Label: model.NewLocalizedText("研究連携", "Research collaboration"),
					Description: model.NewLocalizedText(
						"共同研究や成果発表のご相談。",
						"Discuss collaboration or presenting research outcomes.",
					),
				},
			},
			ConsentText: model.NewLocalizedText(
				"送信によりプライバシーポリシーに同意したものとします。",
				"By submitting you agree to the privacy policy.",
			),
			MinimumLeadHours:   24,
			RecaptchaSiteKey:   "recaptcha-public-key",
			SupportEmail:       "support@example.dev",
			CalendarTimezone:   "Asia/Tokyo",
			GoogleCalendarID:   "primary",
			BookingWindowDays:  30,
			MeetingURLTemplate: "こんにちは {{guest_name}} さん。\n以下のリンクからミーティングにご参加ください: {{meeting_url}}\nよろしくお願いいたします。",
			CreatedAt:          now.AddDate(-1, 0, 0),
			UpdatedAt:          now.Add(-6 * time.Hour),
		},
	}
}

func (r *contactFormSettingsRepository) GetContactFormSettings(ctx context.Context) (*model.ContactFormSettingsV2, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.settings == nil {
		return nil, repository.ErrNotFound
	}
	return cloneContactSettings(r.settings), nil
}

func (r *contactFormSettingsRepository) UpdateContactFormSettings(ctx context.Context, settings *model.ContactFormSettingsV2, expectedUpdatedAt time.Time) (*model.ContactFormSettingsV2, error) {
	_ = ctx

	if settings == nil || settings.ID == 0 {
		return nil, repository.ErrInvalidInput
	}
	if expectedUpdatedAt.IsZero() {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.settings == nil || r.settings.ID != settings.ID {
		return nil, repository.ErrNotFound
	}
	if !r.settings.UpdatedAt.UTC().Equal(expectedUpdatedAt.UTC()) {
		return nil, repository.ErrConflict
	}

	updated := cloneContactSettings(settings)
	updated.CreatedAt = r.settings.CreatedAt
	updated.UpdatedAt = time.Now().UTC()

	r.settings = updated
	return cloneContactSettings(r.settings), nil
}

func cloneContactSettings(settings *model.ContactFormSettingsV2) *model.ContactFormSettingsV2 {
	if settings == nil {
		return nil
	}
	copyTopics := make([]model.ContactTopicV2, len(settings.Topics))
	for i, topic := range settings.Topics {
		copyTopics[i] = model.ContactTopicV2{
			ID: topic.ID,
			Label: model.LocalizedText{
				Ja: topic.Label.Ja,
				En: topic.Label.En,
			},
			Description: model.LocalizedText{
				Ja: topic.Description.Ja,
				En: topic.Description.En,
			},
		}
	}

	return &model.ContactFormSettingsV2{
		ID:                 settings.ID,
		HeroTitle:          model.LocalizedText{Ja: settings.HeroTitle.Ja, En: settings.HeroTitle.En},
		HeroDescription:    model.LocalizedText{Ja: settings.HeroDescription.Ja, En: settings.HeroDescription.En},
		Topics:             copyTopics,
		ConsentText:        model.LocalizedText{Ja: settings.ConsentText.Ja, En: settings.ConsentText.En},
		MinimumLeadHours:   settings.MinimumLeadHours,
		RecaptchaSiteKey:   settings.RecaptchaSiteKey,
		SupportEmail:       settings.SupportEmail,
		CalendarTimezone:   settings.CalendarTimezone,
		GoogleCalendarID:   settings.GoogleCalendarID,
		BookingWindowDays:  settings.BookingWindowDays,
		MeetingURLTemplate: settings.MeetingURLTemplate,
		CreatedAt:          settings.CreatedAt,
		UpdatedAt:          settings.UpdatedAt,
	}
}
