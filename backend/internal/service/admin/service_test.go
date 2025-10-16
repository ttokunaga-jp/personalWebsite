package admin

import (
	"context"
	"testing"
	"time"

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
	require.GreaterOrEqual(t, summary.BlacklistEntries, 1)
}

func newTestService(t *testing.T) Service {
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

	blogs := inmemory.NewBlogRepository()
	meetings := inmemory.NewMeetingRepository()
	bl := inmemory.NewBlacklistRepository()

	svc, err := NewService(adminProjectRepo, adminResearchRepo, blogs, meetings, bl)
	require.NoError(t, err)

	// Seed deterministic state for meeting creations to avoid time.Now drift.
	_, err = svc.CreateMeeting(context.Background(), MeetingInput{
		Name:            "Test User",
		Email:           "test@example.com",
		Datetime:        time.Date(2025, 1, 10, 9, 0, 0, 0, time.UTC),
		DurationMinutes: 30,
		Status:          model.MeetingStatusPending,
	})
	require.NoError(t, err)

	return svc
}
