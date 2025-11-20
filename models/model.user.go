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
	query := `INSERT INTO users (email, password, status, user_level) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, email, status, user_level, premium_expires_at, created_at, updated_at`

	err := db.Get(&user, query, r)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return nil, fmt.Errorf("user with this email already exists")
		}
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}
