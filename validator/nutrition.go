package validator

// CreateNutritionRequest represents the request payload for creating a nutrition intake entry.
type NutritionRequest struct {
	Category     NutritionCategory `json:"category" db:"category" validate:"required,nutrition_category"`
	Fat          float64           `json:"fat" db:"fat" validate:"gte=0,decimal2"`
	Protein      float64           `json:"protein" db:"protein" validate:"gte=0,decimal2"`
	Carbohydrate float64           `json:"carbohydrate" db:"carbohydrate" validate:"gte=0,decimal2"`
	Caloric      float64           `json:"caloric" db:"caloric" validate:"gte=0,decimal2,caloric_calculation"`
	Name         string            `json:"name" db:"name" validate:"required,min=1,max=255"`
}

// NutritionCategory represents nutrition category, like breakfast, lunch, dinner, snack
type NutritionCategory string

const (
	NutritionCategoryBreakfast NutritionCategory = "breakfast"
	NutritionCategoryLunch     NutritionCategory = "lunch"
	NutritionCategoryDinner    NutritionCategory = "dinner"
	NutritionCategorySnack     NutritionCategory = "snack"
)

// ValidNutritionCategories returns a slice of all valid nutrition categories
func ValidNutritionCategories() []NutritionCategory {
	return []NutritionCategory{
		NutritionCategoryBreakfast,
		NutritionCategoryLunch,
		NutritionCategoryDinner,
		NutritionCategorySnack,
	}
}

// IsValid checks if the nutrition category is valid
func (nc NutritionCategory) IsValid() bool {
	validCategories := ValidNutritionCategories()
	for _, valid := range validCategories {
		if nc == valid {
			return true
		}
	}
	return false
}
