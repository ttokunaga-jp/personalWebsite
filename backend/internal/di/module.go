package di

import (
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/calendar"
	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	firestoredb "github.com/takumi/personal-website/internal/infra/firestore"
	"github.com/takumi/personal-website/internal/infra/google"
	mysqlinfra "github.com/takumi/personal-website/internal/infra/mysql"
	"github.com/takumi/personal-website/internal/logging"
	"github.com/takumi/personal-website/internal/mail"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/repository"
	"github.com/takumi/personal-website/internal/repository/provider"
	csrfmgr "github.com/takumi/personal-website/internal/security/csrf"
	"github.com/takumi/personal-website/internal/service"
	adminservice "github.com/takumi/personal-website/internal/service/admin"
	"github.com/takumi/personal-website/internal/service/auth"
	"github.com/takumi/personal-website/internal/telemetry"
)

var Module = fx.Module("di",
	fx.Provide(
		mysqlinfra.NewClient,
		firestoredb.NewClient,
		provideAuthConfig,
		auth.NewJWTIssuer,
		auth.NewStateManager,
		auth.NewGoogleOAuthProvider,
		provideGoogleTokenStore,
		google.NewGmailTokenManager,
		auth.NewService,
		auth.NewAdminService,
		auth.NewJWTVerifier,
		provideTechCatalogRepository,
		provideProfileRepository,
		provideContentProfileRepository,
		provideProjectDocumentRepository,
		provideResearchDocumentRepository,
		provideContactSettingsRepository,
		provideHomePageConfigRepository,
		provider.NewAdminProfileRepository,
		provideProjectRepository,
		provider.NewAdminProjectRepository,
		provideResearchRepository,
		provider.NewAdminResearchRepository,
		provideContactRepository,
		provider.NewAdminContactRepository,
		provideAvailabilityRepository,
		provideBlogRepository,
		provideMeetingReservationRepository,
		provideMeetingNotificationRepository,
		provideBlacklistRepository,
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
		handler.NewAdminAuthHandler,
		handler.NewAdminHandler,
		handler.NewSecurityHandler,
		logging.NewLogger,
		middleware.NewRequestID,
		middleware.NewRequestLogger,
		middleware.NewCharsetMiddleware,
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

func provideGoogleTokenStore(client *firestore.Client, cfg *config.AppConfig) (google.TokenStore, error) {
	if client == nil || cfg == nil {
		log.Printf("google token store: firestore client not available (client=%v cfg=%v)", client != nil, cfg != nil)
		return nil, nil
	}
	if cfg.Auth.StateSecret == "" {
		log.Printf("google token store: state secret not configured")
		return nil, nil
	}
	store, err := google.NewFirestoreTokenStore(client, cfg.Firestore.CollectionPrefix, cfg.Auth.StateSecret)
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

func provideTechCatalogRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.TechCatalogRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewTechCatalogRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewTechCatalogRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewTechCatalogRepository(db, fs, cfg)
	}
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

func provideProfileRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ProfileRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewProfileRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewProfileRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewProfileRepository(db, fs, cfg)
	}
}

func provideContentProfileRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ContentProfileRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewContentProfileRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewContentProfileRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewContentProfileRepository(db, fs, cfg)
	}
}

func provideProjectDocumentRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ProjectDocumentRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewProjectDocumentRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewProjectDocumentRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewProjectDocumentRepository(db, fs, cfg)
	}
}

func provideProjectRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ProjectRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewProjectRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewProjectRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewProjectRepository(db, fs, cfg)
	}
}

func provideResearchDocumentRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ResearchDocumentRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewResearchDocumentRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewResearchDocumentRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewResearchDocumentRepository(db, fs, cfg)
	}
}

func provideResearchRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ResearchRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewResearchRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewResearchRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewResearchRepository(db, fs, cfg)
	}
}

func provideContactSettingsRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ContactFormSettingsRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewContactFormSettingsRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewContactFormSettingsRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewContactFormSettingsRepository(db, fs, cfg)
	}
}

func provideContactRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.ContactRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewContactRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewContactRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewContactRepository(db, fs, cfg)
	}
}

func provideHomePageConfigRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.HomePageConfigRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewHomePageConfigRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewHomePageConfigRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewHomePageConfigRepository(db, fs, cfg)
	}
}

func provideAvailabilityRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.AvailabilityRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewAvailabilityRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewAvailabilityRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewAvailabilityRepository(db, fs, cfg)
	}
}

func provideBlogRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.BlogRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewBlogRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewBlogRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewBlogRepository(db, fs, cfg)
	}
}

func provideMeetingReservationRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.MeetingReservationRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewMeetingReservationRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewMeetingReservationRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewMeetingReservationRepository(db, fs, cfg)
	}
}

func provideMeetingNotificationRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.MeetingNotificationRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewMeetingNotificationRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewMeetingNotificationRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewMeetingNotificationRepository(db, fs, cfg)
	}
}

func provideBlacklistRepository(cfg *config.AppConfig, db *sqlx.DB, fs *firestore.Client) repository.BlacklistRepository {
	driver := normalizedDriver(cfg)
	switch driver {
	case "firestore":
		return provider.NewBlacklistRepository(nil, fs, cfg)
	case "mysql":
		return provider.NewBlacklistRepository(db, nil, cfg)
	default:
		log.Printf("unknown db_driver %q; defaulting to mysql if available", driver)
		return provider.NewBlacklistRepository(db, fs, cfg)
	}
}

func normalizedDriver(cfg *config.AppConfig) string {
	if cfg == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(cfg.DBDriver))
}
