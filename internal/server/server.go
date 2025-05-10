package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/database"
	"github.com/pilegoblin/garbanzo/internal/server/handlers"
	"github.com/pilegoblin/garbanzo/internal/server/middleware"
)

type Server struct {
	Router  *chi.Mux
	port    string
	handler *handlers.HandlerEnv
}

func New(config *config.Config) *Server {
	db := database.New(&config.Database)
	return &Server{
		Router:  chi.NewRouter(),
		port:    config.Server.Port,
		handler: handlers.NewHandlerEnv(db),
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

		// auth routes
		r.Get("/login", handlers.LoginViewHandler)
		r.Get("/auth", handlers.AuthHandler)
		r.Get("/auth/callback", handlers.CallbackHandler)
		r.Get("/logout", handlers.LogoutHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthIDMiddleware)
			r.Get("/", s.handler.IndexViewHandler)
			r.Get("/user/new", handlers.NewUserViewHandler)
			r.Post("/user/new", s.handler.NewUserHandler)
			r.Post("/pod/join", s.handler.JoinPodHandler)
			r.Get("/pod/{podID}", s.handler.PodViewHandler)
			r.Post("/{podID}/{beanID}/post", s.handler.CreatePost)
		})
	})

	slog.Info("Starting server on port " + s.port)
	http.ListenAndServe(":"+s.port, s.Router)
}
