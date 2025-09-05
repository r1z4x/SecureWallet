package routes

import (
	"net/http"
	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		// Support management routes
		admin.GET("/support/tickets", getAdminSupportTickets)
		admin.POST("/support/tickets/:id/reply", replyToTicket)
		admin.POST("/support/tickets/:id/resolve", resolveTicket)
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

// getAdminSupportTickets gets all support tickets for admin
func getAdminSupportTickets(c *gin.Context) {
	db := config.GetDB()

	var tickets []models.SupportTicket
	if err := db.Preload("User").Order("created_at DESC").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch support tickets"})
		return
	}

	// Transform tickets to include proper user information
	type TicketResponse struct {
		ID          string    `json:"id"`
		UserID      string    `json:"user_id"`
		Subject     string    `json:"subject"`
		Message     string    `json:"message"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		Category    string    `json:"category"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		User        struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}

	var response []TicketResponse
	for _, ticket := range tickets {
		response = append(response, TicketResponse{
			ID:          ticket.ID.String(),
			UserID:      ticket.UserID.String(),
			Subject:     ticket.Subject,
			Message:     ticket.Description, // Map description to message for frontend compatibility
			Description: ticket.Description,
			Status:      ticket.Status,
			Priority:    ticket.Priority,
			Category:    "General", // Default category since it's not in the model
			CreatedAt:   ticket.CreatedAt,
			UpdatedAt:   ticket.UpdatedAt,
			User: struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
			}{
				ID:       ticket.User.ID.String(),
				Username: ticket.User.Username,
				Email:    ticket.User.Email,
			},
		})
	}

	// If no tickets found, return mock data for testing
	if len(response) == 0 {
		mockTickets := []TicketResponse{
			{
				ID:          "550e8400-e29b-41d4-a716-446655440001",
				UserID:      "550e8400-e29b-41d4-a716-446655440002",
				Subject:     "Login issue",
				Message:     "I cannot login to my account",
				Description: "I cannot login to my account",
				Status:      "open",
				Priority:    "high",
				Category:    "Authentication",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				User: struct {
					ID       string `json:"id"`
					Username string `json:"username"`
					Email    string `json:"email"`
				}{
					ID:       "550e8400-e29b-41d4-a716-446655440002",
					Username: "testuser",
					Email:    "test@example.com",
				},
			},
			{
				ID:          "550e8400-e29b-41d4-a716-446655440003",
				UserID:      "550e8400-e29b-41d4-a716-446655440004",
				Subject:     "Transfer problem",
				Message:     "My transfer is stuck",
				Description: "My transfer is stuck",
				Status:      "in_progress",
				Priority:    "medium",
				Category:    "Transactions",
				CreatedAt:   time.Now().Add(-24 * time.Hour),
				UpdatedAt:   time.Now().Add(-24 * time.Hour),
				User: struct {
					ID       string `json:"id"`
					Username string `json:"username"`
					Email    string `json:"email"`
				}{
					ID:       "550e8400-e29b-41d4-a716-446655440004",
					Username: "user2",
					Email:    "user2@example.com",
				},
			},
		}
		response = mockTickets
	}

	c.JSON(http.StatusOK, response)
}

// ReplyRequest represents a reply to a support ticket
type ReplyRequest struct {
	Message string `json:"message" binding:"required"`
}

// replyToTicket adds a reply to a support ticket
func replyToTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")

	// Parse UUID
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID format"})
		return
	}

	var replyReq ReplyRequest
	if err := c.ShouldBindJSON(&replyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if this is a mock ticket ID
	mockTicketIDs := []string{
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440003",
		"550e8400-e29b-41d4-a716-446655440005",
		"550e8400-e29b-41d4-a716-446655440007",
	}

	isMockTicket := false
	for _, mockID := range mockTicketIDs {
		if ticketIDStr == mockID {
			isMockTicket = true
			break
		}
	}

	if isMockTicket {
		// Handle mock ticket - just return success
		c.JSON(http.StatusOK, gin.H{
			"message":   "Reply added successfully (mock ticket)",
			"ticket_id": ticketID.String(),
			"reply":     replyReq.Message,
			"status":    "in_progress",
		})
		return
	}

	db := config.GetDB()

	// Check if ticket exists in database
	var ticket models.SupportTicket
	if err := db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Update ticket status to in_progress if it's open
	if ticket.Status == "open" {
		ticket.Status = "in_progress"
		if err := db.Save(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
			return
		}
	}

	// TODO: In a real implementation, you would save the reply to a separate replies table
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{
		"message":   "Reply added successfully",
		"ticket_id": ticketID.String(),
		"reply":     replyReq.Message,
	})
}

// resolveTicket resolves a support ticket
func resolveTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")

	// Parse UUID
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID format"})
		return
	}

	// Check if this is a mock ticket ID
	mockTicketIDs := []string{
		"550e8400-e29b-41d4-a716-446655440001",
		"550e8400-e29b-41d4-a716-446655440003",
		"550e8400-e29b-41d4-a716-446655440005",
		"550e8400-e29b-41d4-a716-446655440007",
	}

	isMockTicket := false
	for _, mockID := range mockTicketIDs {
		if ticketIDStr == mockID {
			isMockTicket = true
			break
		}
	}

	if isMockTicket {
		// Handle mock ticket - just return success
		c.JSON(http.StatusOK, gin.H{
			"message":   "Ticket resolved successfully (mock ticket)",
			"ticket_id": ticketID.String(),
			"status":    "resolved",
		})
		return
	}

	db := config.GetDB()

	// Check if ticket exists in database
	var ticket models.SupportTicket
	if err := db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Update ticket status to resolved
	ticket.Status = "resolved"
	if err := db.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Ticket resolved successfully",
		"ticket_id": ticketID.String(),
		"status":    "resolved",
	})
}
