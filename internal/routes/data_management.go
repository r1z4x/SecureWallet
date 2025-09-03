package routes

import (
	"net/http"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupDataManagementRoutes sets up data management routes
func SetupDataManagementRoutes(router *gin.RouterGroup) {
	data := router.Group("/data")
	{
		data.POST("/init-db", initDatabase)
		data.DELETE("/clear-sample", clearSampleData)
		data.GET("/stats", getDataStats)
		data.POST("/reset-database", resetDatabase)
		data.POST("/force-recreate", forceDatabaseRecreation)
	}
}

// @Summary Initialize database tables
// @Description Initialize database tables using GORM auto-migration
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /data/init-db [post]
func initDatabase(c *gin.Context) {
	dataManager := services.NewSampleDataManager()

	if err := dataManager.InitializeDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database tables initialized successfully",
	})
}

// @Summary Clear sample data
// @Description Clear all sample data
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /data/clear-sample [delete]
func clearSampleData(c *gin.Context) {
	dataManager := services.NewSampleDataManager()

	if err := dataManager.ClearSampleData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clear sample data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sample data cleared successfully",
	})
}

// @Summary Get data statistics
// @Description Get statistics about current data
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /data/stats [get]
func getDataStats(c *gin.Context) {
	dataManager := services.NewSampleDataManager()

	c.JSON(http.StatusOK, gin.H{
		"stats": dataManager.GetSampleDataStats(),
	})
}

// @Summary Reset database completely
// @Description Clear all data, initialize database, and create sample data
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /data/reset-database [post]
func resetDatabase(c *gin.Context) {
	dataManager := services.NewSampleDataManager()

	if err := dataManager.ResetDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reset database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database reset successfully",
		"stats":   dataManager.GetSampleDataStats(),
	})
}

// @Summary Force database recreation
// @Description Completely drop and recreate the database (use when schema incompatibility issues occur)
// @Tags data
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /data/force-recreate [post]
func forceDatabaseRecreation(c *gin.Context) {
	dataManager := services.NewSampleDataManager()

	if err := dataManager.CompleteDatabaseRecreation(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to force database recreation",
			"details": err.Error(),
		})
		return
	}

	// After recreation, initialize the database and create sample data
	if err := dataManager.InitializeDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize database after recreation",
			"details": err.Error(),
		})
		return
	}

	// Create sample data
	if err := dataManager.CreateSampleUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sample users after recreation",
			"details": err.Error(),
		})
		return
	}

	if err := dataManager.CreateSampleWallets(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sample wallets after recreation",
			"details": err.Error(),
		})
		return
	}

	if err := dataManager.CreateSampleTransactions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sample transactions after recreation",
			"details": err.Error(),
		})
		return
	}

	if err := dataManager.CreateSampleLoginHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sample login history after recreation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database force recreation completed successfully",
		"stats":   dataManager.GetSampleDataStats(),
	})
}
