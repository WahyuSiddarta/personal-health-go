package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// PortfolioCash represents a cash portfolio entry
type PortfolioCash struct {
	ID                  int        `db:"id" json:"id"`
	UserID              int        `db:"user_id" json:"user_id"`
	Account             string     `db:"account" json:"account"`
	Bank                string     `db:"bank" json:"bank"`
	Amount              float64    `db:"amount" json:"amount"`
	YieldRate           *float64   `db:"yield_rate" json:"yield_rate"`
	YieldPeriod         string     `db:"yield_period" json:"yield_period"`
	YieldFrequencyType  string     `db:"yield_frequency_type" json:"yield_frequency_type"`
	YieldFrequencyValue int        `db:"yield_frequency_value" json:"yield_frequency_value"`
	YieldPaymentType    string     `db:"yield_payment_type" json:"yield_payment_type"`
	HasMaturity         bool       `db:"has_maturity" json:"has_maturity"`
	MaturityDate        *time.Time `db:"maturity_date" json:"maturity_date"`
	Note                *string    `db:"note" json:"note"`
	Status              string     `db:"status" json:"status"`
	Category            string     `db:"category" json:"category"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deleted_at"`
}

// PortfolioPnlRealizedCash represents a realized PnL entry for a cash portfolio
type PortfolioPnlRealizedCash struct {
	ID              int        `db:"id" json:"id"`
	UserID          int        `db:"user_id" json:"user_id"`
	PortfolioCashID int        `db:"portfolio_cash_id" json:"portfolio_cash_id"`
	Amount          float64    `db:"amount" json:"amount"`
	RealizedAt      time.Time  `db:"realized_at" json:"realized_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

// PnlSummary represents aggregation statistics for realized PnL entries
type PnlSummary struct {
	TotalAmount    float64    `db:"total_amount" json:"total_amount"`
	Count          int        `db:"count" json:"count"`
	AvgAmount      float64    `db:"avg_amount" json:"avg_amount"`
	MaxAmount      float64    `db:"max_amount" json:"max_amount"`
	MinAmount      float64    `db:"min_amount" json:"min_amount"`
	LastRealizedAt *time.Time `db:"last_realized_at" json:"last_realized_at"`
}

// MoveAssetResponse represents the response for moving assets
type MoveAssetResponse struct {
	Source      *PortfolioCash `json:"source"`
	Target      *PortfolioCash `json:"target"`
	MovedAmount float64        `json:"moved_amount"`
}

// RealizeCashPortfolioResponse represents the response for realizing a cash portfolio
type RealizeCashPortfolioResponse struct {
	Portfolio *PortfolioCash            `json:"portfolio"`
	PnL       *PortfolioPnlRealizedCash `json:"pnl"`
}

// PortfolioCashRepository defines all operations related to cash portfolios
type PortfolioCashRepository interface {
	Create(userID int, account string, bank string, amount float64, yieldRate *float64, yieldPeriod string,
		yieldFrequencyType string, yieldFrequencyValue int, yieldPaymentType string, hasMaturity bool,
		maturityDate *time.Time, note *string, category string) (*PortfolioCash, error)
	FindByUserID(userID int) ([]*PortfolioCash, error)
	FindByID(id int, userID int) (*PortfolioCash, error)
	Update(id int, userID int, account string, bank string, amount *float64, yieldRate *float64,
		yieldPeriod string, yieldFrequencyType string, yieldFrequencyValue *int, yieldPaymentType string,
		hasMaturity *bool, maturityDate *time.Time, note *string, status string, category string) (*PortfolioCash, error)
	Delete(id int, userID int) error
	MoveAsset(sourceID int, targetID int, userID int) (*MoveAssetResponse, error)

	RealizeCashPortfolio(userID int, portfolioID int, finalSaldo float64, pnlAmount float64, realizedAt time.Time) (*RealizeCashPortfolioResponse, error)
	CreatePnlEntry(userID, portfolioCashID int, amount float64, realizedAt time.Time) (*PortfolioPnlRealizedCash, error)
	FindPnlByUserID(userID int, limit, offset int) ([]*PortfolioPnlRealizedCash, error)
	FindPnlByPortfolioCashID(portfolioCashID, userID int, limit, offset int) ([]*PortfolioPnlRealizedCash, error)
	FindPnlByID(id, userID int) (*PortfolioPnlRealizedCash, error)
	UpdatePnlEntry(id, userID int, amount *float64, realizedAt *time.Time) (*PortfolioPnlRealizedCash, error)
	DeletePnlEntry(id, userID int) error
	GetPnlSummary(userID int) (*PnlSummary, error)
}

type portfolioCashRepository struct{}

// NewPortfolioCashRepository creates a new portfolio cash repository implementation
func NewPortfolioCashRepository() PortfolioCashRepository {
	return &portfolioCashRepository{}
}

func (r *portfolioCashRepository) getDB() (*sqlx.DB, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db, nil
}

// Create creates a new cash portfolio entry
func (r *portfolioCashRepository) Create(
	userID int,
	account string,
	bank string,
	amount float64,
	yieldRate *float64,
	yieldPeriod string,
	yieldFrequencyType string,
	yieldFrequencyValue int,
	yieldPaymentType string,
	hasMaturity bool,
	maturityDate *time.Time,
	note *string,
	category string,
) (*PortfolioCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	status := "active"
	if yieldPeriod == "" {
		yieldPeriod = "per_tahun"
	}

	query := `
		INSERT INTO portfolio_cash 
		(user_id, account, bank, amount, yield_rate, yield_period, 
		 yield_frequency_type, yield_frequency_value, yield_payment_type, 
		 has_maturity, maturity_date, note, status, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, user_id, account, bank, amount, yield_rate, yield_period, 
			  yield_frequency_type, yield_frequency_value, yield_payment_type, 
			  has_maturity, maturity_date, note, status, category, created_at, updated_at, deleted_at
	`

	var result PortfolioCash
	err = db.QueryRowx(query,
		userID,
		account,
		bank,
		amount,
		yieldRate,
		yieldPeriod,
		yieldFrequencyType,
		yieldFrequencyValue,
		yieldPaymentType,
		hasMaturity,
		maturityDate,
		note,
		status,
		category,
	).StructScan(&result)

	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioCash.Create] Error creating cash portfolio entry")
		return nil, fmt.Errorf("kesalahan membuat portfolio kas: %w", err)
	}

	return &result, nil
}

// FindByUserID retrieves all cash portfolios for a user (excluding soft-deleted)
func (r *portfolioCashRepository) FindByUserID(userID int) ([]*PortfolioCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, account, bank, amount, yield_rate, yield_period, 
		       yield_frequency_type, yield_frequency_value, yield_payment_type, 
		       has_maturity, maturity_date, note, status, category, created_at, updated_at, deleted_at
		FROM portfolio_cash 
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var portfolios []*PortfolioCash
	err = db.Select(&portfolios, query, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("Error finding cash portfolios by user ID")
		return nil, fmt.Errorf("kesalahan mengambil portfolio kas: %w", err)
	}

	return portfolios, nil
}

// FindByID retrieves a specific cash portfolio entry
func (r *portfolioCashRepository) FindByID(id int, userID int) (*PortfolioCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, account, bank, amount, yield_rate, yield_period, 
		       yield_frequency_type, yield_frequency_value, yield_payment_type, 
		       has_maturity, maturity_date, note, status, category, created_at, updated_at, deleted_at
		FROM portfolio_cash 
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	var portfolio PortfolioCash
	err = db.QueryRowx(query, id, userID).StructScan(&portfolio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("Error finding cash portfolio by ID")
		return nil, fmt.Errorf("kesalahan mengambil portfolio kas: %w", err)
	}

	return &portfolio, nil
}

// Update updates a cash portfolio entry (partial updates via COALESCE)
func (r *portfolioCashRepository) Update(
	id int,
	userID int,
	account string,
	bank string,
	amount *float64,
	yieldRate *float64,
	yieldPeriod string,
	yieldFrequencyType string,
	yieldFrequencyValue *int,
	yieldPaymentType string,
	hasMaturity *bool,
	maturityDate *time.Time,
	note *string,
	status string,
	category string,
) (*PortfolioCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE portfolio_cash 
		SET account = COALESCE(NULLIF($1, ''), account),
			bank = COALESCE(NULLIF($2, ''), bank),
			amount = COALESCE($3, amount),
			yield_rate = COALESCE($4, yield_rate),
			yield_period = COALESCE(NULLIF($5, ''), yield_period),
			yield_frequency_type = COALESCE(NULLIF($6, ''), yield_frequency_type),
			yield_frequency_value = COALESCE($7, yield_frequency_value),
			yield_payment_type = COALESCE(NULLIF($8, ''), yield_payment_type),
			has_maturity = COALESCE($9, has_maturity),
			maturity_date = COALESCE($10, maturity_date),
			note = COALESCE($11, note),
			status = COALESCE(NULLIF($12, ''), status),
			category = COALESCE(NULLIF($13, ''), category),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $14 AND user_id = $15 AND deleted_at IS NULL
		RETURNING id, user_id, account, bank, amount, yield_rate, yield_period, 
			  yield_frequency_type, yield_frequency_value, yield_payment_type, 
			  has_maturity, maturity_date, note, status, category, created_at, updated_at, deleted_at
	`

	var result PortfolioCash
	err = db.QueryRowx(query,
		account,
		bank,
		amount,
		yieldRate,
		yieldPeriod,
		yieldFrequencyType,
		yieldFrequencyValue,
		yieldPaymentType,
		hasMaturity,
		maturityDate,
		note,
		status,
		category,
		id,
		userID,
	).StructScan(&result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.Update] Error updating cash portfolio entry")
		return nil, fmt.Errorf("kesalahan memperbarui portfolio kas: %w", err)
	}

	return &result, nil
}

// Delete soft deletes a cash portfolio entry
func (r *portfolioCashRepository) Delete(id int, userID int) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	query := `
		UPDATE portfolio_cash 
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		RETURNING id
	`

	var resultID int
	err = db.QueryRowx(query, id, userID).Scan(&resultID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("portfolio tidak ditemukan")
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.Delete] Error deleting cash portfolio entry")
		return fmt.Errorf("kesalahan menghapus portfolio kas: %w", err)
	}

	return nil
}

// MoveAsset moves assets from source portfolio to target portfolio (transactional)
func (r *portfolioCashRepository) MoveAsset(sourceID int, targetID int, userID int) (*MoveAssetResponse, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	tx, err := db.Beginx()
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioCash.MoveAsset] Error beginning transaction for move asset")
		return nil, fmt.Errorf("kesalahan memulai transaksi: %w", err)
	}

	// Get source portfolio
	sourceQuery := `SELECT * FROM portfolio_cash WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var sourcePortfolio PortfolioCash
	err = tx.QueryRowx(sourceQuery, sourceID, userID).StructScan(&sourcePortfolio)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("portfolio sumber tidak ditemukan")
		}
		return nil, fmt.Errorf("kesalahan menemukan portfolio sumber: %w", err)
	}

	// Check status
	if sourcePortfolio.Status != "maturity" {
		tx.Rollback()
		return nil, fmt.Errorf("portfolio sumber harus memiliki status 'maturity' untuk dipindahkan")
	}

	// Get target portfolio
	targetQuery := `SELECT * FROM portfolio_cash WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var targetPortfolio PortfolioCash
	err = tx.QueryRowx(targetQuery, targetID, userID).StructScan(&targetPortfolio)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("portfolio target tidak ditemukan")
		}
		return nil, fmt.Errorf("kesalahan menemukan portfolio target: %w", err)
	}
	// Check status
	if targetPortfolio.Status != "active" {
		tx.Rollback()
		return nil, fmt.Errorf("portfolio target harus memiliki status 'active' untuk dipindahkan")
	}

	if sourceID == targetID {
		tx.Rollback()
		return nil, fmt.Errorf("portfolio sumber dan target tidak boleh sama")
	}

	// Update target
	newAmount := targetPortfolio.Amount + sourcePortfolio.Amount
	err = tx.QueryRowx(`UPDATE portfolio_cash SET amount = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING *`,
		newAmount, targetID).StructScan(&targetPortfolio)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("kesalahan memperbarui portfolio target: %w", err)
	}

	// Delete source
	_, err = tx.Exec(`UPDATE portfolio_cash SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, sourceID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("kesalahan menghapus portfolio sumber: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("kesalahan melakukan commit: %w", err)
	}

	return &MoveAssetResponse{
		Source:      &sourcePortfolio,
		Target:      &targetPortfolio,
		MovedAmount: sourcePortfolio.Amount,
	}, nil
}

// RealizeCashPortfolio realizes a portfolio and creates PnL (transactional)
func (r *portfolioCashRepository) RealizeCashPortfolio(
	userID int,
	portfolioID int,
	finalSaldo float64,
	pnlAmount float64,
	realizedAt time.Time,
) (*RealizeCashPortfolioResponse, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("kesalahan memulai transaksi: %w", err)
	}

	// Update portfolio
	var portfolio PortfolioCash
	err = tx.QueryRowx(`UPDATE portfolio_cash SET amount = $1, status = 'maturity', updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL RETURNING *`,
		finalSaldo, portfolioID, userID).StructScan(&portfolio)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("portfolio cash tidak ditemukan")
		}
		return nil, fmt.Errorf("kesalahan memperbarui portfolio: %w", err)
	}

	// Create PnL
	var pnl PortfolioPnlRealizedCash
	err = tx.QueryRowx(`INSERT INTO portfolio_pnl_realized_cash (user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, user_id, portfolio_cash_id,amount,realized_at,
		created_at,updated_at,deleted_at`,
		userID, portfolioID, pnlAmount, realizedAt).StructScan(&pnl)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("kesalahan membuat entry PnL: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("kesalahan melakukan commit: %w", err)
	}

	return &RealizeCashPortfolioResponse{
		Portfolio: &portfolio,
		PnL:       &pnl,
	}, nil
}

// CreatePnlEntry creates a new PnL realized cash entry
func (r *portfolioCashRepository) CreatePnlEntry(userID, portfolioCashID int, amount float64, realizedAt time.Time) (*PortfolioPnlRealizedCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	queries := `
INSERT INTO portfolio_pnl_realized_cash (user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at, deleted_at
`

	var result PortfolioPnlRealizedCash
	err = db.QueryRowx(queries, userID, portfolioCashID, amount, realizedAt).StructScan(&result)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioCash.CreatePnlEntry] Error creating PnL realized cash entry")
		return nil, fmt.Errorf("kesalahan membuat entry PnL: %w", err)
	}

	return &result, nil
}

// FindPnlByUserID retrieves all PnL entries for a user with pagination
func (r *portfolioCashRepository) FindPnlByUserID(userID int, limit, offset int) ([]*PortfolioPnlRealizedCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
SELECT id, user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at, deleted_at
FROM portfolio_pnl_realized_cash 
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY realized_at DESC
LIMIT $2 OFFSET $3
`

	var pnlEntries []*PortfolioPnlRealizedCash
	err = db.Select(&pnlEntries, query, userID, limit+1, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioCash.FindPnlByUserID] Error finding PnL entries by user ID")
		return nil, fmt.Errorf("kesalahan mengambil data PnL: %w", err)
	}

	return pnlEntries, nil
}

// FindPnlByPortfolioCashID retrieves all PnL entries for a specific portfolio with pagination
func (r *portfolioCashRepository) FindPnlByPortfolioCashID(portfolioCashID, userID int, limit, offset int) ([]*PortfolioPnlRealizedCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
SELECT id, user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at, deleted_at
FROM portfolio_pnl_realized_cash 
WHERE portfolio_cash_id = $1 AND user_id = $2 AND deleted_at IS NULL
ORDER BY realized_at DESC
LIMIT $3 OFFSET $4
`

	var pnlEntries []*PortfolioPnlRealizedCash
	err = db.Select(&pnlEntries, query, portfolioCashID, userID, limit, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[PortfolioCash.FindPnlByPortfolioCashID] Error finding PnL entries by portfolio cash ID")
		return nil, fmt.Errorf("kesalahan mengambil data PnL: %w", err)
	}

	return pnlEntries, nil
}

// FindPnlByID retrieves a specific PnL entry
func (r *portfolioCashRepository) FindPnlByID(id, userID int) (*PortfolioPnlRealizedCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
SELECT id, user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at, deleted_at
FROM portfolio_pnl_realized_cash 
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
`

	var pnl PortfolioPnlRealizedCash
	err = db.QueryRowx(query, id, userID).StructScan(&pnl)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.FindPnlByID] Error finding PnL entry by ID")
		return nil, fmt.Errorf("kesalahan mengambil data PnL: %w", err)
	}

	return &pnl, nil
}

// UpdatePnlEntry updates a PnL entry (partial updates via COALESCE)
func (r *portfolioCashRepository) UpdatePnlEntry(id, userID int, amount *float64, realizedAt *time.Time) (*PortfolioPnlRealizedCash, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
UPDATE portfolio_pnl_realized_cash 
SET amount = COALESCE($1, amount),
realized_at = COALESCE($2, realized_at),
updated_at = CURRENT_TIMESTAMP
WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL
RETURNING id, user_id, portfolio_cash_id, amount, realized_at, created_at, updated_at, deleted_at
`

	var result PortfolioPnlRealizedCash
	err = db.QueryRowx(query, amount, realizedAt, id, userID).StructScan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.UpdatePnlEntry] Error updating PnL entry")
		return nil, fmt.Errorf("kesalahan memperbarui data PnL: %w", err)
	}

	return &result, nil
}

// DeletePnlEntry soft deletes a PnL entry
func (r *portfolioCashRepository) DeletePnlEntry(id, userID int) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	query := `
UPDATE portfolio_pnl_realized_cash 
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING id
`

	var resultID int
	err = db.QueryRowx(query, id, userID).Scan(&resultID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("data PnL tidak ditemukan")
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.DeletePnlEntry] Error deleting PnL entry")
		return fmt.Errorf("kesalahan menghapus data PnL: %w", err)
	}

	return nil
}

// GetPnlSummary retrieves summary statistics for a user's PnL entries
func (r *portfolioCashRepository) GetPnlSummary(userID int) (*PnlSummary, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	query := `
SELECT 
COALESCE(SUM(amount), 0) as total_amount,
COUNT(*) as count,
COALESCE(AVG(amount), 0) as avg_amount,
COALESCE(MAX(amount), 0) as max_amount,
COALESCE(MIN(amount), 0) as min_amount,
MAX(realized_at) as last_realized_at
FROM portfolio_pnl_realized_cash 
WHERE user_id = $1 AND deleted_at IS NULL
`

	var summary PnlSummary
	err = db.QueryRowx(query, userID).StructScan(&summary)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		Logger.Error().Err(err).Msg("[PortfolioCash.GetPnlSummary] Error getting PnL summary")
		return nil, fmt.Errorf("kesalahan mengambil ringkasan PnL: %w", err)
	}

	return &summary, nil
}
