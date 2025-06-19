package pagecache

import (
	"bytes"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	PagesDir     = "templates/pages"
	FragmentsDir = "templates/fragments"
)

type PageCache struct {
	cache map[string]*template.Template
}

func firstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0]))
}

func NewPageCache() *PageCache {
	functions := template.FuncMap{
		"firstLetter": firstLetter,
	}

	pc := &PageCache{
		cache: make(map[string]*template.Template),
	}

	// First add pages to the cache
	// Pages are meant to be rendered with the base template
	pages, err := fs.Glob(os.DirFS(PagesDir), "*.html")
	if err != nil {
		slog.Error("failed to get pages", "error", err)
		os.Exit(1)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		t, err := template.New(name).Funcs(functions).ParseFiles("templates/base.html")
		if err != nil {
			slog.Error("failed to parse base template", "error", err)
			os.Exit(1)
		}

		t, err = t.ParseGlob(filepath.Join(FragmentsDir, "*.html"))
		if err != nil {
			slog.Error("failed to parse fragments", "error", err)
			os.Exit(1)
		}

		t, err = t.ParseFiles(filepath.Join(PagesDir, page))
		if err != nil {
			slog.Error("failed to parse page", "error", err)
			os.Exit(1)
		}

		pc.cache[name] = t
	}

	// Then add fragments to the cache
	// Fragments are meant to be rendered by themselves or in a page
	fragments, err := fs.Glob(os.DirFS(FragmentsDir), "*.html")
	if err != nil {
		slog.Error("failed to get fragments", "error", err)
		os.Exit(1)
	}

	for _, fragment := range fragments {
		name := filepath.Base(fragment)
		pc.cache[name] = template.Must(
			template.New(name).Funcs(functions).ParseFiles(
				filepath.Join(FragmentsDir, fragment),
			),
		)
	}

	return pc
}

func (pc *PageCache) Render(w http.ResponseWriter, name string, data any) {
	t, ok := pc.cache[name]
	if !ok {
		slog.Error("template not found", "name", name)
		os.Exit(1)
	}

	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("failed to render template", "name", name, "error", err)
		os.Exit(1)
	}
}

func (pc *PageCache) RenderFragment(w http.ResponseWriter, name string, data any) {
	t, ok := pc.cache[name]
	if !ok {
		slog.Error("template not found", "name", name)
		os.Exit(1)
	}

	if err := t.ExecuteTemplate(w, name, data); err != nil {
		slog.Error("failed to render template", "name", name, "error", err)
		os.Exit(1)
	}
}

func (pc *PageCache) FragmentString(name string, data any) string {
	t, ok := pc.cache[name]
	if !ok {
		slog.Error("template not found", "name", name)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		slog.Error("failed to render template", "name", name, "error", err)
		os.Exit(1)
	}

	return buf.String()
}
