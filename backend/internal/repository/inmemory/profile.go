package inmemory

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type profileRepository struct {
	profile *model.Profile
}

func NewProfileRepository() repository.ProfileRepository {
	return &profileRepository{
		profile: defaultProfile,
	}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	return r.profile, nil
}
