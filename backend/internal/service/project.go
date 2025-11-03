package service

import (
	"context"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ProjectService orchestrates retrieval of project aggregates for public and admin flows.
type ProjectService interface {
	ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error)
}

type projectService struct {
	repo repository.ProjectDocumentRepository
}

// NewProjectService constructs a project service backed by the v2 document repository.
func NewProjectService(repo repository.ProjectDocumentRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) ListProjectDocuments(ctx context.Context, includeDrafts bool) ([]model.ProjectDocument, error) {
	projects, err := s.repo.ListProjectDocuments(ctx, includeDrafts)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load project documents", err)
	}
	return projects, nil
}
