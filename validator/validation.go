package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use JSON tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	validate.RegisterValidation("user_status", validateUserStatus)
	validate.RegisterValidation("user_level", validateUserLevel)
}

// validateUserStatus validates user status enum
func validateUserStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{
		string(models.UserStatusActive),
		string(models.UserStatusInactive),
		string(models.UserStatusSuspended),
		string(models.UserStatusBanned),
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// validateUserLevel validates user level enum
func validateUserLevel(fl validator.FieldLevel) bool {
	level := fl.Field().String()
	validLevels := []string{
		string(models.UserLevelFree),
		string(models.UserLevelPremium),
		string(models.UserLevelPremiumPlus),
	}

	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}
	return false
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			message := getValidationMessage(err)

			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
				Tag:     err.Tag(),
			})
		}
	}

	return errors
}

// getValidationMessage returns user-friendly validation messages in Indonesian
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "email":
		return "Silakan masukkan alamat email yang valid"
	case "min":
		if field == "password" {
			return fmt.Sprintf("Password minimal %s karakter", param)
		}
		return fmt.Sprintf("%s minimal %s karakter", field, param)
	case "max":
		return fmt.Sprintf("%s tidak boleh melebihi %s karakter", field, param)
	case "eqfield":
		return fmt.Sprintf("%s harus sama dengan %s", field, param)
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s", field, param)
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s", field, param)
	case "user_status":
		return "Status harus salah satu dari: active, inactive, suspended, banned"
	case "user_level":
		return "Level pengguna harus salah satu dari: free, premium, premium+"
	case "datetime":
		return "Format tanggal tidak valid. Gunakan format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)"
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}

// ValidateRequest is a middleware function that validates request body
func ValidateRequest(requestType interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create new instance of the request type
			req := reflect.New(reflect.TypeOf(requestType).Elem()).Interface()

			// Bind request body to struct
			if err := c.Bind(req); err != nil {
				if helper.Logger != nil {
					helper.Logger.Warn().Err(err).Msg("Format request tidak valid saat binding")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Format request tidak valid", nil)
			}

			// Validate the struct
			if errs := ValidateStruct(req); len(errs) > 0 {
				if helper.Logger != nil {
					helper.Logger.Warn().Int("error_count", len(errs)).Msg("Kesalahan validasi request")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Kesalahan validasi", errs)
			}

			// Store validated request in context
			c.Set("validated_request", req)
			return next(c)
		}
	}
}

// ValidateQuery is a middleware function that validates query parameters
func ValidateQuery(queryType interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create new instance of the query type
			query := reflect.New(reflect.TypeOf(queryType).Elem()).Interface()

			// Bind query parameters to struct
			if err := c.Bind(query); err != nil {
				if helper.Logger != nil {
					helper.Logger.Warn().Err(err).Msg("Parameter query tidak valid saat binding")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Parameter query tidak valid", nil)
			}

			// Validate the struct
			if errs := ValidateStruct(query); len(errs) > 0 {
				if helper.Logger != nil {
					helper.Logger.Warn().Int("error_count", len(errs)).Msg("Kesalahan validasi query")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Kesalahan validasi query", errs)
			}

			// Store validated query in context
			c.Set("validated_query", query)
			return next(c)
		}
	}
}

// GetValidatedRequest retrieves validated request from context
func GetValidatedRequest(c echo.Context) interface{} {
	return c.Get("validated_request")
}

// GetValidatedQuery retrieves validated query from context
func GetValidatedQuery(c echo.Context) interface{} {
	return c.Get("validated_query")
}

// CustomValidation handles custom validation logic for complex scenarios
func CustomValidation(req interface{}) []ValidationError {
	var errors []ValidationError

	// Handle specific custom validations
	switch v := req.(type) {
	case *UpdateUserLevelRequest:
		// If payment data is provided, ensure it's for premium tiers
		if v.PaymentData != nil && v.UserLevel == models.UserLevelFree {
			errors = append(errors, ValidationError{
				Field:   "payment_data",
				Message: "Data pembayaran hanya dapat diberikan untuk tier premium",
				Tag:     "custom",
			})
		}

		// If user level is premium, validate payment data is provided
		if (v.UserLevel == models.UserLevelPremium || v.UserLevel == models.UserLevelPremiumPlus) &&
			v.PaymentData != nil {
			// Additional payment validation logic can be added here
			if v.PaymentData.OriginalPrice <= 0 {
				errors = append(errors, ValidationError{
					Field:   "payment_data.original_price",
					Message: "Harga asli harus lebih besar dari 0",
					Tag:     "custom",
				})
			}
			if v.PaymentData.PaidPrice < 0 {
				errors = append(errors, ValidationError{
					Field:   "payment_data.paid_price",
					Message: "Harga yang dibayar tidak boleh negatif",
					Tag:     "custom",
				})
			}
		}
	}

	return errors
}

// ValidateAndCustom combines struct validation with custom validation
func ValidateAndCustom(req interface{}) []ValidationError {
	// Standard struct validation
	errors := ValidateStruct(req)

	// Add custom validation errors
	customErrors := CustomValidation(req)
	errors = append(errors, customErrors...)

	return errors
}
