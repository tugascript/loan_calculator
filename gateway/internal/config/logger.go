package config

import (
	"log/slog"
	"strings"
)

type LoggerConfig struct {
	level   slog.Level
	service string
}

func mapLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewLoggerConfig(level string, service string) LoggerConfig {
	return LoggerConfig{
		level:   mapLogLevel(level),
		service: service,
	}
}

func (L *LoggerConfig) Level() slog.Level {
	return L.level
}

func (L *LoggerConfig) Service() string {
	return L.service
}
