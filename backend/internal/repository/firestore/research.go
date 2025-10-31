package firestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type researchRepository struct {
	base baseRepository
}

const researchCollection = "research"

type researchDocument struct {
	ID        int64        `firestore:"id"`
	Title     localizedDoc `firestore:"title"`
	Summary   localizedDoc `firestore:"summary"`
	ContentMD localizedDoc `firestore:"contentMd"`
	Year      int          `firestore:"year"`
	Published bool         `firestore:"published"`
	CreatedAt time.Time    `firestore:"createdAt"`
	UpdatedAt time.Time    `firestore:"updatedAt"`
}

func NewResearchRepository(client *firestore.Client, prefix string) repository.ResearchRepository {
	return &researchRepository{base: newBaseRepository(client, prefix)}
}

func (r *researchRepository) ListResearch(ctx context.Context) ([]model.Research, error) {
	docs, err := r.base.collection(researchCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore research: list: %w", err)
	}

	items, err := r.decodeResearch(docs)
	if err != nil {
		return nil, err
	}

	var result []model.Research
	for _, item := range items {
		if !item.Published {
			continue
		}
		result = append(result, model.Research{
			ID:        item.ID,
			Year:      item.Year,
			Title:     fromLocalizedDoc(item.Title),
			Summary:   fromLocalizedDoc(item.Summary),
			ContentMD: fromLocalizedDoc(item.ContentMD),
		})
	}
	return result, nil
}

func (r *researchRepository) ListAdminResearch(ctx context.Context) ([]model.AdminResearch, error) {
	docs, err := r.base.collection(researchCollection).Documents(ctx).GetAll()
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

func (r *researchRepository) GetAdminResearch(ctx context.Context, id int64) (*model.AdminResearch, error) {
	doc, err := r.base.doc(researchCollection, strconv.FormatInt(id, 10)).Get(ctx)
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
	result := mapResearchDocument(payload)
	return &result, nil
}

func (r *researchRepository) CreateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	id, err := nextID(ctx, r.base.client, r.base.prefix, researchCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore research: next id: %w", err)
	}

	now := time.Now().UTC()
	entry := researchDocument{
		ID:        id,
		Title:     toLocalizedDoc(item.Title),
		Summary:   toLocalizedDoc(item.Summary),
		ContentMD: toLocalizedDoc(item.ContentMD),
		Year:      item.Year,
		Published: item.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	docRef := r.base.doc(researchCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("firestore research: create %d: %w", id, err)
	}

	return r.GetAdminResearch(ctx, id)
}

func (r *researchRepository) UpdateAdminResearch(ctx context.Context, item *model.AdminResearch) (*model.AdminResearch, error) {
	if item == nil {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(researchCollection, strconv.FormatInt(item.ID, 10))
	now := time.Now().UTC()

	updates := []firestore.Update{
		{Path: "title", Value: toLocalizedDoc(item.Title)},
		{Path: "summary", Value: toLocalizedDoc(item.Summary)},
		{Path: "contentMd", Value: toLocalizedDoc(item.ContentMD)},
		{Path: "year", Value: item.Year},
		{Path: "published", Value: item.Published},
		{Path: "updatedAt", Value: now},
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore research: update %d: %w", item.ID, err)
	}

	return r.GetAdminResearch(ctx, item.ID)
}

func (r *researchRepository) DeleteAdminResearch(ctx context.Context, id int64) error {
	docRef := r.base.doc(researchCollection, strconv.FormatInt(id, 10))
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
			if id, err := strconv.ParseInt(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		result = append(result, entry)
	}

	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		leftKey := researchSortKey(left)
		rightKey := researchSortKey(right)
		if leftKey != rightKey {
			return leftKey < rightKey
		}
		if left.Year != right.Year {
			return left.Year > right.Year
		}
		return left.ID < right.ID
	})

	return result, nil
}

func researchSortKey(item researchDocument) int {
	return item.Year * 1000
}

func mapResearchDocument(doc researchDocument) model.AdminResearch {
	return model.AdminResearch{
		ID:        doc.ID,
		Title:     fromLocalizedDoc(doc.Title),
		Summary:   fromLocalizedDoc(doc.Summary),
		ContentMD: fromLocalizedDoc(doc.ContentMD),
		Year:      doc.Year,
		Published: doc.Published,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

var _ repository.ResearchRepository = (*researchRepository)(nil)
var _ repository.AdminResearchRepository = (*researchRepository)(nil)
