package handlers

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/pilegoblin/garbanzo/internal/session"
)

// GET /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil {
		redirect(w, "/auth?provider=google")
		return
	}
	email := sess.Values["email"]
	if email == nil {
		redirect(w, "/auth?provider=google")
		return
	}
	redirect(w, "/user")
}

// GET /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	redirect(w, "/")
}

// GET /auth
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

// GET /auth/callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sess, err := session.GetSession(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sess.Values["email"] = user.Email
	sess.Save(r, w)
	redirect(w, "/user")
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
