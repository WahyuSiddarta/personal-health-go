package models

import "fmt"

type UserTarget struct {
	ID               int     `json:"id" db:"id"`
	UserId           int     `json:"user_id" db:"user_id"`
	NutritionCaloric float64 `json:"nutrition_caloric" db:"nutrition_caloric"`
	NutritionProtein float64 `json:"nutrition_protein" db:"nutrition_protein"`
	NutritionCarbs   float64 `json:"nutrition_carbs" db:"nutrition_carbs"`
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
}

// Create creates a new user
func (r *userRepository) FindPersonalTarget(userID int) (*UserTarget, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var user UserTarget
	query := `SELECT 
	id, user_id, nutrition_caloric, nutrition_protein, nutrition_carbs,
	nutrition_fat FROM user_targets WHERE user_id = $1`

	err := db.Get(&user, query, userID)
	return &user, err
}
