package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
)

type stubProfileRepository struct {
	profile *model.Profile
	err     error
}

func (s *stubProfileRepository) GetProfile(context.Context) (*model.Profile, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.profile, nil
}

func TestProfileService_GetProfileSuccess(t *testing.T) {
	t.Parallel()

	expected := &model.Profile{
		Name:        model.NewLocalizedText("氏名", "Name"),
		Title:       model.NewLocalizedText("肩書き", "Title"),
		Affiliation: model.NewLocalizedText("所属", "Affiliation"),
		Lab:         model.NewLocalizedText("研究室", "Lab"),
		Summary:     model.NewLocalizedText("概要", "Summary"),
	}

	service := NewProfileService(&stubProfileRepository{profile: expected})

	actual, err := service.GetProfile(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProfileService_GetProfileError(t *testing.T) {
	t.Parallel()

	service := NewProfileService(&stubProfileRepository{err: errors.New("db failure")})

	profile, err := service.GetProfile(context.Background())
	require.Nil(t, profile)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
	require.Contains(t, appErr.Message, "failed to load profile")
}
