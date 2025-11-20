package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

// Config holds all configuration for our application
type Config struct {
	// Server Configuration
	Port       string
	Env        string
	LogLevel   string
	AppVersion string

	// Sentry Configuration
	SentryDSN string

	// JWT Configuration
	JWT JWTConfig

	// CORS Configuration
	CORS CORSConfig

	// Rate Limiting Configuration
	RateLimit RateLimitConfig

	// Database Configuration
	Database DatabaseConfig
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// Read-Write Database
	RW DatabaseConnection

	// Read-Cache Database
	RC DatabaseConnection
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
	Enabled           bool
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	Enabled        bool
}

// DatabaseConnection holds connection details for a database
type DatabaseConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Schema   string
	MaxCon   int
	MaxIdle  int
}

var config *Config

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env file doesn't exist in production
		if Logger != nil {
			Logger.Warn().Err(err).Msg(".env file not found")
		}
	}

	config = &Config{
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),
		SentryDSN:  getEnv("SENTRY_DSN", ""),
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseCORSOrigins(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")),
			Enabled:        getEnv("CORS_ENABLED", "true") == "true",
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: getEnvAsInt("RATE_LIMIT_RPS", 100),
			BurstSize:         getEnvAsInt("RATE_LIMIT_BURST", 200),
			Enabled:           getEnv("RATE_LIMIT_ENABLED", "true") == "true",
		},
		Database: DatabaseConfig{
			RW: DatabaseConnection{
				Host:     getEnv("DB_RW_HOST", "localhost"),
				Port:     getEnv("DB_RW_PORT", "5432"),
				User:     getEnv("DB_RW_USER", "postgres"),
				Password: getEnv("DB_RW_PASSWORD", ""),
				DBName:   getEnv("DB_RW_NAME", "saham_db"),
				Schema:   getEnv("DB_RW_SCHEMA", "public"),
				MaxCon:   getEnvAsInt("DB_RW_MAX_CONNECTIONS", 25),
				MaxIdle:  getEnvAsInt("DB_RW_MAX_IDLE", 10),
			},
			RC: DatabaseConnection{
				Host:     getEnv("DB_RC_HOST", getEnv("DB_RW_HOST", "localhost")),
				Port:     getEnv("DB_RC_PORT", getEnv("DB_RW_PORT", "5432")),
				User:     getEnv("DB_RC_USER", getEnv("DB_RW_USER", "postgres")),
				Password: getEnv("DB_RC_PASSWORD", getEnv("DB_RW_PASSWORD", "")),
				DBName:   getEnv("DB_RC_NAME", getEnv("DB_RW_NAME", "saham_db")),
				Schema:   getEnv("DB_RC_SCHEMA", getEnv("DB_RW_SCHEMA", "public")),
				MaxCon:   getEnvAsInt("DB_RC_MAX_CONNECTIONS", 25),
				MaxIdle:  getEnvAsInt("DB_RC_MAX_IDLE", 10),
			},
		},
	}

	return config, nil
}

// Get returns the loaded configuration
func Get() *Config {
	if config == nil {
		var err error
		config, err = Load()
		if err != nil {
			if Logger != nil {
				Logger.Fatal().Err(err).Msg("Failed to load configuration")
			} else {
				log.Fatalf("Failed to load configuration: %v", err)
			}
		}
	}
	return config
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a fallback default value
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// parseCORSOrigins parses a comma-separated string of CORS origins into a slice
func parseCORSOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}

	origins := strings.Split(originsStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}
