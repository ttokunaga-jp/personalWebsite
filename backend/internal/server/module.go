package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/middleware"
)

var Module = fx.Module("http",
	fx.Provide(
		newHTTPServer,
	),
)

func newHTTPServer(
	engine *gin.Engine,
	cfg *config.AppConfig,
	healthHandler *handler.HealthHandler,
	profileHandler *handler.ProfileHandler,
	projectHandler *handler.ProjectHandler,
	researchHandler *handler.ResearchHandler,
	contactHandler *handler.ContactHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) *http.Server {
	registerRoutes(engine, healthHandler, profileHandler, projectHandler, researchHandler, contactHandler, jwtMiddleware)

	return &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func Register(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(fmt.Errorf("listen and serve: %w", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
