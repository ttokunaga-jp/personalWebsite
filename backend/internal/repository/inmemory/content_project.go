package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type projectDocumentRepository struct{}

// NewProjectDocumentRepository returns an in-memory repository for project aggregates.
func NewProjectDocumentRepository() repository.ProjectDocumentRepository {
	return &projectDocumentRepository{}
}

func (r *projectDocumentRepository) ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error) {
	_ = ctx

	now := time.Now().UTC()
	start := now.AddDate(-1, 0, 0)
	projects := []model.ProjectDocument{
		{
			ID:   1,
			Slug: "personal-website",
			Title: model.LocalizedText{
				Ja: "個人サイト刷新",
				En: "Personal Website Revamp",
			},
			Summary: model.LocalizedText{
				Ja: "Next.js と Go による再構築",
				En: "Rebuilt with Next.js and Go backend",
			},
			Description: model.LocalizedText{
				Ja: "設計刷新と可観測性の拡充を実施。",
				En: "Implemented a new architecture with improved observability.",
			},
			CoverImageURL: "https://example.dev/assets/projects/pw-cover.png",
			PrimaryLink:   "https://example.dev/projects/personal-website",
			Links: []model.ProjectLink{
				{
					ID:        1,
					ProjectID: 1,
					Type:      model.ProjectLinkTypeRepo,
					Label:     model.NewLocalizedText("GitHub", "GitHub"),
					URL:       "https://github.com/example/personal-website",
					SortOrder: 1,
				},
			},
			Period: model.ProjectPeriod{
				Start: &start,
			},
			Tech: []model.TechMembership{
				{
					MembershipID: 1,
					EntityType:   "project",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          1,
						Slug:        "go",
						DisplayName: "Go",
						Level:       model.TechLevelAdvanced,
						SortOrder:   1,
						Active:      true,
						CreatedAt:   now.AddDate(-3, 0, 0),
						UpdatedAt:   now.Add(-24 * time.Hour),
					},
					SortOrder: 1,
				},
				{
					MembershipID: 2,
					EntityType:   "project",
					EntityID:     1,
					Context:      model.TechContextSupporting,
					Tech: model.TechCatalogEntry{
						ID:          2,
						Slug:        "react",
						DisplayName: "React",
						Level:       model.TechLevelIntermediate,
						SortOrder:   2,
						Active:      true,
						CreatedAt:   now.AddDate(-4, 0, 0),
						UpdatedAt:   now.Add(-48 * time.Hour),
					},
					SortOrder: 2,
				},
			},
			Highlight: true,
			Published: true,
			SortOrder: 1,
			CreatedAt: now.AddDate(-1, 0, 0),
			UpdatedAt: now.Add(-6 * time.Hour),
		},
		{
			ID:   2,
			Slug: "ml-research",
			Title: model.LocalizedText{
				Ja: "ML 研究プロトタイプ",
				En: "ML Research Prototype",
			},
			Summary: model.NewLocalizedText("論文実装の検証", "Validating research ideas"),
			Description: model.NewLocalizedText(
				"学術論文のアイデアを PoC として実装し、推論最適化を評価。",
				"Implemented research paper ideas as a PoC and measured inference optimisations.",
			),
			PrimaryLink: "https://example.dev/projects/ml-research",
			Highlight:   false,
			Published:   false,
			SortOrder:   2,
			CreatedAt:   now.AddDate(0, -2, 0),
			UpdatedAt:   now.Add(-72 * time.Hour),
		},
	}

	if includeDrafts {
		return append([]model.ProjectDocument(nil), projects...), nil
	}

	results := make([]model.ProjectDocument, 0, len(projects))
	for _, p := range projects {
		if p.Published {
			results = append(results, p)
		}
	}
	return results, nil
}
