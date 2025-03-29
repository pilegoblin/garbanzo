package handlers

import (
	"net/http"

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
