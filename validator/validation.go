package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/WahyuSiddarta/be_saham_go/helper"
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
	validate.RegisterValidation("decimal2", validateDecimalPlaces)
	validate.RegisterValidation("nutrition_category", validateNutritionCategory)
	validate.RegisterValidation("caloric_calculation", validateCaloricCalculation)
}

// validateDecimalPlaces validates that a float64 has at most 2 decimal places
func validateDecimalPlaces(fl validator.FieldLevel) bool {
	value := fl.Field().Float()

	// Use string formatting to check decimal places to avoid floating point precision issues
	formatted := fmt.Sprintf("%.10f", value) // Format with 10 decimal places

	// Find the decimal point
	dotIndex := strings.Index(formatted, ".")
	if dotIndex == -1 {
		return true // No decimal point means it's a whole number
	}

	// Count significant decimal places (excluding trailing zeros)
	decimalPart := strings.TrimRight(formatted[dotIndex+1:], "0")

	return len(decimalPart) <= 2
}

// validateNutritionCategory validates that a nutrition category is one of the allowed values
func validateNutritionCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	category := NutritionCategory(value)
	return category.IsValid()
}

// validateCaloricCalculation validates that the caloric value matches the calculation from macronutrients
func validateCaloricCalculation(fl validator.FieldLevel) bool {
	// Get the parent struct
	parent := fl.Parent()

	// Extract field values from the struct
	caloric := fl.Field().Float()
	fat := parent.FieldByName("Fat").Float()
	protein := parent.FieldByName("Protein").Float()
	carbohydrate := parent.FieldByName("Carbohydrate").Float()

	// Use the checkCaloricValue function logic
	expectedCaloric := (fat * 9) + (protein * 4) + (carbohydrate * 4)

	// Allow small floating point tolerance (0.1 calorie difference)
	tolerance := 0.1
	return abs(caloric-expectedCaloric) <= tolerance
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
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

// getValidationMessage returns user-friendly validation messages in English
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Please enter a valid email address"
	case "min":
		if field == "password" {
			return fmt.Sprintf("Password must be at least %s characters", param)
		}
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "decimal2":
		return fmt.Sprintf("%s can only have a maximum of 2 decimal places", field)
	case "nutrition_category":
		return fmt.Sprintf("%s must be one of: breakfast, lunch, dinner, snack", field)
	case "caloric_calculation":
		return fmt.Sprintf("%s does not match the calculated value from macronutrients (fat×9 + protein×4 + carbohydrate×4)", field)
	case "datetime":
		return "Invalid date format. Please use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)"
	default:
		return fmt.Sprintf("%s is invalid", field)
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
					helper.Logger.Warn().Err(err).Msg("Invalid request format during binding")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", nil)
			}

			// Validate the struct
			if errs := ValidateStruct(req); len(errs) > 0 {
				if helper.Logger != nil {
					helper.Logger.Warn().Int("error_count", len(errs)).Msg("Request validation errors")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Validation errors", errs)
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
					helper.Logger.Warn().Err(err).Msg("Invalid query parameters during binding")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", nil)
			}

			// Validate the struct
			if errs := ValidateStruct(query); len(errs) > 0 {
				if helper.Logger != nil {
					helper.Logger.Warn().Int("error_count", len(errs)).Msg("Query validation errors")
				}
				return helper.ErrorResponse(c, http.StatusBadRequest, "Query validation errors", errs)
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
