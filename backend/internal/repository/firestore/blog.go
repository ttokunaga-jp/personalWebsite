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

type blogRepository struct {
	base baseRepository
}

const blogCollection = "blog_posts"

type blogPostDocument struct {
	ID          int64        `firestore:"id"`
	Title       localizedDoc `firestore:"title"`
	Summary     localizedDoc `firestore:"summary"`
	ContentMD   localizedDoc `firestore:"contentMd"`
	Tags        []string     `firestore:"tags"`
	Published   bool         `firestore:"published"`
	PublishedAt *time.Time   `firestore:"publishedAt,omitempty"`
	CreatedAt   time.Time    `firestore:"createdAt"`
	UpdatedAt   time.Time    `firestore:"updatedAt"`
}

func NewBlogRepository(client *firestore.Client, prefix string) repository.BlogRepository {
	return &blogRepository{base: newBaseRepository(client, prefix)}
}

func (r *blogRepository) ListBlogPosts(ctx context.Context) ([]model.BlogPost, error) {
	docs, err := r.base.collection(blogCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("firestore blog: list: %w", err)
	}

	posts, err := r.decodeBlogPosts(docs)
	if err != nil {
		return nil, err
	}

	result := make([]model.BlogPost, 0, len(posts))
	for _, post := range posts {
		result = append(result, mapBlogPostDocument(post))
	}
	return result, nil
}

func (r *blogRepository) GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error) {
	doc, err := r.base.doc(blogCollection, strconv.FormatInt(id, 10)).Get(ctx)
	if err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore blog: get %d: %w", id, err)
	}

	var payload blogPostDocument
	if err := doc.DataTo(&payload); err != nil {
		return nil, fmt.Errorf("firestore blog: decode %d: %w", id, err)
	}
	result := mapBlogPostDocument(payload)
	return &result, nil
}

func (r *blogRepository) CreateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	id, err := nextID(ctx, r.base.client, r.base.prefix, blogCollection)
	if err != nil {
		return nil, fmt.Errorf("firestore blog: next id: %w", err)
	}

	now := time.Now().UTC()
	var publishedAt *time.Time
	if post.PublishedAt != nil && !post.PublishedAt.IsZero() {
		ts := post.PublishedAt.UTC()
		publishedAt = &ts
	}
	entry := blogPostDocument{
		ID:          id,
		Title:       toLocalizedDoc(post.Title),
		Summary:     toLocalizedDoc(post.Summary),
		ContentMD:   toLocalizedDoc(post.ContentMD),
		Tags:        copyStringSlice(post.Tags),
		Published:   post.Published,
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	docRef := r.base.doc(blogCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("firestore blog: create %d: %w", id, err)
	}

	return r.GetBlogPost(ctx, id)
}

func (r *blogRepository) UpdateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	docRef := r.base.doc(blogCollection, strconv.FormatInt(post.ID, 10))
	now := time.Now().UTC()
	updates := []firestore.Update{
		{Path: "title", Value: toLocalizedDoc(post.Title)},
		{Path: "summary", Value: toLocalizedDoc(post.Summary)},
		{Path: "contentMd", Value: toLocalizedDoc(post.ContentMD)},
		{Path: "tags", Value: copyStringSlice(post.Tags)},
		{Path: "published", Value: post.Published},
		{Path: "publishedAt", Value: publishedAtPtr(post.PublishedAt)},
		{Path: "updatedAt", Value: now},
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if notFound(err) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("firestore blog: update %d: %w", post.ID, err)
	}

	return r.GetBlogPost(ctx, post.ID)
}

func (r *blogRepository) DeleteBlogPost(ctx context.Context, id int64) error {
	docRef := r.base.doc(blogCollection, strconv.FormatInt(id, 10))
	if _, err := docRef.Delete(ctx); err != nil {
		if notFound(err) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("firestore blog: delete %d: %w", id, err)
	}
	return nil
}

func (r *blogRepository) decodeBlogPosts(docs []*firestore.DocumentSnapshot) ([]blogPostDocument, error) {
	items := make([]blogPostDocument, 0, len(docs))
	for _, doc := range docs {
		var entry blogPostDocument
		if err := doc.DataTo(&entry); err != nil {
			return nil, fmt.Errorf("firestore blog: decode %s: %w", doc.Ref.ID, err)
		}
		if entry.ID == 0 {
			if id, err := strconv.ParseInt(doc.Ref.ID, 10, 64); err == nil {
				entry.ID = id
			}
		}
		items = append(items, entry)
	}

	sort.Slice(items, func(i, j int) bool {
		left, right := items[i], items[j]
		leftRef := coalesceTime(left.PublishedAt, left.UpdatedAt)
		rightRef := coalesceTime(right.PublishedAt, right.UpdatedAt)
		if !leftRef.Equal(rightRef) {
			return leftRef.After(rightRef)
		}
		return left.ID > right.ID
	})

	return items, nil
}

func coalesceTime(primary *time.Time, fallback time.Time) time.Time {
	if primary != nil {
		return primary.UTC()
	}
	return fallback.UTC()
}

func mapBlogPostDocument(doc blogPostDocument) model.BlogPost {
	var publishedAt *time.Time
	if doc.PublishedAt != nil && !doc.PublishedAt.IsZero() {
		ts := doc.PublishedAt.UTC()
		publishedAt = &ts
	}

	return model.BlogPost{
		ID:          doc.ID,
		Title:       fromLocalizedDoc(doc.Title),
		Summary:     fromLocalizedDoc(doc.Summary),
		ContentMD:   fromLocalizedDoc(doc.ContentMD),
		Tags:        copyStringSlice(doc.Tags),
		Published:   doc.Published,
		PublishedAt: publishedAt,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}

func publishedAtPtr(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	ts := value.UTC()
	return &ts
}

var _ repository.BlogRepository = (*blogRepository)(nil)
