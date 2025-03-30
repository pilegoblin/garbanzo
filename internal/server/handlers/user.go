package handlers

import (
	"log/slog"
	"net/http"

	"github.com/pilegoblin/garbanzo/internal/session"
)

// POST /user/new
func (h *HandlerEnv) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		redirect(w, "/")
		return
	}
	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authID, err := session.GetAuthID(r)
	if err != nil {
		slog.Error("failed to get authID", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email, err := session.GetEmail(r)
	if err != nil {
		slog.Error("failed to get email", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.db.CreateUser(r.Context(), authID, username, email)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.SetUserID(w, r, user.ID)
	w.Header().Set("HX-Location", "/")
}
