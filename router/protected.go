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
	setupFoodNutritionRoutes(protectedGroup)
	setupBodyMeasurementRoutes(protectedGroup)
}

func setupUserRoutes(group *echo.Group) {
	// Define user-related protected routes here
	usersGroup := group.Group("/users")

	// Initialize auth handlers
	userRepo := models.NewUserRepository()
	userHandlers := api.NewUserHandlers(userRepo)
	usersGroup.GET("/personal-target", userHandlers.GetPersonalTarget)
	usersGroup.PUT("/personal-target/nutrition", userHandlers.UpdatePersonalNutritionTarget, validator.ValidateRequest(&validator.PersonalNutritionTargetRequest{}))
	usersGroup.PUT("/personal-target/body-measurement", userHandlers.UpdatePersonalBodyMeasurementTarget, validator.ValidateRequest(&validator.PersonalBodyMeasurementTargetRequest{}))
}

func setupFoodNutritionRoutes(group *echo.Group) {
	// Define user-related protected routes here
	nutritionGroup := group.Group("/food-tracker")

	// Initialize auth handlers
	nutritionRepo := models.NewNutritionRepository()
	nutritionHandler := api.NewNutritionHandlers(nutritionRepo)

	// Daily nutrition intake routes
	nutritionGroup.GET("/today", nutritionHandler.GetTodaysNutritionIntake)
	nutritionGroup.POST("/today", nutritionHandler.AddNutritionIntake, validator.ValidateRequest(&validator.NutritionRequest{}))
	nutritionGroup.PUT("/today/:food_id", nutritionHandler.UpdateNutritionIntake, validator.ValidateRequest(&validator.NutritionRequest{}))
	nutritionGroup.DELETE("/today/:food_id", nutritionHandler.DeleteNutritionIntake)

	// Overview nutrition routes can be added here
	nutritionGroup.GET("/chart", nutritionHandler.GetNutritionChartData)
}

func setupBodyMeasurementRoutes(group *echo.Group) {
	// Define body measurement-related protected routes here
	bodyMeasurementGroup := group.Group("/body-measurements")

	// Initialize body measurement handlers
	bodyMeasurementRepo := models.NewBodyMeasurementRepository()
	bodyMeasurementHandler := api.NewBodyMeasurementHandlers(bodyMeasurementRepo)

	bodyMeasurementGroup.GET("", bodyMeasurementHandler.GetBodyMeasurements, validator.ValidateQuery(&validator.BodyMeasurementRequest{}))
	bodyMeasurementGroup.POST("", bodyMeasurementHandler.AddBodyMeasurement, validator.ValidateRequest(&validator.BodyMeasurementCreateRequest{}))
	bodyMeasurementGroup.PUT("/:measurement_id", bodyMeasurementHandler.UpdateBodyMeasurement, validator.ValidateRequest(&validator.BodyMeasurementCreateRequest{}))
	bodyMeasurementGroup.DELETE("/:measurement_id", bodyMeasurementHandler.DeleteBodyMeasurement)
}
