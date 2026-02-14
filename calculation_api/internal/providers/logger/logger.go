package logger

import (
	"log/slog"
	"os"

	"github.com/tugascript/loan_calculator/calculation_api/internal/config"
)

func DefaultLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))
}

func ConfigLogger(cfg config.LoggerConfig, env string) *slog.Logger {
	var logger *slog.Logger

	if env == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Level(),
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Level(),
		}))
	}

	return logger.With("service", cfg.Service())
}
