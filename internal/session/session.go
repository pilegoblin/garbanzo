package session

import (
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
