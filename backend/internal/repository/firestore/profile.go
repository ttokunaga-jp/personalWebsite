package firestore

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
)

// profileRepository temporarily delegates to the in-memory implementation until Firestore migration is completed.
type profileRepository struct {
	public repository.ProfileRepository
	admin  repository.AdminProfileRepository
}

func NewProfileRepository(client *firestore.Client, prefix string) repository.ProfileRepository {
	_ = client
	_ = prefix

	delegate := inmemory.NewProfileRepository()
	adminDelegate, _ := delegate.(repository.AdminProfileRepository)

	return &profileRepository{
		public: delegate,
		admin:  adminDelegate,
	}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	return r.public.GetProfile(ctx)
}

func (r *profileRepository) GetAdminProfile(ctx context.Context) (*model.AdminProfile, error) {
	if r.admin == nil {
		return nil, repository.ErrNotImplemented
	}
	return r.admin.GetAdminProfile(ctx)
}

func (r *profileRepository) UpdateAdminProfile(ctx context.Context, profile *model.AdminProfile) (*model.AdminProfile, error) {
	if r.admin == nil {
		return nil, repository.ErrNotImplemented
	}
	return r.admin.UpdateAdminProfile(ctx, profile)
}

var _ repository.ProfileRepository = (*profileRepository)(nil)
var _ repository.AdminProfileRepository = (*profileRepository)(nil)
