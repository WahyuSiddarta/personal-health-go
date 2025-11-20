package models

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"
)

var (
	Logger *zerolog.Logger
	DBM    DBManager
)

// SQLTimeFormat : Format string for golang to output SQL standar time
const SQLTimeFormat = "2006-01-02 15:04:05"

// SQLDateFormat :
const SQLDateFormat = "2006-01-02"

// TimeRange :
type TimeRange struct {
	Valid bool
	Start time.Time
	End   time.Time
}

// JSONMap :
type JSONMap map[string]interface{}

// DBManager : sqlx db pointer
type DBManager struct {
	PostgreDBManager DBPointer
}

// GetDB : initialize DB -> pointer is transfered into model from database (localize)
func GetDB() DBManager {
	return DBM
}

type DBPointer struct {
	RW *sqlx.DB
	RC *sqlx.DB
}
