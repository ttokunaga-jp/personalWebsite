package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/takumi/personal-website/internal/config"
)

// NewLogger returns a JSON slog.Logger configured according to application settings.
func NewLogger(cfg *config.AppConfig) *slog.Logger {
	level := slog.LevelInfo
	if cfg != nil {
		switch strings.ToLower(cfg.Logging.Level) {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		case "info":
			level = slog.LevelInfo
		}
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
