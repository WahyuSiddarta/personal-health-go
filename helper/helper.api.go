package helper

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
)

// ErrorResponse sends a standardized error response
func ErrorResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	var builder strings.Builder

	builder.WriteString(`{"success":false,"message":"`)
	builder.WriteString(message)
	builder.WriteString(`"`)

	if data != nil {
		dataJSON, err := sonic.Marshal(data)
		if err == nil {
			builder.WriteString(`,"data":`)
			builder.Write(dataJSON)
		}
	}

	builder.WriteString(`}`)

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(statusCode)
	_, err := c.Response().Write([]byte(builder.String()))
	return err
}

// JsonResponse sends a standardized success response
func JsonResponse(c echo.Context, statusCode int, data interface{}) error {
	var builder strings.Builder
	builder.WriteString(`{"success":true`)

	if data != nil {
		dataJSON, err := sonic.Marshal(data)
		if err == nil {
			builder.WriteString(`,"data":`)
			builder.Write(dataJSON)
		}
	}

	builder.WriteString(`}`)

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(statusCode)
	_, err := c.Response().Write([]byte(builder.String()))

	return err
}

var ERROR struct {
	InvalidRequest string
}
