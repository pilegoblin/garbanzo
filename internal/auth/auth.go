package auth

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pilegoblin/garbanzo/internal/config"
)

const (
	SessionName = "gbzo-session"
)

var store *sessions.CookieStore
var authOnce sync.Once

func SetupAuth(config *config.Config) {
	authOnce.Do(func() {
		setup(config)
	})
}

// Sets up the goth package to do what it needs to do
func setup(config *config.Config) {
	sessionSecret := config.Auth.SessionSecret
	googleKey := config.Auth.GoogleClientID
	googleSecret := config.Auth.GoogleClientSecret

	port := config.Server.Port
	environment := config.Environment

	store = sessions.NewCookieStore([]byte(sessionSecret))

	maxAge := 86400 * 30 // 30 days

	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = environment == "prod"

	gothic.Store = store
	goth.UseProviders(
		google.New(
			googleKey, googleSecret, "http://localhost:"+port+"/auth/callback?provider=google",
		),
	)
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	if store == nil {
		return nil, errors.New("store not initialized")
	}
	session, err := store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}
	return session, nil
}
