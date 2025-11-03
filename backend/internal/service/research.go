package service

import (
	"context"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ResearchService exposes research/blog aggregate queries.
type ResearchService interface {
	ListResearchDocuments(ctx context.Context, includeDrafts bool) ([]model.ResearchDocument, error)
}

type researchService struct {
	repo repository.ResearchDocumentRepository
}

func NewResearchService(repo repository.ResearchDocumentRepository) ResearchService {
	return &researchService{repo: repo}
}

func (s *researchService) ListResearchDocuments(ctx context.Context, includeDrafts bool) ([]model.ResearchDocument, error) {
	research, err := s.repo.ListResearchDocuments(ctx, includeDrafts)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load research documents", err)
	}
	return research, nil
}
