package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ServerConfig struct {
	Port           string   `mapstructure:"port"`
	Mode           string   `mapstructure:"mode"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
	TLSCertFile    string   `mapstructure:"tls_cert_file"`
	TLSKeyFile     string   `mapstructure:"tls_key_file"`
}

type FirestoreConfig struct {
	ProjectID        string `mapstructure:"project_id"`
	DatabaseID       string `mapstructure:"database_id"`
	CollectionPrefix string `mapstructure:"collection_prefix"`
	EmulatorHost     string `mapstructure:"emulator_host"`
}

type AuthConfig struct {
	JWTSecret             string          `mapstructure:"jwt_secret"`
	Issuer                string          `mapstructure:"issuer"`
	Audience              []string        `mapstructure:"audience"`
	ClockSkewSeconds      int64           `mapstructure:"clock_skew_seconds"`
	AccessTokenTTLMinutes int64           `mapstructure:"access_token_ttl_minutes"`
	StateSecret           string          `mapstructure:"state_secret"`
	StateTTLSeconds       int64           `mapstructure:"state_ttl_seconds"`
	Disabled              bool            `mapstructure:"disabled"`
	Admin                 AdminAuthConfig `mapstructure:"admin"`
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
	Timezone         string   `mapstructure:"timezone"`
	SlotDurationMin  int      `mapstructure:"slot_duration_minutes"`
	WorkdayStartHour int      `mapstructure:"workday_start_hour"`
	WorkdayEndHour   int      `mapstructure:"workday_end_hour"`
	HorizonDays      int      `mapstructure:"horizon_days"`
	BufferMinutes    int      `mapstructure:"buffer_minutes"`
	Topics           []string `mapstructure:"topics"`
	RecaptchaSiteKey string   `mapstructure:"recaptcha_site_key"`
	MinimumLeadHours int      `mapstructure:"minimum_lead_hours"`
	ConsentText      string   `mapstructure:"consent_text"`
	SupportEmail     string   `mapstructure:"support_email"`
	CalendarTimezone string   `mapstructure:"calendar_timezone"`
}

type BookingConfig struct {
	CalendarID           string        `mapstructure:"calendar_id"`
	MeetTemplate         string        `mapstructure:"meet_template"`
	NotificationSender   string        `mapstructure:"notification_sender"`
	NotificationReceiver string        `mapstructure:"notification_receiver"`
	AccessTokenEnv       string        `mapstructure:"access_token_env"`
	RequestTimeout       time.Duration `mapstructure:"request_timeout"`
	MaxRetries           int           `mapstructure:"max_retries"`
	InitialBackoff       time.Duration `mapstructure:"initial_backoff"`
	BackoffMultiplier    float64       `mapstructure:"backoff_multiplier"`
	CircuitOpenSeconds   int           `mapstructure:"circuit_open_seconds"`
	CircuitFailureThresh int           `mapstructure:"circuit_failure_threshold"`
}

type SecurityConfig struct {
	EnableCSRF                 bool          `mapstructure:"enable_csrf"`
	CSRFSigningKey             string        `mapstructure:"csrf_signing_key"`
	CSRFTokenTTL               time.Duration `mapstructure:"csrf_token_ttl"`
	CSRFCookieName             string        `mapstructure:"csrf_cookie_name"`
	CSRFCookieDomain           string        `mapstructure:"csrf_cookie_domain"`
	CSRFCookieSecure           bool          `mapstructure:"csrf_cookie_secure"`
	CSRFCookieHTTPOnly         bool          `mapstructure:"csrf_cookie_http_only"`
	CSRFCookieSameSite         string        `mapstructure:"csrf_cookie_same_site"`
	CSRFHeaderName             string        `mapstructure:"csrf_header_name"`
	CSRFExemptPaths            []string      `mapstructure:"csrf_exempt_paths"`
	HTTPSRedirect              bool          `mapstructure:"https_redirect"`
	HSTSMaxAgeSeconds          int           `mapstructure:"hsts_max_age_seconds"`
	ContentSecurityPolicy      string        `mapstructure:"content_security_policy"`
	ReferrerPolicy             string        `mapstructure:"referrer_policy"`
	RateLimitRequestsPerMinute int           `mapstructure:"rate_limit_requests_per_minute"`
	RateLimitBurst             int           `mapstructure:"rate_limit_burst"`
	RateLimitWhitelist         []string      `mapstructure:"rate_limit_whitelist"`
	RateLimitExemptPaths       []string      `mapstructure:"rate_limit_exempt_paths"`
	AllowedOrigins             []string      `mapstructure:"allowed_origins"`
	AllowCredentials           bool          `mapstructure:"allow_credentials"`
}

type MetricsConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Endpoint  string `mapstructure:"endpoint"`
	Namespace string `mapstructure:"namespace"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Driver                 string        `mapstructure:"driver"`
	DSN                    string        `mapstructure:"dsn"`
	Host                   string        `mapstructure:"host"`
	Port                   int           `mapstructure:"port"`
	User                   string        `mapstructure:"user"`
	Name                   string        `mapstructure:"name"`
	InstanceConnectionName string        `mapstructure:"instance_connection_name"`
	Timezone               string        `mapstructure:"timezone"`
	MaxOpenConns           int           `mapstructure:"max_open_conns"`
	MaxIdleConns           int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime        time.Duration `mapstructure:"conn_max_lifetime"`
}

type AppConfig struct {
	Server    ServerConfig      `mapstructure:"server"`
	Firestore FirestoreConfig   `mapstructure:"firestore"`
	Auth      AuthConfig        `mapstructure:"auth"`
	Google    GoogleOAuthConfig `mapstructure:"google"`
	Contact   ContactConfig     `mapstructure:"contact"`
	Booking   BookingConfig     `mapstructure:"booking"`
	Security  SecurityConfig    `mapstructure:"security"`
	Metrics   MetricsConfig     `mapstructure:"metrics"`
	Logging   LoggingConfig     `mapstructure:"logging"`
	Database  DatabaseConfig    `mapstructure:"database"`
	DBDriver  string            `mapstructure:"db_driver"`
}

type AdminAuthConfig struct {
	DefaultRedirectURI    string        `mapstructure:"default_redirect_uri"`
	AllowedDomains        []string      `mapstructure:"allowed_domains"`
	AllowedEmails         []string      `mapstructure:"allowed_emails"`
	SessionCookieName     string        `mapstructure:"session_cookie_name"`
	SessionCookieDomain   string        `mapstructure:"session_cookie_domain"`
	SessionCookiePath     string        `mapstructure:"session_cookie_path"`
	SessionCookieSecure   bool          `mapstructure:"session_cookie_secure"`
	SessionCookieHTTPOnly bool          `mapstructure:"session_cookie_http_only"`
	SessionCookieSameSite string        `mapstructure:"session_cookie_same_site"`
	SessionTTL            time.Duration `mapstructure:"session_ttl"`
	SessionIdleTimeout    time.Duration `mapstructure:"session_idle_timeout"`
	SessionRefreshWindow  time.Duration `mapstructure:"session_refresh_window"`
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
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	v.SetDefault("server.port", "8100")
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.trusted_proxies", []string{})
	v.SetDefault("server.tls_cert_file", "")
	v.SetDefault("server.tls_key_file", "")
	v.SetDefault("firestore.project_id", "")
	v.SetDefault("firestore.database_id", "(default)")
	v.SetDefault("firestore.collection_prefix", "")
	v.SetDefault("firestore.emulator_host", "")
	v.SetDefault("auth.jwt_secret", "local-dev-secret-change-me")
	v.SetDefault("auth.issuer", "personal-website")
	v.SetDefault("auth.audience", []string{"personal-website-admin"})
	v.SetDefault("auth.clock_skew_seconds", 30)
	v.SetDefault("auth.access_token_ttl_minutes", 60)
	v.SetDefault("auth.state_secret", "local-dev-state-secret-change-me")
	v.SetDefault("auth.state_ttl_seconds", 300)
	v.SetDefault("auth.disabled", false)
	v.SetDefault("auth.admin.default_redirect_uri", "/admin/")
	v.SetDefault("auth.admin.allowed_domains", []string{})
	v.SetDefault("auth.admin.allowed_emails", []string{})
	v.SetDefault("auth.admin.session_cookie_name", "ps_admin_session")
	v.SetDefault("auth.admin.session_cookie_domain", "")
	v.SetDefault("auth.admin.session_cookie_path", "/")
	v.SetDefault("auth.admin.session_cookie_secure", true)
	v.SetDefault("auth.admin.session_cookie_http_only", true)
	v.SetDefault("auth.admin.session_cookie_same_site", "strict")
	v.SetDefault("auth.admin.session_ttl", "24h")
	v.SetDefault("auth.admin.session_idle_timeout", "2h")
	v.SetDefault("auth.admin.session_refresh_window", "20m")
	v.SetDefault("google.auth_url", "https://accounts.google.com/o/oauth2/v2/auth")
	v.SetDefault("google.token_url", "https://oauth2.googleapis.com/token")
	v.SetDefault("google.userinfo_url", "https://openidconnect.googleapis.com/v1/userinfo")
	v.SetDefault("google.client_id", "")
	v.SetDefault("google.client_secret", "")
	v.SetDefault("google.redirect_url", "")
	v.SetDefault("google.scopes", []string{
		"openid",
		"email",
		"profile",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/gmail.send",
		"https://www.googleapis.com/auth/calendar.events",
		"https://www.googleapis.com/auth/calendar.readonly",
	})
	v.SetDefault("google.allowed_domains", []string{})
	v.SetDefault("google.allowed_emails", []string{})
	v.SetDefault("contact.timezone", "Asia/Tokyo")
	v.SetDefault("contact.slot_duration_minutes", 30)
	v.SetDefault("contact.workday_start_hour", 9)
	v.SetDefault("contact.workday_end_hour", 18)
	v.SetDefault("contact.horizon_days", 14)
	v.SetDefault("contact.buffer_minutes", 30)
	v.SetDefault("contact.topics", []string{})
	v.SetDefault("contact.recaptcha_site_key", "")
	v.SetDefault("contact.minimum_lead_hours", 24)
	v.SetDefault("contact.consent_text", "")
	v.SetDefault("booking.request_timeout", 8*time.Second)
	v.SetDefault("booking.max_retries", 3)
	v.SetDefault("booking.initial_backoff", 750*time.Millisecond)
	v.SetDefault("booking.backoff_multiplier", 2.0)
	v.SetDefault("booking.circuit_open_seconds", 60)
	v.SetDefault("booking.circuit_failure_threshold", 3)
	v.SetDefault("booking.access_token_env", "")
	v.SetDefault("security.enable_csrf", true)
	v.SetDefault("security.csrf_signing_key", "local-dev-csrf-secret-change-me")
	v.SetDefault("security.csrf_token_ttl", 3600*time.Second)
	v.SetDefault("security.csrf_cookie_name", "ps_csrf")
	v.SetDefault("security.csrf_cookie_domain", "")
	v.SetDefault("security.csrf_cookie_secure", true)
	v.SetDefault("security.csrf_cookie_http_only", true)
	v.SetDefault("security.csrf_cookie_same_site", "strict")
	v.SetDefault("security.csrf_header_name", "X-CSRF-Token")
	v.SetDefault("security.csrf_exempt_paths", []string{"/api/auth/callback"})
	v.SetDefault("security.https_redirect", true)
	v.SetDefault("security.hsts_max_age_seconds", 63072000)
	v.SetDefault("security.content_security_policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'")
	v.SetDefault("security.referrer_policy", "strict-origin-when-cross-origin")
	v.SetDefault("db_driver", "mysql")
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.dsn", "")
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "")
	v.SetDefault("database.name", "")
	v.SetDefault("database.instance_connection_name", "")
	v.SetDefault("database.timezone", "Asia/Tokyo")
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 30*time.Minute)
	v.SetDefault("security.rate_limit_requests_per_minute", 120)
	v.SetDefault("security.rate_limit_burst", 20)
	v.SetDefault("security.rate_limit_whitelist", []string{})
	v.SetDefault("security.rate_limit_exempt_paths", []string{
		"/api/admin/auth/login",
		"/api/admin/auth/callback",
		"/api/admin/auth/session",
	})
	v.SetDefault("security.allowed_origins", []string{})
	v.SetDefault("security.allow_credentials", true)
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.endpoint", "/metrics")
	v.SetDefault("metrics.namespace", "personal_website")
	v.SetDefault("logging.level", "info")

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
