package handlers

import (
	"net/http"

	"github.com/pilegoblin/garbanzo/internal/database"
	"github.com/pilegoblin/garbanzo/internal/pagecache"
)

type HandlerEnv struct {
	db *database.Database
	pc *pagecache.PageCache
}

func NewHandlerEnv(db *database.Database) *HandlerEnv {
	return &HandlerEnv{
		db: db,
		pc: pagecache.NewPageCache(),
	}
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
