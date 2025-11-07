package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/logging"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
	"github.com/takumi/personal-website/internal/service"
	adminsvc "github.com/takumi/personal-website/internal/service/admin"
	"github.com/takumi/personal-website/internal/service/auth"
	"github.com/takumi/personal-website/internal/telemetry"
)

func TestRegisterRoutes(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	engine := gin.New()

	profileSvc := service.NewProfileService(
		inmemory.NewContentProfileRepository(),
		inmemory.NewHomePageConfigRepository(),
	)
	projectSvc := service.NewProjectService(inmemory.NewProjectDocumentRepository())
	researchSvc := service.NewResearchService(inmemory.NewResearchDocumentRepository())
	contactSvc := service.NewContactService(inmemory.NewContactRepository(), inmemory.NewContactFormSettingsRepository())
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

	appCfg := &config.AppConfig{
		Auth: config.AuthConfig{
			Admin: config.AdminAuthConfig{
				SessionCookieName:     "ps_admin_jwt",
				SessionCookiePath:     "/",
				SessionCookieHTTPOnly: true,
				SessionCookieSecure:   true,
				SessionCookieSameSite: "lax",
			},
		},
		Contact: config.ContactConfig{
			Topics:           []string{"Research collaboration"},
			RecaptchaSiteKey: "test-site-key",
			MinimumLeadHours: 48,
			ConsentText:      "Testing purposes only.",
		},
	}

	jwtVerifier := auth.NewJWTVerifier(config.AuthConfig{
		JWTSecret:        "test-secret",
		Issuer:           "personal-website",
		Audience:         []string{"personal-website-admin"},
		ClockSkewSeconds: 30,
		Disabled:         true,
	})
	jwtMiddleware := middleware.NewJWTMiddleware(jwtVerifier, appCfg.Auth)
	adminSvc := &stubAdminService{}

	registerRoutes(
		engine,
		handler.NewHealthHandler(),
		handler.NewProfileHandler(profileSvc),
		handler.NewProjectHandler(projectSvc),
		handler.NewResearchHandler(researchSvc),
		handler.NewContactHandler(contactSvc, availabilitySvc, appCfg),
		handler.NewBookingHandler(&stubBookingService{}),
		handler.NewAuthHandler(&stubAuthService{}),
		handler.NewAdminAuthHandler(&stubAdminAuthService{}, &stubTokenIssuer{}, &stubTokenVerifier{}, appCfg.Auth),
		jwtMiddleware,
		handler.NewAdminHandler(adminSvc),
		middleware.NewAdminGuard(),
		nil,
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

	t.Run("health head route ok", func(t *testing.T) {
		rec := performRequest(engine, http.MethodHead, "/api/health", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Empty(t, rec.Body.String())
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
		require.Contains(t, rec.Body.String(), `"reservation"`)
		require.Contains(t, rec.Body.String(), `"supportEmail"`)
	})

	t.Run("booking lookup route returns reservation", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/contact/bookings/lookup-123", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"lookup-123"`)
	})

	t.Run("profile route returns data", func(t *testing.T) {
		t.Helper()
		rec := performRequest(engine, http.MethodGet, "/api/profile", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, `"data"`)
		require.Contains(t, body, `"home"`)
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

func TestSecurityAndObservabilityFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("csrf issuance and validation", func(t *testing.T) {
		cfg := newSecurityTestConfig()
		cfg.Security.RateLimitRequestsPerMinute = 0
		cfg.Security.RateLimitBurst = 0

		srv := newSecurityTestEngine(t, cfg)

		csrfRec := httptest.NewRecorder()
		csrfReq := httptest.NewRequest(http.MethodGet, "/api/security/csrf", nil)
		srv.ServeHTTP(csrfRec, csrfReq)

		require.Equal(t, http.StatusOK, csrfRec.Code)

		var payload struct {
			Data struct {
				Token     string `json:"token"`
				ExpiresAt string `json:"expires_at"`
			}
		}
		require.NoError(t, json.Unmarshal(csrfRec.Body.Bytes(), &payload))
		require.NotEmpty(t, payload.Data.Token)

		var csrfCookie *http.Cookie
		for _, c := range csrfRec.Result().Cookies() {
			if c.Name == cfg.Security.CSRFCookieName {
				csrfCookie = c
				break
			}
		}
		require.NotNil(t, csrfCookie)

		parts := strings.Split(csrfCookie.Value, ":")
		require.GreaterOrEqual(t, len(parts), 3)
		require.Equal(t, payload.Data.Token, parts[0])

		body := []byte(`{"name":"Tester","email":"tester@example.com","message":"hello"}`)

		// Missing header should be rejected.
		failRec := httptest.NewRecorder()
		failReq := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewReader(body))
		failReq.Header.Set("Content-Type", "application/json")
		failReq.AddCookie(csrfCookie)
		srv.ServeHTTP(failRec, failReq)
		require.Equal(t, http.StatusForbidden, failRec.Code)

		// Proper header + cookie should be accepted.
		successRec := httptest.NewRecorder()
		successReq := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewReader(body))
		successReq.Header.Set("Content-Type", "application/json")
		successReq.Header.Set(cfg.Security.CSRFHeaderName, payload.Data.Token)
		successReq.AddCookie(csrfCookie)
		srv.ServeHTTP(successRec, successReq)

		require.Equal(t, http.StatusAccepted, successRec.Code)
	})

	t.Run("rate limiter returns 429 after threshold", func(t *testing.T) {
		cfg := newSecurityTestConfig()
		cfg.Security.RateLimitRequestsPerMinute = 2
		cfg.Security.RateLimitBurst = 1

		srv := newSecurityTestEngine(t, cfg)

		rec1 := performRequest(srv, http.MethodGet, "/api/contact/availability", nil)
		require.Equal(t, http.StatusOK, rec1.Code)

		rec2 := performRequest(srv, http.MethodGet, "/api/contact/availability", nil)
		require.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})

	t.Run("metrics endpoint exposes prometheus data", func(t *testing.T) {
		cfg := newSecurityTestConfig()
		cfg.Security.RateLimitRequestsPerMinute = 0
		cfg.Security.RateLimitBurst = 0

		srv := newSecurityTestEngine(t, cfg)

		healthRec := performRequest(srv, http.MethodGet, "/api/health", nil)
		require.Equal(t, http.StatusOK, healthRec.Code)

		rec := performRequest(srv, http.MethodGet, "/metrics", nil)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "personal_website_http_request_duration_seconds")
	})

	t.Run("http requests redirect to https when enabled", func(t *testing.T) {
		cfg := newSecurityTestConfig()
		cfg.Security.HTTPSRedirect = true
		cfg.Security.RateLimitRequestsPerMinute = 0
		cfg.Security.RateLimitBurst = 0

		srv := newSecurityTestEngine(t, cfg)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		srv.ServeHTTP(rec, req)

		require.Equal(t, http.StatusPermanentRedirect, rec.Code)
		require.Equal(t, "https://example.com/api/health", rec.Header().Get("Location"))
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

func newSecurityTestConfig() *config.AppConfig {
	return &config.AppConfig{
		Server: config.ServerConfig{
			Mode: gin.TestMode,
		},
		Security: config.SecurityConfig{
			EnableCSRF:                 true,
			CSRFSigningKey:             "test-signing-key",
			CSRFTokenTTL:               time.Hour,
			CSRFCookieName:             "ps_csrf",
			CSRFCookieHTTPOnly:         true,
			CSRFCookieSecure:           false,
			CSRFCookieSameSite:         "lax",
			CSRFHeaderName:             "X-CSRF-Token",
			CSRFExemptPaths:            []string{"/api/auth/callback"},
			ContentSecurityPolicy:      "default-src 'self'",
			ReferrerPolicy:             "no-referrer",
			HTTPSRedirect:              false,
			RateLimitRequestsPerMinute: 0,
			RateLimitBurst:             0,
		},
		Metrics: config.MetricsConfig{
			Enabled:   true,
			Endpoint:  "/metrics",
			Namespace: "personal_website",
		},
		Logging: config.LoggingConfig{Level: "info"},
	}
}

type noopLifecycle struct{}

func (noopLifecycle) Append(fx.Hook) {}

func newSecurityTestEngine(t *testing.T, cfg *config.AppConfig) *gin.Engine {
	t.Helper()

	engine := gin.New()
	engine.Use(gin.Recovery())

	requestID := middleware.NewRequestID()
	logger := logging.NewLogger(cfg)
	requestLogger := middleware.NewRequestLogger(logger)
	securityHeaders := middleware.NewSecurityHeaders(cfg)
	httpsRedirect := middleware.NewHTTPSRedirect(cfg)
	cors := middleware.NewCORSMiddleware(cfg)
	rateLimiter := middleware.NewRateLimiter(noopLifecycle{}, cfg)
	csrfManager := csrfmgr.NewManager(cfg.Security.CSRFSigningKey, cfg.Security.CSRFTokenTTL)
	csrfMiddleware := middleware.NewCSRFMiddleware(cfg, csrfManager)
	metrics := telemetry.NewMetrics(cfg)

	if httpsRedirect != nil {
		engine.Use(httpsRedirect.Handler())
	}
	if requestID != nil {
		engine.Use(requestID.Handler())
	}
	if metrics != nil {
		engine.Use(metrics.Handler())
	}
	if requestLogger != nil {
		engine.Use(requestLogger.Handler())
	}
	if cors != nil {
		engine.Use(cors.Handler())
	}
	if rateLimiter != nil {
		engine.Use(rateLimiter.Handler())
	}
	if securityHeaders != nil {
		engine.Use(securityHeaders.Handler())
	}
	if csrfMiddleware != nil {
		engine.Use(csrfMiddleware.Handler())
	}

	profileSvc := service.NewProfileService(
		inmemory.NewContentProfileRepository(),
		inmemory.NewHomePageConfigRepository(),
	)
	projectSvc := service.NewProjectService(inmemory.NewProjectDocumentRepository())
	researchSvc := service.NewResearchService(inmemory.NewResearchDocumentRepository())
	contactSvc := service.NewContactService(inmemory.NewContactRepository(), inmemory.NewContactFormSettingsRepository())
	availabilitySvc := &stubAvailabilityService{
		response: &model.AvailabilityResponse{
			Timezone:    "Asia/Tokyo",
			GeneratedAt: time.Unix(0, 0),
			Days: []model.AvailabilityDay{{
				Date: "1970-01-01",
				Slots: []model.AvailabilitySlot{{
					Start: time.Unix(0, 0),
					End:   time.Unix(1800, 0),
				}},
			}},
		},
	}

	appCfg := &config.AppConfig{
		Auth: config.AuthConfig{
			Admin: config.AdminAuthConfig{
				SessionCookieName:     "ps_admin_jwt",
				SessionCookiePath:     "/",
				SessionCookieHTTPOnly: true,
				SessionCookieSecure:   true,
				SessionCookieSameSite: "lax",
			},
		},
		Contact: config.ContactConfig{
			Topics:           []string{"Research collaboration"},
			RecaptchaSiteKey: "test-site-key",
			MinimumLeadHours: 48,
			ConsentText:      "Testing purposes only.",
		},
	}

	jwtVerifier := auth.NewJWTVerifier(config.AuthConfig{
		JWTSecret:        "test-jwt-secret",
		Issuer:           "personal-website",
		Audience:         []string{"personal-website-admin"},
		ClockSkewSeconds: 30,
		Disabled:         true,
	})
	jwtMiddleware := middleware.NewJWTMiddleware(jwtVerifier, appCfg.Auth)
	adminSvc := &stubAdminService{}
	securityHandler := handler.NewSecurityHandler(csrfManager, cfg)

	registerRoutes(
		engine,
		handler.NewHealthHandler(),
		handler.NewProfileHandler(profileSvc),
		handler.NewProjectHandler(projectSvc),
		handler.NewResearchHandler(researchSvc),
		handler.NewContactHandler(contactSvc, availabilitySvc, appCfg),
		handler.NewBookingHandler(&stubBookingService{}),
		handler.NewAuthHandler(&stubAuthService{}),
		handler.NewAdminAuthHandler(&stubAdminAuthService{}, &stubTokenIssuer{}, &stubTokenVerifier{}, appCfg.Auth),
		jwtMiddleware,
		handler.NewAdminHandler(adminSvc),
		middleware.NewAdminGuard(),
		securityHandler,
	)

	if metrics != nil {
		metrics.Register(engine)
	}

	return engine
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
		RedirectURI: "/admin/",
	}, nil
}

type stubAdminAuthService struct{}

func (s *stubAdminAuthService) StartLogin(context.Context, string) (*auth.AdminLoginResult, error) {
	return &auth.AdminLoginResult{
		AuthURL: "",
		State:   "admin-stub-state",
	}, nil
}

func (s *stubAdminAuthService) HandleCallback(context.Context, string, string) (*auth.AdminCallbackResult, error) {
	return &auth.AdminCallbackResult{
		Token:        "admin-stub-token",
		ExpiresAt:    456,
		RedirectPath: "/admin/",
	}, nil
}

type stubTokenIssuer struct{}

func (s *stubTokenIssuer) Issue(context.Context, string, string, ...string) (string, time.Time, error) {
	return "stub-issued-token", time.Now().Add(time.Hour), nil
}

type stubTokenVerifier struct{}

func (s *stubTokenVerifier) Verify(context.Context, string) (*auth.Claims, error) {
	return &auth.Claims{
		Subject:   "stub-admin",
		Email:     "stub@example.com",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(time.Hour),
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
		Reservation: model.MeetingReservation{
			ID:                     1,
			LookupHash:             "lookup-123",
			Name:                   "Alan Turing",
			Email:                  "alan@example.com",
			StartAt:                now,
			EndAt:                  now.Add(45 * time.Minute),
			DurationMinutes:        45,
			GoogleEventID:          "evt-123",
			GoogleCalendarStatus:   "confirmed",
			Status:                 model.MeetingReservationStatusPending,
			CreatedAt:              now,
			UpdatedAt:              now,
			LastNotificationSentAt: &now,
		},
		CalendarEventID:  "evt-123",
		SupportEmail:     "support@example.com",
		CalendarTimezone: "UTC",
	}, nil
}

func (s *stubBookingService) LookupReservation(context.Context, string) (*model.BookingResult, error) {
	return s.Book(context.Background(), model.BookingRequest{})
}

func (s *stubBookingService) CancelReservation(context.Context, string, string) (*model.BookingResult, error) {
	return s.Book(context.Background(), model.BookingRequest{})
}

type stubAdminService struct{}

func (s *stubAdminService) GetProfile(context.Context) (*model.AdminProfile, error) {
	now := time.Now().UTC()
	return &model.AdminProfile{
		Name:        model.NewLocalizedText("管理者", "Admin"),
		Title:       model.NewLocalizedText("肩書", "Title"),
		Affiliation: model.NewLocalizedText("所属", "Affiliation"),
		Lab:         model.NewLocalizedText("ラボ", "Lab"),
		Summary:     model.NewLocalizedText("要約", "Summary"),
		Skills: []model.LocalizedText{
			model.NewLocalizedText("Go", "Go"),
			model.NewLocalizedText("React", "React"),
		},
		FocusAreas: []model.LocalizedText{
			model.NewLocalizedText("AI", "AI"),
		},
		UpdatedAt: &now,
	}, nil
}

func (s *stubAdminService) UpdateProfile(ctx context.Context, input adminsvc.ProfileInput) (*model.AdminProfile, error) {
	return &model.AdminProfile{
		Name:        input.Name,
		Title:       input.Title,
		Affiliation: input.Affiliation,
		Lab:         input.Lab,
		Summary:     input.Summary,
		Skills:      input.Skills,
		FocusAreas:  input.FocusAreas,
	}, nil
}

func (s *stubAdminService) ListProjects(context.Context) ([]model.AdminProject, error) {
	return []model.AdminProject{
		{
			ID:          1,
			Title:       model.NewLocalizedText("プロジェクト", "Project"),
			Description: model.NewLocalizedText("説明", "Description"),
			Tech: []model.TechMembership{
				{
					MembershipID: 1,
					EntityType:   "project",
					EntityID:     1,
					Tech: model.TechCatalogEntry{
						ID:          1,
						Slug:        "go",
						DisplayName: "Go",
						Level:       model.TechLevelAdvanced,
						SortOrder:   1,
						Active:      true,
					},
					Context:   model.TechContextPrimary,
					SortOrder: 1,
				},
			},
			LinkURL:   "https://example.com",
			Year:      2024,
			Published: true,
		},
	}, nil
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

func (s *stubAdminService) ListTechCatalog(context.Context, bool) ([]model.TechCatalogEntry, error) {
	return []model.TechCatalogEntry{
		{
			ID:          1,
			Slug:        "go",
			DisplayName: "Go",
			Level:       model.TechLevelAdvanced,
			SortOrder:   1,
			Active:      true,
		},
	}, nil
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

func (s *stubAdminService) ListContactMessages(context.Context) ([]model.ContactMessage, error) {
	now := time.Now().UTC()
	return []model.ContactMessage{
		{
			ID:        "contact-1",
			Name:      "Example",
			Email:     "example@example.com",
			Topic:     "Inquiry",
			Message:   "Hello",
			Status:    model.ContactStatusPending,
			AdminNote: "",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (s *stubAdminService) GetContactMessage(context.Context, string) (*model.ContactMessage, error) {
	return &model.ContactMessage{}, nil
}

func (s *stubAdminService) UpdateContactMessage(context.Context, string, adminsvc.ContactUpdateInput) (*model.ContactMessage, error) {
	return &model.ContactMessage{}, nil
}

func (s *stubAdminService) DeleteContactMessage(context.Context, string) error {
	return nil
}

func (s *stubAdminService) ListBlacklist(context.Context) ([]model.BlacklistEntry, error) {
	return nil, nil
}

func (s *stubAdminService) AddBlacklistEntry(context.Context, adminsvc.BlacklistInput) (*model.BlacklistEntry, error) {
	return &model.BlacklistEntry{}, nil
}

func (s *stubAdminService) UpdateBlacklistEntry(context.Context, int64, adminsvc.BlacklistInput) (*model.BlacklistEntry, error) {
	return &model.BlacklistEntry{}, nil
}

func (s *stubAdminService) RemoveBlacklistEntry(context.Context, int64) error {
	return nil
}

func (s *stubAdminService) IsEmailBlacklisted(context.Context, string) (bool, error) {
	return false, nil
}

func (s *stubAdminService) Summary(context.Context) (*model.AdminSummary, error) {
	now := time.Now().UTC()
	return &model.AdminSummary{
		ProfileUpdatedAt:  &now,
		SkillCount:        3,
		FocusAreaCount:    2,
		PublishedProjects: 1,
		DraftProjects:     0,
		PublishedResearch: 1,
		DraftResearch:     0,
		PendingContacts:   1,
		BlacklistEntries:  1,
	}, nil
}
