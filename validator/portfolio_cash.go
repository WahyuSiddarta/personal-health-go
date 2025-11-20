package validator

import "time"

// CreatePortfolioCashRequest represents request to create a cash portfolio.
type CreatePortfolioCashRequest struct {
	Account             string   `json:"account" validate:"required"`
	Bank                string   `json:"bank" validate:"required"`
	Amount              float64  `json:"amount" validate:"required,min=0"`
	YieldRate           *float64 `json:"yield_rate" validate:"omitempty,min=0,max=100"`
	YieldPeriod         string   `json:"yield_period"`
	YieldFrequencyType  string   `json:"yield_frequency_type" validate:"required,oneof=daily monthly yearly"`
	YieldFrequencyValue int      `json:"yield_frequency_value" validate:"required,min=1"`
	YieldPaymentType    string   `json:"yield_payment_type" validate:"required"`
	HasMaturity         bool     `json:"has_maturity"`
	MaturityDate        *string  `json:"maturity_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Note                *string  `json:"note"`
	Category            string   `json:"category" validate:"required,oneof=liquid time_deposit money_market other"`
}

// ParsedMaturityDate returns the parsed maturity date if provided.
func (r *CreatePortfolioCashRequest) ParsedMaturityDate() (*time.Time, error) {
	return parseDateString(r.MaturityDate)
}

// UpdatePortfolioCashRequest represents request to update cash portfolio data.
type UpdatePortfolioCashRequest struct {
	Account             string   `json:"account"`
	Bank                string   `json:"bank"`
	Amount              *float64 `json:"amount" validate:"omitempty,min=0"`
	YieldRate           *float64 `json:"yield_rate" validate:"omitempty,min=0,max=100"`
	YieldPeriod         string   `json:"yield_period"`
	YieldFrequencyType  string   `json:"yield_frequency_type" validate:"omitempty,oneof=daily monthly yearly"`
	YieldFrequencyValue *int     `json:"yield_frequency_value" validate:"omitempty,min=1"`
	YieldPaymentType    string   `json:"yield_payment_type"`
	HasMaturity         *bool    `json:"has_maturity"`
	MaturityDate        *string  `json:"maturity_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Note                *string  `json:"note"`
	Status              string   `json:"status" validate:"omitempty,oneof=active maturity"`
	Category            string   `json:"category" validate:"omitempty,oneof=liquid time_deposit money_market other"`
}

// ParsedMaturityDate returns the parsed maturity date if provided.
func (r *UpdatePortfolioCashRequest) ParsedMaturityDate() (*time.Time, error) {
	return parseDateString(r.MaturityDate)
}

// MoveAssetRequest represents request to move an asset between portfolios.
type MoveAssetRequest struct {
	SourcePortfolioID int `json:"source_portfolio_id" validate:"required,gt=0"`
	TargetPortfolioID int `json:"target_portfolio_id" validate:"required,gt=0"`
}

// RealizeCashPortfolioRequest represents request to realize a cash portfolio.
type RealizeCashPortfolioRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	FinalSaldo      float64 `json:"final_saldo" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	RealizedAt      string  `json:"realized_at" validate:"required,datetime=2006-01-02"`
}

// CreatePnlRealizedCashRequest represents request to create a realized PnL entry.
type CreatePnlRealizedCashRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	Amount          float64 `json:"amount" validate:"required"`
	Note            *string `json:"note"`
	RealizedAt      *string `json:"realized_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// UpdatePnlRealizedCashRequest represents request to update a realized PnL entry.
type UpdatePnlRealizedCashRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	Amount          float64 `json:"amount" validate:"required"`
	Note            *string `json:"note"`
	RealizedAt      *string `json:"realized_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// PnlListQuery represents query parameters for listing PnL entries.
type PnlListQuery struct {
	Limit  int `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset int `query:"offset" validate:"omitempty,min=0"`
}
