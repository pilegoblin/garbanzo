package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/markbates/goth/gothic"
	"github.com/pilegoblin/garbanzo/internal/session"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	email := sess.Values["email"]
	if email == nil {
		w.Header().Set("Location", "/auth?provider=google")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
	render.JSON(w, r, map[string]any{
		"email": email,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

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
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
