package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type blogRepository struct {
	mu     sync.RWMutex
	posts  []model.BlogPost
	nextID int64
}

func NewBlogRepository() repository.BlogRepository {
	posts := make([]model.BlogPost, len(defaultBlogPosts))
	copy(posts, defaultBlogPosts)
	var maxID int64
	for _, post := range posts {
		if post.ID > maxID {
			maxID = post.ID
		}
	}
	return &blogRepository{
		posts:  posts,
		nextID: maxID + 1,
	}
}

func (r *blogRepository) ListBlogPosts(ctx context.Context) ([]model.BlogPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.BlogPost, len(r.posts))
	for i, post := range r.posts {
		result[i] = copyBlogPost(post)
	}
	return result, nil
}

func (r *blogRepository) GetBlogPost(ctx context.Context, id int64) (*model.BlogPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, post := range r.posts {
		if post.ID == id {
			copied := copyBlogPost(post)
			return &copied, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *blogRepository) CreateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	post.ID = r.nextID
	r.nextID++
	post.CreatedAt = now
	post.UpdatedAt = now
	if post.Published && post.PublishedAt == nil {
		t := now
		post.PublishedAt = &t
	}

	r.posts = append(r.posts, copyBlogPost(*post))
	created := copyBlogPost(*post)
	return &created, nil
}

func (r *blogRepository) UpdateBlogPost(ctx context.Context, post *model.BlogPost) (*model.BlogPost, error) {
	if post == nil {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, existing := range r.posts {
		if existing.ID == post.ID {
			post.CreatedAt = existing.CreatedAt
			if post.Published && post.PublishedAt == nil {
				if existing.PublishedAt != nil {
					t := *existing.PublishedAt
					post.PublishedAt = &t
				} else {
					t := time.Now().UTC()
					post.PublishedAt = &t
				}
			}
			post.UpdatedAt = time.Now().UTC()
			r.posts[idx] = copyBlogPost(*post)
			updated := copyBlogPost(*post)
			return &updated, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *blogRepository) DeleteBlogPost(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, post := range r.posts {
		if post.ID == id {
			r.posts = append(r.posts[:idx], r.posts[idx+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func copyBlogPost(src model.BlogPost) model.BlogPost {
	dst := src
	dst.Tags = append([]string(nil), src.Tags...)
	if src.PublishedAt != nil {
		t := *src.PublishedAt
		dst.PublishedAt = &t
	}
	return dst
}

var _ repository.BlogRepository = (*blogRepository)(nil)
