package handlers

import (
	"net/http"

	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
	"github.com/pilegoblin/garbanzo/internal/session"
)

// GET /
func (h *HandlerEnv) IndexViewHandler(w http.ResponseWriter, r *http.Request) {
	authID, ok := gbzocontext.GetAuthID(r.Context())
	if !ok {
		redirect(w, "/login")
		return
	}

	user, err := h.db.GetUserByAuthID(r.Context(), authID)
	if err != nil {
		redirect(w, "/user/new")
		return
	}

	_, err = session.GetUserID(r)
	if err != nil {
		session.SetUserID(w, r, user.ID)
	}

	renderTemplate(w, user, "index.html")
}

// GET /login
func LoginViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, nil, "login.html")
}

// GET /user/new
func NewUserViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, nil, "new_user.html")
}

// GET /{podID}/{beanID}
func BeanViewHandler(w http.ResponseWriter, r *http.Request) {
	// podID := chi.URLParam(r, "podID")
	// beanID := chi.URLParam(r, "beanID")

	renderTemplate(w, nil, "bean.html")
}
