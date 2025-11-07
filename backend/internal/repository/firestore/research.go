package firestore

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchRepository struct {
	base baseRepository
}

const (
	researchBlogCollection = "research_blog_entries"
	researchEntityType     = "research_blog"
)

type researchDocument struct {
	ID                uint64             `firestore:"id"`
	Slug              string             `firestore:"slug"`
	Kind              string             `firestore:"kind"`
	Title             localizedDoc       `firestore:"title"`
	Overview          localizedDoc       `firestore:"overview"`
	Outcome           localizedDoc       `firestore:"outcome"`
	Outlook           localizedDoc       `firestore:"outlook"`
	ExternalURL       string             `firestore:"externalUrl"`
	HighlightImageURL string             `firestore:"highlightImageUrl"`
	ImageAlt          localizedDoc       `firestore:"imageAlt"`
	PublishedAt       time.Time          `firestore:"publishedAt"`
	IsDraft           bool               `firestore:"isDraft"`
	CreatedAt         time.Time          `firestore:"createdAt"`
	UpdatedAt         time.Time          `firestore:"updatedAt"`
	Tags              []researchTagDoc   `firestore:"tags"`
	Links             []researchLinkDoc  `firestore:"links"`
	Assets            []researchAssetDoc `firestore:"assets"`
	Tech              []researchTechDoc  `firestore:"tech"`
}

type researchTagDoc struct {
	ID        uint64 `firestore:"id"`
	Value     string `firestore:"value"`
	SortOrder int    `firestore:"sortOrder"`
}

type researchLinkDoc struct {
	ID        uint64       `firestore:"id"`
	Type      string       `firestore:"type"`
	Label     localizedDoc `firestore:"label"`
	URL       string       `firestore:"url"`
	SortOrder int          `firestore:"sortOrder"`
}

type researchAssetDoc struct {
	ID        uint64       `firestore:"id"`
	URL       string       `firestore:"url"`
	Caption   localizedDoc `firestore:"caption"`
	SortOrder int          `firestore:"sortOrder"`
}

type researchTechDoc struct {
	ID        uint64 `firestore:"id"`
	TechID    uint64 `firestore:"techId"`
	Context   string `firestore:"context"`
	Note      string `firestore:"note"`
	SortOrder int    `firestore:"sortOrder"`
}

func NewResearchRepository(client *firestore.Client, prefix string) repository.ResearchRepository {
	return &researchRepository{base: newBaseRepository(client, prefix)}
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	docs, err := r.base.collection(researchBlogCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore research: list: %w", err)
	}

	items, err := r.decodeResearch(docs)
	if err != nil {
		return nil, err
	}

	result := make([]model.Research, 0, len(items))
	for _, item := range items {
		if item.IsDraft {
			continue
		}
		id, err := safeUintToInt64(item.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, model.Research{
			ID:        id,
			Year:      item.PublishedAt.Year(),
			Title:     fromLocalizedDoc(item.Title),
			Summary:   fromLocalizedDoc(item.Overview),
			ContentMD: fromLocalizedDoc(item.Outcome),
		})
	}
	return result, nil
}

func (r *researchRepository) ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error) {
	docs, err := r.base.collection(researchBlogCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore research: list admin: %w", err)
	}

	items, err := r.decodeResearch(docs)
	if err != nil {
		return nil, err
	}

	admin := make([]model.AdminResearch, 0, len(items))
	for _, item := range items {
		admin = append(admin, mapResearchDocument(item))
	}
	return admin, nil
}

func (r *researchRepository) GetAdminResearch(ctx context.Context, id uint64) (*model.AdminResearch, error) {
	doc, err := r.base.doc(researchBlogCollection, strconv.FormatUint(id, 10)).Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore research: get %d: %w", id, err)
	}

	var payload researchDocument
	if err := doc.DataTo(&payload); err != nil {
		return nil, fmt.Errorf("firestore research: decode %d: %w", id, err)
	}
	if payload.ID == 0 {
		payload.ID = id
	}
	result := mapResearchDocument(payload)
	return &result, nil
}

func (r *researchRepository) CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	next, err := nextID(ctx, r.base.client, r.base.prefix, researchBlogCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore research: next id: %w", err)
	}
	id := uint64(next)
	now := time.Now().UTC()

	doc := toResearchDocument(id, item, now, now)
	docRef := r.base.doc(researchBlogCollection, strconv.FormatUint(id, 10))
	if _, err := docRef.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("firestore research: create %d: %w", id, err)
	}

	return r.GetAdminResearch(ctx, id)
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}
	if item.ID == 0 {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(researchBlogCollection, strconv.FormatUint(item.ID, 10))
	current, err := docRef.Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore research: update load %d: %w", item.ID, err)
	}

	var existing researchDocument
	if err := current.DataTo(&existing); err != nil {
		return nil, fmt.Errorf("firestore research: update decode %d: %w", item.ID, err)
	}
	if existing.CreatedAt.IsZero() {
		existing.CreatedAt = time.Now().UTC()
	}

	now := time.Now().UTC()
	doc := toResearchDocument(item.ID, item, existing.CreatedAt.UTC(), now)
	if _, err := docRef.Set(ctx, doc); err != nil {
		return nil, fmt.Errorf("firestore research: update %d: %w", item.ID, err)
	}

	return r.GetAdminResearch(ctx, item.ID)
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id uint64) error {
	docRef := r.base.doc(researchBlogCollection, strconv.FormatUint(id, 10))
	if _, err := docRef.Delete(ctx); err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore research: delete %d: %w", id, err)
	}
	return nil
}

func (r *researchRepository) decodeResearch(docs []*firestore.DocumentSnapshot) ([]researchDocument, error) {
	result := make([]researchDocument, 0, len(docs))
	for _, doc := range docs {
		var entry researchDocument
		if err := doc.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore research: decode %s: %w", doc.Ref.ID, err)
		}
		if entry.ID == 0 {
			if id, err := strconv.ParseUint(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		entry.CreatedAt = entry.CreatedAt.UTC()
		entry.UpdatedAt = entry.UpdatedAt.UTC()
		entry.PublishedAt = entry.PublishedAt.UTC()
		result = append(result, entry)
	}

	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		if !left.PublishedAt.Equal(right.PublishedAt) {
			return left.PublishedAt.After(right.PublishedAt)
		}
		return left.ID > right.ID
	})

	return result, nil
}

func toResearchDocument(id uint64, item *model.AdminResearch, createdAt, updatedAt time.Time) researchDocument {
	return researchDocument{
		ID:                id,
		Slug:              strings.TrimSpace(item.Slug),
		Kind:              string(item.Kind),
		Title:             toLocalizedDoc(item.Title),
		Overview:          toLocalizedDoc(item.Overview),
		Outcome:           toLocalizedDoc(item.Outcome),
		Outlook:           toLocalizedDoc(item.Outlook),
		ExternalURL:       strings.TrimSpace(item.ExternalURL),
		HighlightImageURL: strings.TrimSpace(item.HighlightImageURL),
		ImageAlt:          toLocalizedDoc(item.ImageAlt),
		PublishedAt:       item.PublishedAt.UTC(),
		IsDraft:           item.IsDraft,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Tags:              toResearchTagDocs(item.Tags),
		Links:             toResearchLinkDocs(item.Links),
		Assets:            toResearchAssetDocs(item.Assets),
		Tech:              toResearchTechDocs(item.Tech),
	}
}

func toResearchTagDocs(tags []model.ResearchTag) []researchTagDoc {
	if len(tags) == 0 {
		return nil
	}
	result := make([]researchTagDoc, 0, len(tags))
	nextID := uint64(1)
	for _, tag := range tags {
		value := strings.TrimSpace(tag.Value)
		if value == "" {
			continue
		}
		id := tag.ID
		if id == 0 {
			id = nextID
			nextID++
		}
		result = append(result, researchTagDoc{
			ID:        id,
			Value:     value,
			SortOrder: tag.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func toResearchLinkDocs(links []model.ResearchLink) []researchLinkDoc {
	if len(links) == 0 {
		return nil
	}
	result := make([]researchLinkDoc, 0, len(links))
	nextID := uint64(1)
	for _, link := range links {
		url := strings.TrimSpace(link.URL)
		if url == "" {
			continue
		}
		id := link.ID
		if id == 0 {
			id = nextID
			nextID++
		}
		result = append(result, researchLinkDoc{
			ID:        id,
			Type:      string(link.Type),
			Label:     toLocalizedDoc(link.Label),
			URL:       url,
			SortOrder: link.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func toResearchAssetDocs(assets []model.ResearchAsset) []researchAssetDoc {
	if len(assets) == 0 {
		return nil
	}
	result := make([]researchAssetDoc, 0, len(assets))
	nextID := uint64(1)
	for _, asset := range assets {
		url := strings.TrimSpace(asset.URL)
		if url == "" {
			continue
		}
		id := asset.ID
		if id == 0 {
			id = nextID
			nextID++
		}
		result = append(result, researchAssetDoc{
			ID:        id,
			URL:       url,
			Caption:   toLocalizedDoc(asset.Caption),
			SortOrder: asset.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func toResearchTechDocs(tech []model.TechMembership) []researchTechDoc {
	if len(tech) == 0 {
		return nil
	}
	result := make([]researchTechDoc, 0, len(tech))
	nextID := uint64(1)
	for _, membership := range tech {
		if membership.Tech.ID == 0 {
			continue
		}
		id := membership.MembershipID
		if id == 0 {
			id = nextID
			nextID++
		}
		result = append(result, researchTechDoc{
			ID:        id,
			TechID:    membership.Tech.ID,
			Context:   string(membership.Context),
			Note:      strings.TrimSpace(membership.Note),
			SortOrder: membership.SortOrder,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func mapResearchDocument(doc researchDocument) model.AdminResearch {
	admin := model.AdminResearch{
		ID:                doc.ID,
		Slug:              strings.TrimSpace(doc.Slug),
		Kind:              model.ResearchKind(strings.TrimSpace(doc.Kind)),
		Title:             fromLocalizedDoc(doc.Title),
		Overview:          fromLocalizedDoc(doc.Overview),
		Outcome:           fromLocalizedDoc(doc.Outcome),
		Outlook:           fromLocalizedDoc(doc.Outlook),
		ExternalURL:       strings.TrimSpace(doc.ExternalURL),
		HighlightImageURL: strings.TrimSpace(doc.HighlightImageURL),
		ImageAlt:          fromLocalizedDoc(doc.ImageAlt),
		PublishedAt:       doc.PublishedAt.UTC(),
		IsDraft:           doc.IsDraft,
		CreatedAt:         doc.CreatedAt.UTC(),
		UpdatedAt:         doc.UpdatedAt.UTC(),
	}

	if len(doc.Tags) > 0 {
		admin.Tags = make([]model.ResearchTag, 0, len(doc.Tags))
		for _, tag := range doc.Tags {
			admin.Tags = append(admin.Tags, model.ResearchTag{
				ID:        tag.ID,
				EntryID:   doc.ID,
				Value:     strings.TrimSpace(tag.Value),
				SortOrder: tag.SortOrder,
			})
		}
	}

	if len(doc.Links) > 0 {
		admin.Links = make([]model.ResearchLink, 0, len(doc.Links))
		for _, link := range doc.Links {
			admin.Links = append(admin.Links, model.ResearchLink{
				ID:        link.ID,
				EntryID:   doc.ID,
				Type:      model.ResearchLinkType(strings.TrimSpace(link.Type)),
				Label:     fromLocalizedDoc(link.Label),
				URL:       strings.TrimSpace(link.URL),
				SortOrder: link.SortOrder,
			})
		}
	}

	if len(doc.Assets) > 0 {
		admin.Assets = make([]model.ResearchAsset, 0, len(doc.Assets))
		for _, asset := range doc.Assets {
			admin.Assets = append(admin.Assets, model.ResearchAsset{
				ID:        asset.ID,
				EntryID:   doc.ID,
				URL:       strings.TrimSpace(asset.URL),
				Caption:   fromLocalizedDoc(asset.Caption),
				SortOrder: asset.SortOrder,
			})
		}
	}

	if len(doc.Tech) > 0 {
		admin.Tech = make([]model.TechMembership, 0, len(doc.Tech))
		for _, membership := range doc.Tech {
			admin.Tech = append(admin.Tech, model.TechMembership{
				MembershipID: membership.ID,
				EntityType:   researchEntityType,
				EntityID:     doc.ID,
				Tech: model.TechCatalogEntry{
					ID: membership.TechID,
				},
				Context:   model.TechContext(strings.TrimSpace(membership.Context)),
				Note:      strings.TrimSpace(membership.Note),
				SortOrder: membership.SortOrder,
			})
		}
	}

	return admin
}

func safeUintToInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("value %d exceeds int64 range", value)
	}
	return int64(value), nil
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
