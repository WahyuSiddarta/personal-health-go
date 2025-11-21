package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// AuthHandlers contains all authentication-related handlers
type NutritionHandlers struct {
	repo models.NutritionRepository
}

// NewNutritionHandlers creates a new instance of nutrition handlers
func NewNutritionHandlers(repo models.NutritionRepository) *NutritionHandlers {
	return &NutritionHandlers{repo: repo}
}

// / Overview Nutrition Handlers
func (h *NutritionHandlers) GetNutritionChartData(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	chartData, err := h.repo.GetNutritionChartData(userId)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetNutritionChartData] Failed to get nutrition chart data")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get nutrition chart data", nil)
	}

	Logger.Info().Msgf("[GetNutritionChartData] Retrieved nutrition chart data for user %d", userId)
	return helper.JsonResponse(c, http.StatusOK, chartData)
}

// / Daily Nutrition Intake Handlers
func (h *NutritionHandlers) GetTodaysNutritionIntake(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	userIntakes, err := h.repo.FindUserTodayIntake(userId)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error().Err(err).Msg("[GetTodaysNutritionIntake] Failed to get today's nutrition intake")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get today's nutrition intake", nil)
	}

	Logger.Info().Msgf("[GetTodaysNutritionIntake] Retrieved %d nutrition intake records for user %d", len(userIntakes), userId)
	return helper.JsonResponse(c, http.StatusOK, userIntakes)
}

func (h *NutritionHandlers) AddNutritionIntake(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic
	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[AddNutritionIntake] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	req, ok := validatedRequest.(*validator.NutritionRequest)
	if !ok {
		Logger.Error().Msg("[AddNutritionIntake] Failed to cast validated request to NutritionRequest")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	nutritionTracker := &models.NutritionTracker{
		UserId:       userId,
		Category:     models.NutritionCategory(req.Category),
		Fat:          req.Fat,
		Protein:      req.Protein,
		Carbohydrate: req.Carbohydrate,
		Caloric:      req.Caloric,
		Name:         req.Name,
	}

	err := h.repo.AddTodayIntake(nutritionTracker)
	if err != nil {
		Logger.Error().Err(err).Msg("[AddTodayIntake] Failed to add nutrition intake")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to add nutrition intake", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, nutritionTracker)
}

func (h *NutritionHandlers) UpdateNutritionIntake(c echo.Context) error {
	// Implementation for updating nutrition intake goes here
	var userId int = 1 // Replace with actual user ID retrieval logic
	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[UpdateNutritionIntake] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	req, ok := validatedRequest.(*validator.NutritionRequest)
	if !ok {
		Logger.Error().Msg("[UpdateNutritionIntake] Failed to cast validated request to NutritionRequest")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}
	foodId := c.Param("food_id")
	if foodId == "" {
		Logger.Error().Msg("[UpdateNutritionIntake] food_id parameter is missing")
		return helper.ErrorResponse(c, http.StatusBadRequest, "food_id parameter is required", nil)
	}

	foodIdInt, err := strconv.Atoi(foodId)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateNutritionIntake] Invalid food_id format")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid food_id format", nil)
	}

	nutritionTracker := &models.NutritionTracker{
		UserId:       userId,
		FoodId:       foodIdInt,
		Category:     models.NutritionCategory(req.Category),
		Fat:          req.Fat,
		Protein:      req.Protein,
		Carbohydrate: req.Carbohydrate,
		Caloric:      req.Caloric,
		Name:         req.Name,
	}

	err = h.repo.UpdateTodayIntake(nutritionTracker)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateNutritionIntake] Failed to update nutrition intake")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update nutrition intake", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nutritionTracker)
}

func (h *NutritionHandlers) DeleteNutritionIntake(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic
	foodId := c.Param("food_id")
	if foodId == "" {
		Logger.Error().Msg("[DeleteNutritionIntake] food_id parameter is missing")
		return helper.ErrorResponse(c, http.StatusBadRequest, "food_id parameter is required", nil)
	}

	foodIdInt, err := strconv.Atoi(foodId)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteNutritionIntake] Invalid food_id format")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid food_id format", nil)
	}

	err = h.repo.DeleteTodayIntake(userId, foodIdInt)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteNutritionIntake] Failed to delete nutrition intake")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete nutrition intake", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]string{"message": "Nutrition intake deleted successfully"})
}
