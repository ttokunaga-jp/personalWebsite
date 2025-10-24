package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service"
	adminsvc "github.com/takumi/personal-website/internal/service/admin"
	"github.com/takumi/personal-website/internal/service/auth"
)

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	engine := gin.New()

	profileSvc := service.NewProfileService(inmemory.NewProfileRepository())
	projectSvc := service.NewProjectService(inmemory.NewProjectRepository())
	researchSvc := service.NewResearchService(inmemory.NewResearchRepository())
	contactSvc := service.NewContactService(inmemory.NewContactRepository())
	availabilitySvc := &stubAvailabilityService{
		response: &model.AvailabilityResponse{
			Timezone:    "Asia/Tokyo",
			GeneratedAt: time.Unix(0, 0),
			Days: []model.AvailabilityDay{
				{
					Date: "1970-01-01",
					Slots: []model.AvailabilitySlot{
						{
							Start: time.Unix(0, 0),
							End:   time.Unix(1800, 0),
						},
					},
				},
			},
		},
	}

	jwtVerifier := auth.NewJWTVerifier(config.AuthConfig{
		JWTSecret:        "test-secret",
		Issuer:           "personal-website",
		Audience:         []string{"personal-website-admin"},
		ClockSkewSeconds: 30,
		Disabled:         true,
	})
	jwtMiddleware := middleware.NewJWTMiddleware(jwtVerifier)
	adminSvc := &stubAdminService{}

	registerRoutes(
		engine,
		handler.NewHealthHandler(),
		handler.NewProfileHandler(profileSvc),
		handler.NewProjectHandler(projectSvc),
		handler.NewResearchHandler(researchSvc),
		handler.NewContactHandler(contactSvc, availabilitySvc),
		handler.NewBookingHandler(&stubBookingService{}),
		handler.NewAuthHandler(&stubAuthService{}),
		jwtMiddleware,
		handler.NewAdminHandler(adminSvc),
		middleware.NewAdminGuard(),
	)

	t.Run("health route ok", func(t *testing.T) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"status":"ok"`)
	})

	t.Run("booking route schedules meeting", func(t *testing.T) {
		t.Helper()
		start := time.Now().Add(2 * time.Hour).UTC()
		body, err := json.Marshal(model.BookingRequest{
			Name:            "Alan Turing",
			Email:           "alan@example.com",
			StartTime:       start,
			DurationMinutes: 45,
			Agenda:          "Discuss computation theory",
		})
		require.NoError(t, err)

		rec := performRequest(engine, http.MethodPost, "/api/contact/bookings", body)
		require.Equal(t, http.StatusCreated, rec.Code)
		require.Contains(t, rec.Body.String(), `"calendarEventId"`)
	})

	t.Run("profile route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/profile", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("projects route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/projects", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("research route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/research", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("contact route accepts payload", func(t *testing.T) {
		t.Helper()
		body, err := json.Marshal(model.ContactRequest{
			Name:    "Ada Lovelace",
			Email:   "ada@example.com",
			Message: "I'd like to learn more about your research.",
		})
		require.NoError(t, err)

		rec := performRequest(engine, http.MethodPost, "/api/contact", body)
		require.Equal(t, http.StatusAccepted, rec.Code)
		require.Contains(t, rec.Body.String(), `"data"`)
	})

	t.Run("availability route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/contact/availability", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "Asia/Tokyo")
	})

	t.Run("auth login route responds", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/auth/login", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"state"`)
	})

	t.Run("auth callback route responds", func(t *testing.T) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/callback?state=stub&code=ok", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"token"`)
	})
}

func performRequest(engine *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

type stubAuthService struct{}

func (s *stubAuthService) StartLogin(context.Context, string) (*auth.LoginResult, error) {
	return &auth.LoginResult{
		AuthURL: "",
		State:   "stub-state",
	}, nil
}

func (s *stubAuthService) HandleCallback(context.Context, string, string) (*auth.CallbackResult, error) {
	return &auth.CallbackResult{
		Token:       "stub-token",
		ExpiresAt:   123,
		RedirectURI: "/admin",
	}, nil
}

type stubAvailabilityService struct {
	response *model.AvailabilityResponse
}

func (s *stubAvailabilityService) GetAvailability(context.Context, service.AvailabilityOptions) (*model.AvailabilityResponse, error) {
	return s.response, nil
}

type stubBookingService struct{}

func (s *stubBookingService) Book(context.Context, model.BookingRequest) (*model.BookingResult, error) {
	now := time.Now().UTC()
	return &model.BookingResult{
		Meeting: model.Meeting{
			ID:              1,
			Name:            "Alan Turing",
			Email:           "alan@example.com",
			Datetime:        now,
			DurationMinutes: 45,
			MeetURL:         "https://meet.example.com/test",
			CalendarEventID: "evt-123",
			Status:          model.MeetingStatusPending,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		CalendarEventID: "evt-123",
	}, nil
}

type stubAdminService struct{}

func (s *stubAdminService) ListProjects(context.Context) ([]model.AdminProject, error) {
	return nil, nil
}

func (s *stubAdminService) GetProject(context.Context, int64) (*model.AdminProject, error) {
	return &model.AdminProject{}, nil
}

func (s *stubAdminService) CreateProject(context.Context, adminsvc.ProjectInput) (*model.AdminProject, error) {
	return &model.AdminProject{}, nil
}

func (s *stubAdminService) UpdateProject(context.Context, int64, adminsvc.ProjectInput) (*model.AdminProject, error) {
	return &model.AdminProject{}, nil
}

func (s *stubAdminService) DeleteProject(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) ListResearch(context.Context) ([]model.AdminResearch, error) {
	return nil, nil
}

func (s *stubAdminService) GetResearch(context.Context, int64) (*model.AdminResearch, error) {
	return &model.AdminResearch{}, nil
}

func (s *stubAdminService) CreateResearch(context.Context, adminsvc.ResearchInput) (*model.AdminResearch, error) {
	return &model.AdminResearch{}, nil
}

func (s *stubAdminService) UpdateResearch(context.Context, int64, adminsvc.ResearchInput) (*model.AdminResearch, error) {
	return &model.AdminResearch{}, nil
}

func (s *stubAdminService) DeleteResearch(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) ListBlogPosts(context.Context) ([]model.BlogPost, error) {
	return nil, nil
}

func (s *stubAdminService) GetBlogPost(context.Context, int64) (*model.BlogPost, error) {
	return &model.BlogPost{}, nil
}

func (s *stubAdminService) CreateBlogPost(context.Context, adminsvc.BlogPostInput) (*model.BlogPost, error) {
	return &model.BlogPost{}, nil
}

func (s *stubAdminService) UpdateBlogPost(context.Context, int64, adminsvc.BlogPostInput) (*model.BlogPost, error) {
	return &model.BlogPost{}, nil
}

func (s *stubAdminService) DeleteBlogPost(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) ListMeetings(context.Context) ([]model.Meeting, error) {
	return nil, nil
}

func (s *stubAdminService) GetMeeting(context.Context, int64) (*model.Meeting, error) {
	return &model.Meeting{}, nil
}

func (s *stubAdminService) CreateMeeting(context.Context, adminsvc.MeetingInput) (*model.Meeting, error) {
	return &model.Meeting{}, nil
}

func (s *stubAdminService) UpdateMeeting(context.Context, int64, adminsvc.MeetingInput) (*model.Meeting, error) {
	return &model.Meeting{}, nil
}

func (s *stubAdminService) DeleteMeeting(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) ListBlacklist(context.Context) ([]model.BlacklistEntry, error) {
	return nil, nil
}

func (s *stubAdminService) AddBlacklistEntry(context.Context, adminsvc.BlacklistInput) (*model.BlacklistEntry, error) {
	return &model.BlacklistEntry{}, nil
}

func (s *stubAdminService) RemoveBlacklistEntry(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) IsEmailBlacklisted(context.Context, string) (bool, error) {
	return false, nil
}

func (s *stubAdminService) Summary(context.Context) (*model.AdminSummary, error) {
	return &model.AdminSummary{}, nil
}
