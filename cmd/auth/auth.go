package main

import (
	"github.com/banana-slugg/garbanzo/pkg/util"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// Sets up the goth package to do what it needs to do
func SetupAuth() {
	sessionKey := util.GetEnvVarOrPanic("SESSION_SECRET")
	googleKey := util.GetEnvVarOrPanic("GOOGLE_KEY")
	googleSecret := util.GetEnvVarOrPanic("GOOGLE_SECRET")

	port := util.GetEnvVarOrDefault("PORT", "8080")
	environment := util.GetEnvVarOrDefault("ENVIRONMENT", "dev")

	isProd := false
	if environment == "prod" {
		isProd = true
	}
	maxAge := 86400 * 30 // 30 days

	store := sessions.NewCookieStore([]byte(sessionKey))
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
