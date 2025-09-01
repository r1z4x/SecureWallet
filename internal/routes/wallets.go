package routes

import (
	"fmt"
	"net/http"
	"strings"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes sets up wallet routes
func SetupWalletRoutes(router *gin.RouterGroup) {
	wallets := router.Group("/wallets")
	{
		wallets.GET("/", middleware.AuthMiddleware(), getWallets)
		wallets.GET("/balance", middleware.AuthMiddleware(), getBalance)
		wallets.POST("/deposit", middleware.AuthMiddleware(), deposit)
		wallets.POST("/transfer", middleware.AuthMiddleware(), transfer)
		wallets.GET("/:id", middleware.AuthMiddleware(), getWallet)
		wallets.POST("/", middleware.AuthMiddleware(), createWallet)
		wallets.PUT("/:id", middleware.AuthMiddleware(), updateWallet)
		wallets.DELETE("/:id", middleware.AuthMiddleware(), deleteWallet)
	}
}

// DepositRequest represents deposit request data
type DepositRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// TransferRequest represents transfer request data
type TransferRequest struct {
	Recipient   string  `json:"recipient" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// Transfer fee constants
const (
	TRANSFER_FEE_PERCENTAGE = 0.01   // 1% transfer fee
	MIN_TRANSFER_FEE        = 1.0    // Minimum $1 fee
	MAX_TRANSFER_FEE        = 50.0   // Maximum $50 fee
	MIN_TRANSFER_AMOUNT     = 1.0    // Minimum transfer amount
	MAX_TRANSFER_AMOUNT     = 1000.0 // Maximum transfer amount
)

// deposit handles wallet deposit
func deposit(c *gin.Context) {
	var depositReq DepositRequest
	if err := c.ShouldBindJSON(&depositReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Start transaction
	tx := db.Begin()

	// Update wallet balance
	if err := tx.Model(&userWallet).Update("balance", userWallet.Balance+depositReq.Amount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		WalletID:    userWallet.ID,
		Type:        "DEPOSIT",
		Amount:      depositReq.Amount,
		Currency:    userWallet.Currency,
		Description: depositReq.Description,
		Status:      "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Get updated wallet
	var updatedWallet models.Wallet
	db.First(&updatedWallet, userWallet.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deposit successful",
		"wallet":  updatedWallet,
		"transaction": gin.H{
			"id":          transaction.ID,
			"amount":      transaction.Amount,
			"description": transaction.Description,
			"status":      transaction.Status,
		},
	})
}

// transfer handles wallet transfer between users
func transfer(c *gin.Context) {
	var transferReq TransferRequest
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SECURE: Validate transfer amount
	if transferReq.Amount < MIN_TRANSFER_AMOUNT || transferReq.Amount > MAX_TRANSFER_AMOUNT {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Transfer amount must be between $%.2f and $%.2f", MIN_TRANSFER_AMOUNT, MAX_TRANSFER_AMOUNT),
		})
		return
	}

	// SECURE: Validate recipient email format
	if !strings.Contains(transferReq.Recipient, "@") || len(transferReq.Recipient) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient email format"})
		return
	}

	// SECURE: Validate description length
	if len(transferReq.Description) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description too long (max 255 characters)"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	// Get sender's wallet
	var senderWallet models.Wallet
	if err := db.Where("user_id = ?", currentUser.ID).First(&senderWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender wallet not found"})
		return
	}

	// Calculate transfer fee
	transferFee := transferReq.Amount * TRANSFER_FEE_PERCENTAGE
	if transferFee < MIN_TRANSFER_FEE {
		transferFee = MIN_TRANSFER_FEE
	} else if transferFee > MAX_TRANSFER_FEE {
		transferFee = MAX_TRANSFER_FEE
	}

	totalAmount := transferReq.Amount + transferFee

	// Check if sender has sufficient balance (including fee)
	if senderWallet.Balance < totalAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Insufficient balance",
			"details": gin.H{
				"transfer_amount": transferReq.Amount,
				"transfer_fee":    transferFee,
				"total_amount":    totalAmount,
				"current_balance": senderWallet.Balance,
			},
		})
		return
	}

	// Find recipient by email
	var recipient models.User
	if err := db.Where("email = ?", transferReq.Recipient).First(&recipient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	// Get recipient's wallet
	var recipientWallet models.Wallet
	if err := db.Where("user_id = ?", recipient.ID).First(&recipientWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient wallet not found"})
		return
	}

	// Prevent self-transfer
	if senderWallet.ID == recipientWallet.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to yourself"})
		return
	}

	// Start database transaction
	tx := db.Begin()

	// Deduct from sender's wallet (amount + fee)
	if err := tx.Model(&senderWallet).Update("balance", senderWallet.Balance-totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender wallet"})
		return
	}

	// Add to recipient's wallet
	if err := tx.Model(&recipientWallet).Update("balance", recipientWallet.Balance+transferReq.Amount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipient wallet"})
		return
	}

	// Create outgoing transaction record for sender (including fee)
	outgoingTransaction := models.Transaction{
		WalletID:    senderWallet.ID,
		Type:        "TRANSFER",
		Amount:      totalAmount,
		Currency:    senderWallet.Currency,
		Description: transferReq.Description + " (to " + recipient.Username + ") + $" + fmt.Sprintf("%.2f", transferFee) + " fee",
		Status:      "completed",
	}

	if err := tx.Create(&outgoingTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create outgoing transaction"})
		return
	}

	// Create incoming transaction record for recipient
	incomingTransaction := models.Transaction{
		WalletID:    recipientWallet.ID,
		Type:        "TRANSFER",
		Amount:      transferReq.Amount,
		Currency:    recipientWallet.Currency,
		Description: transferReq.Description + " (from " + currentUser.Username + ")",
		Status:      "completed",
	}

	if err := tx.Create(&incomingTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incoming transaction"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transfer"})
		return
	}

	// Get updated wallets
	var updatedSenderWallet, updatedRecipientWallet models.Wallet
	db.First(&updatedSenderWallet, senderWallet.ID)
	db.First(&updatedRecipientWallet, recipientWallet.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"sender_wallet": gin.H{
			"balance":  updatedSenderWallet.Balance,
			"currency": updatedSenderWallet.Currency,
		},
		"recipient": gin.H{
			"username": recipient.Username,
			"email":    recipient.Email,
		},
		"transfer": gin.H{
			"amount":       transferReq.Amount,
			"transfer_fee": transferFee,
			"total_amount": totalAmount,
			"description":  transferReq.Description,
			"status":       "completed",
		},
	})
}

// getBalance gets the current user's wallet balance
func getBalance(c *gin.Context) {
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

	// Get transaction count
	var transactionCount int64
	db.Model(&models.Transaction{}).Where("wallet_id = ?", userWallet.ID).Count(&transactionCount)

	c.JSON(http.StatusOK, gin.H{
		"balance":           userWallet.Balance,
		"currency":          userWallet.Currency,
		"transaction_count": transactionCount,
	})
}

// getWallets gets all wallets for the current user
func getWallets(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	var wallets []models.Wallet
	if err := db.Where("user_id = ?", currentUser.ID).Find(&wallets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// getWallet gets a specific wallet
func getWallet(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get wallet", "id": id})
}

// createWallet creates a new wallet
func createWallet(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create wallet"})
}

// updateWallet updates a wallet
func updateWallet(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update wallet", "id": id})
}

// deleteWallet deletes a wallet
func deleteWallet(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete wallet", "id": id})
}
