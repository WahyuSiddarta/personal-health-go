package models

import "fmt"

type UserTarget struct {
	TargetId                    int     `json:"target_id" db:"target_id"`
	UserId                      int     `json:"user_id" db:"user_id"`
	NutritionCaloric            float64 `json:"nutrition_caloric" db:"nutrition_caloric"`
	NutritionProtein            float64 `json:"nutrition_protein" db:"nutrition_protein"`
	NutritionCarbs              float64 `json:"nutrition_carbohydrate" db:"nutrition_carbohydrate"`
	NutritionFat                float64 `json:"nutrition_fat" db:"nutrition_fat"`
	BodyWeight                  float64 `json:"bodyweight" db:"bodyweight"`
	ViceralFat                  float64 `json:"viceral_fat" db:"viceral_fat"`
	FatPercentage               float64 `json:"fat_percentage" db:"fat_percentage"`
	WeeklyExerciseMinutes       int     `json:"weekly_exercise_minutes" db:"weekly_exercise_minutes"`
	WeeklyExcerciseSessions     int     `json:"weekly_exercise_sessions" db:"weekly_exercise_sessions"`
	WeeklyExcerciseCaloric      int     `json:"weekly_exercise_caloric" db:"weekly_exercise_caloric"`
	WeeklyWeightLiftingSessions int     `json:"weekly_weight_lifting_sessions" db:"weekly_weight_lifting_sessions"`
	WeeklyCardioMinutes         int     `json:"weekly_cardio_minutes" db:"weekly_cardio_minutes"`
}

// userRepository implements UserRepository interface
type userRepository struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindPersonalTarget(userID int) (*UserTarget, error)
	UpdatePersonalNutritionTarget(userTarget *UserTarget) error
	UpdatePersonalBodyMeasurementTarget(userTarget *UserTarget) error
	UpdatePersonalExerciseTarget(userTarget *UserTarget) error
}

// FindPersonalTarget retrieves the personal target for a given user ID
func (r *userRepository) FindPersonalTarget(userID int) (*UserTarget, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var user UserTarget
	query := `SELECT 
	target_id, user_id, nutrition_caloric, nutrition_protein, nutrition_carbohydrate,
	weekly_exercise_minutes, weekly_exercise_sessions, weekly_exercise_caloric, weekly_weight_lifting_sessions, weekly_cardio_minutes,
	nutrition_fat, bodyweight, viceral_fat, fat_percentage FROM users_target WHERE user_id = $1`

	err := db.Get(&user, query, userID)
	return &user, err
}

// UpdatePersonalBodyMeasurementTarget updates the personal target for a user
func (r *userRepository) UpdatePersonalBodyMeasurementTarget(userTarget *UserTarget) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE users_target SET 
	bodyweight = $1, viceral_fat = $2, fat_percentage = $3
	WHERE user_id = $4`
	Logger.Debug().Msgf("Executing query: %s with values %+v", query, userTarget)
	_, err := db.Exec(query, userTarget.BodyWeight, userTarget.ViceralFat, userTarget.FatPercentage, userTarget.UserId)
	return err
}

// UpdatePersonalTarget updates the personal target for a user
func (r *userRepository) UpdatePersonalNutritionTarget(userTarget *UserTarget) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE users_target SET 
	nutrition_caloric = $1, nutrition_protein = $2, nutrition_carbohydrate = $3,
	nutrition_fat = $4 WHERE user_id = $5`

	_, err := db.Exec(query, userTarget.NutritionCaloric, userTarget.NutritionProtein,
		userTarget.NutritionCarbs, userTarget.NutritionFat, userTarget.UserId)
	return err
}

// UpdatePersonalExerciseTarget updates the personal target for a user
func (r *userRepository) UpdatePersonalExerciseTarget(userTarget *UserTarget) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE users_target SET 
	weekly_exercise_minutes = $1, weekly_exercise_sessions = $2, weekly_exercise_caloric = $3, weekly_weight_lifting_sessions = $4, weekly_cardio_minutes = $5
	WHERE user_id = $6`
	Logger.Debug().Msgf("Executing query: %s with values %+v", query, userTarget)
	_, err := db.Exec(query, userTarget.WeeklyExerciseMinutes, userTarget.WeeklyExcerciseSessions, userTarget.WeeklyExcerciseCaloric, userTarget.WeeklyWeightLiftingSessions, userTarget.WeeklyCardioMinutes, userTarget.UserId)
	return err
}
