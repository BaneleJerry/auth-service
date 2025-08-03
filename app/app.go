package app

import (
	controllers "github.com/BaneleJerry/auth-service/handlers"
	"github.com/BaneleJerry/auth-service/services"
	"gorm.io/gorm"
)

// App holds your controllers
type App struct {
	Controllers Handlers
}

// Handlers groups your HTTP handlers/controllers
type Handlers struct {
	AuthHandler *controllers.AuthHandler
}

// InitApp wires everything using shared dependencies (like DB)
func InitApp(db *gorm.DB) *App {
	// Initialize services
	authService := services.NewAuthService(db)

	// Initialize handlers
	authHandler := controllers.NewAuthHandler(authService)

	// Group controllers
	ctrls := Handlers{
		AuthHandler: authHandler,
	}

	return &App{
		Controllers: ctrls,
	}
}
