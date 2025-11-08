package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/handler"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/telemetry"
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
	bookingHandler *handler.BookingHandler,
	authHandler *handler.AuthHandler,
	adminAuthHandler *handler.AdminAuthHandler,
	sessionMiddleware *middleware.AdminSessionMiddleware,
	adminHandler *handler.AdminHandler,
	adminGuard *middleware.AdminGuard,
	securityHandler *handler.SecurityHandler,
	metrics *telemetry.Metrics,
) *http.Server {
	registerRoutes(engine, healthHandler, profileHandler, projectHandler, researchHandler, contactHandler, bookingHandler, authHandler, adminAuthHandler, sessionMiddleware, adminHandler, adminGuard, securityHandler)
	if metrics != nil {
		metrics.Register(engine)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:           engine,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if cfg != nil && cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
		if err != nil {
			panic(fmt.Errorf("load tls certificate: %w", err))
		}
		srv.TLSConfig = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{certificate},
		}
	}

	return srv
}

func Register(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				var err error
				if srv.TLSConfig != nil {
					err = srv.ListenAndServeTLS("", "")
				} else {
					err = srv.ListenAndServe()
				}
				if err != nil && err != http.ErrServerClosed {
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
