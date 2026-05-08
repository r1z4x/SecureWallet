package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes sets up wallet routes
func SetupWalletRoutes(router *gin.RouterGroup) {
	wallets := router.Group("/wallets")
	wallets.Use(middleware.AuthMiddleware())
	{
		wallets.GET("/", middleware.RequirePermission(models.PermWalletRead), getWallets)
		wallets.GET("/balance", middleware.RequirePermission(models.PermWalletRead), getBalance)
		wallets.POST("/deposit", middleware.RequirePermission(models.PermWalletWrite), deposit)
		wallets.POST("/transfer", middleware.RequirePermission(models.PermTransferWrite), middleware.IdempotencyMiddleware("transfer"), transfer)
		wallets.GET("/:id", middleware.RequirePermission(models.PermWalletRead), getWallet)
		wallets.POST("/", middleware.RequirePermission(models.PermWalletWrite), createWallet)
		wallets.PUT("/:id", middleware.RequirePermission(models.PermWalletWrite), updateWallet)
		wallets.DELETE("/:id", middleware.RequirePermission(models.PermWalletDelete), deleteWallet)
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
// @Summary Deposit funds
// @Description Deposit funds into the authenticated user's wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param body body DepositRequest true "Deposit data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Security BearerAuth
// @Router /wallets/deposit [post]
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
	audit := services.NewAuditLogger().WithGinContext(c)

	svc := services.NewTransferService()
	result, err := svc.Deposit(services.DepositRequest{
		UserID:      currentUser.ID,
		Amount:      depositReq.Amount,
		Description: depositReq.Description,
	})
	if err != nil {
		if errors.Is(err, services.ErrWalletNotFound) {
			audit.Log(currentUser.ID, "DEPOSIT", "wallet", "Wallet not found", services.AuditResultFailure)
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		if errors.Is(err, services.ErrInvalidAmount) {
			audit.Log(currentUser.ID, "DEPOSIT", "wallet", err.Error(), services.AuditResultFailure)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		audit.Log(currentUser.ID, "DEPOSIT", "wallet", "Deposit failed", services.AuditResultFailure)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process deposit"})
		return
	}

	audit.Log(currentUser.ID, "DEPOSIT", "wallet", fmt.Sprintf("Deposited %.2f %s", result.Amount, result.TransactionID), services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deposit successful",
		"wallet": gin.H{
			"id":      result.WalletID,
			"balance": result.NewBalance,
		},
		"transaction": gin.H{
			"id":     result.TransactionID,
			"amount": result.Amount,
		},
	})
}

// transfer handles wallet transfer between users
// @Summary Transfer funds
// @Description Transfer funds to another user's wallet by email. Requires Idempotency-Key header for safe retries.
// @Tags wallets
// @Accept json
// @Produce json
// @Param Idempotency-Key header string false "Idempotency key for safe retries"
// @Param body body TransferRequest true "Transfer data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 409 {object} gin.H
// @Security BearerAuth
// @Router /wallets/transfer [post]
func transfer(c *gin.Context) {
	var transferReq TransferRequest
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if transferReq.Amount < MIN_TRANSFER_AMOUNT || transferReq.Amount > MAX_TRANSFER_AMOUNT {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Transfer amount must be between $%.2f and $%.2f", MIN_TRANSFER_AMOUNT, MAX_TRANSFER_AMOUNT),
		})
		return
	}

	if !strings.Contains(transferReq.Recipient, "@") || len(transferReq.Recipient) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient email format"})
		return
	}

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
	audit := services.NewAuditLogger().WithGinContext(c)

	idempotencyKey := c.GetHeader("Idempotency-Key")

	svc := services.NewTransferService()
	result, err := svc.Transfer(services.TransferRequest{
		SenderUserID:   currentUser.ID,
		RecipientEmail: transferReq.Recipient,
		Amount:         transferReq.Amount,
		Description:    transferReq.Description,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInsufficientBalance):
			fee := svc.CalculateFee(transferReq.Amount)
			audit.Log(currentUser.ID, "TRANSFER", "wallet", fmt.Sprintf("Insufficient balance for transfer to %s", transferReq.Recipient), services.AuditResultFailure)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient balance",
				"details": gin.H{
					"transfer_amount": transferReq.Amount,
					"transfer_fee":    fee,
					"total_amount":    transferReq.Amount + fee,
				},
			})
		case errors.Is(err, services.ErrWalletNotFound):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", "Wallet not found", services.AuditResultFailure)
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		case errors.Is(err, services.ErrRecipientNotFound):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", fmt.Sprintf("Recipient not found: %s", transferReq.Recipient), services.AuditResultFailure)
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		case errors.Is(err, services.ErrSelfTransfer):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", "Self-transfer attempted", services.AuditResultDenied)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to yourself"})
		case errors.Is(err, services.ErrDuplicateRequest):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", "Duplicate request", services.AuditResultDenied)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Duplicate request"})
		case errors.Is(err, services.ErrCurrencyMismatch):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", "Currency mismatch", services.AuditResultFailure)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Currency mismatch between sender and recipient"})
		case errors.Is(err, services.ErrInvalidAmount):
			audit.Log(currentUser.ID, "TRANSFER", "wallet", err.Error(), services.AuditResultFailure)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("Transfer error: %v", err)
			audit.Log(currentUser.ID, "TRANSFER", "wallet", "Transfer failed", services.AuditResultFailure)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transfer"})
		}
		return
	}

	audit.Log(currentUser.ID, "TRANSFER", "wallet", fmt.Sprintf("Transferred %.2f to %s", result.Amount, transferReq.Recipient), services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"sender_wallet": gin.H{
			"id":      result.SenderWalletID,
			"balance": result.SenderBalance,
		},
		"recipient_wallet": gin.H{
			"id":      result.RecipientWalletID,
			"balance": result.RecipientBalance,
		},
		"transfer": gin.H{
			"id":           result.TransactionID,
			"amount":       result.Amount,
			"transfer_fee": result.Fee,
			"total_amount": result.TotalDeducted,
			"description":  transferReq.Description,
			"status":       "completed",
		},
	})
}

// getBalance gets the current user's wallet balance
// @Summary Get wallet balance
// @Description Get the authenticated user's wallet balance and transaction count
// @Tags wallets
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Security BearerAuth
// @Router /wallets/balance [get]
func getBalance(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	var userWallet models.Wallet
	if err := db.Select("id", "balance", "currency").Where("user_id = ?", currentUser.ID).First(&userWallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	var transactionCount int64
	db.Model(&models.Transaction{}).Where("wallet_id = ?", userWallet.ID).Count(&transactionCount)

	c.JSON(http.StatusOK, gin.H{
		"balance":           userWallet.Balance,
		"currency":          userWallet.Currency,
		"transaction_count": transactionCount,
	})
}

// getWallets gets all wallets for the current user
// @Summary List wallets
// @Description Get all wallets owned by the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Success 200 {array} models.Wallet
// @Failure 401 {object} gin.H
// @Security BearerAuth
// @Router /wallets [get]
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
// @Summary Get wallet by ID
// @Description Get a specific wallet by ID (must be owned by authenticated user)
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Security BearerAuth
// @Router /wallets/{id} [get]
func getWallet(c *gin.Context) {
	id := c.Param("id")

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	var wallet models.Wallet
	if err := db.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	audit := services.NewAuditLogger().WithGinContext(c)
	audit.Log(currentUser.ID, "WALLET_READ", "wallet", fmt.Sprintf("Accessed wallet %s", id), services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{
		"wallet": gin.H{
			"id":         wallet.ID,
			"balance":    wallet.Balance,
			"currency":   wallet.Currency,
			"created_at": wallet.CreatedAt,
			"updated_at": wallet.UpdatedAt,
		},
	})
}

// createWallet creates a new wallet
// @Summary Create wallet
// @Description Create a new wallet for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Success 201 {object} gin.H
// @Security BearerAuth
// @Router /wallets [post]
func createWallet(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create wallet"})
}

// updateWallet updates a wallet
// @Summary Update wallet
// @Description Update a wallet's attributes (must be owned by authenticated user)
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /wallets/{id} [put]
func updateWallet(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update wallet", "id": id})
}

// deleteWallet deletes a wallet
// @Summary Delete wallet
// @Description Delete a wallet (must be owned by authenticated user, must have zero balance)
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /wallets/{id} [delete]
func deleteWallet(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete wallet", "id": id})
}
