package handlers

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"time"

	"github.com/pilegoblin/garbanzo/db/sqlc"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/pagecache"
	"github.com/pilegoblin/garbanzo/internal/switchboard"
)

type HandlerEnv struct {
	query       *sqlc.Queries
	pc          *pagecache.PageCache
	switchboard *switchboard.Switchboard
	origins     []string
}

type MessageData struct {
	AuthorAvatarURL string        `json:"author_avatar_url"`
	AuthorID        int64         `json:"author_id"`
	AuthorUsername  string        `json:"author_username"`
	AuthorUserColor string        `json:"author_user_color"`
	Content         string        `json:"content"`
	CreatedAt       time.Time     `json:"created_at"`
	ID              string        `json:"id"`
	Action          MessageAction `json:"action"`
	Editable        bool          `json:"editable"`
}

type MessageAction string

const (
	MessageActionNew    MessageAction = "new"
	MessageActionEdit   MessageAction = "edit"
	MessageActionDelete MessageAction = "delete"
)

type BeanWithMessages struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	PodID    int64         `json:"pod_id"`
	PodName  string        `json:"pod_name"`
	Messages []MessageData `json:"messages"`
}

func New(config *config.Config, queries *sqlc.Queries) *HandlerEnv {
	return &HandlerEnv{
		query:       queries,
		pc:          pagecache.NewPageCache(),
		switchboard: switchboard.New(),
		origins: []string{
			"http://localhost:8080",
			"https://" + config.Server.Host,
			"https://www." + config.Server.Host,
		},
	}
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func createUserColor(username string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(username))
	source := rand.NewSource(int64(hasher.Sum64()))
	r := rand.New(source)
	num := r.Intn(16777215)
	return fmt.Sprintf("%06x", num)
}
