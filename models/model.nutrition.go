package models

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
)

type NutritionTracker struct {
	UserId       int               `json:"user_id" db:"user_id"`
	FoodId       int               `json:"food_id" db:"food_id"`
	Category     NutritionCategory `json:"category" db:"category"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	Fat          float64           `json:"fat" db:"fat"`
	Protein      float64           `json:"protein" db:"protein"`
	Carbohydrate float64           `json:"carbohydrate" db:"carbohydrate"`
	Caloric      float64           `json:"caloric" db:"caloric"`
	Name         string            `json:"name" db:"name"`
}

// MarshalJSON : Overloads NutritionTracker
func (a NutritionTracker) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(struct {
		UserId       int               `json:"user_id"`
		FoodId       int               `json:"food_id"`
		Category     NutritionCategory `json:"category"`
		CreatedAt    string            `json:"created_at"`
		Fat          float64           `json:"fat"`
		Protein      float64           `json:"protein"`
		Carbohydrate float64           `json:"carbohydrate"`
		Caloric      float64           `json:"caloric"`
		Name         string            `json:"name"`
	}{
		UserId:       a.UserId,
		FoodId:       a.FoodId,
		Category:     a.Category,
		CreatedAt:    a.CreatedAt.Format(time.RFC3339),
		Fat:          a.Fat,
		Protein:      a.Protein,
		Carbohydrate: a.Carbohydrate,
		Caloric:      a.Caloric,
		Name:         a.Name,
	})
}

// NutritionCategory represents nutrition category, like breakfast, lunch, dinner, snack
type NutritionCategory string

const (
	NutritionCategoryBreakfast NutritionCategory = "breakfast"
	NutritionCategoryLunch     NutritionCategory = "lunch"
	NutritionCategoryDinner    NutritionCategory = "dinner"
	NutritionCategorySnack     NutritionCategory = "snack"
)

// nutritionRepository implements NutritionRepository interface
type nutritionRepository struct{}

// NewNutritionRepository creates a new nutrition repository
func NewNutritionRepository() NutritionRepository {
	return &nutritionRepository{}
}

// NutritionRepository defines the interface for user data operations
type NutritionRepository interface {
	DeleteTodayIntake(userID, foodId int) error
	UpdateTodayIntake(nutritionTracker *NutritionTracker) error
	AddTodayIntake(nutritionTracker *NutritionTracker) error
	FindUserTodayIntake(userID int) ([]NutritionTracker, error)
}

// DeleteTodayIntake deletes today's food intake for a user
func (r *nutritionRepository) DeleteTodayIntake(userID, foodId int) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `DELETE FROM users_food_intake 
	WHERE user_id = $1 
	AND food_id = $2`
	_, err := db.Exec(query, userID, foodId)
	return err
}

// UpdateTodayIntake updates today's food intake for a user
func (r *nutritionRepository) UpdateTodayIntake(nutritionTracker *NutritionTracker) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE users_food_intake SET 
	fat = $1, protein = $2, carbohydrate = $3,
	caloric = $4, name = $5 WHERE user_id = $6 AND food_id = $7`

	_, err := db.Exec(query, nutritionTracker.Fat, nutritionTracker.Protein,
		nutritionTracker.Carbohydrate, nutritionTracker.Caloric,
		nutritionTracker.Name, nutritionTracker.UserId, nutritionTracker.FoodId)
	return err
}

// AddTodayIntake adds today's food intake for a user
func (r *nutritionRepository) AddTodayIntake(nutritionTracker *NutritionTracker) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `INSERT INTO users_food_intake 
	(user_id, category, created_at, fat, protein, carbohydrate, caloric, name) 
	VALUES ($1, $2, NOW(), $3, $4, $5, $6, $7)`

	_, err := db.Exec(query,
		nutritionTracker.UserId,
		nutritionTracker.Category, nutritionTracker.Fat,
		nutritionTracker.Protein, nutritionTracker.Carbohydrate,
		nutritionTracker.Caloric, nutritionTracker.Name)
	return err
}

// FindUserTodayIntake retrieves the personal target for a given user ID
func (r *nutritionRepository) FindUserTodayIntake(userID int) ([]NutritionTracker, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var users []NutritionTracker
	query := `SELECT 
	user_id, food_id, category, created_at, fat,
	protein, carbohydrate, caloric, name 
 FROM users_food_intake 
 WHERE user_id = $1 
 AND created_at AT TIME ZONE 'Asia/Jakarta' >= CURRENT_DATE AT TIME ZONE 'Asia/Jakarta'
 AND created_at AT TIME ZONE 'Asia/Jakarta' < (CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Jakarta'`

	err := db.Select(&users, query, userID)
	return users, err
}
