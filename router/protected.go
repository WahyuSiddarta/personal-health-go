package router

import (
	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(apiGroup *echo.Group) {
	protectedGroup := apiGroup.Group("/protected")
	// protectedGroup.Use(middleware.JWTMiddleware())
	setupUserRoutes(protectedGroup)

}

func setupUserRoutes(group *echo.Group) {
	// Define user-related protected routes here
	usersGroup := group.Group("/users")

	// Initialize auth handlers
	userRepo := models.NewUserRepository()
	userHandlers := api.NewUserHandlers(userRepo)

	usersGroup.GET("/personal-target", userHandlers.GetPersonalTarget)
	usersGroup.PUT("/personal-target/nutrition", userHandlers.UpdatePersonalNutritionTarget, validator.ValidateRequest(&validator.PersonalNutritionTargetRequest{}))
}
