package service

import (
	"context"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ProfileService describes operations for fetching profile information.
type ProfileService interface {
	GetProfile(ctx context.Context) (*model.Profile, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(ctx context.Context) (*model.Profile, error) {
	profile, err := s.repo.GetProfile(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, 500, "failed to load profile", err)
	}
	return profile, nil
}
