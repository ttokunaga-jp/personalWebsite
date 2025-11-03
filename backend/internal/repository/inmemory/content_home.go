package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type homePageConfigRepository struct{}

// NewHomePageConfigRepository returns an in-memory home page configuration repository.
func NewHomePageConfigRepository() repository.HomePageConfigRepository {
	return &homePageConfigRepository{}
}

func (r *homePageConfigRepository) GetHomePageConfig(ctx context.Context) (*model.HomePageConfigDocument, error) {
	_ = ctx

	now := time.Now().UTC()
	config := &model.HomePageConfigDocument{
		ID:        1,
		ProfileID: 1,
		HeroSubtitle: model.NewLocalizedText(
			"AI × Web エンジニアリング",
			"AI × Web Engineering",
		),
		QuickLinks: []model.HomeQuickLink{
			{
				ID:          1,
				ConfigID:    1,
				Section:     "profile",
				Label:       model.NewLocalizedText("プロフィール", "Profile"),
				Description: model.NewLocalizedText("研究と開発のバックグラウンド", "Research and product background"),
				CTA:         model.NewLocalizedText("詳しく見る", "View details"),
				TargetURL:   "/profile",
				SortOrder:   1,
			},
			{
				ID:          2,
				ConfigID:    1,
				Section:     "projects",
				Label:       model.NewLocalizedText("プロジェクト", "Projects"),
				Description: model.NewLocalizedText("最近の取り組みを紹介", "Latest initiatives"),
				CTA:         model.NewLocalizedText("見る", "Explore"),
				TargetURL:   "/projects",
				SortOrder:   2,
			},
		},
		ChipSources: []model.HomeChipSource{
			{
				ID:        1,
				ConfigID:  1,
				Source:    "affiliation",
				Label:     model.NewLocalizedText("所属", "Affiliations"),
				Limit:     3,
				SortOrder: 1,
			},
			{
				ID:        2,
				ConfigID:  1,
				Source:    "tech",
				Label:     model.NewLocalizedText("注目技術", "Featured Tech"),
				Limit:     6,
				SortOrder: 2,
			},
		},
		UpdatedAt: now.Add(-12 * time.Hour),
	}

	return config, nil
}
