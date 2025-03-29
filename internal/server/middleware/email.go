package middleware

import (
	"net/http"

	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
	"github.com/pilegoblin/garbanzo/internal/session"
)

func EmailMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, err := session.GetEmail(r)
		if err != nil {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		ctx := gbzocontext.SetEmail(r.Context(), email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
