package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// AuthHandlers contains all authentication-related handlers
type AuthHandlers struct {
	repo models.UserRepository
}

// NewAuthHandlers creates a new instance of auth handlers
func NewAuthHandlers(repo models.UserRepository) *AuthHandlers {
	return &AuthHandlers{repo: repo}
}

// convertPaymentData converts payment request data to model
func convertPaymentData(reqPaymentData *validator.PaymentDataRequest) (*models.PaymentData, error) {
	if reqPaymentData == nil {
		return nil, nil
	}

	paymentData := &models.PaymentData{
		OriginalPrice:  reqPaymentData.OriginalPrice,
		PaidPrice:      reqPaymentData.PaidPrice,
		DiscountAmount: reqPaymentData.DiscountAmount,
		DiscountReason: reqPaymentData.DiscountReason,
		PaymentMethod:  reqPaymentData.PaymentMethod,
		Notes:          reqPaymentData.Notes,
	}

	// Parse payment date if provided
	if reqPaymentData.PaymentDate != nil {
		paymentDate, err := time.Parse(time.RFC3339, *reqPaymentData.PaymentDate)
		if err != nil {
			return nil, err
		}
		paymentData.PaymentDate = &paymentDate
	}

	return paymentData, nil
}

// Login handles user authentication
func (h *AuthHandlers) Login(c echo.Context) error {
	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.LoginRequest)

	// Find user by email
	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		middleware.CaptureException(c, err)
		Logger.Error().Err(err).Str("email", req.Email).Msg("[Login] Gagal mencari pengguna saat login")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	if user == nil {
		err := errors.New("[login] user not found")
		middleware.CaptureException(c, err)
		return helper.ErrorResponse(c, http.StatusBadRequest, "Email atau password tidak valid", nil)
	}

	// Validate password
	if err := h.repo.ValidatePassword(req.Password, user.Password); err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Email atau password tidak valid", nil)
	}

	// Check if user account is active
	if user.Status != models.UserStatusActive {
		return helper.ErrorResponse(c, http.StatusForbidden, "Akses akun ditolak", nil)
	}

	middleware.SetUserContext(c, user.ID, user.Email)
	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		middleware.CaptureError(c, err,
			map[string]string{"action": "generate_token"},
			nil,
		)
		Logger.Error().Err(err).Int("user_id", user.ID).Msg("[Login] Gagal membuat token")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	// Remove password from response
	user.Password = ""
	return helper.JsonResponse(c, http.StatusOK, validator.LoginData{
		User:  user,
		Token: token,
	})
}

// Register handles user registration
func (h *AuthHandlers) Register(c echo.Context) error {
	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.RegisterRequest)

	// Check if user already exists
	existingUser, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		middleware.CaptureException(c, err)
		Logger.Error().Err(err).Str("email", req.Email).Msg("[Register] Gagal memeriksa pengguna saat registrasi")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	if existingUser != nil {
		return helper.ErrorResponse(c, http.StatusConflict, "Email sudah terdaftar", nil)
	}

	// Set default values if not provided
	status := models.UserStatusActive
	if req.Status != nil {
		status = *req.Status
	}

	userLevel := models.UserLevelFree
	if req.UserLevel != nil {
		userLevel = *req.UserLevel
	}

	createReq := &models.CreateUserRequest{
		Email:     req.Email,
		Password:  req.Password,
		Status:    status,
		UserLevel: userLevel,
	}

	// Create new user
	newUser, err := h.repo.Create(createReq)
	if err != nil {
		middleware.CaptureError(c, err,
			map[string]string{"action": "CreateUser"},
			map[string]interface{}{"email": req.Email},
		)
		Logger.Error().Err(err).Str("email", req.Email).Msg("[Register] Gagal membuat pengguna saat registrasi")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(newUser.ID)
	if err != nil {
		middleware.CaptureError(c, err,
			map[string]string{"action": "generate_token"},
			nil,
		)
		Logger.Error().Err(err).Int("user_id", newUser.ID).Msg("[Register] Gagal membuat token untuk pengguna baru")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", nil)
	}

	Logger.Info().Str("email", req.Email).Int("user_id", newUser.ID).Msg("[Register] Pengguna baru berhasil didaftarkan")
	middleware.SetUserContext(c, newUser.ID, newUser.Email)

	return helper.JsonResponse(c, http.StatusCreated, validator.RegisterData{
		User:  newUser,
		Token: token,
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandlers) GetProfile(c echo.Context) error {
	// Get authenticated user from middleware
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, authUser)
}

// UpdateUserLevel handles updating user subscription level (admin only)
func (h *AuthHandlers) UpdateUserLevel(c echo.Context) error {
	// Get user ID from path parameter
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
	}

	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.UpdateUserLevelRequest)

	// Validate custom logic
	if errs := validator.CustomValidation(req); len(errs) > 0 {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Kesalahan validasi", errs)
	}

	// Get admin user ID for audit trail
	adminUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", nil)
	}

	// Convert payment data if provided
	paymentData, err := convertPaymentData(req.PaymentData)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal pembayaran tidak valid", nil)
	}

	// Update user level
	result, err := h.repo.UpdateUserLevel(userID, req.UserLevel, paymentData, &adminUser.ID)
	if err != nil {
		Logger.Error().Err(err).Int("user_id", userID).Str("new_level", string(req.UserLevel)).Msg("[UpdateUserLevel] Gagal memperbarui level pengguna")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui level pengguna", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// UpdateUserStatus handles updating user account status (admin only)
func (h *AuthHandlers) UpdateUserStatus(c echo.Context) error {
	// Get user ID from path parameter
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
	}

	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.UpdateUserStatusRequest)

	// Get admin user for logging
	adminUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", nil)
	}

	// Update user status
	updatedUser, err := h.repo.UpdateUserStatus(userID, req.Status)
	if err != nil {
		Logger.Error().Err(err).Int("user_id", userID).Str("new_status", string(req.Status)).Msg("[UpdateUserStatus] Error updating user status")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error updating user status", nil)
	}

	Logger.Info().Int("user_id", userID).Str("new_status", string(req.Status)).Int("admin_id", adminUser.ID).Msg("[UpdateUserStatus] User status updated by admin")

	return helper.JsonResponse(c, http.StatusOK, updatedUser)
}

// GetAllUsers returns paginated list of users with optional filters (admin only)
func (h *AuthHandlers) GetAllUsers(c echo.Context) error {
	// Get validated query parameters
	query := validator.GetValidatedQuery(c).(*validator.GetUsersQuery)

	// Set defaults
	page := query.Page
	if page <= 0 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get users with filters
	result, err := h.repo.GetAllUsers(page, limit, query.Status, query.UserLevel, query.EmailFilter)
	if err != nil {
		Logger.Error().Err(err).Msg("[GetAllUsers] Error fetching users")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error fetching users", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, result)
}

// GetExpiredUsers returns users with expired premium subscriptions (admin only)
func (h *AuthHandlers) GetExpiredUsers(c echo.Context) error {
	// Get expired users
	expiredUsers, err := h.repo.GetExpiredUsers()
	if err != nil {
		Logger.Error().Err(err).Msg("[GetExpiredUsers] Error fetching expired users")
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error fetching expired users", nil)
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"expired_users": expiredUsers,
		"count":         len(expiredUsers),
	})
}
