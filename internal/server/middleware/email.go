package middleware

import (
	"net/http"

	"github.com/pilegoblin/garbanzo/internal/session"
)

func EmailMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetEmail(r)
		if err != nil {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}
