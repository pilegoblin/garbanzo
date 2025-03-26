// This application will be the RESTful endpoint responsible for providing and verifying auth to users

package main

import (
	"log/slog"
	"os"

	"github.com/pilegoblin/garbanzo/internal/auth"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/database"
	"github.com/pilegoblin/garbanzo/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	config := config.New()
	auth.SetupAuth(config)
	db := database.New(&config.Database)
	db.Migrate()
	server := server.New(config)
	server.Run()
}
