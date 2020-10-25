package config

import (
	"log"
	"os"
)

type Config struct {
	JwtSigningSecret []byte
	MongoDBUrl       string
	Host             string
}

func NewConfig() *Config {
	jwtPanicMessage := "JWT_TOKEN is required to start. Please supply a passphrase."
	config := &Config{
		JwtSigningSecret: []byte(getEnvOrPanic("JWT_TOKEN", jwtPanicMessage)),
		MongoDBUrl:       getEnvWithFallback("MONGO_DB_URL", "mongodb://localhost:27017"),
		Host:             getEnvWithFallback("HOST", ":9002"),
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
		log.Panic(message)
	}
	return value
}
