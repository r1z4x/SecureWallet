package routes

import (
	"net/http"
	"securewallet/internal/middleware"
	"securewallet/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupSecurityRoutes sets up security monitoring routes
func SetupSecurityRoutes(router *gin.RouterGroup) {
	security := router.Group("/security")
	{
		// Security detection endpoints (admin only)
		security.GET("/idor/stats", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), getIDORStats)
		security.GET("/alerts", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), getSecurityAlerts)
		security.PUT("/alerts/:id/status", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), updateAlertStatus)
		security.POST("/users/:id/reset-attempts", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), resetUserAttempts)
		security.POST("/cleanup", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), cleanupSecurityData)
	}
}

// getIDORStats returns IDOR detection statistics
func getIDORStats(c *gin.Context) {
	securityDetector := services.NewSecurityDetector()
	stats := securityDetector.GetIDORStats()

	c.JSON(http.StatusOK, gin.H{
		"idor_detection_stats": stats,
		"message":              "IDOR detection statistics retrieved successfully",
	})
}

// getSecurityAlerts returns security alerts with filtering
func getSecurityAlerts(c *gin.Context) {
	alertType := c.Query("type") // IDOR, RATE_LIMIT, etc.
	status := c.Query("status")  // OPEN, INVESTIGATING, RESOLVED
	limitStr := c.Query("limit") // Number of alerts to return

	var limit int
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	securityDetector := services.NewSecurityDetector()
	alerts, err := securityDetector.GetSecurityAlerts(alertType, status, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve security alerts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"security_alerts": alerts,
		"count":           len(alerts),
		"filters": gin.H{
			"type":   alertType,
			"status": status,
			"limit":  limit,
		},
	})
}

// updateAlertStatus updates the status of a security alert
func updateAlertStatus(c *gin.Context) {
	alertID := c.Param("id")

	var req struct {
		Status     string `json:"status" binding:"required"`
		ResolvedBy string `json:"resolved_by"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	securityDetector := services.NewSecurityDetector()
	err := securityDetector.UpdateAlertStatus(alertID, req.Status, req.ResolvedBy)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update alert status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Alert status updated successfully",
		"alert_id":   alertID,
		"new_status": req.Status,
	})
}

// resetUserAttempts resets security attempt counters for a user
func resetUserAttempts(c *gin.Context) {
	userID := c.Param("id")

	securityDetector := services.NewSecurityDetector()
	securityDetector.ResetUserAttempts(userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "User attempt counters reset successfully",
		"user_id": userID,
	})
}

// cleanupSecurityData cleans up old security data
func cleanupSecurityData(c *gin.Context) {
	securityDetector := services.NewSecurityDetector()
	securityDetector.CleanupOldData()

	c.JSON(http.StatusOK, gin.H{
		"message":   "Security data cleanup completed successfully",
		"timestamp": "now",
	})
}
