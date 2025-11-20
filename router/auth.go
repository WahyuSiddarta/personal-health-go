package router

import (
	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// setupPublicRoutes configures public routes (no authentication required)
func (r *Router) setupAuthRoutes(apiGroup *echo.Group) {

	// Initialize auth handlers
	userRepo := models.NewUserRepository()
	authHandlers := api.NewAuthHandlers(userRepo)

	// Authentication routes (no auth required)
	authGroup := apiGroup.Group("/auth")

	// Login endpoint - accessible at /api/public/auth/login
	authGroup.POST("/login", authHandlers.Login, validator.ValidateRequest(&validator.LoginRequest{}))

	// Register endpoint - accessible at /api/public/auth/register
	authGroup.POST("/register", authHandlers.Register, validator.ValidateRequest(&validator.RegisterRequest{}))
}
