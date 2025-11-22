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

type NutritionChartData struct {
	UserId       int     `json:"user_id" db:"user_id"`
	Fat          float64 `json:"total_fat" db:"total_fat"`
	Protein      float64 `json:"total_protein" db:"total_protein"`
	Carbohydrate float64 `json:"total_carbohydrate" db:"total_carbohydrate"`
	Caloric      float64 `json:"total_caloric" db:"total_caloric"`
	Period       string  `json:"time_slice_label" db:"time_slice_label"`
}

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

	GetNutritionChartData(userID int) ([]NutritionChartData, error)
}

// / Overview Nutrition Handlers
func (r *nutritionRepository) GetNutritionChartData(userID int) ([]NutritionChartData, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var query string = `SELECT
    TO_CHAR(time_slice_date, 'YYYY-MM-DD') AS time_slice_label,
    SUM(fat) AS total_fat,
    SUM(protein) AS total_protein,
    SUM(carbohydrate) AS total_carbohydrate,
    SUM(caloric) AS total_caloric
FROM (
    SELECT
        DATE(created_at AT TIME ZONE 'Asia/Makassar') AS time_slice_date,
        fat,
        protein,
        carbohydrate,
        caloric
    FROM public.users_food_intake
    WHERE user_id = $1
      AND created_at AT TIME ZONE 'Asia/Makassar' >= (NOW() AT TIME ZONE 'Asia/Makassar') - INTERVAL '366 days'
) AS T
GROUP BY time_slice_date
ORDER BY time_slice_date DESC`

	var chartData []NutritionChartData
	err := db.Select(&chartData, query, userID)
	return chartData, err
}

// / Daily Nutrition Intake Handlers
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
 AND created_at AT TIME ZONE 'Asia/Makassar' >= CURRENT_DATE AT TIME ZONE 'Asia/Makassar'
 AND created_at AT TIME ZONE 'Asia/Makassar' < (CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Makassar'`

	err := db.Select(&users, query, userID)
	return users, err
}
