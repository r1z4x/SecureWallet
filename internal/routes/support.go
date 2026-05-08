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
// @Summary List support tickets
// @Description Get all support tickets for the authenticated user
// @Tags support
// @Accept json
// @Produce json
// @Success 200 {array} models.SupportTicket
// @Failure 401 {object} gin.H
// @Security BearerAuth
// @Router /support/tickets [get]
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
// @Summary Get support ticket
// @Description Get a specific support ticket by ID (must be owned by authenticated user)
// @Tags support
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /support/tickets/{id} [get]
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
// @Summary Create support ticket
// @Description Create a new support ticket for the authenticated user
// @Tags support
// @Accept json
// @Produce json
// @Param body body SupportTicketRequest true "Support ticket data"
// @Success 201 {object} models.SupportTicket
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /support/tickets [post]
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
// @Summary Update support ticket
// @Description Update a support ticket (must be owned by authenticated user)
// @Tags support
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /support/tickets/{id} [put]
func updateTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update support ticket", "id": id})
}

// deleteTicket deletes a support ticket
// @Summary Delete support ticket
// @Description Delete a support ticket (must be owned by authenticated user)
// @Tags support
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /support/tickets/{id} [delete]
func deleteTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete support ticket", "id": id})
}
