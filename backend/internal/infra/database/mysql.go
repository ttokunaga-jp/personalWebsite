package database

import (
	"context"
	"fmt"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
)

// NewMySQLClient initialises an sqlx DB handle that is ready for dependency injection.
// It returns nil when no DSN is provided, allowing tests to run without a database.
func NewMySQLClient(lc fx.Lifecycle, cfg *config.AppConfig) (*sqlx.DB, error) {
	dsn := cfg.Database.DSN
	if dsn == "" {
		fmt.Println("mysql client: DSN not provided; database disabled")
		return nil, nil
	}

	fmt.Printf("mysql client: connecting to %s\n", dsn)

	normalizedDSN, err := normalizeDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("normalize mysql dsn: %w", err)
	}

	db, err := sqlx.Open("mysql", normalizedDSN)
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

func normalizeDSN(dsn string) (string, error) {
	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		return "", err
	}

	if cfg.Params == nil {
		cfg.Params = make(map[string]string)
	}

	if _, ok := cfg.Params["charset"]; !ok {
		cfg.Params["charset"] = "utf8mb4"
	}

	if cfg.Collation == "" {
		cfg.Collation = "utf8mb4_unicode_ci"
	}

	if !cfg.ParseTime {
		cfg.ParseTime = true
	}

	return cfg.FormatDSN(), nil
}
