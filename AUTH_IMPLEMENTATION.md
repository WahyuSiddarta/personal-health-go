# Authentication System Implementation

This Go implementation provides the same authentication functionality as the JavaScript version you provided, with the following features:

## Features Implemented

### 1. User Management

- User registration and login
- Password hashing with bcrypt (cost 12)
- JWT token generation and validation
- User profile management
- User status management (active, inactive, suspended, banned)
- Premium subscription levels (free, premium, premium+)

### 2. Authentication & Authorization

- JWT-based authentication
- Token validation middleware
- Role-based access control (admin features)
- Premium subscription validation
- Optional authentication middleware

### 3. API Endpoints

#### Public Endpoints (No Authentication Required)

- `POST /api/public/auth/login` - User login
- `POST /api/public/auth/register` - User registration

#### Protected Endpoints (Authentication Required)

- `GET /api/protected/user/profile` - Get current user profile
- `POST /api/protected/user/logout` - User logout

#### Admin Endpoints (Admin Access Required)

- `GET /api/protected/admin/users` - Get all users with pagination and filters
- `PUT /api/protected/admin/users/:id/level` - Update user subscription level
- `PUT /api/protected/admin/users/:id/status` - Update user status
- `GET /api/protected/admin/users/expired` - Get users with expired subscriptions
- `POST /api/protected/admin/users/downgrade-expired` - Downgrade expired users

## File Structure

```
├── models/
│   ├── user.go                 # User model with database operations
│   └── model.common.go         # Common model utilities
├── api/
│   ├── auth.go                 # Authentication handlers
│   └── api.common.go           # Common API utilities
├── middleware/
│   ├── auth.go                 # JWT authentication middleware
│   └── [other middleware files]
├── validation/
│   └── validation.go           # Request validation logic
├── router/
│   ├── public.go               # Public routes setup
│   ├── protected.go            # Protected routes setup
│   └── router.go               # Main router configuration
├── config/
│   └── config.go               # Configuration management
└── .env.example                # Environment variables example
```

## Environment Variables

Key environment variables needed for authentication:

```bash
# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# Database connection details (already configured)
DB_RW_HOST=localhost
DB_RW_PORT=5432
DB_RW_USER=postgres
DB_RW_PASSWORD=your_password
DB_RW_NAME=saham_db
```

## Usage Examples

### 1. User Registration

```bash
curl -X POST http://localhost:8080/api/public/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/api/public/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/protected/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Admin - Update User Level

```bash
curl -X PUT http://localhost:8080/api/protected/admin/users/123/level \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_level": "premium",
    "payment_data": {
      "original_price": 100.00,
      "paid_price": 80.00,
      "discount_amount": 20.00,
      "discount_reason": "Early bird discount",
      "payment_method": "manual"
    }
  }'
```

## Key Features from JavaScript Implementation

✅ **Password Hashing**: Uses bcrypt with salt rounds 12  
✅ **JWT Token Generation**: Configurable expiration time  
✅ **User Status Management**: active, inactive, suspended, banned  
✅ **Premium Subscription Logic**: Handles extension vs new subscription  
✅ **Payment Records**: Tracks premium subscription payments  
✅ **Admin Functions**: User management and expired user handling  
✅ **Validation**: Comprehensive request/response validation  
✅ **Error Handling**: Structured error responses  
✅ **Logging**: Integrated with existing logger

## Middleware Usage

- `middleware.RequireAuth()` - Requires valid JWT authentication
- `middleware.RequirePremium()` - Requires active premium subscription
- `middleware.RequirePremiumPlus()` - Requires active premium+ subscription
- `middleware.AdminRequired()` - Requires admin privileges
- `middleware.OptionalAuth()` - Optional authentication (doesn't fail if no auth)

## Next Steps

1. Set the `JWT_SECRET` environment variable to a secure secret
2. Update admin email list in `middleware/auth.go` AdminRequired() function
3. Test the endpoints with your existing database
4. Implement any additional business logic specific to your application

The implementation maintains the same structure and functionality as your JavaScript version while leveraging Go's type safety and performance benefits.
