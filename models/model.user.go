package models

import "fmt"

type UserTarget struct {
	TargetId         int     `json:"target_id" db:"target_id"`
	UserId           int     `json:"user_id" db:"user_id"`
	NutritionCaloric float64 `json:"nutrition_caloric" db:"nutrition_caloric"`
	NutritionProtein float64 `json:"nutrition_protein" db:"nutrition_protein"`
	NutritionCarbs   float64 `json:"nutrition_carbohydrate" db:"nutrition_carbohydrate"`
	NutritionFat     float64 `json:"nutrition_fat" db:"nutrition_fat"`
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
	nutrition_fat FROM users_target WHERE user_id = $1`

	err := db.Get(&user, query, userID)
	return &user, err
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
