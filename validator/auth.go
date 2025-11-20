package validator

import (
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
)

// LoginRequest represents login request payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// RegisterRequest represents registration request payload.
type RegisterRequest struct {
	Email           string             `json:"email" validate:"required,email"`
	Password        string             `json:"password" validate:"required,min=6"`
	ConfirmPassword string             `json:"confirmPassword" validate:"required,eqfield=Password"`
	Status          *models.UserStatus `json:"status,omitempty" validate:"omitempty,user_status"`
	UserLevel       *models.UserLevel  `json:"user_level,omitempty" validate:"omitempty,user_level"`
}

// UpdateUserLevelRequest represents request to update user level.
type UpdateUserLevelRequest struct {
	UserLevel   models.UserLevel    `json:"user_level" validate:"required,user_level"`
	PaymentData *PaymentDataRequest `json:"payment_data,omitempty"`
}

// PaymentDataRequest represents payment data nested inside other requests.
type PaymentDataRequest struct {
	OriginalPrice  float64  `json:"original_price" validate:"required,gt=0"`
	PaidPrice      float64  `json:"paid_price" validate:"required,gte=0"`
	DiscountAmount *float64 `json:"discount_amount,omitempty" validate:"omitempty,gte=0"`
	DiscountReason *string  `json:"discount_reason,omitempty" validate:"omitempty,max=255"`
	PaymentMethod  *string  `json:"payment_method,omitempty" validate:"omitempty,max=50"`
	PaymentDate    *string  `json:"payment_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Notes          *string  `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// UpdateUserStatusRequest represents request to update user status.
type UpdateUserStatusRequest struct {
	Status models.UserStatus `json:"status" validate:"required,user_status"`
}

// GetUsersQuery represents query parameters for fetching users.
type GetUsersQuery struct {
	Page        int                `query:"page" validate:"omitempty,min=1"`
	Limit       int                `query:"limit" validate:"omitempty,min=1,max=100"`
	Status      *models.UserStatus `query:"status" validate:"omitempty,user_status"`
	UserLevel   *models.UserLevel  `query:"user_level" validate:"omitempty,user_level"`
	EmailFilter *string            `query:"email_filter" validate:"omitempty,max=100"`
}

// LoginData contains login response data.
type LoginData struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// RegisterData contains registration response data.
type RegisterData struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// ProfileData contains authenticated user profile data.
type ProfileData struct {
	User *middleware.AuthUser `json:"user"`
}

// ValidationError represents a field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}
