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

func (c *Config) Port() int64 {
	return c.port
}

func (c *Config) ENV() string {
	return c.env
}

func (c *Config) DatabaseURL() string {
	return c.databaseURL
}

func (c *Config) Logger() LoggerConfig {
	return c.logger
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
