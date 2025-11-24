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
type ExcerciseHandlers struct {
	repo models.ExcerciseRecordRepository
}

// NewExcerciseHandlers creates a new instance of exercise handlers
func NewExcerciseHandlers(repo models.ExcerciseRecordRepository) *ExcerciseHandlers {
	return &ExcerciseHandlers{repo: repo}
}

// Get user excercise records
func (h *ExcerciseHandlers) GetUserExcercises(c echo.Context) error {
	var userId int = 1 // Replace with actual user ID retrieval logic

	// Get validated request from middleware
	req := validator.GetValidatedQuery(c).(*validator.ExerciseRequest)
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := config.Get().PaginationDefaultPageSize
	userMeasurements, err := h.repo.GetByUserId(userId, limit, page)
	if err != nil && err != sql.ErrNoRows {
		Logger.Error().Err(err).Msg("[GetUserExcercises] Failed to get today's nutrition intake")
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
func (h *ExcerciseHandlers) AddExercise(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[AddExercise] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	measurementRequest := validatedRequest.(*validator.ExcerciseMutationRequest)
	var Minute *int
	if measurementRequest.Minute != nil && *measurementRequest.Minute != 0 {
		Minute = measurementRequest.Minute
	}

	newMeasurement := models.ExcerciseRecord{
		Minute:  Minute,
		Caloric: measurementRequest.Caloric,
		Type:    measurementRequest.Type,
	}

	err := h.repo.Create(userId, newMeasurement)
	if err != nil {
		Logger.Error().Err(err).Msg("[AddBodyMeasurement] Failed to add body measurement")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to add body measurement", nil)
	}

	Logger.Info().Msgf("[AddBodyMeasurement] Added new body measurement for user %d", userId)
	return helper.JsonResponse(c, http.StatusCreated, map[string]string{"message": "Body measurement added successfully"})
}
func (h *ExcerciseHandlers) UpdateExercise(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	excerciseIdParam := c.Param("exercise_id")
	excerciseId, err := strconv.Atoi(excerciseIdParam)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateExercise] Invalid exercise ID parameter")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid exercise ID", nil)
	}

	// Get validated request from middleware
	validatedRequest := validator.GetValidatedRequest(c)
	if validatedRequest == nil {
		Logger.Error().Msg("[UpdateExercise] No validated request found in context")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	measurementRequest := validatedRequest.(*validator.ExcerciseMutationRequest)
	var Minute *int
	if measurementRequest.Minute != nil && *measurementRequest.Minute != 0 {
		Minute = measurementRequest.Minute
	}

	updatedMeasurement := &models.ExcerciseRecord{
		Minute:  Minute,
		Caloric: measurementRequest.Caloric,
		Type:    measurementRequest.Type,
	}

	err = h.repo.Update(userId, excerciseId, updatedMeasurement)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateExercise] Failed to update exercise record")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update exercise record", nil)
	}

	Logger.Info().Msgf("[UpdateExercise] Updated exercise record %d for user %d", excerciseId, userId)
	return helper.JsonResponse(c, http.StatusOK, map[string]string{"message": "Exercise record updated successfully"})
}

func (h *ExcerciseHandlers) DeleteExercise(c echo.Context) error {
	var userId int = 1
	// Replace with actual user ID retrieval logic

	exerciseIdParam := c.Param("exercise_id")
	exerciseId, err := strconv.Atoi(exerciseIdParam)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteExercise] Invalid exercise ID parameter")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid exercise ID", nil)
	}

	err = h.repo.Delete(userId, exerciseId)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteExercise] Failed to delete exercise record")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete exercise record", nil)
	}

	Logger.Info().Msgf("[DeleteExercise] Deleted exercise record %d for user %d", exerciseId, userId)
	return helper.JsonResponse(c, http.StatusOK, map[string]string{"message": "Exercise record deleted successfully"})
}
