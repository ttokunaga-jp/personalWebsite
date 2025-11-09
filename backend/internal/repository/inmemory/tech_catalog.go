package inmemory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type techCatalogRepository struct {
	mu      sync.RWMutex
	seq     uint64
	entries []model.TechCatalogEntry
}

// NewTechCatalogRepository returns a seeded in-memory tech catalog repository.
func NewTechCatalogRepository() repository.TechCatalogRepository {
	now := time.Now().UTC()
	entries := []model.TechCatalogEntry{
		{
			ID:          1,
			Slug:        "go",
			DisplayName: "Go",
			Category:    "backend",
			Level:       model.TechLevelAdvanced,
			Icon:        "🐹",
			SortOrder:   1,
			Active:      true,
			CreatedAt:   now.Add(-365 * 24 * time.Hour),
			UpdatedAt:   now.Add(-48 * time.Hour),
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
			CreatedAt:   now.Add(-400 * 24 * time.Hour),
			UpdatedAt:   now.Add(-24 * time.Hour),
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
			CreatedAt:   now.Add(-600 * 24 * time.Hour),
			UpdatedAt:   now.Add(-72 * time.Hour),
		},
	}

	var seq uint64
	for _, entry := range entries {
		if entry.ID > seq {
			seq = entry.ID
		}
	}

	return &techCatalogRepository{
		seq:     seq,
		entries: entries,
	}
}

func (r *techCatalogRepository) ListTechCatalog(ctx context.Context, includeInactive bool) ([]model.TechCatalogEntry, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]model.TechCatalogEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		if includeInactive || entry.Active {
			results = append(results, copyTechCatalogEntry(entry))
		}
	}
	return results, nil
}

func (r *techCatalogRepository) GetTechCatalogEntry(ctx context.Context, id uint64) (*model.TechCatalogEntry, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.entries {
		if entry.ID == id {
			result := copyTechCatalogEntry(entry)
			return &result, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *techCatalogRepository) CreateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error) {
	_ = ctx

	if entry == nil {
		return nil, repository.ErrInvalidInput
	}

	slug := strings.TrimSpace(entry.Slug)
	displayName := strings.TrimSpace(entry.DisplayName)
	if slug == "" || displayName == "" {
		return nil, repository.ErrInvalidInput
	}
	if !isValidTechLevel(entry.Level) {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.entries {
		if strings.EqualFold(existing.Slug, slug) {
			return nil, repository.ErrDuplicate
		}
	}

	r.seq++
	now := time.Now().UTC()
	created := model.TechCatalogEntry{
		ID:          r.seq,
		Slug:        slug,
		DisplayName: displayName,
		Category:    strings.TrimSpace(entry.Category),
		Level:       entry.Level,
		Icon:        strings.TrimSpace(entry.Icon),
		SortOrder:   entry.SortOrder,
		Active:      entry.Active,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.entries = append(r.entries, created)
	result := copyTechCatalogEntry(created)
	return &result, nil
}

func (r *techCatalogRepository) UpdateTechCatalogEntry(ctx context.Context, entry *model.TechCatalogEntry) (*model.TechCatalogEntry, error) {
	_ = ctx

	if entry == nil || entry.ID == 0 {
		return nil, repository.ErrInvalidInput
	}

	slug := strings.TrimSpace(entry.Slug)
	displayName := strings.TrimSpace(entry.DisplayName)
	if slug == "" || displayName == "" {
		return nil, repository.ErrInvalidInput
	}
	if !isValidTechLevel(entry.Level) {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var updated *model.TechCatalogEntry
	for index, existing := range r.entries {
		if existing.ID != entry.ID {
			if strings.EqualFold(existing.Slug, slug) {
				return nil, repository.ErrDuplicate
			}
			continue
		}

		now := time.Now().UTC()
		newEntry := model.TechCatalogEntry{
			ID:          existing.ID,
			Slug:        slug,
			DisplayName: displayName,
			Category:    strings.TrimSpace(entry.Category),
			Level:       entry.Level,
			Icon:        strings.TrimSpace(entry.Icon),
			SortOrder:   entry.SortOrder,
			Active:      entry.Active,
			CreatedAt:   existing.CreatedAt,
			UpdatedAt:   now,
		}
		r.entries[index] = newEntry
		copy := copyTechCatalogEntry(newEntry)
		updated = &copy
		break
	}

	if updated == nil {
		return nil, repository.ErrNotFound
	}
	return updated, nil
}

func copyTechCatalogEntry(entry model.TechCatalogEntry) model.TechCatalogEntry {
	result := entry
	return result
}

func isValidTechLevel(level model.TechLevel) bool {
	switch strings.TrimSpace(string(level)) {
	case string(model.TechLevelBeginner), string(model.TechLevelIntermediate), string(model.TechLevelAdvanced):
		return true
	default:
		return false
	}
}
