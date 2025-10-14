package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN              string        `mapstructure:"dsn"`
	MaxOpenConns     int           `mapstructure:"max_open_conns"`
	MaxIdleConns     int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
	SkipPingOnStart  bool          `mapstructure:"skip_ping_on_start"`
	MigrationsEnable bool          `mapstructure:"migrations_enable"`
}

type AuthConfig struct {
	JWTSecret        string   `mapstructure:"jwt_secret"`
	Issuer           string   `mapstructure:"issuer"`
	Audience         []string `mapstructure:"audience"`
	ClockSkewSeconds int64    `mapstructure:"clock_skew_seconds"`
	Disabled         bool     `mapstructure:"disabled"`
}

type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

var Module = fx.Module("config",
	fx.Provide(load),
)

func load() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "release")
	v.SetDefault("database.dsn", "")
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 30*time.Minute)
	v.SetDefault("database.skip_ping_on_start", true)
	v.SetDefault("database.migrations_enable", false)
	v.SetDefault("auth.jwt_secret", "local-dev-secret-change-me")
	v.SetDefault("auth.issuer", "personal-website")
	v.SetDefault("auth.audience", []string{"personal-website-admin"})
	v.SetDefault("auth.clock_skew_seconds", 30)
	v.SetDefault("auth.disabled", false)

	if err := v.ReadInConfig(); err != nil {
		// Ignore missing file; rely on defaults/env overrides.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
