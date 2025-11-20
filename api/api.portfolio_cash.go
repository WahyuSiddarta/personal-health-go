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

// PortfolioCashHandlers contains all portfolio cash-related handlers
type PortfolioCashHandlers struct {
	repo models.PortfolioCashRepository
}

// NewPortfolioCashHandlers creates a new instance of portfolio cash handlers
func NewPortfolioCashHandlers(repo models.PortfolioCashRepository) *PortfolioCashHandlers {
	return &PortfolioCashHandlers{repo: repo}
}

// CreateCashPortfolio handles creating a new cash portfolio
func (h *PortfolioCashHandlers) CreateCashPortfolio(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.CreatePortfolioCashRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	// Parse maturity date
	maturityDate, err := req.ParsedMaturityDate()
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal jatuh tempo tidak valid", nil)
	}

	result, err := h.repo.Create(
		userID,
		req.Account,
		req.Bank,
		req.Amount,
		req.YieldRate,
		req.YieldPeriod,
		req.YieldFrequencyType,
		req.YieldFrequencyValue,
		req.YieldPaymentType,
		req.HasMaturity,
		maturityDate,
		req.Note,
		req.Category,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreateCashPortfolio] Error creating cash portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "CreateCashPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, result)
}

// GetMyCashPortfolios retrieves all cash portfolios for the authenticated user
func (h *PortfolioCashHandlers) GetMyCashPortfolios(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolios, err := h.repo.FindByUserID(userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetMyCashPortfolios] Error fetching cash portfolios")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetMyCashPortfolios"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if portfolios == nil {
		portfolios = []*models.PortfolioCash{}
	}

	return helper.JsonResponse(c, http.StatusOK, portfolios)
}

// UpdateCashPortfolio handles updating a cash portfolio entry
func (h *PortfolioCashHandlers) UpdateCashPortfolio(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdatePortfolioCashRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	// Parse maturity date
	maturityDate, err := req.ParsedMaturityDate()
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal jatuh tempo tidak valid", nil)
	}

	result, err := h.repo.Update(
		portfolioID,
		userID,
		req.Account,
		req.Bank,
		req.Amount,
		req.YieldRate,
		req.YieldPeriod,
		req.YieldFrequencyType,
		req.YieldFrequencyValue,
		req.YieldPaymentType,
		req.HasMaturity,
		maturityDate,
		req.Note,
		req.Status,
		req.Category,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdateCashPortfolio] Error updating cash portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdateCashPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Portofolio tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// DeleteCashPortfolio handles deleting a cash portfolio entry
func (h *PortfolioCashHandlers) DeleteCashPortfolio(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	err = h.repo.Delete(portfolioID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeleteCashPortfolio] Error deleting cash portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "DeleteCashPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}

// MoveAsset handles moving assets between portfolios
func (h *PortfolioCashHandlers) MoveAsset(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.MoveAssetRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	result, err := h.repo.MoveAsset(req.SourcePortfolioID, req.TargetPortfolioID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[MoveAsset] Error moving asset")
		middleware.CaptureError(c, err, map[string]string{"handler": "MoveAsset"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// RealizeCashPortfolio handles realizing a cash portfolio
func (h *PortfolioCashHandlers) RealizeCashPortfolio(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.RealizeCashPortfolioRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	// Parse realized at date
	realizedAt, err := time.Parse("2006-01-02", req.RealizedAt)
	if err != nil {
		Logger.Error().Err(err).Msg("[RealizeCashPortfolio] Invalid realized at date format")
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal realisasi tidak valid", nil)
	}

	result, err := h.repo.RealizeCashPortfolio(
		userID,
		req.PortfolioCashID,
		req.FinalSaldo,
		req.Amount,
		realizedAt,
	)
	if err != nil {
		Logger.Error().Err(err).Msg("[RealizeCashPortfolio] Error realizing cash portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "RealizeCashPortfolio"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// GetPnlRealizedCash retrieves all PnL entries for the authenticated user
func (h *PortfolioCashHandlers) GetPnlRealizedCash(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	limit, offset := parseLimitOffset(c)

	pnlEntries, err := h.repo.FindPnlByUserID(userID, limit, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetPnlRealizedCash] Error fetching PnL entries")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetPnlRealizedCash"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if pnlEntries == nil {
		pnlEntries = []*models.PortfolioPnlRealizedCash{}
	}

	hasNextData := false
	if len(pnlEntries) > limit {
		pnlEntries = pnlEntries[:limit]
		hasNextData = true
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"entries": pnlEntries,
		"pagination": map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"hasNextData": hasNextData,
		},
	})
}

// GetPnlByPortfolioCashID retrieves all PnL entries for a specific portfolio
func (h *PortfolioCashHandlers) GetPnlByPortfolioCashID(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	portfolioID, err := strconv.Atoi(c.Param("portfolioId"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID portofolio tidak valid", nil)
	}

	limit, offset := parseLimitOffset(c)

	pnlEntries, err := h.repo.FindPnlByPortfolioCashID(portfolioID, userID, limit, offset)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetPnlByPortfolioCashID] Error fetching PnL entries by portfolio")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetPnlByPortfolioCashID"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if pnlEntries == nil {
		pnlEntries = []*models.PortfolioPnlRealizedCash{}
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"entries": pnlEntries,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetPnlById retrieves a specific PnL entry by ID
func (h *PortfolioCashHandlers) GetPnlById(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	pnlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID PnL tidak valid", nil)
	}

	result, err := h.repo.FindPnlByID(pnlID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetPnlById] Error fetching PnL entry")
		middleware.CaptureError(c, err, map[string]string{"handler": "GetPnlById"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Entri PnL tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// CreatePnlRealizedCash handles creating a new PnL entry
func (h *PortfolioCashHandlers) CreatePnlRealizedCash(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.CreatePnlRealizedCashRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	realizedAt := time.Now()
	if req.RealizedAt != nil && *req.RealizedAt != "" {
		t, err := time.Parse(time.RFC3339, *req.RealizedAt)
		if err != nil {
			middleware.CaptureError(c, err, map[string]string{"handler": "CreatePnlRealizedCash", "error_type": "invalid_date_format"}, nil)
			return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal realisasi tidak valid", nil)
		}
		realizedAt = t
	}

	result, err := h.repo.CreatePnlEntry(userID, req.PortfolioCashID, req.Amount, realizedAt)
	if err != nil {
		Logger.Error().Err(err).Msg("[CreatePnlRealizedCash] Error creating PnL entry")
		middleware.CaptureError(c, err, map[string]string{"handler": "CreatePnlRealizedCash"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusCreated, result)
}

// UpdatePnlRealizedCash handles updating a PnL entry
func (h *PortfolioCashHandlers) UpdatePnlRealizedCash(c echo.Context) error {
	req := validator.GetValidatedRequest(c).(*validator.UpdatePnlRealizedCashRequest)

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	pnlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID PnL tidak valid", nil)
	}

	var realizedAt *time.Time
	if req.RealizedAt != nil && *req.RealizedAt != "" {
		t, err := time.Parse(time.RFC3339, *req.RealizedAt)
		if err != nil {
			middleware.CaptureError(c, err, map[string]string{"handler": "UpdatePnlRealizedCash", "error_type": "invalid_date_format"}, nil)
			return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal realisasi tidak valid", nil)
		}
		realizedAt = &t
	}

	result, err := h.repo.UpdatePnlEntry(pnlID, userID, &req.Amount, realizedAt)
	if err != nil {
		Logger.Error().Err(err).Msg("[UpdatePnlRealizedCash] Error updating PnL entry")
		middleware.CaptureError(c, err, map[string]string{"handler": "UpdatePnlRealizedCash"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	if result == nil {
		return helper.ErrorResponse(c, http.StatusNotFound, "Entri PnL tidak ditemukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// DeletePnlRealizedCash handles deleting a PnL entry
func (h *PortfolioCashHandlers) DeletePnlRealizedCash(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna belum diautentikasi", nil)
	}

	pnlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID PnL tidak valid", nil)
	}

	err = h.repo.DeletePnlEntry(pnlID, userID)
	if err != nil {
		Logger.Error().Err(err).Msg("[DeletePnlRealizedCash] Error deleting PnL entry")
		middleware.CaptureError(c, err, map[string]string{"handler": "DeletePnlRealizedCash"}, nil)
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, nil)
}
