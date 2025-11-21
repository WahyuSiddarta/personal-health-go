package validator

// PersonalNutritionTargetRequest represents the request payload for updating personal nutrition targets.
type PersonalNutritionTargetRequest struct {
	NutritionCaloric float64 `json:"nutrition_caloric" validate:"gte=0,decimal2"`
	NutritionProtein float64 `json:"nutrition_protein" validate:"gte=0,decimal2"`
	NutritionCarbs   float64 `json:"nutrition_carbohydrate" validate:"gte=0,decimal2"`
	NutritionFat     float64 `json:"nutrition_fat" validate:"gte=0,decimal2"`
}

// PersonalBodyMeasurementTargetRequest represents the request payload for updating personal nutrition targets.
type PersonalBodyMeasurementTargetRequest struct {
	BodyWeight    float64 `json:"body_weight" validate:"gte=0,decimal2"`
	ViceralFat    float64 `json:"viceral_fat" validate:"gte=0,decimal2"`
	FatPercentage float64 `json:"fat_percentage" validate:"gte=0,decimal2"`
}
