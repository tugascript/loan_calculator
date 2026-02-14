package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	port        int64
	env         string
	databaseURL string
	logger      LoggerConfig
}

func (C *Config) Port() int64 {
	return C.port
}

func (C *Config) ENV() string {
	return C.env
}

func (C *Config) DatabaseURL() string {
	return C.databaseURL
}

func (C *Config) Logger() LoggerConfig {
	return C.logger
}

var variables = [5]string{
	"PORT",
	"ENV",
	"DATABASE_URL",
	"LOG_LEVEL",
	"SERVICE_NAME",
}

var numericVariables = [1]string{
	"PORT",
}

func NewConfig(logger *slog.Logger, envFilePath string) *Config {
	err := godotenv.Load(envFilePath)
	if err != nil {
		logger.Debug("Error loading .env file", "error", err)
	}

	variablesMap := make(map[string]string)
	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			logger.Error(variable + " is not set")
			panic(variable + " is not set")
		}
		variablesMap[variable] = value
	}

	intMap := make(map[string]int64)
	for _, numeric := range numericVariables {
		value, err := strconv.ParseInt(variablesMap[numeric], 10, 0)
		if err != nil {
			logger.Error(numeric + " is not an integer")
			panic(numeric + " is not an integer")
		}
		intMap[numeric] = value
	}

	return &Config{
		port:        intMap["PORT"],
		env:         variablesMap["ENV"],
		databaseURL: variablesMap["DATABASE_URL"],
		logger:      NewLoggerConfig(variablesMap["LOG_LEVEL"], variablesMap["SERVICE_NAME"]),
	}
}
