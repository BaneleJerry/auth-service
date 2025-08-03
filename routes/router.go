package routes

import (
	"github.com/BaneleJerry/auth-service/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(handlers *app.Handlers) *chi.Mux {

	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Define your routes here
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", handlers.AuthHandler.RegisterUser)
	})

	return router
}
