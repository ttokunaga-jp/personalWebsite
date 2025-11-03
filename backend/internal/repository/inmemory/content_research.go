package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchDocumentRepository struct{}

// NewResearchDocumentRepository returns an in-memory research/blog repository.
func NewResearchDocumentRepository() repository.ResearchDocumentRepository {
	return &researchDocumentRepository{}
}

func (r *researchDocumentRepository) ListResearchDocuments(ctx context.Context, includeDrafts bool) ([]model.ResearchDocument, error) {
	_ = ctx

	now := time.Now().UTC()

	documents := []model.ResearchDocument{
		{
			ID:   1,
			Slug: "nlp-observability",
			Kind: model.ResearchKindResearch,
			Title: model.NewLocalizedText(
				"自然言語処理システムの可観測性",
				"Observability Strategies for NLP Systems",
			),
			Overview: model.NewLocalizedText(
				"推論パイプラインの計測基盤構築について紹介。",
				"Discusses a metrics pipeline for inference workloads.",
			),
			Outcome: model.NewLocalizedText(
				"誤差要因の可視化により MTBF が 30% 改善。",
				"Reduced MTBF by 30% through improved error attribution.",
			),
			Outlook: model.NewLocalizedText(
				"構造化ログとトレーシングの自動化を進める。",
				"Automating structured logging and tracing will be the next step.",
			),
			ExternalURL:       "https://example.dev/blog/nlp-observability",
			HighlightImageURL: "https://example.dev/assets/articles/nlp.png",
			ImageAlt:          model.NewLocalizedText("NLP 可観測性", "NLP Observability"),
			PublishedAt:       now.AddDate(0, -1, 0),
			UpdatedAt:         now.Add(-24 * time.Hour),
			Tags: []model.ResearchTag{
				{ID: 1, EntryID: 1, Value: "observability", SortOrder: 1},
				{ID: 2, EntryID: 1, Value: "nlp", SortOrder: 2},
			},
			Links: []model.ResearchLink{
				{
					ID:        1,
					EntryID:   1,
					Type:      model.ResearchLinkTypeSlides,
					Label:     model.NewLocalizedText("発表資料", "Slides"),
					URL:       "https://speakerdeck.com/example/nlp-observability",
					SortOrder: 1,
				},
			},
			Assets: []model.ResearchAsset{
				{
					ID:        1,
					EntryID:   1,
					URL:       "https://example.dev/assets/articles/nlp-chart.png",
					Caption:   model.NewLocalizedText("可視化ダッシュボード", "Observability dashboard"),
					SortOrder: 1,
				},
			},
			Tech: []model.TechMembership{
				{
					MembershipID: 1,
					EntityType:   "research_blog",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          3,
						Slug:        "python",
						DisplayName: "Python",
						Level:       model.TechLevelAdvanced,
						SortOrder:   3,
						Active:      true,
						CreatedAt:   now.AddDate(-5, 0, 0),
						UpdatedAt:   now.Add(-48 * time.Hour),
					},
					SortOrder: 1,
				},
			},
			IsDraft: false,
		},
		{
			ID:          2,
			Slug:        "ui-review-2024",
			Kind:        model.ResearchKindBlog,
			Title:       model.NewLocalizedText("UI レビューの振り返り", "2024 UI Review Retrospective"),
			Overview:    model.NewLocalizedText("UI 改善施策の裏側を記録。", "Notes on the UI refresh programme."),
			ExternalURL: "https://example.dev/blog/ui-review-2024",
			PublishedAt: now.AddDate(0, 0, -10),
			UpdatedAt:   now.Add(-3 * time.Hour),
			IsDraft:     true,
		},
	}

	if includeDrafts {
		return append([]model.ResearchDocument(nil), documents...), nil
	}

	results := make([]model.ResearchDocument, 0, len(documents))
	for _, doc := range documents {
		if !doc.IsDraft {
			results = append(results, doc)
		}
	}
	return results, nil
}
