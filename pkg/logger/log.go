package logger

import (
	"log/slog"
	"os"
	"sport-assistance/pkg/configs"
	"strings"
)

func New(cfg configs.LoggerConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	base := slog.New(handler)

	// глобальные поля для всех логов
	return base.With(
		"pid", os.Getpid(),
	)
}

// parseLevel — строка → slog.Level
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
