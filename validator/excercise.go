package validator

// ExerciseRequest
type ExerciseRequest struct {
	Page int `json:"page" validate:"omitempty,gte=1"`
}
type ExcerciseMutationRequest struct {
	Minute  *int   `json:"minute,omitempty" db:"minute" validate:"omitempty,gt=1"`
	Caloric int    `json:"caloric" validate:"required,gt=1"`
	Type    string `json:"type" db:"type" validate:"required,oneof=HIT WeightLifting Cardio"`
}
