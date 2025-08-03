package routes

import (
	"net/http"

	"github.com/BaneleJerry/auth-service/app"
	customMiddleware "github.com/BaneleJerry/auth-service/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(handlers *app.Handlers) *chi.Mux {
	router := chi.NewRouter()

	// Public middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Public routes - no auth required
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.AuthHandler.RegisterUser)
		r.Post("/login", handlers.AuthHandler.LoginUser)
		r.Post("/refresh-token", handlers.AuthHandler.RefreshTokenHandler)

	})

	// Protected routes - require JWT auth
	router.Route("/api", func(r chi.Router) {
		r.Use(customMiddleware.JWTAuthMiddleware)

		r.Delete("/logout", handlers.AuthHandler.LogoutHandler)
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

	})

	return router
}
