package handlers

import (
	"net/http"
	"strings"
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

func renderTemplateRaw(w http.ResponseWriter, data any, file string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New("templates")
	t.Funcs(template.FuncMap{
		"firstLetter": firstLetter,
	})
	t, err := t.ParseFiles("templates/" + file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, file, data)
}

func renderTemplate(w http.ResponseWriter, data any, files ...string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	files = append([]string{"base.html"}, files...)
	for i, file := range files {
		files[i] = "templates/" + file
	}
	t := template.New("templates")
	t.Funcs(template.FuncMap{
		"firstLetter": firstLetter,
	})
	tmpl, err := t.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func firstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0]))
}
