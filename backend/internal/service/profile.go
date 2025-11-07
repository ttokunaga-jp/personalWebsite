package service

import (
	"context"
	"errors"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service/support"
)

// ProfileService exposes profile retrieval for the public API (v2 schema).
type ProfileService interface {
	GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error)
}

type profileService struct {
	repo         repository.ContentProfileRepository
	homeRepo     repository.HomePageConfigRepository
	fallbackRepo repository.ContentProfileRepository
	fallbackHome repository.HomePageConfigRepository
}

// NewProfileService returns a service backed by the v2 content repository.
func NewProfileService(repo repository.ContentProfileRepository, homeRepo repository.HomePageConfigRepository) ProfileService {
	return &profileService{
		repo:         repo,
		homeRepo:     homeRepo,
		fallbackRepo: inmemory.NewContentProfileRepository(),
		fallbackHome: inmemory.NewHomePageConfigRepository(),
	}
}

func (s *profileService) GetProfileDocument(ctx context.Context) (*model.ProfileDocument, error) {
	document, err := s.repo.GetProfileDocument(ctx)
	if err != nil {
		var ok bool
		document, ok = s.tryFallbackDocument(ctx, err)
		if !ok {
			return nil, support.MapRepositoryError(err, "profile document")
		}
	}
	if document == nil {
		return nil, errs.New(errs.CodeInternal, 500, "profile document is empty", nil)
	}

	if s.homeRepo != nil {
		homeConfig, err := s.homeRepo.GetHomePageConfig(ctx)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) || support.ShouldFallback(err) {
				if applied := s.applyFallbackHome(ctx, document); applied {
					return document, nil
				}
				return document, nil
			}
			return nil, support.MapRepositoryError(err, "home page config")
		} else {
			document.Home = homeConfig
		}
	}

	if document.Home == nil {
		s.applyFallbackHome(ctx, document)
	}

	return document, nil
}

func (s *profileService) tryFallbackDocument(ctx context.Context, cause error) (*model.ProfileDocument, bool) {
	if s.fallbackRepo == nil || !support.ShouldFallback(cause) {
		return nil, false
	}

	fallbackDoc, err := s.fallbackRepo.GetProfileDocument(ctx)
	if err != nil || fallbackDoc == nil {
		return nil, false
	}

	return fallbackDoc, true
}

func (s *profileService) applyFallbackHome(ctx context.Context, document *model.ProfileDocument) bool {
	if document == nil || s.fallbackHome == nil {
		return false
	}
	home, err := s.fallbackHome.GetHomePageConfig(ctx)
	if err != nil || home == nil {
		return false
	}
	document.Home = home
	return true
}
