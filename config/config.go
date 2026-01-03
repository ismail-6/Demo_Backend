package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	Port         string
	Environment  string
	DatabaseType string // "postgres" or "sqlite"
}

func LoadConfig() *Config {
	// Default to development settings
	config := &Config{
		DatabaseURL:  "learning_app.db",
		Port:         "8080",
		Environment:  "development",
		DatabaseType: "sqlite",
	}

	// Check for production database URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
		config.DatabaseType = "postgres"
		config.Environment = "production"
	}

	// Check for custom port
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	return config
}
