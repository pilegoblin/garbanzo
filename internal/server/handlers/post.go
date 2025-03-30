package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pilegoblin/garbanzo/internal/session"
)

func (h *HandlerEnv) CreatePost(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	beanID, err := strconv.Atoi(chi.URLParam(r, "beanID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserID(r)
	if err != nil {
		redirect(w, "/login")
		return
	}

	post, err := h.db.CreatePost(r.Context(), userID, content, beanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, post, "post.html")
}
