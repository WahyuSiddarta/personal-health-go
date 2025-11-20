package api

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

type API struct {
	Router       *echo.Echo
	ServerIP     string
	ServerStatus string
}

// getUserIDFromContext extracts the user ID from the echo context
func getUserIDFromContext(c echo.Context) (int, error) {
	if userIDValue := c.Get("user_id"); userIDValue != nil {
		if userID, ok := userIDValue.(int); ok {
			return userID, nil
		}
		return 0, fmt.Errorf("user ID has invalid type stored for key user_id")
	}

	if userIDValue := c.Get("userID"); userIDValue != nil {
		if userID, ok := userIDValue.(int); ok {
			return userID, nil
		}
		return 0, fmt.Errorf("user ID has invalid type stored for key userID")
	}

	return 0, fmt.Errorf("user ID not found in context")
}

// parseLimitOffset extracts pagination params from query with sane defaults.
func parseLimitOffset(c echo.Context) (int, int) {
	limit := 10
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
