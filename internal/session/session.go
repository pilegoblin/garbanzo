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

func GetAuthID(r *http.Request) (string, error) {
	session, err := GetSession(r)
	if err != nil {
		return "", err
	}
	authID, ok := session.Values["authID"].(string)
	if !ok {
		return "", errors.New("authID not found in session")
	}
	return authID, nil
}

func SetUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Values["userID"] = userID
	return session.Save(r, w)
}

func GetUserID(r *http.Request) (int64, error) {
	session, err := GetSession(r)
	if err != nil {
		return 0, err
	}
	userID, ok := session.Values["userID"].(int64)
	if !ok {
		return 0, errors.New("userID not found in session")
	}
	return userID, nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
