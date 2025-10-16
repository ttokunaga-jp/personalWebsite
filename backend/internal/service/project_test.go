package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
)

type stubProjectRepository struct {
	projects []model.Project
	err      error
}

func (s *stubProjectRepository) ListProjects(context.Context) ([]model.Project, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]model.Project(nil), s.projects...), nil
}

func TestProjectService_ListProjectsSuccess(t *testing.T) {
	t.Parallel()

	expected := []model.Project{
		{ID: 1, Title: model.NewLocalizedText("プロジェクト1", "Project 1")},
		{ID: 2, Title: model.NewLocalizedText("プロジェクト2", "Project 2")},
	}

	service := NewProjectService(&stubProjectRepository{projects: expected})

	projects, err := service.ListProjects(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, projects)
}

func TestProjectService_ListProjectsError(t *testing.T) {
	t.Parallel()

	service := NewProjectService(&stubProjectRepository{err: errors.New("db failure")})

	projects, err := service.ListProjects(context.Background())
	require.Nil(t, projects)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
	require.Contains(t, appErr.Message, "failed to load projects")
}
