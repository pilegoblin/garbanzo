package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/pilegoblin/garbanzo/ent"
	"github.com/pilegoblin/garbanzo/ent/user"
	"github.com/pilegoblin/garbanzo/internal/config"

	_ "github.com/lib/pq"
)

type db struct {
	client *ent.Client
}

func New(config *config.DatabaseConfig) *db {
	client, err := ent.Open("postgres", config.DatabaseURL)
	if err != nil {
		slog.Error("failed to open postgres client", "error", err)
		os.Exit(1)
	}
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
