package router

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/labstack/echo/v4"
)

// setupPublicRoutes configures public routes (no authentication required)
func (r *Router) setupPublicRoutes(apiGroup *echo.Group) {

	rpub := apiGroup.Group("/public")
	rpub.GET("/health", func(c echo.Context) error {
		healthData := map[string]interface{}{
			"service":     "be_saham_go",
			"version":     config.Get().AppVersion,
			"environment": config.Get().Env,
			"rate_limit":  "enabled",
		}
		return helper.JsonResponse(c, http.StatusOK, healthData)
	})

	// TEST endpoint - accessible at /api/public/test
	rpub.GET("/test", r.API.Test)

	// Test panic recovery - accessible at /api/public/test-panic (for testing only)
	rpub.GET("/test-panic", func(c echo.Context) error {
		// This endpoint intentionally panics to test the recover middleware
		panic("Test panic for middleware testing")
	})
}
