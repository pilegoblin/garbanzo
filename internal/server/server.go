package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/database"
	"github.com/pilegoblin/garbanzo/internal/server/handlers"
	"github.com/pilegoblin/garbanzo/internal/server/middleware"
)

type Server struct {
	Router  *chi.Mux
	DB      *pgxpool.Pool
	port    string
	handler *handlers.HandlerEnv
}

func New(config *config.Config) *Server {
	db := database.New(&config.Database)
	q := database.NewQueries(db)
	return &Server{
		Router:  chi.NewRouter(),
		DB:      db,
		port:    config.Server.Port,
		handler: handlers.New(config, q),
	}
}

func (s *Server) Run() {

	s.Router.Group(func(r chi.Router) {
		// middleware
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		r.Use(middleware.Logger)
		r.Use(chimiddleware.Recoverer)

		// public
		r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
			fs := http.FileServer(http.Dir("./public/"))
			w.Header().Add("Cache-Control", "no-cache")
			http.StripPrefix("/public/", fs).ServeHTTP(w, r)
		})

		// auth routes
		r.Get("/login", s.handler.LoginViewHandler)
		r.Get("/auth", handlers.AuthHandler)
		r.Get("/auth/callback", handlers.CallbackHandler)
		r.Get("/logout", handlers.LogoutHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthIDMiddleware)
			r.Get("/", s.handler.IndexViewHandler)
			r.Get("/user/new", s.handler.NewUserViewHandler)
			r.Post("/user/new", s.handler.NewUserHandler)
			r.Post("/pod/join", s.handler.JoinPodHandler)
			r.Get("/pod/{podID}", s.handler.PodViewHandler)
			r.Get("/ws/{podID}/{beanID}", s.handler.WebsocketHandler)
			r.Get("/messages/edit/{messageID}", s.handler.EditMessageViewHandler)
			r.Post("/messages/edit/{messageID}", s.handler.EditMessageHandler)
		})
	})

	slog.Info("Starting server on port " + s.port)

	if err := http.ListenAndServe(":"+s.port, s.Router); err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
