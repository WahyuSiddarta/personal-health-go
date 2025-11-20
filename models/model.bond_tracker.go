package models

import "time"

type BondTracker struct {
	BondId      string     `db:"bond_id"`
	MarketPrice float64    `db:"market_price"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
