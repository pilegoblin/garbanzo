package handlers

import (
	"net/http"

	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
)

// GET /
func (h *HandlerEnv) IndexViewHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := gbzocontext.GetEmail(r.Context())
	if !ok {
		redirect(w, "/")
		return
	}
	user, err := h.db.GetUser(r.Context(), email)
	if err != nil {
		redirect(w, "/user/create")
		return
	}

	renderTemplate(w, user, "index.html")
}

func LoginViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, nil, "login.html")
}

// GET /user/create
func CreateUserViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, nil, "new_user.html")
}

// GET /{podID}/{beanID}
func BeanViewHandler(w http.ResponseWriter, r *http.Request) {
	// podID := chi.URLParam(r, "podID")
	// beanID := chi.URLParam(r, "beanID")

	renderTemplate(w, nil, "bean.html")
}
