package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type stubContentProfileRepository struct {
	document *model.ProfileDocument
	err      error
}

func (s *stubContentProfileRepository) GetProfileDocument(context.Context) (*model.ProfileDocument, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.document, nil
}

func TestProfileService_GetProfileDocumentSuccess(t *testing.T) {
	t.Parallel()

	expected := &model.ProfileDocument{
		ID:          1,
		DisplayName: "Takumi",
		Headline:    model.NewLocalizedText("エンジニア", "Engineer"),
		Summary:     model.NewLocalizedText("紹介", "Summary"),
		Theme: model.ProfileTheme{
			Mode:        model.ProfileThemeModeLight,
			AccentColor: "#ff9900",
		},
		UpdatedAt: time.Now(),
	}

	svc := NewProfileService(&stubContentProfileRepository{document: expected})

	actual, err := svc.GetProfileDocument(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProfileService_GetProfileDocumentError(t *testing.T) {
	t.Parallel()

	svc := NewProfileService(&stubContentProfileRepository{err: repository.ErrNotFound})

	document, err := svc.GetProfileDocument(context.Background())
	require.Nil(t, document)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeNotFound, appErr.Code)
}

func TestProfileService_EmptyProfileReturnsError(t *testing.T) {
	t.Parallel()

	svc := NewProfileService(&stubContentProfileRepository{document: nil})

	doc, err := svc.GetProfileDocument(context.Background())
	require.Nil(t, doc)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
}
