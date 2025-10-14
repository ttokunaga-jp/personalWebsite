package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
	"github.com/takumi/personal-website/internal/di"
	"github.com/takumi/personal-website/internal/server"
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

func newEngine(cfg *config.AppConfig) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return engine
}
