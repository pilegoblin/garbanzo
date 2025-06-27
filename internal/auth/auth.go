package auth

import (
	"log/slog"
	"sync"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/session"
)

var authOnce sync.Once

func SetupAuth(config *config.Config) {
	authOnce.Do(func() {
		setup(config)
	})
}

// Sets up the goth package to do what it needs to do
func setup(config *config.Config) {
	googleKey := config.Auth.GoogleClientID
	googleSecret := config.Auth.GoogleClientSecret
	port := config.Server.Port
	host := config.Server.Host
	environment := config.Environment

	var callbackURL string
	if environment == "prod" {
		callbackURL = "https://" + host + "/auth/callback?provider=google"
	} else {
		callbackURL = "http://localhost:" + port + "/auth/callback?provider=google"
	}

	store, err := session.GetSessionStore()
	if err != nil {
		panic(err)
	}

	// Add validation
	if store == nil {
		panic("session store is nil")
	}

	slog.Info("Setting up Goth with session store", "callbackURL", callbackURL)
	gothic.Store = store
	goth.UseProviders(
		google.New(
			googleKey, googleSecret, callbackURL,
		),
	)
}
