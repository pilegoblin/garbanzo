// This application will be the RESTful endpoint responsible for providing and verifying auth to users

package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/pilegoblin/garbanzo/internal/auth"
	"github.com/pilegoblin/garbanzo/internal/server"
)

func main() {
	godotenv.Load(".env")
	auth.SetupAuth()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	server := server.New()
	server.Run()
}
