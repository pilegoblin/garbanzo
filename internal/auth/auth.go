package auth

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pilegoblin/garbanzo/internal/util"
)

var store *sessions.CookieStore
var authOnce sync.Once

func SetupAuth() {
	authOnce.Do(setup)
}

// Sets up the goth package to do what it needs to do
func setup() {
	sessionSecret := util.GetEnvVarOrPanic("SESSION_SECRET")
	googleKey := util.GetEnvVarOrPanic("GOOGLE_KEY")
	googleSecret := util.GetEnvVarOrPanic("GOOGLE_SECRET")

	port := util.GetEnvVarOrDefault("PORT", "8080")
	environment := util.GetEnvVarOrDefault("ENVIRONMENT", "dev")

	store = sessions.NewCookieStore([]byte(sessionSecret))

	isProd := false
	if environment == "prod" {
		isProd = true
	}
	maxAge := 86400 * 30 // 30 days

	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

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
	session, err := store.Get(r, "gbzo-session")
	if err != nil {
		return nil, err
	}
	return session, nil
}
