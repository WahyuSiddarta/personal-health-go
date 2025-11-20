package validator

// CreatePortfolioBondRequest represents the payload required to add a bond to a user portfolio.
type CreatePortfolioBondRequest struct {
	BondId          string  `db:"bond_id" json:"bond_id"`
	Name            *string `json:"name" validate:"omitempty"`
	PurchasePrice   float64 `json:"purchase_price" validate:"required,min=0"`
	CouponRate      float64 `json:"coupon_rate" validate:"required,min=0"`
	CouponFrequency string  `json:"coupon_frequency" validate:"required,oneof=monthly quarterly semi-annual annual"`
	NextCouponDate  *string `json:"next_coupon_date" validate:"omitempty,datetime=2006-01-02"`
	MaturityDate    string  `json:"maturity_date" validate:"required,datetime=2006-01-02"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
	SecondaryMarket bool    `db:"secondary_market" json:"secondary_market"`
	Note            *string `json:"note"`
}

// UpdatePortfolioBondRequest represents partial updates for a bond.
type UpdatePortfolioBondRequest struct {
	ISIN            *string  `json:"isin"`
	Name            *string  `json:"name"`
	PurchasePrice   *float64 `json:"purchase_price" validate:"omitempty,min=0"`
	FaceValue       *float64 `json:"face_value" validate:"omitempty,min=0"`
	CouponRate      *float64 `json:"coupon_rate" validate:"omitempty,min=0"`
	CouponFrequency *string  `json:"coupon_frequency" validate:"omitempty,oneof=monthly quarterly semi-annual annual"`
	NextCouponDate  *string  `json:"next_coupon_date" validate:"omitempty,datetime=2006-01-02"`
	MaturityDate    *string  `json:"maturity_date" validate:"omitempty,datetime=2006-01-02"`
	Quantity        *int     `json:"quantity" validate:"omitempty,min=1"`
	Issuer          *string  `json:"issuer"`
	Note            *string  `json:"note"`
	Status          *string  `json:"status" validate:"omitempty,oneof=active inactive matured sold"`
	MarketPrice     *float64 `json:"market_price" validate:"omitempty,min=0"`
}

// UpdateMarketPriceOverrideRequest captures data for market price overrides.
type UpdateMarketPriceOverrideRequest struct {
	MarketPriceOverride float64 `json:"market_price_override" validate:"required,min=0"`
}

// CreateCouponRequest represents the payload for creating coupon records.
type CreateCouponRequest struct {
	PortfolioBondID int     `json:"portfolio_bond_id" validate:"required,gt=0"`
	CouponNumber    int     `json:"coupon_number" validate:"required,gt=0"`
	PaymentDate     string  `json:"payment_date" validate:"required,datetime=2006-01-02"`
	Amount          float64 `json:"amount" validate:"required,min=0"`
	Status          string  `json:"status" validate:"omitempty,oneof=pending received missed"`
	Note            *string `json:"note"`
}

// UpdateCouponRequest represents partial updates to coupon records.
type UpdateCouponRequest struct {
	CouponNumber *int     `json:"coupon_number" validate:"omitempty,gt=0"`
	PaymentDate  *string  `json:"payment_date" validate:"omitempty,datetime=2006-01-02"`
	Amount       *float64 `json:"amount" validate:"omitempty,min=0"`
	Status       *string  `json:"status" validate:"omitempty,oneof=pending received missed"`
	Note         *string  `json:"note"`
}

// CreateRealizedBondRequest represents the payload for recording realized bonds.
type CreateRealizedBondRequest struct {
	PortfolioBondID      int      `json:"portfolio_bond_id" validate:"required,gt=0"`
	RealizedPrice        float64  `json:"realized_price" validate:"required"`
	TotalCouponsReceived *float64 `json:"total_coupons_received" validate:"omitempty,min=0"`
	RealizedDate         *string  `json:"realized_date" validate:"omitempty,datetime=2006-01-02"`
	Note                 *string  `json:"note"`
}

// UpdateRealizedBondRequest represents partial updates to realized bond records.
type UpdateRealizedBondRequest struct {
	RealizedPrice        *float64 `json:"realized_price" validate:"omitempty"`
	TotalCouponsReceived *float64 `json:"total_coupons_received" validate:"omitempty,min=0"`
	RealizedDate         *string  `json:"realized_date" validate:"omitempty,datetime=2006-01-02"`
	Note                 *string  `json:"note"`
}
