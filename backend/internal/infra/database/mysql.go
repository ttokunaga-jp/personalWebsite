package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/takumi/personal-website/internal/config"
)

// NewMySQLClient initialises an sqlx DB handle that is ready for dependency injection.
// It returns nil when no DSN is provided, allowing tests to run without a database.
func NewMySQLClient(lc fx.Lifecycle, cfg *config.AppConfig) (*sqlx.DB, error) {
	dsn := cfg.Database.DSN
	if dsn == "" {
		return nil, nil
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	if cfg.Database.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	}
	if cfg.Database.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	}

	if !cfg.Database.SkipPingOnStart {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("ping mysql: %w", err)
		}
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}
