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
	BodyWeight    float64 `json:"bodyweight" validate:"gte=0,decimal2"`
	ViceralFat    float64 `json:"viceral_fat" validate:"gte=0,decimal2"`
	FatPercentage float64 `json:"fat_percentage" validate:"gte=0,decimal2"`
}

type PersonalExerciseTargetRequest struct {
	WeeklyExerciseMinutes       int `json:"weekly_exercise_minutes" validate:"gte=0"`
	WeeklyExcerciseSessions     int `json:"weekly_exercise_sessions" validate:"gte=0"`
	WeeklyExcerciseCaloric      int `json:"weekly_exercise_caloric" validate:"gte=0"`
	WeeklyWeightLiftingSessions int `json:"weekly_weight_lifting_sessions" validate:"gte=0"`
	WeeklyCardioMinutes         int `json:"weekly_cardio_minutes" validate:"gte=0"`
}
