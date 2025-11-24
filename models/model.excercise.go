package models

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
)

type ExcerciseRecord struct {
	ExcerciseId int       `json:"excercise_id" db:"excercise_id"`
	UserId      int       `json:"user_id" db:"user_id"`
	RecordAt    time.Time `json:"record_at" db:"record_at"`
	Minute      *int      `json:"minute,omitempty" db:"minute"`
	Caloric     int       `json:"caloric" db:"caloric"`
	Type        string    `json:"type" db:"type"`
}

// MarshalJSON : Overloads BodyMeasurement
func (a ExcerciseRecord) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(struct {
		ExcerciseId int    `json:"excercise_id" db:"excercise_id"`
		UserId      int    `json:"user_id" db:"user_id"`
		RecordAt    string `json:"record_at" db:"record_at"`
		Minute      *int   `json:"minute,omitempty" db:"minute"`
		Caloric     int    `json:"caloric" db:"caloric"`
		Type        string `json:"type" db:"type"`
	}{
		ExcerciseId: a.ExcerciseId,
		UserId:      a.UserId,
		RecordAt:    a.RecordAt.Format(time.RFC3339),
		Minute:      a.Minute,
		Caloric:     a.Caloric,
		Type:        a.Type,
	})
}

// excerciseRecord implements excerciseRecord interface
type excerciseRecordRepository struct{}

// NewexcerciseRecord creates a new excerciseRecord repository
func NewexcerciseRecordRepository() ExcerciseRecordRepository {
	return &excerciseRecordRepository{}
}

// excerciseRecord defines the interface for user data operations
type ExcerciseRecordRepository interface {
	Create(userId int, data ExcerciseRecord) error
	GetByUserId(userId, limit, page int) ([]ExcerciseRecord, error)
	Update(userId int, excerciseId int, data *ExcerciseRecord) error
	Delete(userId int, excerciseId int) error
}

// Delete ExcerciseRecord for a user
func (r *excerciseRecordRepository) Delete(userID, excerciseId int) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE excercise_record SET deleted_at = NOW() 
	WHERE user_id = $1 
	AND excercise_id = $2`
	_, err := db.Exec(query, userID, excerciseId)
	return err
}

// Update ExcerciseRecord for a user
func (r *excerciseRecordRepository) Update(userId int, excerciseId int, data *ExcerciseRecord) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE excercise_record SET
	minute = $1, caloric = $2, type = $3
	WHERE user_id = $4 AND excercise_id = $5`

	_, err := db.Exec(query, data.Minute, data.Caloric, data.Type, userId, excerciseId)
	return err
}

// AddTodayIntake adds today's food intake for a user
func (r *excerciseRecordRepository) Create(userId int, data ExcerciseRecord) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `INSERT INTO excercise_record
	(user_id, minute, caloric, type, record_at)
	VALUES ($1, $2, $3, $4, NOW())`

	_, err := db.Exec(query,
		userId,
		data.Minute, data.Caloric,
		data.Type)

	return err
}

func (r *excerciseRecordRepository) GetByUserId(userID, limit, page int) ([]ExcerciseRecord, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	offset := (page - 1) * limit
	limit = limit + 1

	var records []ExcerciseRecord
	query := `SELECT 
	excercise_id, user_id, minute, caloric, type, record_at 
 FROM excercise_record 
 WHERE user_id = $1 AND is_deleted = null ORDER BY record_at DESC LIMIT $2 OFFSET $3`
	err := db.Select(&records, query, userID, limit, offset)
	return records, err
}
