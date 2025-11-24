package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// AuthHandlers contains all authentication-related handlers
type BodyMeasurementHandlers struct {
	repo models.BodyMeasurementRepository
}

// NewBodyMeasurementHandlers creates a new instance of body measurement handlers
func NewBodyMeasurementHandlers(repo models.BodyMeasurementRepository) *BodyMeasurementHandlers {
	return &BodyMeasurementHandlers{repo: repo}
}

// / Daily Nutrition Intake Handlers
func (h *BodyMeasurementHandlers) GetBodyMeasurements(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	// Get validated request from middleware
	req := validator.GetValidatedQuery(c).(*validator.BodyMeasurementRequest)
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := config.Get().PaginationDefaultPageSize
	userMeasurements, err := h.repo.GetByUserId(userId, limit, page)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error().Err(err).Msg("[GetBodyMeasurements] Failed to get today's nutrition intake")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to get today's nutrition intake", nil)
	}

	hasNext := false
	if len(userMeasurements) > limit {
		hasNext = true
		userMeasurements = userMeasurements[:limit]
	}

	response := map[string]interface{}{
		"measurements": userMeasurements,
		"nextPage":     hasNext,
	}

	return helper.JsonResponse(c, http.StatusOK, response)
}

func (h *BodyMeasurementHandlers) AddBodyMeasurement(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[AddBodyMeasurement] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	measurementRequest := validatedRequest.(*validator.BodyMeasurementCreateRequest)

	var viceralFat, fatPercentage, nickCm, waistCm *float64
	if measurementRequest.ViceralFat != 0 {
		viceralFat = &measurementRequest.ViceralFat
	}
	if measurementRequest.FatPercentage != 0 {
		fatPercentage = &measurementRequest.FatPercentage
	}
	if measurementRequest.NickCm != 0 {
		nickCm = &measurementRequest.NickCm
	}
	if measurementRequest.WaistCm != 0 {
		waistCm = &measurementRequest.WaistCm
	}

	newMeasurement := models.BodyMeasurement{
		Bodyweight:    measurementRequest.Bodyweight,
		ViceralFat:    viceralFat,
		FatPercentage: fatPercentage,
		NickCm:        nickCm,
		WaistCm:       waistCm,
	}

	err := h.repo.Create(userId, newMeasurement)
	if err != nil {
		Logger.Error().Err(err).Msg("[AddBodyMeasurement] Failed to add body measurement")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to add body measurement", nil)
	}

	Logger.Info().Msgf("[AddBodyMeasurement] Added new body measurement for user %d", userId)
	return helper.JsonResponse(c, http.StatusCreated, map[string]string{"message": "Body measurement added successfully"})
}

func (h *BodyMeasurementHandlers) UpdateBodyMeasurement(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	measurementIdParam := c.Param("measurement_id")
	measurementId, err := strconv.Atoi(measurementIdParam)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateBodyMeasurement] Invalid measurement ID parameter")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid measurement ID", nil)
	}

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[UpdateBodyMeasurement] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	measurementRequest := validatedRequest.(*validator.BodyMeasurementCreateRequest)

	var viceralFat, fatPercentage, nickCm, waistCm *float64
	if measurementRequest.ViceralFat != 0 {
		viceralFat = &measurementRequest.ViceralFat
	}
	if measurementRequest.FatPercentage != 0 {
		fatPercentage = &measurementRequest.FatPercentage
	}
	if measurementRequest.NickCm != 0 {
		nickCm = &measurementRequest.NickCm
	}
	if measurementRequest.WaistCm != 0 {
		waistCm = &measurementRequest.WaistCm
	}

	updatedMeasurement := &models.BodyMeasurement{
		Bodyweight:    measurementRequest.Bodyweight,
		ViceralFat:    viceralFat,
		FatPercentage: fatPercentage,
		NickCm:        nickCm,
		WaistCm:       waistCm,
	}

	err = h.repo.Update(userId, measurementId, updatedMeasurement)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateBodyMeasurement] Failed to update body measurement")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update body measurement", nil)
	}

	Logger.Info().Msgf("[UpdateBodyMeasurement] Updated body measurement %d for user %d", measurementId, userId)
	return helper.JsonResponse(c, http.StatusOK, map[string]string{"message": "Body measurement updated successfully"})
}

func (h *BodyMeasurementHandlers) DeleteBodyMeasurement(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	measurementIdParam := c.Param("measurement_id")
	measurementId, err := strconv.Atoi(measurementIdParam)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteBodyMeasurement] Invalid measurement ID parameter")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid measurement ID", nil)
	}

	err = h.repo.Delete(userId, measurementId)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteBodyMeasurement] Failed to delete body measurement")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete body measurement", nil)
	}

	Logger.Info().Msgf("[DeleteBodyMeasurement] Deleted body measurement %d for user %d", measurementId, userId)
	return helper.JsonResponse(c, http.StatusOK, map[string]string{"message": "Body measurement deleted successfully"})
}
