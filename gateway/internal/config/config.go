package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	port                     int64
	env                      string
	logger                   LoggerConfig
	calculationDomain        string
	calculationTLSServerName string // optional: for SNI / cert verification in production
	calculationTLSCertPath   string // optional: for cert verification in production
}

func (c *Config) Port() int64 {
	return c.port
}

func (c *Config) ENV() string {
	return c.env
}

func (c *Config) CalculationDomain() string {
	return c.calculationDomain
}

func (c *Config) CalculationTLSServerName() string {
	return c.calculationTLSServerName
}

func (c *Config) CalculationTLSCertPath() string {
	return c.calculationTLSCertPath
}

func (c *Config) Logger() LoggerConfig {
	return c.logger
}

var variables = [5]string{
	"PORT",
	"ENV",
	"LOG_LEVEL",
	"SERVICE_NAME",
	"CALCULATION_DOMAIN",
}

// Optional: set for production TLS (e.g. when backend uses different hostname for cert)
var optionalVariables = [2]string{
	"CALCULATION_SERVICE_TLS_SERVER_NAME",
	"CALCULATION_SERVICE_TLS_CERT_PATH",
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

	optionalMap := make(map[string]string)
	for _, variable := range optionalVariables {
		optionalMap[variable] = os.Getenv(variable)
	}

	return &Config{
		port:                     intMap["PORT"],
		env:                      variablesMap["ENV"],
		calculationDomain:        variablesMap["CALCULATION_DOMAIN"],
		calculationTLSServerName: optionalMap["CALCULATION_SERVICE_TLS_SERVER_NAME"],
		calculationTLSCertPath:   optionalMap["CALCULATION_SERVICE_TLS_CERT_PATH"],
		logger:                   NewLoggerConfig(variablesMap["LOG_LEVEL"], variablesMap["SERVICE_NAME"]),
	}
}
