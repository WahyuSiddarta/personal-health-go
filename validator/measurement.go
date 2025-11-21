package validator

// BodyMeasurementRequest
type BodyMeasurementRequest struct {
	Page int `json:"page" validate:"omitempty,gte=1"`
}

type BodyMeasurementCreateRequest struct {
	Bodyweight    float64 `json:"bodyweight" validate:"required,gt=0,decimal2"`
	ViceralFat    float64 `json:"viceral_fat,omitempty" validate:"omitempty,gt=0,decimal2"`
	FatPercentage float64 `json:"fat_percentage,omitempty" validate:"omitempty,gt=0,decimal2"`
	NickCm        float64 `json:"nick_cm,omitempty" validate:"omitempty,gt=0,decimal2"`
	WaistCm       float64 `json:"waist_cm,omitempty" validate:"omitempty,gt=0,decimal2"`
}
