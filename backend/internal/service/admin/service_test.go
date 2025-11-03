package admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/inmemory"
)

func TestService_CreateProjectAndList(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	ctx := context.Background()
	input := ProjectInput{
		Title:       model.NewLocalizedText("新規プロジェクト", "New Project"),
		Description: model.NewLocalizedText("説明", "Description"),
		TechStack:   []string{"Go", "React"},
		LinkURL:     "https://example.com",
		Year:        2025,
		Published:   true,
	}

	created, err := svc.CreateProject(ctx, input)
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	require.Equal(t, 2025, created.Year)

	projects, err := svc.ListProjects(ctx)
	require.NoError(t, err)
	var found bool
	for _, project := range projects {
		if project.ID == created.ID {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestService_CreateProjectValidation(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.CreateProject(context.Background(), ProjectInput{Year: 2025})
	require.Error(t, err)
}

func TestService_AddBlacklistEntryDuplicate(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	ctx := context.Background()

	entry, err := svc.AddBlacklistEntry(ctx, BlacklistInput{Email: "duplicate@example.com", Reason: "test"})
	require.NoError(t, err)
	require.NotZero(t, entry.ID)

	_, err = svc.AddBlacklistEntry(ctx, BlacklistInput{Email: "duplicate@example.com"})
	require.Error(t, err)
}

func TestService_Summary(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	summary, err := svc.Summary(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, summary.PublishedProjects, 1)
	require.GreaterOrEqual(t, summary.SkillCount, 1)
	require.GreaterOrEqual(t, summary.BlacklistEntries, 1)
}

func TestService_UpdateProfileNormalisesInput(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	ctx := context.Background()
	input := ProfileInput{
		Name:        model.NewLocalizedText(" 名前 ", " Name "),
		Title:       model.NewLocalizedText("肩書", "Title"),
		Affiliation: model.NewLocalizedText("所属", "Affiliation"),
		Lab:         model.NewLocalizedText("ラボ", "Lab"),
		Summary:     model.NewLocalizedText(" 要約 ", " Summary "),
		Skills: []model.LocalizedText{
			{Ja: " Go ", En: " Go "},
			{},
		},
		FocusAreas: []model.LocalizedText{
			{Ja: " AI ", En: " AI "},
			{Ja: "", En: ""},
		},
	}

	profile, err := svc.UpdateProfile(ctx, input)
	require.NoError(t, err)
	require.Len(t, profile.Skills, 1)
	require.Equal(t, "Go", profile.Skills[0].Ja)
	require.Len(t, profile.FocusAreas, 1)
	require.Equal(t, "AI", profile.FocusAreas[0].Ja)
	require.NotNil(t, profile.UpdatedAt)
}

func TestService_UpdateContactMessageInvalidStatus(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.UpdateContactMessage(context.Background(), "contact-1", ContactUpdateInput{Status: "unknown"})
	require.Error(t, err)
}

func newTestService(t *testing.T) Service {
	profileRepo := inmemory.NewProfileRepository()
	adminProfileRepo, ok := profileRepo.(repository.AdminProfileRepository)
	if !ok {
		t.Fatalf("profile repository missing admin interface")
	}

	projectRepo := inmemory.NewProjectRepository()
	adminProjectRepo, ok := projectRepo.(repository.AdminProjectRepository)
	if !ok {
		t.Fatalf("project repository missing admin interface")
	}

	researchRepo := inmemory.NewResearchRepository()
	adminResearchRepo, ok := researchRepo.(repository.AdminResearchRepository)
	if !ok {
		t.Fatalf("research repository missing admin interface")
	}

	contactRepo := inmemory.NewContactRepository()
	adminContactRepo, ok := contactRepo.(repository.AdminContactRepository)
	if !ok {
		t.Fatalf("contact repository missing admin interface")
	}

	bl := inmemory.NewBlacklistRepository()

	svc, err := NewService(adminProfileRepo, adminProjectRepo, adminResearchRepo, adminContactRepo, bl)
	require.NoError(t, err)

	return svc
}
