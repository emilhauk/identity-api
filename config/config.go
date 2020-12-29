package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	MongoDBUrl   string
	Host         string
	LogLevel     logrus.Level
	KeyStorePath string
	DefaultKeyId string
}

func NewConfig() *Config {
	config := &Config{
		MongoDBUrl:   getEnvWithFallback("MONGO_DB_URL", "mongodb://localhost:27017"),
		Host:         getEnvWithFallback("HOST", ":9002"),
		LogLevel:     getLogLevel(),
		KeyStorePath: getEnvOrFatal("KEY_STORE", "KEY_STORE path is required to start"),
		DefaultKeyId: getEnvOrFatal("DEFAULT_KEY", "DEFAULT_KEY id is required to start"),
	}
	return config
}

func getEnvWithFallback(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getEnvOrFatal(key string, message string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logrus.Fatalln(message)
	}
	return value
}

func getLogLevel() logrus.Level {
	if strings.ToLower(getEnvWithFallback("DEBUG", "false")) != "false" {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}
