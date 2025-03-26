package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pilegoblin/garbanzo/ent"
	"github.com/pilegoblin/garbanzo/ent/user"
	"github.com/pilegoblin/garbanzo/internal/config"
)

type db struct {
	client *ent.Client
}

func New(config *config.DatabaseConfig) *db {
	database, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	drv := entsql.OpenDB(dialect.Postgres, database)
	client := ent.NewClient(ent.Driver(drv))
	return &db{
		client: client,
	}
}

func (d *db) Migrate() {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		slog.Error("failed to create schema", "error", err)
		os.Exit(1)
	}
}

func (d *db) GetUser(ctx context.Context, email string) (*ent.User, error) {
	user, err := d.client.User.Query().Where(user.Email(email)).First(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *db) CreateUser(ctx context.Context, user *ent.User) (*ent.User, error) {
	user, err := d.client.User.Create().SetEmail(user.Email).Save(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
