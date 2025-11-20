package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger provides request logging middleware
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Log request details
			if Logger != nil {
				duration := time.Since(start)

				// Log successful requests
				if err == nil {
					Logger.Info().
						Str("method", c.Request().Method).
						Str("path", c.Request().URL.Path).
						Str("remote_ip", c.RealIP()).
						Int("status", c.Response().Status).
						Dur("duration", duration).
						Str("user_agent", c.Request().UserAgent()).
						Msg("Request processed")
				} else {
					// Log errors
					Logger.Error().
						Str("method", c.Request().Method).
						Str("path", c.Request().URL.Path).
						Str("remote_ip", c.RealIP()).
						Dur("duration", duration).
						Err(err).
						Msg("Request failed")
				}
			}

			return err
		}
	}
}

// HealthCheckLogger provides lighter logging for health check endpoints
func HealthCheckLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Only log health checks at debug level to avoid spam
			if Logger != nil && c.Request().URL.Path == "/api/public/health" {
				duration := time.Since(start)
				Logger.Debug().
					Str("method", c.Request().Method).
					Str("path", c.Request().URL.Path).
					Int("status", c.Response().Status).
					Dur("duration", duration).
					Msg("Health check")
			}

			return err
		}
	}
}
