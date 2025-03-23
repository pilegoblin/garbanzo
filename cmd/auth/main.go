// This application will be the RESTful endpoint responsible for providing and verifying auth to users

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"github.com/pilegoblin/garbanzo/internal/util"
)

func init() {
	godotenv.Load(".env")

	store = sessions.NewCookieStore([]byte(util.GetEnvVarOrPanic("SESSION_SECRET")))
}

func main() {
	SetupAuth()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		sess := GetSession(r)
		email := sess.Values["email"]
		if email == nil {
			w.Header().Set("Location", "/auth?provider=google")
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
		render.JSON(w, r, sess.Values["email"])

	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		sess := GetSession(r)
		sess.Values["email"] = user.Email
		sess.Save(r, w)
		render.JSON(w, r, "Hello "+user.Email)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})
	log.Println("we up")

	port := util.GetEnvVarOrDefault("PORT", "8080")
	http.ListenAndServe(":"+port, r)
}
