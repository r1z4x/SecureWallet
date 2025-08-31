package routes

import (
	"net/http"
	"strconv"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupTransactionRoutes sets up transaction routes
func SetupTransactionRoutes(router *gin.RouterGroup) {
	transactions := router.Group("/transactions")
	{
		transactions.GET("/", middleware.AuthMiddleware(), getTransactions)
		transactions.GET("/:id", middleware.AuthMiddleware(), getTransaction)
		transactions.POST("/", middleware.AuthMiddleware(), createTransaction)
		transactions.PUT("/:id", middleware.AuthMiddleware(), updateTransaction)
		transactions.DELETE("/:id", middleware.AuthMiddleware(), deleteTransaction)
	}
}

// getTransactions gets all transactions for the current user
func getTransactions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	// Get user's wallet
	var userWallet models.Wallet
	if err := db.Where("user_id = ?", currentUser.ID).First(&userWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Get limit parameter
	limitStr := c.Query("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get transactions for the user's wallet
	var transactions []models.Transaction
	if err := db.Where("wallet_id = ?", userWallet.ID).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// getTransaction gets a specific transaction
func getTransaction(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get transaction", "id": id})
}

// createTransaction creates a new transaction
func createTransaction(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create transaction"})
}

// updateTransaction updates a transaction
func updateTransaction(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update transaction", "id": id})
}

// deleteTransaction deletes a transaction
func deleteTransaction(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete transaction", "id": id})
}
