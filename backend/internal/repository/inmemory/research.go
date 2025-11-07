package inmemory

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

const researchEntityType = "research_blog"

type researchRepository struct {
	mu       sync.RWMutex
	research []model.AdminResearch
	nextID   uint64
}

func NewResearchRepository() repository.ResearchRepository {
	return newResearchRepository()
}

func newResearchRepository() *researchRepository {
	items := make([]model.AdminResearch, len(defaultAdminResearch))
	for i, entry := range defaultAdminResearch {
		items[i] = copyAdminResearch(entry)
	}

	var maxID uint64
	for _, item := range items {
		if item.ID > maxID {
			maxID = item.ID
		}
		setTechEntity(&item)
	}

	return &researchRepository{
		research: items,
		nextID:   maxID + 1,
	}
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Research, 0, len(r.research))
	for _, item := range r.research {
		if item.IsDraft {
			continue
		}
		id, err := safeUintToInt64(item.ID)
		if err != nil {
			return nil, err
		}
		year := item.PublishedAt.Year()
		if year <= 0 {
			year = item.UpdatedAt.Year()
		}
		result = append(result, model.Research{
			ID:        id,
			Year:      year,
			Title:     item.Title,
			Summary:   item.Overview,
			ContentMD: item.Outcome,
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

func (r *researchRepository) GetAdminResearch(ctx context.Context, id uint64) (*model.AdminResearch, error) {
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
	if item.PublishedAt.IsZero() {
		item.PublishedAt = now
	}
	item.ID = r.nextID
	r.nextID++
	item.CreatedAt = now
	item.UpdatedAt = now
	setTechEntity(item)

	copied := copyAdminResearch(*item)
	r.research = append(r.research, copied)
	created := copyAdminResearch(copied)
	return &created, nil
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.research {
		if existing.ID != item.ID {
			continue
		}

		item.CreatedAt = existing.CreatedAt
		if item.PublishedAt.IsZero() {
			item.PublishedAt = existing.PublishedAt
		}
		item.UpdatedAt = time.Now().UTC()
		setTechEntity(item)

		r.research[idx] = copyAdminResearch(*item)
		updated := copyAdminResearch(r.research[idx])
		return &updated, nil
	}

	return nil, repository.ErrNotFound
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id uint64) error {
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
	dst.Title = copyLocalized(src.Title)
	dst.Overview = copyLocalized(src.Overview)
	dst.Outcome = copyLocalized(src.Outcome)
	dst.Outlook = copyLocalized(src.Outlook)
	dst.ImageAlt = copyLocalized(src.ImageAlt)
	dst.Tags = cloneResearchTags(src.Tags)
	dst.Links = cloneResearchLinks(src.Links)
	dst.Assets = cloneResearchAssets(src.Assets)
	dst.Tech = cloneTechMemberships(src.Tech)
	return dst
}

func copyLocalized(src model.LocalizedText) model.LocalizedText {
	return model.LocalizedText{
		Ja: src.Ja,
		En: src.En,
	}
}

func cloneResearchTags(src []model.ResearchTag) []model.ResearchTag {
	if len(src) == 0 {
		return nil
	}
	result := make([]model.ResearchTag, len(src))
	copy(result, src)
	return result
}

func cloneResearchLinks(src []model.ResearchLink) []model.ResearchLink {
	if len(src) == 0 {
		return nil
	}
	result := make([]model.ResearchLink, len(src))
	for i, link := range src {
		result[i] = model.ResearchLink{
			ID:        link.ID,
			EntryID:   link.EntryID,
			Type:      link.Type,
			Label:     copyLocalized(link.Label),
			URL:       link.URL,
			SortOrder: link.SortOrder,
		}
	}
	return result
}

func cloneResearchAssets(src []model.ResearchAsset) []model.ResearchAsset {
	if len(src) == 0 {
		return nil
	}
	result := make([]model.ResearchAsset, len(src))
	for i, asset := range src {
		result[i] = model.ResearchAsset{
			ID:        asset.ID,
			EntryID:   asset.EntryID,
			URL:       asset.URL,
			Caption:   copyLocalized(asset.Caption),
			SortOrder: asset.SortOrder,
		}
	}
	return result
}

func cloneTechMemberships(src []model.TechMembership) []model.TechMembership {
	if len(src) == 0 {
		return nil
	}
	result := make([]model.TechMembership, len(src))
	for i, membership := range src {
		result[i] = model.TechMembership{
			MembershipID: membership.MembershipID,
			EntityType:   membership.EntityType,
			EntityID:     membership.EntityID,
			Tech: model.TechCatalogEntry{
				ID:          membership.Tech.ID,
				Slug:        membership.Tech.Slug,
				DisplayName: membership.Tech.DisplayName,
				Category:    membership.Tech.Category,
				Level:       membership.Tech.Level,
				Icon:        membership.Tech.Icon,
				SortOrder:   membership.Tech.SortOrder,
				Active:      membership.Tech.Active,
				CreatedAt:   membership.Tech.CreatedAt,
				UpdatedAt:   membership.Tech.UpdatedAt,
			},
			Context:   membership.Context,
			Note:      membership.Note,
			SortOrder: membership.SortOrder,
		}
	}
	return result
}

func setTechEntity(item *model.AdminResearch) {
	if len(item.Tech) == 0 {
		return
	}
	for idx := range item.Tech {
		item.Tech[idx].EntityType = researchEntityType
		item.Tech[idx].EntityID = item.ID
	}
}

func safeUintToInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, repository.ErrInvalidInput
	}
	return int64(value), nil
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
