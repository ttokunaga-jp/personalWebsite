package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/di"
	"github.com/takumi/personal-website/internal/middleware"
	"github.com/takumi/personal-website/internal/server"
	"github.com/takumi/personal-website/internal/telemetry"
)

type Application struct {
	fx *fx.App
}

func New() *Application {
	app := fx.New(
		config.Module,
		di.Module,
		server.Module,
		fx.Provide(
			newEngine,
		),
		fx.Invoke(server.Register),
	)

	return &Application{fx: app}
}

func (a *Application) Start(ctx context.Context) error {
	startCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := a.fx.Start(startCtx); err != nil {
		return fmt.Errorf("fx start: %w", err)
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	stopCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := a.fx.Stop(stopCtx); err != nil {
		return fmt.Errorf("fx stop: %w", err)
	}

	return nil
}

func newEngine(
	cfg *config.AppConfig,
	requestID *middleware.RequestID,
	requestLogger *middleware.RequestLogger,
	charset *middleware.Charset,
	securityHeaders *middleware.SecurityHeaders,
	httpsRedirect *middleware.HTTPSRedirect,
	cors *middleware.CORSMiddleware,
	csrf *middleware.CSRFMiddleware,
	rateLimiter *middleware.RateLimiter,
	metrics *telemetry.Metrics,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	if cfg != nil {
		if len(cfg.Server.TrustedProxies) == 0 {
			_ = engine.SetTrustedProxies(nil)
		} else {
			if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
				panic(fmt.Errorf("set trusted proxies: %w", err))
			}
		}
	}

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
	if charset != nil {
		engine.Use(charset.Handler())
	}
	if securityHeaders != nil {
		engine.Use(securityHeaders.Handler())
	}
	if csrf != nil {
		engine.Use(csrf.Handler())
	}

	return engine
}
