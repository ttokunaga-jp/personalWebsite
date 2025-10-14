package service

import (
	"context"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ProjectService orchestrates project listing logic.
type ProjectService interface {
	ListProjects(ctx context.Context) ([]model.Project, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) ListProjects(ctx context.Context) ([]model.Project, error) {
	projects, err := s.repo.ListProjects(ctx)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load projects", err)
	}
	return projects, nil
}
