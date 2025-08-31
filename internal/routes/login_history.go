package routes

import (
	"fmt"
	"net/http"

	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupLoginHistoryRoutes sets up login history routes
func SetupLoginHistoryRoutes(router *gin.RouterGroup) {
	loginHistory := router.Group("/login-history")
	{
		loginHistory.GET("/", middleware.AuthMiddleware(), getLoginHistory)
		loginHistory.GET("/recent", middleware.AuthMiddleware(), getRecentLoginHistory)
	}
}

// getLoginHistory gets login history for the current user
func getLoginHistory(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	limit := 50 // Default limit

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseIntParam(limitStr, 50); err == nil {
			limit = parsedLimit
		}
	}

	loginHistoryService := services.NewLoginHistoryService()
	history, err := loginHistoryService.GetLoginHistory(currentUser.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch login history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login_history": history,
		"total":         len(history),
	})
}

// getRecentLoginHistory gets recent login history for the current user
func getRecentLoginHistory(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	limit := 10 // Default limit for recent

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := parseIntParam(limitStr, 10); err == nil {
			limit = parsedLimit
		}
	}

	loginHistoryService := services.NewLoginHistoryService()
	history, err := loginHistoryService.GetLoginHistory(currentUser.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent login history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recent_login_history": history,
		"total":                len(history),
	})
}

// Helper function to parse integer parameters
func parseIntParam(param string, defaultValue int) (int, error) {
	// This is a simplified version - in production you'd want proper validation
	if param == "" {
		return defaultValue, nil
	}

	// Parse the string to int
	var result int
	_, err := fmt.Sscanf(param, "%d", &result)
	if err != nil {
		return defaultValue, err
	}

	// Validate the result
	if result <= 0 {
		return defaultValue, nil
	}

	return result, nil
}
