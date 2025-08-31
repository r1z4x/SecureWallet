package routes

import (
	"net/http"
	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupSupportRoutes sets up support routes
func SetupSupportRoutes(router *gin.RouterGroup) {
	support := router.Group("/support")
	{
		support.GET("/tickets", middleware.AuthMiddleware(), getTickets)
		support.GET("/tickets/:id", middleware.AuthMiddleware(), getTicket)
		support.POST("/tickets", middleware.AuthMiddleware(), createTicket)
		support.PUT("/tickets/:id", middleware.AuthMiddleware(), updateTicket)
		support.DELETE("/tickets/:id", middleware.AuthMiddleware(), deleteTicket)
	}
}

// getTickets gets all support tickets for the current user
func getTickets(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	var tickets []models.SupportTicket
	if err := db.Where("user_id = ?", currentUser.ID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// getTicket gets a specific support ticket
func getTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get support ticket", "id": id})
}

// SupportTicketRequest represents support ticket request data
type SupportTicketRequest struct {
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"message" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
}

// createTicket creates a new support ticket
func createTicket(c *gin.Context) {
	var ticketReq SupportTicketRequest
	if err := c.ShouldBindJSON(&ticketReq); err != nil {
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

	ticket := models.SupportTicket{
		UserID:      currentUser.ID,
		Subject:     ticketReq.Subject,
		Description: ticketReq.Description,
		Priority:    ticketReq.Priority,
		Status:      "open",
	}

	if err := db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// updateTicket updates a support ticket
func updateTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update support ticket", "id": id})
}

// deleteTicket deletes a support ticket
func deleteTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete support ticket", "id": id})
}
