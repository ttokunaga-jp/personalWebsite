package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchRepository struct {
	mu       sync.RWMutex
	research []model.AdminResearch
	nextID   int64
}

func NewResearchRepository() repository.ResearchRepository {
	return newResearchRepository()
}

func newResearchRepository() *researchRepository {
	items := make([]model.AdminResearch, len(defaultAdminResearch))
	copy(items, defaultAdminResearch)
	var maxID int64
	for _, r := range items {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	return &researchRepository{
		research: items,
		nextID:   maxID + 1,
	}
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Research
	for _, item := range r.research {
		if !item.Published {
			continue
		}
		result = append(result, model.Research{
			ID:        item.ID,
			Title:     item.Title,
			Summary:   item.Summary,
			ContentMD: item.ContentMD,
			Year:      item.Year,
		})
	}
	return result, nil
}

func (r *researchRepository) ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.AdminResearch, len(r.research))
	for i, item := range r.research {
		result[i] = copyAdminResearch(item)
	}
	return result, nil
}

func (r *researchRepository) GetAdminResearch(ctx context.Context, id int64) (*model.AdminResearch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.research {
		if item.ID == id {
			copied := copyAdminResearch(item)
			return &copied, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *researchRepository) CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	item.ID = r.nextID
	r.nextID++
	item.CreatedAt = now
	item.UpdatedAt = now

	r.research = append(r.research, copyAdminResearch(*item))
	created := copyAdminResearch(*item)
	return &created, nil
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.research {
		if existing.ID == item.ID {
			item.CreatedAt = existing.CreatedAt
			item.UpdatedAt = time.Now().UTC()
			r.research[idx] = copyAdminResearch(*item)
			updated := copyAdminResearch(*item)
			return &updated, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, item := range r.research {
		if item.ID == id {
			r.research = append(r.research[:idx], r.research[idx+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func copyAdminResearch(src model.AdminResearch) model.AdminResearch {
	dst := src
	return dst
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
