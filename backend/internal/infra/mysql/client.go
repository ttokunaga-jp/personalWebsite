package mysql

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/takumi/personal-website/internal/config"
)

// NewClient initialises a sqlx-backed MySQL client derived from application configuration.
func NewClient(lc fx.Lifecycle, cfg *config.AppConfig) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("mysql client: missing app config")
	}

	if !strings.EqualFold(cfg.DBDriver, "mysql") {
		return nil, nil
	}

	dbCfg := cfg.Database
	if driver := strings.TrimSpace(dbCfg.Driver); driver != "" && !strings.EqualFold(driver, "mysql") {
		log.Printf("mysql client: driver %q is not supported, skipping connection", driver)
		return nil, nil
	}

	dsn := strings.TrimSpace(dbCfg.DSN)
	if dsn == "" {
		user := firstNonEmpty(dbCfg.User, os.Getenv("DB_USER"))
		password := firstNonEmpty(os.Getenv("DB_PASSWORD"))
		name := firstNonEmpty(dbCfg.Name, os.Getenv("DB_NAME"))

		if user == "" || password == "" || name == "" {
			log.Println("mysql client: missing DB credentials, skipping connection")
			return nil, nil
		}

		location := url.QueryEscape(firstNonEmpty(dbCfg.Timezone, "Asia/Tokyo"))
		instance := firstNonEmpty(dbCfg.InstanceConnectionName, os.Getenv("DB_INSTANCE_CONNECTION_NAME"))

		var address string
		switch {
		case instance != "":
			address = fmt.Sprintf("unix(/cloudsql/%s)", instance)
		case dbCfg.Host != "":
			address = fmt.Sprintf("tcp(%s:%d)", dbCfg.Host, valueOrDefault(dbCfg.Port, 3306))
		default:
			address = "tcp(127.0.0.1:3306)"
		}

		dsn = fmt.Sprintf("%s:%s@%s/%s?parseTime=true&charset=utf8mb4&loc=%s", user, password, address, name, location)
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql client: open: %w", err)
	}

	db.SetMaxOpenConns(valueOrDefault(dbCfg.MaxOpenConns, 10))
	db.SetMaxIdleConns(valueOrDefault(dbCfg.MaxIdleConns, 5))
	db.SetConnMaxLifetime(valueOrDefaultDuration(dbCfg.ConnMaxLifetime, 30*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql client: ping: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func valueOrDefault(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func valueOrDefaultDuration(value, fallback time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return fallback
}
