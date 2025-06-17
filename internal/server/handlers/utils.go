package handlers

import (
	"net/http"

	"github.com/pilegoblin/garbanzo/db/sqlc"
	"github.com/pilegoblin/garbanzo/internal/pagecache"
)

type HandlerEnv struct {
	query *sqlc.Queries
	pc    *pagecache.PageCache
}

func NewHandlerEnv(queries *sqlc.Queries) *HandlerEnv {
	return &HandlerEnv{
		query: queries,
		pc:    pagecache.NewPageCache(),
	}
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
