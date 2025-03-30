package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pilegoblin/garbanzo/ent"
	"github.com/pilegoblin/garbanzo/ent/pod"
	"github.com/pilegoblin/garbanzo/ent/user"
	"github.com/pilegoblin/garbanzo/internal/config"
)

type Database struct {
	client *ent.Client
}

func New(config *config.DatabaseConfig) *Database {
	database, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	drv := entsql.OpenDB(dialect.Postgres, database)
	client := ent.NewClient(ent.Driver(drv))
	return &Database{
		client: client,
	}
}

func (d *Database) Migrate() {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		slog.Error("failed to create schema", "error", err)
		os.Exit(1)
	}
}

func (d *Database) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	user, err := d.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Database) GetUserByAuthID(ctx context.Context, authID string) (*ent.User, error) {
	user, err := d.client.User.Query().Where(user.AuthID(authID)).First(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Database) CreateUser(ctx context.Context, authID, username, email string) (*ent.User, error) {
	user, err := d.client.User.Create().SetAuthID(authID).SetUsername(username).SetEmail(email).Save(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Database) GetPosts(ctx context.Context) ([]*ent.Post, error) {
	posts, err := d.client.Post.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (d *Database) CreatePost(ctx context.Context, userID int, content string, beanID int) (*ent.Post, error) {
	post, err := d.client.Post.Create().SetContent(content).SetUserID(userID).SetBeanID(beanID).Save(ctx)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// PODS

// check if user is in pod
func (d *Database) IsUserInPod(ctx context.Context, userID int, podID int) (bool, error) {
	userInPod, err := d.client.User.Query().
		Where(user.ID(userID)).
		QueryJoinedPods().
		Where(pod.ID(podID)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return userInPod, nil
}
