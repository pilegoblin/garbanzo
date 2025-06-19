package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pilegoblin/garbanzo/db/sqlc"
	"github.com/pilegoblin/garbanzo/internal/pagecache"
	"github.com/pilegoblin/garbanzo/internal/switchboard"
)

type HandlerEnv struct {
	query       *sqlc.Queries
	pc          *pagecache.PageCache
	upgrader    *websocket.Upgrader
	switchboard *switchboard.Switchboard
}

type FullMessage struct {
	AuthorAvatarURL string    `json:"author_avatar_url"`
	AuthorID        int       `json:"author_id"`
	AuthorUsername  string    `json:"author_username"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	ID              int       `json:"id"`
}

type BeanWithMessages struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	PodID    int64         `json:"pod_id"`
	PodName  string        `json:"pod_name"`
	Messages []FullMessage `json:"messages"`
}

func NewHandlerEnv(queries *sqlc.Queries) *HandlerEnv {
	return &HandlerEnv{
		query: queries,
		pc:    pagecache.NewPageCache(),
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// allowedOrigins := []string{
				// 	"http://localhost:8080",
				// }
				// return slices.Contains(allowedOrigins, r.Header.Get("Origin"))
				return true
			},
		},
		switchboard: switchboard.New(),
	}
}

func redirect(w http.ResponseWriter, path string) {
	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
