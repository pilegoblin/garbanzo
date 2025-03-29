package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
)

func (h *HandlerEnv) CreatePost(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	email, ok := gbzocontext.GetEmail(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	beanID := chi.URLParam(r, "beanID")
	if beanID == "" {
		http.Error(w, "Bean ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUser(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.db.CreatePost(r.Context(), user, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, post, "post.html")
}
