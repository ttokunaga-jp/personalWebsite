package di

import (
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/infra/database"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/repository/inmemory"
	"github.com/takumi/personal-website/internal/service"
	"github.com/takumi/personal-website/internal/service/auth"
)

var Module = fx.Module("di",
	fx.Provide(
		database.NewMySQLClient,
		auth.NewJWTVerifier,
		inmemory.NewProfileRepository,
		inmemory.NewProjectRepository,
		inmemory.NewResearchRepository,
		inmemory.NewContactRepository,
		service.NewProfileService,
		service.NewProjectService,
		service.NewResearchService,
		service.NewContactService,
		handler.NewHealthHandler,
		handler.NewProfileHandler,
		handler.NewProjectHandler,
		handler.NewResearchHandler,
		handler.NewContactHandler,
		middleware.NewJWTMiddleware,
	),
)
