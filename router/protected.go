package router

import (
	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(apiGroup *echo.Group) {
	// Initialize auth handlers
	userRepo := models.NewUserRepository()
	authHandlers := api.NewAuthHandlers(userRepo)

	userGroup := apiGroup.Group("/users")
	userGroup.Use(middleware.RequireAuth()) // Add authentication middleware to protected routes
	userGroup.GET("/profile", authHandlers.GetProfile)
	portfolioGroup := userGroup.Group("/portfolio")
	setupCashPortfolioRoutes(portfolioGroup) // Setup CashPortfolio routes (includes PnL)
	setupBondPortfolioRoutes(portfolioGroup) // Setup BondPortfolio routes

	// Setup admin routes
	setupAdminRoutes(apiGroup, authHandlers)

}

// setupAdminRoutes configures admin routes (admin authentication required)
func setupAdminRoutes(rprotected *echo.Group, authHandlers *api.AuthHandlers) {
	adminGroup := rprotected.Group("/admin")
	adminGroup.Use(middleware.AdminRequired())

	// User management endpoints, accessible at /api/admin/users
	usersGroup := adminGroup.Group("/users")
	usersGroup.GET("", authHandlers.GetAllUsers, validator.ValidateQuery(&validator.GetUsersQuery{}))
	usersGroup.PUT("/:id/level", authHandlers.UpdateUserLevel, validator.ValidateRequest(&validator.UpdateUserLevelRequest{}))
	usersGroup.PUT("/:id/status", authHandlers.UpdateUserStatus, validator.ValidateRequest(&validator.UpdateUserStatusRequest{}))
	usersGroup.GET("/expired", authHandlers.GetExpiredUsers)
}

// setupCashPortfolioRoutes configures portfolio cash routes
func setupCashPortfolioRoutes(portfolioGroup *echo.Group) {
	// Initialize portfolio handlers
	portfolioRepo := models.NewPortfolioCashRepository()
	portfolioHandlers := api.NewPortfolioCashHandlers(portfolioRepo)

	// Create cash portfolio - accessible at /api/users/portfolio/cash
	cashGroup := portfolioGroup.Group("/cash")
	cashGroup.POST("", portfolioHandlers.CreateCashPortfolio, validator.ValidateRequest(&validator.CreatePortfolioCashRequest{}))
	cashGroup.GET("", portfolioHandlers.GetMyCashPortfolios)
	cashGroup.PUT("/:id", portfolioHandlers.UpdateCashPortfolio, validator.ValidateRequest(&validator.UpdatePortfolioCashRequest{}))
	cashGroup.DELETE("/:id", portfolioHandlers.DeleteCashPortfolio)
	cashGroup.POST("/move", portfolioHandlers.MoveAsset, validator.ValidateRequest(&validator.MoveAssetRequest{}))
	cashGroup.POST("/realize", portfolioHandlers.RealizeCashPortfolio, validator.ValidateRequest(&validator.RealizeCashPortfolioRequest{}))

	// PnL sub-routes under portfolio /api/users/portfolio/cash/pnl
	pnlGroup := cashGroup.Group("/pnl")
	pnlGroup.GET("", portfolioHandlers.GetPnlRealizedCash)
	pnlGroup.GET("/portfolio/:portfolioId", portfolioHandlers.GetPnlByPortfolioCashID)
	pnlGroup.GET("/:id", portfolioHandlers.GetPnlById)
	pnlGroup.POST("", portfolioHandlers.CreatePnlRealizedCash, validator.ValidateRequest(&validator.CreatePnlRealizedCashRequest{}))
	pnlGroup.PUT("/:id", portfolioHandlers.UpdatePnlRealizedCash, validator.ValidateRequest(&validator.UpdatePnlRealizedCashRequest{}))
	pnlGroup.DELETE("/:id", portfolioHandlers.DeletePnlRealizedCash)
}

// setupBondPortfolioRoutes configures portfolio bond routes
func setupBondPortfolioRoutes(portfolioGroup *echo.Group) {
	// Initialize portfolio bond handlers
	bondRepo := models.NewPortfolioBondRepository()
	couponRepo := models.NewPortfolioBondCouponRepository()
	realizedRepo := models.NewPortfolioBondRealizedRepository()
	portfolioBondHandlers := api.NewPortfolioBondHandlers(bondRepo, couponRepo, realizedRepo)

	bondGroup := portfolioGroup.Group("/bond")
	bondGroup.POST("", portfolioBondHandlers.CreateBondPortfolio)
	bondGroup.GET("", portfolioBondHandlers.GetMyBondPortfolios)

	bondGroup.PUT("/:portfolioId", portfolioBondHandlers.UpdateBondPortfolio)
	bondGroup.DELETE("/:portfolioId", portfolioBondHandlers.DeleteBondPortfolio)

	bondGroup.PUT("/:portfolioId/market-price-override", portfolioBondHandlers.UpdateMarketPriceOverride, validator.ValidateRequest(&validator.UpdateMarketPriceOverrideRequest{}))

}
