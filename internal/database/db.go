package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pilegoblin/garbanzo/db/sqlc"
	"github.com/pilegoblin/garbanzo/internal/config"
)

func New(config *config.DatabaseConfig) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	return db
}

func NewQueries(db *pgxpool.Pool) *sqlc.Queries {
	return sqlc.New(db)
}
