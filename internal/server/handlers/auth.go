package handlers

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/pilegoblin/garbanzo/internal/session"
)

// GET /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.Logout(w, r)
	redirect(w, "/login")
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
	sess.Values["authID"] = user.UserID
	sess.Values["email"] = user.Email
	sess.Save(r, w)
	redirect(w, "/")
}
