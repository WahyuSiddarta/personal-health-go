package models

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
)

type BodyMeasurement struct {
	MeasurementId int       `json:"measurement_id" db:"measurement_id"`
	UserId        int       `json:"user_id" db:"user_id"`
	Bodyweight    float64   `json:"bodyweight" db:"bodyweight"`
	ViceralFat    *float64  `json:"viceral_fat,omitempty" db:"viceral_fat"`
	FatPercentage *float64  `json:"fat_percentage,omitempty" db:"fat_percentage"`
	NickCm        *float64  `json:"nick_cm,omitempty" db:"nick_cm"`
	WaistCm       *float64  `json:"waist_cm,omitempty" db:"waist_cm"`
	MeasuredAt    time.Time `json:"measured_at" db:"measured_at"`
}

// MarshalJSON : Overloads BodyMeasurement
func (a BodyMeasurement) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(struct {
		MeasurementId int      `json:"measurement_id" db:"measurement_id"`
		UserId        int      `json:"user_id" db:"user_id"`
		Bodyweight    float64  `json:"bodyweight" db:"bodyweight"`
		ViceralFat    *float64 `json:"viceral_fat,omitempty" db:"viceral_fat"`
		FatPercentage *float64 `json:"fat_percentage,omitempty" db:"fat_percentage"`
		NickCm        *float64 `json:"nick_cm,omitempty" db:"nick_cm"`
		WaistCm       *float64 `json:"waist_cm,omitempty" db:"waist_cm"`
		MeasuredAt    string   `json:"measured_at" db:"measured_at"`
	}{
		MeasurementId: a.MeasurementId,
		UserId:        a.UserId,
		Bodyweight:    a.Bodyweight,
		ViceralFat:    a.ViceralFat,
		FatPercentage: a.FatPercentage,
		NickCm:        a.NickCm,
		WaistCm:       a.WaistCm,
		MeasuredAt:    a.MeasuredAt.Format(time.RFC3339),
	})
}

// bodyMeasurement implements bodyMeasurement interface
type bodyMeasurementRepository struct{}

// NewBodyMeasurement creates a new bodyMeasurement repository
func NewBodyMeasurementRepository() BodyMeasurementRepository {
	return &bodyMeasurementRepository{}
}

// BodyMeasurement defines the interface for user data operations
type BodyMeasurementRepository interface {
	Create(userId int, data BodyMeasurement) error
	GetByUserId(userId, limit, page int) ([]BodyMeasurement, error)
	Update(userId int, measurementId int, data *BodyMeasurement) error
	Delete(userId int, measurementId int) error
}

// DeleteTodayIntake deletes today's food intake for a user
func (r *bodyMeasurementRepository) Delete(userID, measurementId int) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `DELETE FROM body_measurement 
	WHERE user_id = $1 
	AND measurement_id = $2`
	_, err := db.Exec(query, userID, measurementId)
	return err
}

// Update updates a body measurement for a user
func (r *bodyMeasurementRepository) Update(userId int, measurementId int, data *BodyMeasurement) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `UPDATE body_measurement SET
	bodyweight = $1, viceral_fat = $2, fat_percentage = $3,
	nick_cm = $4, waist_cm = $5 WHERE user_id = $6 AND measurement_id = $7`

	_, err := db.Exec(query, data.Bodyweight, data.ViceralFat,
		data.FatPercentage, data.NickCm,
		data.WaistCm, userId, measurementId)

	return err
}

// AddTodayIntake adds today's food intake for a user
func (r *bodyMeasurementRepository) Create(userId int, data BodyMeasurement) error {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := `INSERT INTO body_measurement
	(user_id, bodyweight, viceral_fat, fat_percentage, nick_cm, waist_cm, measured_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err := db.Exec(query,
		userId,
		data.Bodyweight, data.ViceralFat,
		data.FatPercentage, data.NickCm,
		data.WaistCm)
	return err
}

func (r *bodyMeasurementRepository) GetByUserId(userID, limit, page int) ([]BodyMeasurement, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	offset := (page - 1) * limit
	limit = limit + 1

	var measurements []BodyMeasurement
	query := `SELECT 
	measurement_id, user_id, bodyweight, viceral_fat,
	fat_percentage, nick_cm, waist_cm, measured_at 
 FROM body_measurement 
 WHERE user_id = $1 ORDER BY measured_at DESC LIMIT $2 OFFSET $3`

	err := db.Select(&measurements, query, userID, limit, offset)
	return measurements, err
}
