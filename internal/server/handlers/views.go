package handlers

import (
	"html/template"
	"net/http"
)

func IndexViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func CreateUserViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_user.html", nil)
}

func UserViewHandler(w http.ResponseWriter, r *http.Request) {

}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/" + name))
	tmpl.Execute(w, data)
}
