package service

import (
	"context"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/service/support"
)

// ProfileService exposes profile retrieval for the public API (v2 schema).
type ProfileService interface {
	GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error)
}

type profileService struct {
	repo repository.ContentProfileRepository
}

// NewProfileService returns a service backed by the v2 content repository.
func NewProfileService(repo repository.ContentProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	document, err := s.repo.GetProfileDocument(ctx)
	if err != nil {
		return nil, support.MapRepositoryError(err, "profile document")
	}
	if document == nil {
		return nil, errs.New(errs.CodeInternal, 500, "profile document is empty", nil)
	}
	return document, nil
}
