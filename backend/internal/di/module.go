package di

import (
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/infra/database"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/repository/provider"
	"github.com/takumi/personal-website/internal/service"
	adminservice "github.com/takumi/personal-website/internal/service/admin"
	"github.com/takumi/personal-website/internal/service/auth"
)

var Module = fx.Module("di",
	fx.Provide(
		database.NewMySQLClient,
		provideAuthConfig,
		auth.NewJWTIssuer,
		auth.NewStateManager,
		auth.NewGoogleOAuthProvider,
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
		service.NewProfileService,
		service.NewProjectService,
		service.NewResearchService,
		service.NewContactService,
		service.NewAvailabilityService,
		adminservice.NewService,
		handler.NewHealthHandler,
		handler.NewProfileHandler,
		handler.NewProjectHandler,
		handler.NewResearchHandler,
		handler.NewContactHandler,
		handler.NewAuthHandler,
		handler.NewAdminHandler,
		middleware.NewJWTMiddleware,
		middleware.NewAdminGuard,
	),
)

func provideAuthConfig(cfg *config.AppConfig) config.AuthConfig {
	if cfg == nil {
		return config.AuthConfig{}
	}
	return cfg.Auth
}
