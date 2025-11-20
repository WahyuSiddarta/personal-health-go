package middleware

import (
	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ConfigureCORS configures and returns CORS middleware based on configuration
func ConfigureCORS() echo.MiddlewareFunc {
	corsConfig := config.Get().CORS

	if !corsConfig.Enabled {
		// Return a no-op middleware if CORS is disabled
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	})
}

// LogCORSStatus logs the current CORS configuration status
func LogCORSStatus() {
	corsConfig := config.Get().CORS

	if corsConfig.Enabled {
		Logger.Info().
			Strs("allowed_origins", corsConfig.AllowedOrigins).
			Msg("CORS middleware enabled")
	} else {
		Logger.Info().Msg("CORS middleware disabled")
	}
}
