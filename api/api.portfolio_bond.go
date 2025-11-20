package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// PortfolioBondHandlers contains handlers for bond portfolios and related entities.
type PortfolioBondHandlers struct {
	bondRepo     models.PortfolioBondRepository
	couponRepo   models.PortfolioBondCouponRepository
	realizedRepo models.PortfolioBondRealizedRepository
}

// NewPortfolioBondHandlers constructs a new PortfolioBondHandlers instance.
func NewPortfolioBondHandlers(
	bondRepo models.PortfolioBondRepository,
	couponRepo models.PortfolioBondCouponRepository,
	realizedRepo models.PortfolioBondRealizedRepository,
) *PortfolioBondHandlers {
	return &PortfolioBondHandlers{
		bondRepo:     bondRepo,
		couponRepo:   couponRepo,
		realizedRepo: realizedRepo,
	}
}

// CreateBondPortfolio creates a new bond entry for the authenticated user.
func (h *PortfolioBondHandlers) CreateBondPortfolio(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.CreatePortfolioBondRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	var nextCouponDate *time.Time
	if req.NextCouponDate != nil {
		nextCouponDate, err = validator.ParseDate(req.NextCouponDate)
		if err != nil {
			Logger.Error().Err(err).Msg("[CreateBondPortfolio] Invalid next coupon date")
			return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal kupon berikutnya tidak valid", nil)
		}
	}

	maturityDate, err := validator.ParseDate(&req.MaturityDate)
	if err != nil || maturityDate == nil {
		Logger.Error().Err(err).Msg("[CreateBondPortfolio] Invalid maturity date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal jatuh tempo tidak valid", nil)
	}

	payload := models.PortfolioBondCreateRequest{
		BondId:          req.BondId,
		Name:            req.Name,
		PurchasePrice:   req.PurchasePrice,
		CouponRate:      req.CouponRate,
		CouponFrequency: req.CouponFrequency,
		NextCouponDate:  nextCouponDate,
		MaturityDate:    *maturityDate,
		Quantity:        req.Quantity,
		SecondaryMarket: req.SecondaryMarket,
		Note:            req.Note,
	}

	result, err := h.bondRepo.Create(userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreateBondPortfolio] Error creating bond portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "CreateBondPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, result)
}

// GetMyBondPortfolios returns all bond portfolios enriched with potential gain.
func (h *PortfolioBondHandlers) GetMyBondPortfolios(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolios, err := h.bondRepo.FindByUserID(userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetMyBondPortfolios] Error fetching bond portfolios")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetMyBondPortfolios"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if portfolios == nil {
		portfolios = []*models.PortfolioBond{}
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{"portfolios": portfolios})
}

// UpdateMarketPriceOverride updates override values for the bond market price.
func (h *PortfolioBondHandlers) UpdateMarketPriceOverride(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdateMarketPriceOverrideRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("portfolioId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	if _, err := h.bondRepo.UpdateMarketPriceOverride(portfolioID, userID, req.MarketPriceOverride); err != nil {
		Logger.Error().Err(err).Msg("[UpdateMarketPriceOverride] Error updating override")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdateMarketPriceOverride"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}

//// checker

// UpdateBondPortfolio handles updating a bond portfolio entry.
func (h *PortfolioBondHandlers) UpdateBondPortfolio(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdatePortfolioBondRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("portfolioId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	nextCouponDate, err := validator.ParseDate(req.NextCouponDate)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateBondPortfolio] Invalid next coupon date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal kupon berikutnya tidak valid", nil)
	}

	maturityDate, err := validator.ParseDate(req.MaturityDate)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateBondPortfolio] Invalid maturity date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal jatuh tempo tidak valid", nil)
	}

	payload := models.PortfolioBondUpdateRequest{
		ISIN:            req.ISIN,
		Name:            req.Name,
		PurchasePrice:   req.PurchasePrice,
		FaceValue:       req.FaceValue,
		CouponRate:      req.CouponRate,
		CouponFrequency: req.CouponFrequency,
		NextCouponDate:  nextCouponDate,
		MaturityDate:    maturityDate,
		Quantity:        req.Quantity,
		Issuer:          req.Issuer,
		Note:            req.Note,
		Status:          req.Status,
		MarketPrice:     req.MarketPrice,
	}

	result, err := h.bondRepo.Update(portfolioID, userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateBondPortfolio] Error updating bond portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdateBondPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Portofolio tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// DeleteBondPortfolio removes a bond portfolio entry.
func (h *PortfolioBondHandlers) DeleteBondPortfolio(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("portfolioId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	err = h.bondRepo.Delete(portfolioID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteBondPortfolio] Error deleting bond portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "DeleteBondPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}

// CreateCoupon adds a coupon payment record to a bond.
func (h *PortfolioBondHandlers) CreateCoupon(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.CreateCouponRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	paymentDate, err := validator.ParseDate(&req.PaymentDate)
	if err != nil || paymentDate == nil {
		Logger.Error().Err(err).Msg("[CreateCoupon] Invalid payment date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal pembayaran tidak valid", nil)
	}

	payload := models.PortfolioBondCouponCreateRequest{
		PortfolioBondID: req.PortfolioBondID,
		CouponNumber:    req.CouponNumber,
		PaymentDate:     *paymentDate,
		Amount:          req.Amount,
		Status:          req.Status,
		Note:            req.Note,
	}

	result, err := h.couponRepo.Create(userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreateCoupon] Error creating coupon")
		middleware.CaptureError(c, err, map[string]string{"handler": "CreateCoupon"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, result)
}

// GetCouponsByBond returns all coupons for a specific bond.
func (h *PortfolioBondHandlers) GetCouponsByBond(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioBondID, err := strconv.Atoi(c.Param("portfolioBondId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portfolio tidak valid", nil)
	}

	coupons, err := h.couponRepo.FindByPortfolioBondID(userID, portfolioBondID)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetCouponsByBond] Error fetching coupons")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetCouponsByBond"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{"coupons": coupons})
}

// UpdateCoupon modifies an existing coupon record.
func (h *PortfolioBondHandlers) UpdateCoupon(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdateCouponRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	couponID, err := strconv.Atoi(c.Param("couponId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID kupon tidak valid", nil)
	}

	paymentDate, err := validator.ParseDate(req.PaymentDate)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateCoupon] Invalid payment date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal pembayaran tidak valid", nil)
	}

	payload := models.PortfolioBondCouponUpdateRequest{
		CouponNumber: req.CouponNumber,
		PaymentDate:  paymentDate,
		Amount:       req.Amount,
		Status:       req.Status,
		Note:         req.Note,
	}

	result, err := h.couponRepo.Update(couponID, userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateCoupon] Error updating coupon")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdateCoupon"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Kupon tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// DeleteCoupon removes a coupon record.
func (h *PortfolioBondHandlers) DeleteCoupon(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	couponID, err := strconv.Atoi(c.Param("couponId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID kupon tidak valid", nil)
	}

	err = h.couponRepo.Delete(couponID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteCoupon] Error deleting coupon")
		middleware.CaptureError(c, err, map[string]string{"handler": "DeleteCoupon"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}

// CreateRealizedBond records a realized bond entry.
func (h *PortfolioBondHandlers) CreateRealizedBond(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.CreateRealizedBondRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	realizedDate, err := validator.ParseDate(req.RealizedDate)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreateRealizedBond] Invalid realized date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal realisasi tidak valid", nil)
	}

	payload := models.PortfolioBondRealizedCreateRequest{
		PortfolioBondID:      req.PortfolioBondID,
		RealizedPrice:        req.RealizedPrice,
		TotalCouponsReceived: req.TotalCouponsReceived,
		RealizedDate:         realizedDate,
		Note:                 req.Note,
	}

	result, err := h.realizedRepo.Create(userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreateRealizedBond] Error creating realized bond")
		middleware.CaptureError(c, err, map[string]string{"handler": "CreateRealizedBond"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, result)
}

// GetRealizedBonds returns all realized entries for a user with pagination.
func (h *PortfolioBondHandlers) GetRealizedBonds(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	limit, offset := parseLimitOffset(c)

	entries, err := h.realizedRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetRealizedBonds] Error fetching realized bonds")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetRealizedBonds"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	hasNextData := false
	if len(entries) > limit {
		entries = entries[:limit]
		hasNextData = true
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"pagination": map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"hasNextData": hasNextData,
		},
	})
}

// GetRealizedBondsByPortfolioId returns realized entries for a specific bond.
func (h *PortfolioBondHandlers) GetRealizedBondsByPortfolioId(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioBondID, err := strconv.Atoi(c.Param("portfolioBondId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portfolio obligasi tidak valid", nil)
	}

	limit, offset := parseLimitOffset(c)

	entries, err := h.realizedRepo.FindByPortfolioBondID(userID, portfolioBondID, limit, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetRealizedBondsByPortfolioId] Error fetching realized bonds")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetRealizedBondsByPortfolioId"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	hasNextData := false
	if len(entries) > limit {
		entries = entries[:limit]
		hasNextData = true
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"pagination": map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"hasNextData": hasNextData,
		},
	})
}

// UpdateRealizedBond updates a realized bond entry.
func (h *PortfolioBondHandlers) UpdateRealizedBond(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdateRealizedBondRequest)
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	realizedID, err := strconv.Atoi(c.Param("realizedId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID data realisasi tidak valid", nil)
	}

	realizedDate, err := validator.ParseDate(req.RealizedDate)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateRealizedBond] Invalid realized date")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal realisasi tidak valid", nil)
	}

	payload := models.PortfolioBondRealizedUpdateRequest{
		RealizedPrice:        req.RealizedPrice,
		TotalCouponsReceived: req.TotalCouponsReceived,
		RealizedDate:         realizedDate,
		Note:                 req.Note,
	}

	result, err := h.realizedRepo.Update(realizedID, userID, payload)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateRealizedBond] Error updating realized bond")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdateRealizedBond"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Data obligasi terealisasi tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// DeleteRealizedBond removes a realized bond entry.
func (h *PortfolioBondHandlers) DeleteRealizedBond(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	realizedID, err := strconv.Atoi(c.Param("realizedId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID data realisasi tidak valid", nil)
	}

	err = h.realizedRepo.Delete(realizedID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteRealizedBond] Error deleting realized bond")
		middleware.CaptureError(c, err, map[string]string{"handler": "DeleteRealizedBond"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}
