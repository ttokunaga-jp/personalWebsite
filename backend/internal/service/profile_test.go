package service

import (
	"context"
	"testing"
	"time"

	mysqlerr "github.com/go-sql-driver/mysql"
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

type stubHomePageConfigRepository struct {
	document *model.HomePageConfigDocument
	err      error
}

func (s *stubHomePageConfigRepository) GetHomePageConfig(context.Context) (*model.HomePageConfigDocument, error) {
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

	homeConfig := &model.HomePageConfigDocument{
		ID:           10,
		ProfileID:    expected.ID,
		HeroSubtitle: model.NewLocalizedText("ホーム", "Home"),
		QuickLinks:   []model.HomeQuickLink{},
		ChipSources:  []model.HomeChipSource{},
		UpdatedAt:    time.Now(),
	}

	svc := NewProfileService(
		&stubContentProfileRepository{document: expected},
		&stubHomePageConfigRepository{document: homeConfig},
	)

	actual, err := svc.GetProfileDocument(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.Equal(t, homeConfig, actual.Home)
}

func TestProfileService_GetProfileDocumentError(t *testing.T) {
	t.Parallel()

	svc := NewProfileService(
		&stubContentProfileRepository{err: repository.ErrNotFound},
		&stubHomePageConfigRepository{},
	)

	document, err := svc.GetProfileDocument(context.Background())
	require.NoError(t, err)
	require.NotNil(t, document)
	require.NotEmpty(t, document.DisplayName)
	require.NotNil(t, document.Home)
}

func TestProfileService_EmptyProfileReturnsError(t *testing.T) {
	t.Parallel()

	svc := NewProfileService(
		&stubContentProfileRepository{document: nil},
		&stubHomePageConfigRepository{},
	)

	doc, err := svc.GetProfileDocument(context.Background())
	require.Nil(t, doc)
	require.Error(t, err)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInternal, appErr.Code)
}

func TestProfileService_HomeConfigNotFoundIsIgnored(t *testing.T) {
	t.Parallel()

	profile := &model.ProfileDocument{
		ID:          1,
		DisplayName: "Takumi",
		Headline:    model.NewLocalizedText("エンジニア", "Engineer"),
		UpdatedAt:   time.Now(),
	}

	svc := NewProfileService(
		&stubContentProfileRepository{document: profile},
		&stubHomePageConfigRepository{err: repository.ErrNotFound},
	)

	actual, err := svc.GetProfileDocument(context.Background())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Home)
}

func TestProfileService_HomeConfigErrorPropagates(t *testing.T) {
	t.Parallel()

	profile := &model.ProfileDocument{
		ID:          1,
		DisplayName: "Takumi",
		Headline:    model.NewLocalizedText("エンジニア", "Engineer"),
		UpdatedAt:   time.Now(),
	}

	svc := NewProfileService(
		&stubContentProfileRepository{document: profile},
		&stubHomePageConfigRepository{err: repository.ErrInvalidInput},
	)

	document, err := svc.GetProfileDocument(context.Background())
	require.Error(t, err)
	require.Nil(t, document)

	appErr := errs.From(err)
	require.Equal(t, errs.CodeInvalidInput, appErr.Code)
}

func TestProfileService_FallbackOnMissingTable(t *testing.T) {
	t.Parallel()

	mysqlErr := &mysqlerr.MySQLError{
		Number:  1146,
		Message: "Table 'personal_website.profiles' doesn't exist",
	}

	svc := NewProfileService(
		&stubContentProfileRepository{err: mysqlErr},
		&stubHomePageConfigRepository{err: mysqlErr},
	)

	doc, err := svc.GetProfileDocument(context.Background())
	require.NoError(t, err)
	require.NotNil(t, doc)
	require.NotEmpty(t, doc.DisplayName)
	require.NotNil(t, doc.Home)
}
