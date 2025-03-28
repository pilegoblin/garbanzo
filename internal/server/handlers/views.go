package handlers

import (
	"html/template"
	"net/http"

	gbzocontext "github.com/pilegoblin/garbanzo/internal/context"
)

// GET /
func IndexViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

// GET /user/create
func CreateUserViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_user.html", nil)
}

// GET /user
func (h *HandlerEnv) UserViewHandler(w http.ResponseWriter, r *http.Request) {
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
	renderTemplate(w, "user.html", user)
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/"+name))
	tmpl.Execute(w, data)
}
