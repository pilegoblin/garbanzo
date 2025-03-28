package middleware

import (
	"context"
	"net/http"

	"github.com/pilegoblin/garbanzo/internal/session"
)

type ContextKey string

const emailContextKey ContextKey = "email"

func EmailMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, err := session.GetEmail(r)
		if err != nil {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
		ctx := context.WithValue(r.Context(), emailContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(emailContextKey).(string)
	return email, ok
}
