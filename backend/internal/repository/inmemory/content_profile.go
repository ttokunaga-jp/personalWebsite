package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contentProfileRepository struct{}

// NewContentProfileRepository returns an in-memory ContentProfileRepository seeded from default fixtures.
func NewContentProfileRepository() repository.ContentProfileRepository {
	return &contentProfileRepository{}
}

func (r *contentProfileRepository) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	_ = ctx

	// Convert the default admin profile into the richer ProfileDocument structure.
	doc := &model.ProfileDocument{
		ID:          1,
		DisplayName: defaultAdminProfile.Name.En,
		Headline:    defaultAdminProfile.Title,
		Summary:     defaultAdminProfile.Summary,
		AvatarURL:   "https://example.dev/avatar.png",
		Location:    model.NewLocalizedText("東京", "Tokyo"),
		Theme: model.ProfileTheme{
			Mode:        model.ProfileThemeModeLight,
			AccentColor: "#3b82f6",
		},
		Lab: model.ProfileLab{
			Name:    defaultAdminProfile.Lab,
			Advisor: model.NewLocalizedText("指導教員", "Advisor"),
			Room:    model.NewLocalizedText("4F 研究室", "Lab 4F"),
			URL:     "https://example.dev/lab",
		},
		Affiliations: []model.ProfileAffiliation{
			{
				ID:          1,
				ProfileID:   1,
				Kind:        model.ProfileAffiliationKindAffiliation,
				Name:        "Example University",
				URL:         "https://example.dev",
				Description: model.NewLocalizedText("研究員", "Researcher"),
				StartedAt:   time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:   1,
			},
		},
		Communities: []model.ProfileAffiliation{
			{
				ID:          2,
				ProfileID:   1,
				Kind:        model.ProfileAffiliationKindCommunity,
				Name:        "Open Source Guild",
				URL:         "https://oss.example",
				Description: model.NewLocalizedText("OSS コミュニティ", "OSS community"),
				StartedAt:   time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				SortOrder:   1,
			},
		},
		WorkHistory: []model.ProfileWorkExperience{
			{
				ID:           1,
				ProfileID:    1,
				Organization: model.NewLocalizedText("Example Corp", "Example Corp"),
				Role:         model.NewLocalizedText("フルスタックエンジニア", "Full-stack Engineer"),
				Summary:      model.NewLocalizedText("AI と Web を横断", "Bridging AI and web"),
				StartedAt:    time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
				ExternalURL:  "https://example.dev/company",
				SortOrder:    1,
			},
		},
		TechSections: []model.ProfileTechSection{
			{
				ID:         1,
				ProfileID:  1,
				Title:      model.NewLocalizedText("スキルセット", "Skill Set"),
				Layout:     "chips",
				Breakpoint: "lg",
				SortOrder:  1,
				Members: []model.TechMembership{
					{
						MembershipID: 1,
						EntityType:   "profile_section",
						EntityID:     1,
						Tech: model.TechCatalogEntry{
							ID:          1,
							Slug:        "go",
							DisplayName: "Go",
							Category:    "backend",
							Level:       model.TechLevelAdvanced,
							SortOrder:   1,
							Active:      true,
							CreatedAt:   time.Now().Add(-365 * 24 * time.Hour),
							UpdatedAt:   time.Now(),
						},
						Context:   model.TechContextPrimary,
						SortOrder: 1,
					},
					{
						MembershipID: 2,
						EntityType:   "profile_section",
						EntityID:     1,
						Tech: model.TechCatalogEntry{
							ID:          2,
							Slug:        "react",
							DisplayName: "React",
							Category:    "frontend",
							Level:       model.TechLevelIntermediate,
							SortOrder:   2,
							Active:      true,
							CreatedAt:   time.Now().Add(-400 * 24 * time.Hour),
							UpdatedAt:   time.Now(),
						},
						Context:   model.TechContextSupporting,
						SortOrder: 2,
					},
				},
			},
		},
		SocialLinks: []model.ProfileSocialLink{
			{
				ID:        1,
				ProfileID: 1,
				Provider:  model.ProfileSocialProviderGitHub,
				Label:     model.NewLocalizedText("GitHub", "GitHub"),
				URL:       "https://github.com/example",
				IsFooter:  true,
				SortOrder: 1,
			},
			{
				ID:        2,
				ProfileID: 1,
				Provider:  model.ProfileSocialProviderLinkedIn,
				Label:     model.NewLocalizedText("LinkedIn", "LinkedIn"),
				URL:       "https://linkedin.com/in/example",
				IsFooter:  true,
				SortOrder: 2,
			},
		},
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	return doc, nil
}

var _ repository.ContentProfileRepository = (*contentProfileRepository)(nil)
