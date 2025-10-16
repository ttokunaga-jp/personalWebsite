package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
)

type stubResearchRepository struct {
	research []model.Research
	err      error
}

func (s *stubResearchRepository) ListResearch(context.Context) ([]model.Research, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]model.Research(nil), s.research...), nil
}

func TestResearchService_ListResearchSuccess(t *testing.T) {
	t.Parallel()

	expected := []model.Research{
		{ID: 1, Title: model.NewLocalizedText("研究1", "Research 1")},
		{ID: 2, Title: model.NewLocalizedText("研究2", "Research 2")},
	}

	service := NewResearchService(&stubResearchRepository{research: expected})

	research, err := service.ListResearch(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, research)
}

func TestResearchService_ListResearchError(t *testing.T) {
	t.Parallel()

	service := NewResearchService(&stubResearchRepository{err: errors.New("db failure")})

	research, err := service.ListResearch(context.Background())
	require.Nil(t, research)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
	require.Contains(t, appErr.Message, "failed to load research")
}
