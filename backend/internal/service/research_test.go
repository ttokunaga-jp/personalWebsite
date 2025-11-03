package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
)

type stubResearchDocumentRepository struct {
	research []model.ResearchDocument
	err      error
}

func (s *stubResearchDocumentRepository) ListResearchDocuments(context.Context, bool) ([]model.ResearchDocument, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]model.ResearchDocument(nil), s.research...), nil
}

func TestResearchService_ListResearchSuccess(t *testing.T) {
	t.Parallel()

	expected := []model.ResearchDocument{
		{ID: 1, Slug: "research-1", Title: model.NewLocalizedText("研究1", "Research 1")},
		{ID: 2, Slug: "research-2", Title: model.NewLocalizedText("研究2", "Research 2")},
	}

	service := NewResearchService(&stubResearchDocumentRepository{research: expected})

	research, err := service.ListResearchDocuments(context.Background(), false)
	require.NoError(t, err)
	require.Equal(t, expected, research)
}

func TestResearchService_ListResearchError(t *testing.T) {
	t.Parallel()

	service := NewResearchService(&stubResearchDocumentRepository{err: errors.New("db failure")})

	research, err := service.ListResearchDocuments(context.Background(), false)
	require.Nil(t, research)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
	require.Contains(t, appErr.Message, "failed to load research documents")
}
