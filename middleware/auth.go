package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/config"
	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTClaims represents the claims stored in JWT token
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthUser represents authenticated user data stored in context
type AuthUser struct {
	ID               int               `json:"id"`
	Email            string            `json:"email"`
	Status           models.UserStatus `json:"status"`
	UserLevel        models.UserLevel  `json:"user_level"`
	PremiumExpiresAt *time.Time        `json:"premium_expires_at,omitempty"`
}

// GenerateToken generates a JWT token for the given user
func GenerateToken(userID int) (string, error) {
	cfg := config.Get()

	// Parse expires duration
	expiresIn, err := time.ParseDuration(cfg.JWT.ExpiresIn)
	if err != nil {
		expiresIn = 24 * time.Hour // fallback to 24 hours
	}

	// Create claims
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.Get()

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return token, nil
}

// AuthMiddleware provides JWT-based authentication middleware
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			token, err := ExtractTokenFromHeader(authHeader)
			if err != nil {
				Logger.Error().Err(err).Msg("[AuthMiddleware] Invalid authorization header")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Otorisasi diperlukan", nil)
			}

			// Validate token
			claims, err := ValidateToken(token)
			if err != nil {
				Logger.Error().Err(err).Msg("[AuthMiddleware] Token validation failed")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid", nil)
			}

			// Get user from database to ensure user still exists and get current data
			userRepo := models.NewUserRepository()
			user, err := userRepo.FindByID(claims.UserID)
			if err != nil {
				Logger.Error().Err(err).Msg("[AuthMiddleware] Gagal mengambil data pengguna")
				return helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengguna", nil)
			}

			if user == nil {
				Logger.Warn().Msg("[AuthMiddleware] Pengguna tidak ditemukan untuk token")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna tidak ditemukan", nil)
			}

			// Check if user is active
			if user.Status != models.UserStatusActive {
				return helper.ErrorResponse(c, http.StatusForbidden, "Account access denied", nil)
			}

			// Store authenticated user data in context
			authUser := &AuthUser{
				ID:               user.ID,
				Email:            user.Email,
				Status:           user.Status,
				UserLevel:        user.UserLevel,
				PremiumExpiresAt: user.PremiumExpiresAt,
			}

			c.Set("user", authUser)
			c.Set("user_id", user.ID)

			// Set user context in Sentry for error tracking
			SetUserContext(c, user.ID, user.Email)

			return next(c)
		}
	}
}

// RequireAuth returns a middleware that requires authentication
func RequireAuth() echo.MiddlewareFunc {
	return AuthMiddleware()
}

// RequirePremium returns a middleware that requires premium subscription
func RequirePremium() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authenticated user from context
			authUser, ok := c.Get("user").(*AuthUser)
			if !ok {
				Logger.Warn().Msg("[RequirePremium] Missing authenticated user in context")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", nil)
			}

			// Check if user has premium subscription
			if authUser.UserLevel == models.UserLevelFree {
				return helper.ErrorResponse(c, http.StatusForbidden, "Premium subscription required", nil)
			}

			// Check if premium subscription has expired
			if authUser.PremiumExpiresAt != nil && authUser.PremiumExpiresAt.Before(time.Now()) {
				return helper.ErrorResponse(c, http.StatusForbidden, "Premium subscription expired", nil)
			}

			return next(c)
		}
	}
}

// RequirePremiumPlus returns a middleware that requires premium+ subscription
func RequirePremiumPlus() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authenticated user from context
			authUser, ok := c.Get("user").(*AuthUser)
			if !ok {
				Logger.Warn().Msg("[RequirePremiumPlus] Missing authenticated user in context")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", nil)
			}

			// Check if user has premium+ subscription
			if authUser.UserLevel != models.UserLevelPremiumPlus {
				Logger.Warn().Str("user_level", string(authUser.UserLevel)).Msg("[RequirePremiumPlus] Non-premium+ user attempted premium+ route")
				return helper.ErrorResponse(c, http.StatusForbidden, "Premium+ subscription required", nil)
			}

			// Check if premium+ subscription has expired
			if authUser.PremiumExpiresAt != nil && authUser.PremiumExpiresAt.Before(time.Now()) {
				Logger.Warn().Str("user_level", string(authUser.UserLevel)).Msg("[RequirePremiumPlus] Premium+ subscription expired")

				return helper.ErrorResponse(c, http.StatusForbidden, "Premium+ subscription expired", nil)
			}

			return next(c)
		}
	}
}

// GetAuthUser retrieves authenticated user from context
func GetAuthUser(c echo.Context) (*AuthUser, error) {
	user, ok := c.Get("user").(*AuthUser)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

// GetUserID retrieves authenticated user ID from context
func GetUserID(c echo.Context) (int, error) {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// OptionalAuth middleware that tries to authenticate but doesn't fail if no auth provided
func OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No auth provided, continue without setting user context
				return next(c)
			}

			token, err := ExtractTokenFromHeader(authHeader)
			if err != nil {
				// Invalid auth format, continue without setting user context
				return next(c)
			}

			// Validate token
			claims, err := ValidateToken(token)
			if err != nil {
				// Invalid token, continue without setting user context
				return next(c)
			}

			// Get user from database
			userRepo := models.NewUserRepository()
			user, err := userRepo.FindByID(claims.UserID)
			if err != nil || user == nil {
				// User not found, continue without setting user context
				return next(c)
			}

			// Check if user is active
			if user.Status != models.UserStatusActive {
				// User not active, continue without setting user context
				return next(c)
			}

			// Store authenticated user data in context
			authUser := &AuthUser{
				ID:               user.ID,
				Email:            user.Email,
				Status:           user.Status,
				UserLevel:        user.UserLevel,
				PremiumExpiresAt: user.PremiumExpiresAt,
			}

			c.Set("user", authUser)
			c.Set("user_id", user.ID)

			return next(c)
		}
	}
}

// AdminRequired middleware that requires admin privileges
func AdminRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// For now, this is a placeholder. In a real implementation,
			// you would check for admin role in the user model or JWT claims

			authUser, err := GetAuthUser(c)
			if err != nil {
				Logger.Warn().Err(err).Msg("[AdminRequired] Authentication required")
				return helper.ErrorResponse(c, http.StatusUnauthorized, "Authentication required", nil)
			}

			// TODO: Implement proper admin role checking
			// For now, we'll use a simple check based on user ID or email
			// In production, you should have a proper role system

			adminEmails := []string{"admin@example.com", "superadmin@example.com"}
			isAdmin := false
			for _, adminEmail := range adminEmails {
				if authUser.Email == adminEmail {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				Logger.Warn().Str("email", authUser.Email).Msg("[AdminRequired] Non-admin user attempted admin route")
				return helper.ErrorResponse(c, http.StatusForbidden, "Admin access required", nil)
			}

			return next(c)
		}
	}
}
