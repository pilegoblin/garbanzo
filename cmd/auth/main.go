// This application will be the RESTful endpoint responsible for providing and verifying auth to users

package main

import (
	"github.com/joho/godotenv"
	"github.com/pilegoblin/garbanzo/internal/auth"
	"github.com/pilegoblin/garbanzo/internal/server"
)

func init() {

}

func main() {
	godotenv.Load(".env")
	auth.SetupAuth()

	server := server.New()
	server.Run()
}
