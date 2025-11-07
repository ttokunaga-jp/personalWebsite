package service

import (
	"context"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service/support"
)

// ProjectService orchestrates retrieval of project aggregates for public and admin flows.
type ProjectService interface {
	ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error)
}

type projectService struct {
	repo     repository.ProjectDocumentRepository
	fallback repository.ProjectDocumentRepository
}

// NewProjectService constructs a project service backed by the v2 document repository.
func NewProjectService(repo repository.ProjectDocumentRepository, fallbacks ...repository.ProjectDocumentRepository) ProjectService {
	var fallback repository.ProjectDocumentRepository
	for _, candidate := range fallbacks {
		if candidate != nil {
			fallback = candidate
			break
		}
	}
	if fallback == nil {
		fallback = inmemory.NewProjectDocumentRepository()
	}

	return &projectService{
		repo:     repo,
		fallback: fallback,
	}
}

func (s *projectService) ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error) {
	projects, err := s.repo.ListProjectDocuments(ctx, includeDrafts)
	if err != nil {
		if s.fallback != nil && support.ShouldFallback(err) {
			if fallbackProjects, fallbackErr := s.fallback.ListProjectDocuments(ctx, includeDrafts); fallbackErr == nil {
				return fallbackProjects, nil
			}
		}
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load project documents", err)
	}
	return projects, nil
}
