package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Auth        AuthConfig
	Database    DatabaseConfig
	Server      ServerConfig
	Environment string
}

type AuthConfig struct {
	SessionSecret      string
	GoogleClientID     string
	GoogleClientSecret string
}

type DatabaseConfig struct {
	DatabaseURL string
}

type ServerConfig struct {
	Port string
}

func New() *Config {
	return &Config{
		Auth: AuthConfig{
			SessionSecret:      os.Getenv("SESSION_SECRET"),
			GoogleClientID:     os.Getenv("GOOGLE_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_SECRET"),
		},
		Database: DatabaseConfig{
			DatabaseURL: os.Getenv("DATABASE_URL"),
		},
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Environment: os.Getenv("ENVIRONMENT"),
	}
}
