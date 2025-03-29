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
	// Db, config can be added here
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

		// index route

		// auth routes
		r.Get("/login", handlers.LoginViewHandler)
		r.Get("/auth", handlers.AuthHandler)
		r.Get("/auth/callback", handlers.CallbackHandler)
		r.Get("/logout", handlers.LogoutHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.EmailMiddleware)
			r.Get("/", s.handler.IndexViewHandler)
			// user routes
			r.Get("/user/create", handlers.CreateUserViewHandler)
			r.Post("/user/create", s.handler.CreateUserHandler)
		})

	})

	slog.Info("Starting server on port " + s.port)
	http.ListenAndServe(":"+s.port, s.Router)
}
