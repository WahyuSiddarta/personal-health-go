package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID               int        `json:"id" db:"id"`
	Email            string     `json:"email" db:"email"`
	Password         string     `json:"-" db:"password"` // Never include password in JSON responses
	Status           UserStatus `json:"status" db:"status"`
	UserLevel        UserLevel  `json:"user_level" db:"user_level"`
	PremiumExpiresAt *time.Time `json:"premium_expires_at,omitempty" db:"premium_expires_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

// Scan implements the sql.Scanner interface
func (us *UserStatus) Scan(value interface{}) error {
	if value == nil {
		*us = UserStatusActive
		return nil
	}
	if str, ok := value.(string); ok {
		*us = UserStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into UserStatus", value)
}

// Value implements the driver.Valuer interface
func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
}

// UserLevel represents subscription level
type UserLevel string

const (
	UserLevelFree        UserLevel = "free"
	UserLevelPremium     UserLevel = "premium"
	UserLevelPremiumPlus UserLevel = "premium+"
	UserLevelAdmin       UserLevel = "admin"
)

// Scan implements the sql.Scanner interface
func (ul *UserLevel) Scan(value interface{}) error {
	if value == nil {
		*ul = UserLevelFree
		return nil
	}
	if str, ok := value.(string); ok {
		*ul = UserLevel(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into UserLevel", value)
}

// Value implements the driver.Valuer interface
func (ul UserLevel) Value() (driver.Value, error) {
	return string(ul), nil
}

// PaymentRecord represents a payment record for premium subscriptions
type PaymentRecord struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	SubscriptionType   UserLevel `json:"subscription_type" db:"subscription_type"`
	OriginalPrice      float64   `json:"original_price" db:"original_price"`
	PaidPrice          float64   `json:"paid_price" db:"paid_price"`
	DiscountAmount     float64   `json:"discount_amount" db:"discount_amount"`
	DiscountReason     *string   `json:"discount_reason,omitempty" db:"discount_reason"`
	PaymentMethod      string    `json:"payment_method" db:"payment_method"`
	PaymentStatus      string    `json:"payment_status" db:"payment_status"`
	PaymentDate        time.Time `json:"payment_date" db:"payment_date"`
	ExpiresAt          time.Time `json:"expires_at" db:"expires_at"`
	Notes              *string   `json:"notes,omitempty" db:"notes"`
	ProcessedByAdminID *int      `json:"processed_by_admin_id,omitempty" db:"processed_by_admin_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// UserWithPayment represents user data with payment information
type UserWithPayment struct {
	User          *User          `json:"user"`
	PaymentRecord *PaymentRecord `json:"payment_record,omitempty"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Status    UserStatus `json:"status,omitempty"`
	UserLevel UserLevel  `json:"user_level,omitempty"`
}

// PaymentData represents payment information for premium subscriptions
type PaymentData struct {
	OriginalPrice  float64    `json:"original_price"`
	PaidPrice      float64    `json:"paid_price"`
	DiscountAmount *float64   `json:"discount_amount,omitempty"`
	DiscountReason *string    `json:"discount_reason,omitempty"`
	PaymentMethod  *string    `json:"payment_method,omitempty"`
	PaymentDate    *time.Time `json:"payment_date,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Basic CRUD operations
	FindByEmail(email string) (*User, error)
	FindByID(id int) (*User, error)
	Create(req *CreateUserRequest) (*User, error)

	// Password operations
	ValidatePassword(plainPassword, hashedPassword string) error
	UpdatePassword(userID int, newPassword string) (*User, error)

	// User level and status management
	UpdateUserLevel(userID int, userLevel UserLevel, paymentData *PaymentData, processedByAdminID *int) (*UserWithPayment, error)
	UpdateUserStatus(userID int, status UserStatus) (*User, error)

	// Bulk operations
	GetAllUsers(page, limit int, status *UserStatus, userLevel *UserLevel, emailFilter *string) (*UsersResponse, error)
	DowngradeExpiredUsers() (*DowngradeResponse, error)
	GetExpiredUsers() ([]*User, error)
}

// UsersResponse represents paginated users response
type UsersResponse struct {
	Users      []*User         `json:"users"`
	Pagination *PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	CurrentPage int  `json:"current_page"`
	HasMore     bool `json:"has_more"`
	Limit       int  `json:"limit"`
}

// DowngradeResponse represents the response from downgrading expired users
type DowngradeResponse struct {
	DowngradedCount int     `json:"downgraded_count"`
	DowngradedUsers []*User `json:"downgraded_users"`
}

// userRepository implements UserRepository interface
type userRepository struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	// Cost 10 = ~65ms, Cost 11 = ~130ms, Cost 12 = ~400ms
	// Using cost 10 for better performance while maintaining good security
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// FindByEmail finds a user by email address
func (r *userRepository) FindByEmail(email string) (*User, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var user User
	query := `SELECT id, email, password, status, user_level, premium_expires_at, created_at, updated_at 
			  FROM users WHERE email = $1`

	err := db.Get(&user, query, email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return &user, nil
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(id int) (*User, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var user User
	query := `SELECT id, email, status, user_level, premium_expires_at, created_at, updated_at 
			  FROM users WHERE id = $1`

	err := db.Get(&user, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	return &user, nil
}

// Create creates a new user
func (r *userRepository) Create(req *CreateUserRequest) (*User, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Set defaults if not provided
	if req.Status == "" {
		req.Status = UserStatusActive
	}
	if req.UserLevel == "" {
		req.UserLevel = UserLevelFree
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	var user User
	query := `INSERT INTO users (email, password, status, user_level) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, email, status, user_level, premium_expires_at, created_at, updated_at`

	err = db.Get(&user, query, req.Email, hashedPassword, req.Status, req.UserLevel)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return nil, fmt.Errorf("user with this email already exists")
		}
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

// ValidatePassword checks if a plain password matches the hashed password
func (r *userRepository) ValidatePassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// UpdatePassword updates user's password
func (r *userRepository) UpdatePassword(userID int, newPassword string) (*User, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	var user User
	query := `UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $2 
			  RETURNING id, email, status, user_level, premium_expires_at, created_at, updated_at`

	err = db.Get(&user, query, hashedPassword, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating password: %w", err)
	}

	return &user, nil
}

// UpdateUserLevel updates user level and handles premium subscription logic
func (r *userRepository) UpdateUserLevel(userID int, userLevel UserLevel, paymentData *PaymentData, processedByAdminID *int) (*UserWithPayment, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Validate user level
	validLevels := []UserLevel{UserLevelFree, UserLevelPremium, UserLevelPremiumPlus}
	isValidLevel := false
	for _, level := range validLevels {
		if userLevel == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		return nil, fmt.Errorf("invalid user level")
	}

	// If payment data is provided for premium tiers, ensure it's valid
	if paymentData != nil && userLevel == UserLevelFree {
		return nil, fmt.Errorf("payment data can only be provided for premium tiers")
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var premiumExpiresAt *time.Time

	// Get current user data to check existing subscription
	var currentUser User
	err = tx.Get(&currentUser, "SELECT user_level, premium_expires_at FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Set expiration dates based on user level
	if userLevel == UserLevelPremium || userLevel == UserLevelPremiumPlus {
		now := time.Now()

		// Check if user has a valid (non-expired) premium subscription to extend from
		hasValidSubscription := (currentUser.UserLevel == UserLevelPremium || currentUser.UserLevel == UserLevelPremiumPlus) &&
			currentUser.PremiumExpiresAt != nil &&
			currentUser.PremiumExpiresAt.After(now)

		// Start from current expiration if valid, otherwise start from now
		var baseDate time.Time
		if hasValidSubscription {
			baseDate = *currentUser.PremiumExpiresAt
		} else {
			baseDate = now
		}

		// Add duration based on subscription type
		if userLevel == UserLevelPremium {
			baseDate = baseDate.AddDate(0, 1, 0) // Add 1 month
		} else {
			baseDate = baseDate.AddDate(1, 0, 0) // Add 1 year
		}

		premiumExpiresAt = &baseDate

	}

	// Update user level
	var updatedUser User
	query := `UPDATE users SET user_level = $1, premium_expires_at = $2, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $3 
			  RETURNING id, email, user_level, premium_expires_at, status, created_at, updated_at`

	err = tx.Get(&updatedUser, query, userLevel, premiumExpiresAt, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating user level: %w", err)
	}

	// Create payment record if payment data is provided
	var paymentRecord *PaymentRecord
	if paymentData != nil && premiumExpiresAt != nil {
		paymentMethod := "manual"
		if paymentData.PaymentMethod != nil {
			paymentMethod = *paymentData.PaymentMethod
		}

		paymentDate := time.Now()
		if paymentData.PaymentDate != nil {
			paymentDate = *paymentData.PaymentDate
		}

		discountAmount := 0.0
		if paymentData.DiscountAmount != nil {
			discountAmount = *paymentData.DiscountAmount
		}

		paymentQuery := `INSERT INTO payment_records 
						(user_id, subscription_type, original_price, paid_price, discount_amount, 
						 discount_reason, payment_method, payment_status, payment_date, expires_at, 
						 notes, processed_by_admin_id) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
						RETURNING *`

		var pr PaymentRecord
		err = tx.Get(&pr, paymentQuery, userID, userLevel, paymentData.OriginalPrice,
			paymentData.PaidPrice, discountAmount, paymentData.DiscountReason,
			paymentMethod, "completed", paymentDate, *premiumExpiresAt,
			paymentData.Notes, processedByAdminID)
		if err != nil {
			return nil, fmt.Errorf("error creating payment record: %w", err)
		}
		paymentRecord = &pr
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &UserWithPayment{
		User:          &updatedUser,
		PaymentRecord: paymentRecord,
	}, nil
}

// UpdateUserStatus updates user status
func (r *userRepository) UpdateUserStatus(userID int, status UserStatus) (*User, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	validStatuses := []UserStatus{UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusBanned}
	isValidStatus := false
	for _, s := range validStatuses {
		if status == s {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return nil, fmt.Errorf("invalid user status")
	}

	var user User
	query := `UPDATE users SET status = $1, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $2 
			  RETURNING id, email, status, user_level, premium_expires_at, created_at, updated_at`

	err := db.Get(&user, query, status, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating user status: %w", err)
	}

	return &user, nil
}

// GetAllUsers retrieves paginated users with optional filters
func (r *userRepository) GetAllUsers(page, limit int, status *UserStatus, userLevel *UserLevel, emailFilter *string) (*UsersResponse, error) {
	db := GetDB().PostgreDBManager.RC
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	offset := (page - 1) * limit
	args := []interface{}{}
	argCount := 0

	// Build base query
	query := `SELECT id, email, status, user_level, premium_expires_at, created_at, updated_at 
			  FROM users WHERE 1=1`

	// Add status filter if provided
	if status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *status)
	}

	// Add user level filter if provided
	if userLevel != nil {
		argCount++
		query += fmt.Sprintf(" AND user_level = $%d", argCount)
		args = append(args, *userLevel)
	}

	// Add email filter if provided
	if emailFilter != nil {
		argCount++
		query += fmt.Sprintf(" AND email ILIKE $%d", argCount)
		args = append(args, *emailFilter+"%")
	}

	// Add ordering and pagination
	argCount++
	args = append(args, limit+1) // Fetch one extra to check if there's more data
	argCount++
	args = append(args, offset)
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount-1, argCount)

	var users []*User
	err := db.Select(&users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	// Check if there are more results
	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit] // Remove the extra record
	}

	return &UsersResponse{
		Users: users,
		Pagination: &PaginationInfo{
			CurrentPage: page,
			HasMore:     hasMore,
			Limit:       limit,
		},
	}, nil
}

// DowngradeExpiredUsers downgrades users whose premium subscription has expired
func (r *userRepository) DowngradeExpiredUsers() (*DowngradeResponse, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var users []*User
	query := `UPDATE users 
			  SET user_level = 'free', premium_expires_at = NULL, updated_at = CURRENT_TIMESTAMP 
			  WHERE user_level IN ('premium', 'premium+') 
			  AND premium_expires_at IS NOT NULL 
			  AND premium_expires_at <= CURRENT_TIMESTAMP
			  RETURNING id, email, user_level, status, premium_expires_at, created_at, updated_at`

	err := db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("error downgrading expired users: %w", err)
	}

	return &DowngradeResponse{
		DowngradedCount: len(users),
		DowngradedUsers: users,
	}, nil
}

// GetExpiredUsers returns users whose premium subscription has expired
func (r *userRepository) GetExpiredUsers() ([]*User, error) {
	db := GetDB().PostgreDBManager.RC
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var users []*User
	query := `SELECT id, email, user_level, premium_expires_at, status, created_at, updated_at
			  FROM users 
			  WHERE user_level IN ('premium', 'premium+') 
			  AND premium_expires_at IS NOT NULL 
			  AND premium_expires_at <= CURRENT_TIMESTAMP
			  ORDER BY premium_expires_at DESC`

	err := db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching expired users: %w", err)
	}

	return users, nil
}
