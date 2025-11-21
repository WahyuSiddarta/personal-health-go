package api

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

type UsersHandlers struct {
	repo models.UserRepository
}

// NewAuthHandlers creates a new instance of auth handlers
func NewUserHandlers(repo models.UserRepository) *UsersHandlers {
	return &UsersHandlers{repo: repo}
}

func (h *UsersHandlers) GetPersonalTarget(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	userTarget, err := h.repo.FindPersonalTarget(userId)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetPersonalTarget - FindPersonalTarget] Failed to get personal target")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get personal target", nil)
	}
	return helper.JsonResponse(c, http.StatusOK, userTarget)
}

func (h *UsersHandlers) UpdatePersonalBodyMeasurementTarget(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[UpdatePersonalBodyMeasurementTarget] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	req, ok := validatedRequest.(*validator.PersonalBodyMeasurementTargetRequest)
	if !ok {
		Logger.Error().Msg("[UpdatePersonalBodyMeasurementTarget] Failed to cast validated request to PersonalBodyMeasurementTargetRequest")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}
	userTarget := &models.UserTarget{
		UserId:        userId,
		BodyWeight:    req.BodyWeight,
		ViceralFat:    req.ViceralFat,
		FatPercentage: req.FatPercentage,
	}

	err := h.repo.UpdatePersonalBodyMeasurementTarget(userTarget)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdatePersonalBodyMeasurementTarget] Failed to update personal body measurement target")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update personal body measurement target", nil)
	}
	return helper.JsonResponse(c, http.StatusOK, userTarget)
}
func (h *UsersHandlers) UpdatePersonalNutritionTarget(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[UpdatePersonalNutritionTarget] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	req, ok := validatedRequest.(*validator.PersonalNutritionTargetRequest)
	if !ok {
		Logger.Error().Msg("[UpdatePersonalNutritionTarget] Failed to cast validated request to PersonalNutritionTargetRequest")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	userTarget := &models.UserTarget{
		UserId:           userId,
		NutritionCaloric: req.NutritionCaloric,
		NutritionProtein: req.NutritionProtein,
		NutritionCarbs:   req.NutritionCarbs,
		NutritionFat:     req.NutritionFat,
	}

	err := h.repo.UpdatePersonalNutritionTarget(userTarget)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdatePersonalNutritionTarget] Failed to update personal nutrition target")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update personal nutrition target", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, userTarget)
}
