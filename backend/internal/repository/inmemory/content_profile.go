package inmemory

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type contentProfileRepository struct{}

// NewContentProfileRepository returns an in-memory ContentProfileRepository seeded from default fixtures.
func NewContentProfileRepository() repository.ContentProfileRepository {
	return &contentProfileRepository{}
}

func (r *contentProfileRepository) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	_ = ctx
	return cloneAdminProfile(defaultAdminProfile), nil
}

var _ repository.ContentProfileRepository = (*contentProfileRepository)(nil)
