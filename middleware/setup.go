package middleware

import (
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

// SetupGlobalMiddleware configures all global middleware for the Echo instance
func SetupGlobalMiddleware(e *echo.Echo) {
	// Add panic recovery middleware (should be first for safety)
	e.Use(Recover())

	// Add Sentry middleware
	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: false, // Already handled by custom Recover() middleware
	}))

	// Add request logging middleware
	e.Use(RequestLogger())

	// Add CORS middleware
	e.Use(ConfigureCORS())

	// Log middleware setup status
	LogCORSStatus()

	Logger.Info().Msg("Global middleware configured: Panic Recovery, Sentry, Request Logging, CORS")
}

// SetupAPIMiddleware configures middleware specifically for API routes
func SetupAPIMiddleware(apiGroup *echo.Group) {
	// Add rate limiting to all API routes
	rateLimiter := NewRateLimiter()
	apiGroup.Use(rateLimiter.Middleware())

	// Log rate limiting status
	LogRateLimitStatus()

	Logger.Info().Msg("API middleware configured: Rate Limiting")
}
