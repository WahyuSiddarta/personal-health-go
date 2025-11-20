package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/config"
	database "github.com/WahyuSiddarta/be_saham_go/db"
	exLogger "github.com/WahyuSiddarta/be_saham_go/logger"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/router"

	"github.com/rs/zerolog"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

var (
	// Global instances
	Logger *zerolog.Logger
)

// handleCriticalError logs the error and exits the application gracefully
func handleCriticalError(logger *zerolog.Logger, operation string, err error) {
	if logger == nil {
		// If logger is not initialized, use a basic logger
		basicLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		basicLogger.Fatal().Err(err).Msgf("Critical error during %s", operation)
	} else {
		logger.Fatal().Err(err).Msgf("Critical error during %s", operation)
	}
	os.Exit(1)
}

// initializeSystem handles all system-wide initialization
func initializeSystem() *api.API {
	// Initialize logger first
	Logger = exLogger.InitLogger()
	exLogger.DistrubuteLogger(Logger)
	Logger.Info().Msg("System initialization started - 1 / 7 - Logger initialized")

	// Get configuration
	// var err error

	configStruct, err := config.Load()
	if err != nil {
		handleCriticalError(Logger, "loading configuration", err)
	}

	if configStruct == nil {
		handleCriticalError(Logger, "loading configuration 2", err)
	}

	Logger.Info().Msg("System initialization started - 2 / 7 - Configuration loaded")
	loglevel, err := zerolog.ParseLevel(configStruct.LogLevel)
	if err != nil {
		Logger.Warn().Err(err).Msg("Invalid log level in config, defaulting to info")
		loglevel, _ = zerolog.ParseLevel("info")
	}
	Logger.Level(loglevel)

	// Initialize Sentry
	if configStruct.SentryDSN != "" {
		// Set sampling rate based on environment
		// Development: 100% sampling (capture all transactions)
		// Production: 10% sampling (reduce server load)
		tracesSampleRate := 1.0 // Default to 100% for development
		if configStruct.Env == "production" {
			tracesSampleRate = 0.1 // 10% for production
		}

		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              configStruct.SentryDSN,
			Environment:      configStruct.Env,
			Release:          configStruct.AppVersion,
			EnableTracing:    true,
			TracesSampleRate: tracesSampleRate,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				// Ignore health check requests to reduce noise
				if event.Request != nil {
					if strings.Contains(event.Request.URL, "/health") {
						return nil
					}
				}
				return event
			},
		}); err != nil {
			Logger.Warn().Err(err).Msg("Sentry initialization failed")
		} else {
			Logger.Info().Msgf("System initialization started - 3 / 7 - Sentry initialized (TracesSampleRate: %.1f%%)", tracesSampleRate*100)
		}
	} else {
		Logger.Info().Msg("Sentry DSN not configured, skipping Sentry initialization")
	}

	// Create Echo instance
	Logger.Info().Msg("System initialization started - 4 / 7 - Echo instance created")
	echoInstance := echo.New()

	// Setup global middleware
	Logger.Info().Msg("System initialization started - 5 / 7 - Middleware initialized")
	middleware.SetupGlobalMiddleware(echoInstance)

	// Initialize databases and validate connections
	Logger.Info().Msg("System initialization started - 6 / 7 - Database initialized")
	var (
		dbManager models.DBManager
		wg        sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		pDBManagerRW := database.PSQLGetDBReadWrite()
		if pDBManagerRW == nil {
			handleCriticalError(Logger, "PostgreSQL database readwrite initialization failed", fmt.Errorf("PostgreSQL DB connection is nil"))
		} else {
			dbManager.PostgreDBManager.RW = pDBManagerRW
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pDBManagerRC := database.PSQLGetDBReadCache()
		if pDBManagerRC == nil {
			handleCriticalError(Logger, "PostgreSQL database read-cache initialization failed", fmt.Errorf("PostgreSQL DB connection is nil"))
		} else {
			dbManager.PostgreDBManager.RC = pDBManagerRC
		}
	}()

	wg.Wait()
	models.DBM = dbManager

	// Create API instance with all dependencies
	Logger.Info().Msg("System initialization started - 7 / 7 - API instance created")
	apiInstance := &api.API{
		Router: echoInstance,
	}
	Logger.Info().Msg("System initialization completed")
	return apiInstance
}

func main() {
	runtime.GOMAXPROCS(2 * runtime.NumCPU())
	fmt.Println("VCPU Proc :", runtime.NumCPU())

	// Initialize all system components
	apiInstance := initializeSystem()

	// Create router and setup routes
	r := router.New(apiInstance, Logger)
	go func() {
		if err := r.SetupRoutes(); err != nil {
			Logger.Fatal().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)
	<-quit

	Logger.Info().Msg("Application gracefully shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiInstance.Router.Shutdown(ctx); err != nil {
		Logger.Fatal().Err(err).Msg("Error during shutdown")
	}

	// Flush Sentry
	middleware.FlushSentry(5)
	Logger.Info().Msg("Application fully shutdown")
}
