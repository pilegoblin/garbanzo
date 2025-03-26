package session

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var storeOnce sync.Once

func GetSession(r *http.Request) (*sessions.Session, error) {
	storeOnce.Do(func() {
		store = sessions.NewCookieStore([]byte("secret"))
	})

	session, err := store.Get(r, "gbzo-session")
	if err != nil {
		return nil, err
	}
	return session, nil
}

func GetEmail(r *http.Request) (string, error) {
	session, err := GetSession(r)
	if err != nil {
		return "", err
	}
	email, ok := session.Values["email"].(string)
	if !ok {
		return "", errors.New("email not found in session")
	}
	return email, nil
}
