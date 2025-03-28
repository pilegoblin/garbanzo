package handlers

import (
	"log/slog"
	"net/http"

	"github.com/pilegoblin/garbanzo/internal/session"
)

// /user/create
func (h *HandlerEnv) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		redirect(w, "/user")
		return
	}
	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	email, err := session.GetEmail(r)
	if err != nil {
		slog.Error("failed to get email", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.db.CreateUser(r.Context(), email, username)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Location", "/user")
}
