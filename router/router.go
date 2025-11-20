package router

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/WahyuSiddarta/be_saham_go/middleware"

	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

// Router handles all route setup and configuration
type Router struct {
	API *api.API
}

// New creates a new Router instance
func New(apiInstance *api.API, logger *zerolog.Logger) *Router {
	return &Router{
		API: apiInstance,
	}
}

// SetupRoutes configures all routes and starts the server
func (r *Router) SetupRoutes() error {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error().Interface("panic", err).Msg("Panic occurred during route setup")
			panic(err)
		}
	}()

	// Setup API group with middleware
	apiGroup := r.API.Router.Group("/api")

	// Apply middleware specific to API routes
	middleware.SetupAPIMiddleware(apiGroup) // Setup all route groups
	r.setupAuthRoutes(apiGroup)
	r.setupPublicRoutes(apiGroup)
	r.setupProtectedRoutes(apiGroup) // Future routes that require authentication

	Logger.Info().Msg("All routes configured successfully")

	// Start server
	port := config.Get().Port
	if port == "" {
		port = ":1328" // fallback default
	}
	if port[0] != ':' {
		port = ":" + port
	}
	Logger.Info().Msgf("Starting server on port %s", port)
	err := r.API.Router.Start(port)
	if err != nil && err != http.ErrServerClosed {
		Logger.Fatal().Err(err).Msgf("Failed to start server on port %s", port)
		return err
	}
	return nil
}
