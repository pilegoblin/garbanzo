package middleware

import (
	"net/http"

	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
	"github.com/pilegoblin/garbanzo/internal/session"
)

func AuthIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authID, err := session.GetAuthID(r)
		if err != nil {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		ctx := gbzocontext.SetAuthID(r.Context(), authID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
