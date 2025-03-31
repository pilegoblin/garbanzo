package handlers

import (
	"net/http"
	"text/template"

	"github.com/pilegoblin/garbanzo/internal/database"
)

type HandlerEnv struct {
	db *database.Database
}

func NewHandlerEnv(db *database.Database) *HandlerEnv {
	return &HandlerEnv{
		db: db,
	}
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func renderTemplate(w http.ResponseWriter, data any, files ...string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	files = append([]string{"base.html"}, files...)
	for i, file := range files {
		files[i] = "templates/" + file
	}
	tmpl := template.Must(template.ParseFiles(files...))
	tmpl.Execute(w, data)
}

func renderTemplateRaw(w http.ResponseWriter, data any, files ...string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for i, file := range files {
		files[i] = "templates/" + file
	}
	tmpl := template.Must(template.ParseFiles(files...))
	tmpl.Execute(w, data)
}
