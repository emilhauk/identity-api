package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	JwtSigningSecret []byte
	MongoDBUrl       string
	Host             string
	LogLevel         logrus.Level
}

func NewConfig() *Config {
	jwtPanicMessage := "JWT_TOKEN is required to start. Please supply a passphrase."
	config := &Config{
		JwtSigningSecret: []byte(getEnvOrPanic("JWT_TOKEN", jwtPanicMessage)),
		MongoDBUrl:       getEnvWithFallback("MONGO_DB_URL", "mongodb://localhost:27017"),
		Host:             getEnvWithFallback("HOST", ":9002"),
		LogLevel:         getLogLevel(),
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

func getEnvOrPanic(key string, message string) string {
	value := os.Getenv("JWT_TOKEN")
	if len(value) == 0 {
		logrus.Panic(message)
	}
	return value
}

func getLogLevel() logrus.Level {
	if strings.ToLower(getEnvWithFallback("DEBUG", "false")) != "false" {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}
