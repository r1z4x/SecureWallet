package routes

import (
	"net/http"
	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes sets up admin routes
func SetupAdminRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	{
		// SECURE: Add authentication and admin authorization middleware
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())

		admin.GET("/dashboard", getDashboard)
		admin.GET("/users", getAdminUsers)
		admin.GET("/transactions", getAdminTransactions)
		admin.POST("/users/:id/disable", disableUser)
		admin.POST("/users/:id/enable", enableUser)
		admin.GET("/settings", getSystemSettings)
		admin.POST("/settings", saveSystemSettings)
	}
}

// getDashboard gets admin dashboard
func getDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Admin dashboard"})
}

// getAdminUsers gets all users for admin
func getAdminUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all users for admin"})
}

// getAdminTransactions gets all transactions for admin (system-wide)
func getAdminTransactions(c *gin.Context) {
	db := config.GetDB()

	// Get limit parameter
	limitStr := c.Query("limit")
	limit := 100 // default limit for admin
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get ALL transactions from the system (not just current user's)
	var transactions []models.Transaction
	if err := db.Preload("Wallet.User"). // Include user information
						Order("created_at DESC").
						Limit(limit).
						Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// Transform transactions to include transaction_type field for frontend compatibility
	type TransactionResponse struct {
		ID              string    `json:"id"`
		WalletID        string    `json:"wallet_id"`
		TransactionType string    `json:"transaction_type"`
		Amount          float64   `json:"amount"`
		Currency        string    `json:"currency"`
		Description     string    `json:"description"`
		Status          string    `json:"status"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		Wallet          struct {
			User struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"wallet"`
	}

	var response []TransactionResponse
	for _, t := range transactions {
		response = append(response, TransactionResponse{
			ID:              t.ID.String(),
			WalletID:        t.WalletID.String(),
			TransactionType: t.Type,
			Amount:          t.Amount,
			Currency:        t.Currency,
			Description:     t.Description,
			Status:          t.Status,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
			Wallet: struct {
				User struct {
					Email string `json:"email"`
				} `json:"user"`
			}{
				User: struct {
					Email string `json:"email"`
				}{
					Email: t.Wallet.User.Email,
				},
			},
		})
	}

	c.JSON(http.StatusOK, response)
}

// disableUser disables a user
func disableUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Disable user", "id": id})
}

// enableUser enables a user
func enableUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Enable user", "id": id})
}

// getSystemSettings gets current system settings
func getSystemSettings(c *gin.Context) {
	// TODO: Implement getting settings from database
	settings := gin.H{
		"security": gin.H{
			"twoFactorEnabled": true,
			"sessionTimeout":   30,
			"passwordPolicy": gin.H{
				"minLength":           8,
				"requireUppercase":    true,
				"requireLowercase":    true,
				"requireNumbers":      true,
				"requireSpecialChars": true,
			},
		},
		"transactionLimits": gin.H{
			"dailyTransferLimit":   10000,
			"monthlyTransferLimit": 50000,
			"minTransferAmount":    1,
		},
	}
	c.JSON(http.StatusOK, settings)
}

// saveSystemSettings saves system settings
func saveSystemSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
		return
	}

	// TODO: Implement saving settings to database
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message":  "Settings saved successfully",
		"settings": settings,
	})
}
