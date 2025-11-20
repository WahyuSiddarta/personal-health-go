package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PortfolioBond represents a bond holding within a user's portfolio.
type PortfolioBond struct {
	ID                      int        `db:"id" json:"id"`
	BondId                  string     `db:"bond_id" json:"bond_id"`
	UserID                  int        `db:"user_id" json:"user_id"`
	Name                    string     `db:"name" json:"name"`
	PurchasePrice           float64    `db:"purchase_price" json:"purchase_price"`
	CouponRate              float64    `db:"coupon_rate" json:"coupon_rate"`
	CouponFrequency         string     `db:"coupon_frequency" json:"coupon_frequency"`
	NextCouponDate          *time.Time `db:"next_coupon_date" json:"next_coupon_date"`
	MaturityDate            *time.Time `db:"maturity_date" json:"maturity_date"`
	Quantity                int        `db:"quantity" json:"quantity"`
	Status                  string     `db:"status" json:"status"`
	Note                    *string    `db:"note" json:"note"`
	MarketPrice             *float64   `db:"market_price" json:"market_price"`
	MarketPriceOverride     *float64   `db:"market_price_override" json:"market_price_override"`
	MarketPriceOverrideDate *time.Time `db:"market_price_override_date" json:"market_price_override_date"`
	SecondaryMarket         bool       `db:"secondary_market" json:"secondary_market"`
	CreatedAt               time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt               *time.Time `db:"deleted_at" json:"-"`
}

// PortfolioBondCreateRequest captures fields needed when creating a bond.
type PortfolioBondCreateRequest struct {
	BondId          string     `db:"bond_id" json:"bond_id"`
	Name            *string    `json:"name" validate:"omitempty"`
	PurchasePrice   float64    `json:"purchase_price" validate:"required,min=0"`
	CouponRate      float64    `json:"coupon_rate" validate:"required,min=0"`
	CouponFrequency string     `json:"coupon_frequency" validate:"required,oneof=monthly quarterly semi-annual annual"`
	NextCouponDate  *time.Time `json:"next_coupon_date" validate:"omitempty,datetime=2006-01-02"`
	MaturityDate    time.Time  `json:"maturity_date" validate:"required,datetime=2006-01-02"`
	Quantity        int        `json:"quantity" validate:"required,min=1"`
	SecondaryMarket bool       `json:"secondary_market" validate:"required,boolean"`
	Note            *string    `json:"note"`
}

// PortfolioBondUpdateRequest captures partial updates to an existing bond.
type PortfolioBondUpdateRequest struct {
	ISIN            *string    `json:"isin"`
	Name            *string    `json:"name"`
	PurchasePrice   *float64   `json:"purchase_price"`
	FaceValue       *float64   `json:"face_value"`
	CouponRate      *float64   `json:"coupon_rate"`
	CouponFrequency *string    `json:"coupon_frequency"`
	NextCouponDate  *time.Time `json:"next_coupon_date"`
	MaturityDate    *time.Time `json:"maturity_date"`
	Quantity        *int       `json:"quantity"`
	Issuer          *string    `json:"issuer"`
	Note            *string    `json:"note"`
	Status          *string    `json:"status"`
	MarketPrice     *float64   `json:"market_price"`
}

// PortfolioBondWithPotentialGain enriches a bond with matrix calculations.
type PortfolioBondWithPotentialGain struct {
	PortfolioBond
	TotalCouponsReceived float64 `json:"total_coupons_received"`
	PotentialGain        float64 `json:"potential_gain"`
	MarketPriceType      string  `json:"market_price_type"`
}

// PortfolioBondCoupon represents coupon payment records for bonds.
type PortfolioBondCoupon struct {
	ID              int        `db:"id" json:"id"`
	UserID          int        `db:"user_id" json:"user_id"`
	PortfolioBondID int        `db:"portfolio_bond_id" json:"portfolio_bond_id"`
	CouponNumber    int        `db:"coupon_number" json:"coupon_number"`
	PaymentDate     time.Time  `db:"payment_date" json:"payment_date"`
	Amount          float64    `db:"amount" json:"amount"`
	Status          string     `db:"status" json:"status"`
	Note            *string    `db:"note" json:"note"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

// PortfolioBondCouponCreateRequest defines the payload for new coupons.
type PortfolioBondCouponCreateRequest struct {
	PortfolioBondID int       `json:"portfolio_bond_id"`
	CouponNumber    int       `json:"coupon_number"`
	PaymentDate     time.Time `json:"payment_date"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	Note            *string   `json:"note"`
}

// PortfolioBondCouponUpdateRequest allows partial updates to coupon records.
type PortfolioBondCouponUpdateRequest struct {
	CouponNumber *int       `json:"coupon_number"`
	PaymentDate  *time.Time `json:"payment_date"`
	Amount       *float64   `json:"amount"`
	Status       *string    `json:"status"`
	Note         *string    `json:"note"`
}

// PortfolioBondRealized represents bonds that have been realized/closed.
type PortfolioBondRealized struct {
	ID                   int        `db:"id" json:"id"`
	UserID               int        `db:"user_id" json:"user_id"`
	PortfolioBondID      int        `db:"portfolio_bond_id" json:"portfolio_bond_id"`
	RealizedPrice        float64    `db:"realized_price" json:"realized_price"`
	TotalCouponsReceived float64    `db:"total_coupons_received" json:"total_coupons_received"`
	RealizedDate         time.Time  `db:"realized_date" json:"realized_date"`
	Note                 *string    `db:"note" json:"note"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt            *time.Time `db:"deleted_at" json:"deleted_at"`
}

// PortfolioBondRealizedCreateRequest captures required fields for realized bonds.
type PortfolioBondRealizedCreateRequest struct {
	PortfolioBondID      int        `json:"portfolio_bond_id"`
	RealizedPrice        float64    `json:"realized_price"`
	TotalCouponsReceived *float64   `json:"total_coupons_received"`
	RealizedDate         *time.Time `json:"realized_date"`
	Note                 *string    `json:"note"`
}

// PortfolioBondRealizedUpdateRequest allows partial updates of realized entries.
type PortfolioBondRealizedUpdateRequest struct {
	RealizedPrice        *float64   `json:"realized_price"`
	TotalCouponsReceived *float64   `json:"total_coupons_received"`
	RealizedDate         *time.Time `json:"realized_date"`
	Note                 *string    `json:"note"`
}

// PortfolioBondRepository defines all operations over bond portfolios.
type PortfolioBondRepository interface {
	Create(userID int, payload PortfolioBondCreateRequest) (*PortfolioBond, error)
	FindByUserID(userID int) ([]*PortfolioBond, error)
	FindByID(id int, userID int) (*PortfolioBond, error)
	Update(id int, userID int, payload PortfolioBondUpdateRequest) (*PortfolioBond, error)
	Delete(id int, userID int) error

	UpdateMarketPriceOverride(id int, userID int, marketPrice float64) (*PortfolioBond, error)
	FindByUserIDWithPotentialGain(userID int) ([]*PortfolioBondWithPotentialGain, error)
}

// PortfolioBondCouponRepository defines operations for bond coupons.
type PortfolioBondCouponRepository interface {
	Create(userID int, payload PortfolioBondCouponCreateRequest) (*PortfolioBondCoupon, error)
	FindByPortfolioBondID(userID int, portfolioBondID int) ([]*PortfolioBondCoupon, error)
	FindByID(id int, userID int) (*PortfolioBondCoupon, error)
	Update(id int, userID int, payload PortfolioBondCouponUpdateRequest) (*PortfolioBondCoupon, error)
	Delete(id int, userID int) error
	GetTotalReceived(portfolioBondID int, userID int) (float64, error)
}

// PortfolioBondRealizedRepository defines operations for realized bonds.
type PortfolioBondRealizedRepository interface {
	Create(userID int, payload PortfolioBondRealizedCreateRequest) (*PortfolioBondRealized, error)
	FindByUserID(userID int, limit, offset int) ([]*PortfolioBondRealized, error)
	FindByPortfolioBondID(userID int, portfolioBondID int, limit, offset int) ([]*PortfolioBondRealized, error)
	FindByID(id int, userID int) (*PortfolioBondRealized, error)
	Update(id int, userID int, payload PortfolioBondRealizedUpdateRequest) (*PortfolioBondRealized, error)
	Delete(id int, userID int) error
}

type portfolioBondRepository struct{}

func NewPortfolioBondRepository() PortfolioBondRepository {
	return &portfolioBondRepository{}
}

func (r *portfolioBondRepository) getDB() (*sqlx.DB, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db, nil
}

func (r *portfolioBondRepository) Create(userID int, payload PortfolioBondCreateRequest) (*PortfolioBond, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const insertQuery = `
	INSERT INTO portfolio_bond (
	bond_id, user_id, name, purchase_price, coupon_rate,
  coupon_frequency, next_coupon_date, maturity_date, quantity, status,
  note, secondary_market, created_at, updated_at, deleted_at )
	VALUES (
	$1, $2, $3, $4, $5,
	$6, $7, $8, $9, 'active',
	$10, $11, NOW(), NOW(), NULL
	)
	RETURNING id, user_id, isin, name, purchase_price, face_value, coupon_rate,
		coupon_frequency, next_coupon_date, maturity_date, quantity, issuer,
		note, status, market_price, market_price_override, market_price_override_date,
		created_at, updated_at, deleted_at`

	var bond PortfolioBond
	err = db.Get(&bond, insertQuery,
		payload.BondId,
		userID,
		payload.Name,
		payload.PurchasePrice,
		payload.CouponRate,
		payload.CouponFrequency,
		payload.NextCouponDate,
		payload.MaturityDate,
		payload.Quantity,
		payload.Note,
		payload.Note,
		payload.SecondaryMarket,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBond.Create] Error creating bond entry")
		return nil, fmt.Errorf("kesalahan membuat portfolio obligasi: %w", err)
	}

	return &bond, nil
}

func (r *portfolioBondRepository) FindByUserID(userID int) ([]*PortfolioBond, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT 
		a.id, a.bond_id, a.user_id, a.name, a.purchase_price,
		a.coupon_rate, a.coupon_frequency, a.next_coupon_date, a.maturity_date, a.quantity, 
		a.status, a.note, b.market_price, a.market_price_override, a.market_price_override_date,
		a.created_at, a.updated_at, a.deleted_at, a.secondary_market 
	FROM portfolio_bond a LEFT JOIN bond_tracker b ON a.bond_id = b.bond_id 
	WHERE a.user_id = $1 AND a.deleted_at IS NULL AND b.deleted_at IS NULL 
	ORDER BY a.created_at DESC`

	var bonds []*PortfolioBond
	err = db.Select(&bonds, query, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBond.FindByUserID] Error querying bonds")
		return nil, fmt.Errorf("kesalahan mengambil portfolio obligasi: %w", err)
	}

	return bonds, nil
}

func (r *portfolioBondRepository) FindByID(id int, userID int) (*PortfolioBond, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, isin, name, purchase_price, face_value, coupon_rate,
		coupon_frequency, next_coupon_date, maturity_date, quantity, issuer,
		note, status, market_price, market_price_override, market_price_override_date,
		created_at, updated_at, deleted_at
	FROM portfolio_bond
	WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var bond PortfolioBond
	err = db.Get(&bond, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBond.FindByID] Error querying bond")
		return nil, fmt.Errorf("kesalahan mengambil portfolio obligasi: %w", err)
	}

	return &bond, nil
}

func (r *portfolioBondRepository) Update(id int, userID int, payload PortfolioBondUpdateRequest) (*PortfolioBond, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	UPDATE portfolio_bond
	SET isin = COALESCE(NULLIF($1, ''), isin),
		name = COALESCE(NULLIF($2, ''), name),
		purchase_price = COALESCE($3, purchase_price),
		face_value = COALESCE($4, face_value),
		coupon_rate = COALESCE($5, coupon_rate),
		coupon_frequency = COALESCE(NULLIF($6, ''), coupon_frequency),
		next_coupon_date = COALESCE($7, next_coupon_date),
		maturity_date = COALESCE($8, maturity_date),
		quantity = COALESCE($9, quantity),
		issuer = COALESCE(NULLIF($10, ''), issuer),
		note = COALESCE($11, note),
		status = COALESCE(NULLIF($12, ''), status),
		market_price = COALESCE($13, market_price),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $14 AND user_id = $15 AND deleted_at IS NULL
	RETURNING id, user_id, isin, name, purchase_price, face_value, coupon_rate,
		coupon_frequency, next_coupon_date, maturity_date, quantity, issuer,
		note, status, market_price, market_price_override, market_price_override_date,
		created_at, updated_at, deleted_at`

	var bond PortfolioBond
	err = db.Get(&bond, query,
		payload.ISIN,
		payload.Name,
		payload.PurchasePrice,
		payload.FaceValue,
		payload.CouponRate,
		payload.CouponFrequency,
		payload.NextCouponDate,
		payload.MaturityDate,
		payload.Quantity,
		payload.Issuer,
		payload.Note,
		payload.Status,
		payload.MarketPrice,
		id,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBond.Update] Error updating bond")
		return nil, fmt.Errorf("kesalahan memperbarui portfolio obligasi: %w", err)
	}

	return &bond, nil
}

func (r *portfolioBondRepository) Delete(id int, userID int) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	const query = `
	DELETE FROM portfolio_bond
	WHERE id = $1 AND user_id = $2
	RETURNING id`

	var deletedID int
	err = db.Get(&deletedID, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBond.Delete] Error deleting bond")
		return fmt.Errorf("kesalahan menghapus portfolio obligasi: %w", err)
	}

	return nil
}

func (r *portfolioBondRepository) UpdateMarketPriceOverride(id int, userID int, marketPrice float64) (*PortfolioBond, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	UPDATE portfolio_bond
	SET 
		market_price_override = $1,
		market_price_override_date = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	RETURNING id, user_id, isin, name, purchase_price, face_value, coupon_rate,
		coupon_frequency, next_coupon_date, maturity_date, quantity, issuer,
		note, status, market_price, market_price_override, market_price_override_date,
		created_at, updated_at, deleted_at`

	var bond PortfolioBond
	err = db.Get(&bond, query, marketPrice, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBond.UpdateMarketPriceOverride] Error updating override")
		return nil, fmt.Errorf("kesalahan memperbarui harga pasar obligasi: %w", err)
	}

	return &bond, nil
}

func (r *portfolioBondRepository) FindByUserIDWithPotentialGain(userID int) ([]*PortfolioBondWithPotentialGain, error) {
	bonds, err := r.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(bonds) == 0 {
		return []*PortfolioBondWithPotentialGain{}, nil
	}

	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const couponQuery = `
	SELECT portfolio_bond_id, COALESCE(SUM(amount), 0) AS total_coupons
	FROM portfolio_bond_coupons
	WHERE user_id = $1 AND status = 'received'
	GROUP BY portfolio_bond_id`

	rows, err := db.Queryx(couponQuery, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBond.FindByUserIDWithPotentialGain] Coupon aggregation failed")
		return nil, fmt.Errorf("kesalahan mengambil kupon obligasi: %w", err)
	}
	defer rows.Close()

	couponMap := make(map[int]float64)
	for rows.Next() {
		var id int
		var total float64
		if err := rows.Scan(&id, &total); err != nil {
			Logger.Error().Err(err).Msg("[PortfolioBond.FindByUserIDWithPotentialGain] Coupon scan failed")
			return nil, fmt.Errorf("kesalahan membaca kupon obligasi: %w", err)
		}
		couponMap[id] = total
	}

	var enriched []*PortfolioBondWithPotentialGain
	for _, bond := range bonds {
		marketPrice := 0.0
		if bond.MarketPrice != nil {
			marketPrice = *bond.MarketPrice
		}
		quantity := float64(bond.Quantity)
		if quantity == 0 {
			quantity = 1
		}
		totalCoupons := couponMap[bond.ID]
		potentialGain := marketPrice*quantity + totalCoupons - bond.PurchasePrice
		marketPriceType := "market_tracking"
		if bond.MarketPriceOverride != nil {
			marketPriceType = "user_override"
		}
		enriched = append(enriched, &PortfolioBondWithPotentialGain{
			PortfolioBond:        *bond,
			TotalCouponsReceived: totalCoupons,
			PotentialGain:        potentialGain,
			MarketPriceType:      marketPriceType,
		})
	}

	return enriched, nil
}

type portfolioBondCouponRepository struct{}

func NewPortfolioBondCouponRepository() PortfolioBondCouponRepository {
	return &portfolioBondCouponRepository{}
}

func (r *portfolioBondCouponRepository) getDB() (*sqlx.DB, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db, nil
}

func (r *portfolioBondCouponRepository) Create(userID int, payload PortfolioBondCouponCreateRequest) (*PortfolioBondCoupon, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	if payload.Status == "" {
		payload.Status = "pending"
	}

	const query = `
	INSERT INTO portfolio_bond_coupons (
		portfolio_bond_id, user_id, coupon_number, payment_date, amount, status, note
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, user_id, portfolio_bond_id, coupon_number, payment_date, amount, status, note, created_at, updated_at, deleted_at`

	var coupon PortfolioBondCoupon
	err = db.Get(&coupon, query,
		payload.PortfolioBondID,
		userID,
		payload.CouponNumber,
		payload.PaymentDate,
		payload.Amount,
		payload.Status,
		payload.Note,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.Create] Error inserting coupon")
		return nil, fmt.Errorf("kesalahan membuat kupon obligasi: %w", err)
	}

	return &coupon, nil
}

func (r *portfolioBondCouponRepository) FindByPortfolioBondID(userID int, portfolioBondID int) ([]*PortfolioBondCoupon, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, portfolio_bond_id, coupon_number, payment_date, amount, status, note, created_at, updated_at, deleted_at
	FROM portfolio_bond_coupons
	WHERE portfolio_bond_id = $1 AND user_id = $2 AND deleted_at IS NULL
	ORDER BY payment_date DESC`

	var coupons []*PortfolioBondCoupon
	err = db.Select(&coupons, query, portfolioBondID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.FindByPortfolioBondID] Error querying coupons")
		return nil, fmt.Errorf("kesalahan mengambil kupon obligasi: %w", err)
	}

	return coupons, nil
}

func (r *portfolioBondCouponRepository) FindByID(id int, userID int) (*PortfolioBondCoupon, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, portfolio_bond_id, coupon_number, payment_date, amount, status, note, created_at, updated_at, deleted_at
	FROM portfolio_bond_coupons
	WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var coupon PortfolioBondCoupon
	err = db.Get(&coupon, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.FindByID] Error querying coupon")
		return nil, fmt.Errorf("kesalahan mengambil kupon obligasi: %w", err)
	}

	return &coupon, nil
}

func (r *portfolioBondCouponRepository) Update(id int, userID int, payload PortfolioBondCouponUpdateRequest) (*PortfolioBondCoupon, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	UPDATE portfolio_bond_coupons
	SET coupon_number = COALESCE($1, coupon_number),
		payment_date = COALESCE($2, payment_date),
		amount = COALESCE($3, amount),
		status = COALESCE(NULLIF($4, ''), status),
		note = COALESCE($5, note),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6 AND user_id = $7 AND deleted_at IS NULL
	RETURNING id, user_id, portfolio_bond_id, coupon_number, payment_date, amount, status, note, created_at, updated_at, deleted_at`

	var coupon PortfolioBondCoupon
	err = db.Get(&coupon, query,
		payload.CouponNumber,
		payload.PaymentDate,
		payload.Amount,
		payload.Status,
		payload.Note,
		id,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.Update] Error updating coupon")
		return nil, fmt.Errorf("kesalahan memperbarui kupon obligasi: %w", err)
	}

	return &coupon, nil
}

func (r *portfolioBondCouponRepository) Delete(id int, userID int) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	const query = `
	DELETE FROM portfolio_bond_coupons
	WHERE id = $1 AND user_id = $2
	RETURNING id`

	var deletedID int
	err = db.Get(&deletedID, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.Delete] Error deleting coupon")
		return fmt.Errorf("kesalahan menghapus kupon obligasi: %w", err)
	}

	return nil
}

func (r *portfolioBondCouponRepository) GetTotalReceived(portfolioBondID int, userID int) (float64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	const query = `
	SELECT COALESCE(SUM(amount), 0) AS total
	FROM portfolio_bond_coupons
	WHERE portfolio_bond_id = $1 AND user_id = $2 AND status = 'received'`

	var total float64
	err = db.Get(&total, query, portfolioBondID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondCoupon.GetTotalReceived] Error summing coupons")
		return 0, fmt.Errorf("kesalahan menghitung total kupon: %w", err)
	}

	return total, nil
}

type portfolioBondRealizedRepository struct{}

func NewPortfolioBondRealizedRepository() PortfolioBondRealizedRepository {
	return &portfolioBondRealizedRepository{}
}

func (r *portfolioBondRealizedRepository) getDB() (*sqlx.DB, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db, nil
}

func (r *portfolioBondRealizedRepository) Create(userID int, payload PortfolioBondRealizedCreateRequest) (*PortfolioBondRealized, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	realizedDate := time.Now()
	if payload.RealizedDate != nil {
		realizedDate = *payload.RealizedDate
	}
	totalCoupons := 0.0
	if payload.TotalCouponsReceived != nil {
		totalCoupons = *payload.TotalCouponsReceived
	}

	const query = `
	INSERT INTO portfolio_bond_realized (
		user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note, created_at, updated_at, deleted_at`

	var realized PortfolioBondRealized
	err = db.Get(&realized, query,
		userID,
		payload.PortfolioBondID,
		payload.RealizedPrice,
		totalCoupons,
		realizedDate,
		payload.Note,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.Create] Error inserting realized bond")
		return nil, fmt.Errorf("kesalahan membuat entry obligasi terealisasi: %w", err)
	}

	return &realized, nil
}

func (r *portfolioBondRealizedRepository) FindByUserID(userID int, limit, offset int) ([]*PortfolioBondRealized, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note, created_at, updated_at, deleted_at
	FROM portfolio_bond_realized
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY realized_date DESC
	LIMIT $2 OFFSET $3`

	var realizations []*PortfolioBondRealized
	err = db.Select(&realizations, query, userID, limit+1, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.FindByUserID] Error querying realized bonds")
		return nil, fmt.Errorf("kesalahan mengambil data obligasi terealisasi: %w", err)
	}

	return realizations, nil
}

func (r *portfolioBondRealizedRepository) FindByPortfolioBondID(userID int, portfolioBondID int, limit, offset int) ([]*PortfolioBondRealized, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note, created_at, updated_at, deleted_at
	FROM portfolio_bond_realized
	WHERE user_id = $1 AND portfolio_bond_id = $2 AND deleted_at IS NULL
	ORDER BY realized_date DESC
	LIMIT $3 OFFSET $4`

	var realizations []*PortfolioBondRealized
	err = db.Select(&realizations, query, userID, portfolioBondID, limit+1, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.FindByPortfolioBondID] Error querying realized bonds")
		return nil, fmt.Errorf("kesalahan mengambil data obligasi terealisasi: %w", err)
	}

	return realizations, nil
}

func (r *portfolioBondRealizedRepository) FindByID(id int, userID int) (*PortfolioBondRealized, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	SELECT id, user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note, created_at, updated_at, deleted_at
	FROM portfolio_bond_realized
	WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var realized PortfolioBondRealized
	err = db.Get(&realized, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.FindByID] Error querying realized bond")
		return nil, fmt.Errorf("kesalahan mengambil data obligasi terealisasi: %w", err)
	}

	return &realized, nil
}

func (r *portfolioBondRealizedRepository) Update(id int, userID int, payload PortfolioBondRealizedUpdateRequest) (*PortfolioBondRealized, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `
	UPDATE portfolio_bond_realized
	SET realized_price = COALESCE($1, realized_price),
		total_coupons_received = COALESCE($2, total_coupons_received),
		realized_date = COALESCE($3, realized_date),
		note = COALESCE($4, note),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $5 AND user_id = $6 AND deleted_at IS NULL
	RETURNING id, user_id, portfolio_bond_id, realized_price, total_coupons_received, realized_date, note, created_at, updated_at, deleted_at`

	var realized PortfolioBondRealized
	err = db.Get(&realized, query,
		payload.RealizedPrice,
		payload.TotalCouponsReceived,
		payload.RealizedDate,
		payload.Note,
		id,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.Update] Error updating realized bond")
		return nil, fmt.Errorf("kesalahan memperbarui data obligasi terealisasi: %w", err)
	}

	return &realized, nil
}

func (r *portfolioBondRealizedRepository) Delete(id int, userID int) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	const query = `
	DELETE FROM portfolio_bond_realized
	WHERE id = $1 AND user_id = $2
	RETURNING id`

	var deletedID int
	err = db.Get(&deletedID, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		Logger.Error().Err(err).Msg("[PortfolioBondRealized.Delete] Error deleting realized bond")
		return fmt.Errorf("kesalahan menghapus data obligasi terealisasi: %w", err)
	}

	return nil
}
