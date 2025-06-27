package main

import (
	"log/slog"
	"os"

	"github.com/pilegoblin/garbanzo/internal/auth"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/server"
	"github.com/pilegoblin/garbanzo/internal/session"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	config := config.New()
	session.SetupSessionStore(config)
	auth.SetupAuth(config)
	server := server.New(config)
	defer server.DB.Close()
	server.Run()
}
