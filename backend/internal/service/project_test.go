package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
)

type stubProjectDocumentRepository struct {
	projects []model.ProjectDocument
	err      error
}

func (s *stubProjectDocumentRepository) ListProjectDocuments(context.Context, bool) ([]model.ProjectDocument, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]model.ProjectDocument(nil), s.projects...), nil
}

func TestProjectService_ListProjectsSuccess(t *testing.T) {
	t.Parallel()

	expected := []model.ProjectDocument{
		{ID: 1, Slug: "alpha", Title: model.NewLocalizedText("プロジェクト1", "Project 1")},
		{ID: 2, Slug: "beta", Title: model.NewLocalizedText("プロジェクト2", "Project 2")},
	}

	service := NewProjectService(&stubProjectDocumentRepository{projects: expected})

	projects, err := service.ListProjectDocuments(context.Background(), false)
	require.NoError(t, err)
	require.Equal(t, expected, projects)
}

func TestProjectService_ListProjectsError(t *testing.T) {
	t.Parallel()

	service := NewProjectService(&stubProjectDocumentRepository{err: errors.New("db failure")})

	projects, err := service.ListProjectDocuments(context.Background(), false)
	require.Nil(t, projects)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
	require.Contains(t, appErr.Message, "failed to load project documents")
}
