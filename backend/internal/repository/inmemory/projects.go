package inmemory

import (
	"context"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type projectRepository struct {
	projects []model.Project
}

func NewProjectRepository() repository.ProjectRepository {
	return &projectRepository{
		projects: defaultProjects,
	}
}

func (r *projectRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	return r.projects, nil
}
