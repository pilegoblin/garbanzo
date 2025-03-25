package middleware

import (
	"log/slog"
	"net/http"

	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			slog.Info("Request Info",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
				slog.String("host", r.Host),
				slog.String("request_ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Float64("duration_ms", float64(time.Since(t).Microseconds())/1000),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}
