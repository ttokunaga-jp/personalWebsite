package inmemory

import (
	"context"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type techCatalogRepository struct{}

// NewTechCatalogRepository returns a seeded in-memory tech catalog repository.
func NewTechCatalogRepository() repository.TechCatalogRepository {
	return &techCatalogRepository{}
}

func (r *techCatalogRepository) ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error) {
	_ = ctx

	data := []model.TechCatalogEntry{
		{
			ID:          1,
			Slug:        "go",
			DisplayName: "Go",
			Category:    "backend",
			Level:       model.TechLevelAdvanced,
			Icon:        "🐹",
			SortOrder:   1,
			Active:      true,
			CreatedAt:   time.Now().Add(-365 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			ID:          2,
			Slug:        "react",
			DisplayName: "React",
			Category:    "frontend",
			Level:       model.TechLevelIntermediate,
			Icon:        "⚛️",
			SortOrder:   2,
			Active:      true,
			CreatedAt:   time.Now().Add(-400 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          3,
			Slug:        "python",
			DisplayName: "Python",
			Category:    "ml",
			Level:       model.TechLevelAdvanced,
			Icon:        "🐍",
			SortOrder:   3,
			Active:      false,
			CreatedAt:   time.Now().Add(-600 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-72 * time.Hour),
		},
	}

	if includeInactive {
		return append([]model.TechCatalogEntry(nil), data...), nil
	}

	active := make([]model.TechCatalogEntry, 0, len(data))
	for _, entry := range data {
		if entry.Active {
			active = append(active, entry)
		}
	}
	return active, nil
}
