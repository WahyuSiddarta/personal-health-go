package api

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
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
