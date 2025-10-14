package service

import (
	"context"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ResearchService exposes research related queries.
type ResearchService interface {
	ListResearch(ctx context.Context) ([]model.Research, error)
}

type researchService struct {
	repo repository.ResearchRepository
}

func NewResearchService(repo repository.ResearchRepository) ResearchService {
	return &researchService{repo: repo}
}

func (s *researchService) ListResearch(ctx context.Context) ([]model.Research, error) {
	research, err := s.repo.ListResearch(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load research", err)
	}
	return research, nil
}
