package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pilegoblin/garbanzo/internal/config"
	"github.com/pilegoblin/garbanzo/internal/server/handlers"
	"github.com/pilegoblin/garbanzo/internal/server/middleware"
)

type Server struct {
	Router *chi.Mux
	port   string
	// Db, config can be added here
}

func New(config *config.ServerConfig) *Server {
	return &Server{
		Router: chi.NewRouter(),
		port:   config.Port,
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

		// routes
		r.Get("/", handlers.MainHandler)
		r.Get("/auth", handlers.AuthHandler)
		r.Get("/auth/callback", handlers.CallbackHandler)
		r.Get("/logout", handlers.LogoutHandler)
	})

	slog.Info("Starting server on port " + s.port)
	http.ListenAndServe(":"+s.port, s.Router)
}
