package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "gbzo-session")
	return session
}
