package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contactFormSettingsRepository struct{}

// NewContactFormSettingsRepository returns an in-memory contact form settings repository.
func NewContactFormSettingsRepository() repository.ContactFormSettingsRepository {
	return &contactFormSettingsRepository{}
}

func (r *contactFormSettingsRepository) GetContactFormSettings(ctx context.Context) (*model.ContactFormSettingsV2, error) {
	_ = ctx

	now := time.Now().UTC()
	settings := &model.ContactFormSettingsV2{
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
		MinimumLeadHours:  24,
		RecaptchaSiteKey:  "recaptcha-public-key",
		SupportEmail:      "support@example.dev",
		CalendarTimezone:  "Asia/Tokyo",
		GoogleCalendarID:  "primary",
		BookingWindowDays: 30,
		CreatedAt:         now.AddDate(-1, 0, 0),
		UpdatedAt:         now.Add(-6 * time.Hour),
	}

	return settings, nil
}
