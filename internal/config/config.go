package config

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Auth        AuthConfig
	Database    DatabaseConfig
	Server      ServerConfig
	Environment string `env:"ENVIRONMENT"`
}

type AuthConfig struct {
	SessionSecret      string `env:"SESSION_SECRET"`
	GoogleClientID     string `env:"GOOGLE_ID"`
	GoogleClientSecret string `env:"GOOGLE_SECRET"`
}

type DatabaseConfig struct {
	DatabaseURL string `env:"DATABASE_URL"`
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

func New() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		os.Exit(1)
	}
	return &cfg
}
