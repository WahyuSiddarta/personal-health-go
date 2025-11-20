package validator

// PersonalNutritionTargetRequest represents the request payload for updating personal nutrition targets.
type PersonalNutritionTargetRequest struct {
	NutritionCaloric float64 `json:"nutrition_caloric" validate:"gte=0,decimal2"`
	NutritionProtein float64 `json:"nutrition_protein" validate:"gte=0,decimal2"`
	NutritionCarbs   float64 `json:"nutrition_carbohydrate" validate:"gte=0,decimal2"`
	NutritionFat     float64 `json:"nutrition_fat" validate:"gte=0,decimal2"`
}
