package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pilegoblin/garbanzo/ent"
	"github.com/pilegoblin/garbanzo/ent/bean"
	"github.com/pilegoblin/garbanzo/ent/pod"
	"github.com/pilegoblin/garbanzo/ent/post"
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
	user, err := d.client.User.Query().Where(user.AuthID(authID)).WithJoinedPods().First(ctx)
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

func (d *Database) CreatePost(ctx context.Context, userID int, beanID int, content string) (*ent.Post, error) {
	// get pod from beanid
	p, err := d.client.Bean.Query().Where(bean.ID(beanID)).QueryPod().First(ctx)
	if err != nil {
		return nil, err
	}

	userInPod, err := d.CheckUserInPod(ctx, userID, p.ID)
	if err != nil {
		return nil, err
	}

	if !userInPod {
		return nil, errors.New("user is not in pod")
	}

	newPost, err := d.client.Post.Create().SetContent(content).SetUserID(userID).SetBeanID(beanID).Save(ctx)
	if err != nil {
		return nil, err
	}
	fullPost, err := d.client.Post.Query().Where(post.ID(newPost.ID)).WithUser().First(ctx)
	if err != nil {
		return nil, err
	}
	return fullPost, nil
}

func (d *Database) CheckUserInPod(ctx context.Context, userID, podID int) (bool, error) {
	userInPod, err := d.client.User.Query().Where(user.ID(userID)).QueryJoinedPods().Where(pod.ID(podID)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return userInPod, nil
}

func (d *Database) CreatePod(ctx context.Context, userID int, name string) (*ent.Pod, error) {
	inviteCode := uuid.New().String()
	pod, err := d.client.Pod.Create().SetPodName(name).SetOwnerID(userID).SetInviteCode(inviteCode).Save(ctx)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (d *Database) CreateBean(ctx context.Context, podID int, name string) (*ent.Bean, error) {
	bean, err := d.client.Bean.Create().SetName(name).SetPodID(podID).Save(ctx)
	if err != nil {
		return nil, err
	}
	return bean, nil
}

func (d *Database) JoinPod(ctx context.Context, userID int, inviteCode string) (*ent.Pod, error) {
	newPod, err := d.client.Pod.Query().Where(pod.InviteCode(inviteCode)).First(ctx)
	if err != nil {
		return nil, err
	}
	_, err = d.client.User.UpdateOneID(userID).AddJoinedPods(newPod).Save(ctx)
	if err != nil {
		return nil, err
	}

	return newPod, nil
}

func (d *Database) GetBeans(ctx context.Context, userID, podID int) ([]*ent.Bean, error) {
	userInPod, err := d.CheckUserInPod(ctx, userID, podID)
	if err != nil {
		return nil, err
	}

	if !userInPod {
		return nil, errors.New("user is not in pod")
	}
	// get all beans and their posts with their edges
	beans, err := d.client.Bean.Query().
		Where(bean.HasPodWith(pod.ID(podID))).
		WithPosts(func(q *ent.PostQuery) {
			q.WithUser()
		}).
		WithPod().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return beans, nil
}
