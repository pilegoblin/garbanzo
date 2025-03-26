package handlers

import (
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
