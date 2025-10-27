package di

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/calendar"
	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/infra/database"
	"github.com/takumi/personal-website/internal/infra/google"
	"github.com/takumi/personal-website/internal/logging"
	"github.com/takumi/personal-website/internal/mail"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/repository/provider"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
	"github.com/takumi/personal-website/internal/service"
	adminservice "github.com/takumi/personal-website/internal/service/admin"
	"github.com/takumi/personal-website/internal/service/auth"
	"github.com/takumi/personal-website/internal/telemetry"
)

var Module = fx.Module("di",
	fx.Provide(
		database.NewMySQLClient,
		provideAuthConfig,
		auth.NewJWTIssuer,
		auth.NewStateManager,
		auth.NewGoogleOAuthProvider,
		provideGoogleTokenStore,
		google.NewGmailTokenManager,
		auth.NewService,
		auth.NewJWTVerifier,
		provider.NewProfileRepository,
		provider.NewProjectRepository,
		provider.NewAdminProjectRepository,
		provider.NewResearchRepository,
		provider.NewAdminResearchRepository,
		provider.NewContactRepository,
		provider.NewAvailabilityRepository,
		provider.NewBlogRepository,
		provider.NewMeetingRepository,
		provider.NewBlacklistRepository,
		provideHTTPClient,
		provideGoogleTokenProvider,
		provideCalendarClient,
		provideGmailClient,
		service.NewProfileService,
		service.NewProjectService,
		service.NewResearchService,
		service.NewContactService,
		service.NewAvailabilityService,
		service.NewBookingService,
		adminservice.NewService,
		handler.NewHealthHandler,
		handler.NewProfileHandler,
		handler.NewProjectHandler,
		handler.NewResearchHandler,
		handler.NewContactHandler,
		handler.NewBookingHandler,
		handler.NewAuthHandler,
		handler.NewAdminHandler,
		handler.NewSecurityHandler,
		logging.NewLogger,
		middleware.NewRequestID,
		middleware.NewRequestLogger,
		middleware.NewSecurityHeaders,
		middleware.NewHTTPSRedirect,
		middleware.NewCORSMiddleware,
		middleware.NewRateLimiter,
		middleware.NewJWTMiddleware,
		middleware.NewAdminGuard,
		middleware.NewCSRFMiddleware,
		provideCSRFManager,
		telemetry.NewMetrics,
	),
)

func provideAuthConfig(cfg *config.AppConfig) config.AuthConfig {
	if cfg == nil {
		return config.AuthConfig{}
	}
	return cfg.Auth
}

func provideHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

func provideGoogleTokenStore(db *sqlx.DB, cfg *config.AppConfig) (google.TokenStore, error) {
	if db == nil || cfg == nil {
		log.Printf("google token store: database not available (db=%v cfg=%v)", db != nil, cfg != nil)
		return nil, nil
	}
	if cfg.Auth.StateSecret == "" {
		log.Printf("google token store: state secret not configured")
		return nil, nil
	}
	store, err := google.NewMySQLTokenStore(db, cfg.Auth.StateSecret)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func provideGoogleTokenProvider(client *http.Client, cfg *config.AppConfig, store google.TokenStore) google.TokenProvider {
	var providers []google.TokenProvider
	if cfg != nil && store != nil {
		if provider, err := google.NewRefreshingTokenProvider(cfg.Google, store, client); err == nil {
			providers = append(providers, provider)
		} else {
			log.Printf("refreshing token provider disabled: %v", err)
		}
	}

	envVar := ""
	if cfg != nil {
		envVar = cfg.Booking.AccessTokenEnv
	}
	if strings.TrimSpace(envVar) == "" {
		envVar = "GOOGLE_GMAIL_TOKEN"
	}
	if strings.TrimSpace(envVar) != "" {
		providers = append(providers, &google.EnvTokenProvider{EnvVar: envVar})
	}

	return google.NewFallbackTokenProvider(providers...)
}

func provideCalendarClient(client *http.Client, provider google.TokenProvider, cfg *config.AppConfig) calendar.Client {
	return google.NewCalendarAPIClient(client, provider, cfg.Contact.Timezone)
}

func provideGmailClient(client *http.Client, provider google.TokenProvider) mail.Client {
	return google.NewGmailAPIClient(client, provider)
}

func provideCSRFManager(cfg *config.AppConfig) *csrfmgr.Manager {
	if cfg == nil || !cfg.Security.EnableCSRF {
		return nil
	}
	return csrfmgr.NewManager(cfg.Security.CSRFSigningKey, cfg.Security.CSRFTokenTTL)
}
