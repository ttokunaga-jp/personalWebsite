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
	JWTSecret             string   `mapstructure:"jwt_secret"`
	Issuer                string   `mapstructure:"issuer"`
	Audience              []string `mapstructure:"audience"`
	ClockSkewSeconds      int64    `mapstructure:"clock_skew_seconds"`
	AccessTokenTTLMinutes int64    `mapstructure:"access_token_ttl_minutes"`
	StateSecret           string   `mapstructure:"state_secret"`
	StateTTLSeconds       int64    `mapstructure:"state_ttl_seconds"`
	Disabled              bool     `mapstructure:"disabled"`
}

type GoogleOAuthConfig struct {
	ClientID       string   `mapstructure:"client_id"`
	ClientSecret   string   `mapstructure:"client_secret"`
	RedirectURL    string   `mapstructure:"redirect_url"`
	Scopes         []string `mapstructure:"scopes"`
	AuthURL        string   `mapstructure:"auth_url"`
	TokenURL       string   `mapstructure:"token_url"`
	UserInfoURL    string   `mapstructure:"userinfo_url"`
	AllowedDomains []string `mapstructure:"allowed_domains"`
	AllowedEmails  []string `mapstructure:"allowed_emails"`
}

type ContactConfig struct {
	Timezone         string `mapstructure:"timezone"`
	SlotDurationMin  int    `mapstructure:"slot_duration_minutes"`
	WorkdayStartHour int    `mapstructure:"workday_start_hour"`
	WorkdayEndHour   int    `mapstructure:"workday_end_hour"`
	HorizonDays      int    `mapstructure:"horizon_days"`
	BufferMinutes    int    `mapstructure:"buffer_minutes"`
}

type AppConfig struct {
	Server   ServerConfig      `mapstructure:"server"`
	Database DatabaseConfig    `mapstructure:"database"`
	Auth     AuthConfig        `mapstructure:"auth"`
	Google   GoogleOAuthConfig `mapstructure:"google"`
	Contact  ContactConfig     `mapstructure:"contact"`
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

	v.SetDefault("server.port", "8100")
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
	v.SetDefault("auth.access_token_ttl_minutes", 60)
	v.SetDefault("auth.state_secret", "local-dev-state-secret-change-me")
	v.SetDefault("auth.state_ttl_seconds", 300)
	v.SetDefault("auth.disabled", false)
	v.SetDefault("google.auth_url", "https://accounts.google.com/o/oauth2/v2/auth")
	v.SetDefault("google.token_url", "https://oauth2.googleapis.com/token")
	v.SetDefault("google.userinfo_url", "https://openidconnect.googleapis.com/v1/userinfo")
	v.SetDefault("google.scopes", []string{
		"openid",
		"email",
		"profile",
	})
	v.SetDefault("google.allowed_domains", []string{})
	v.SetDefault("google.allowed_emails", []string{})
	v.SetDefault("contact.timezone", "Asia/Tokyo")
	v.SetDefault("contact.slot_duration_minutes", 30)
	v.SetDefault("contact.workday_start_hour", 9)
	v.SetDefault("contact.workday_end_hour", 18)
	v.SetDefault("contact.horizon_days", 14)
	v.SetDefault("contact.buffer_minutes", 30)

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
