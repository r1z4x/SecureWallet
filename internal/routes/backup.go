package routes

import (
	"net/http"

	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupBackupRoutes sets up backup routes
func SetupBackupRoutes(router *gin.RouterGroup) {
	backup := router.Group("/backup")
	backup.Use(middleware.AuthMiddleware())
	backup.Use(middleware.RequirePermission(models.PermBackupRead))
	{
		backup.GET("/", listBackups)
		backup.GET("/:filename", getBackupInfo)
		backup.GET("/stats", getBackupStats)
		backup.GET("/config", getBackupConfig)
	}
}

// listBackups lists all available backups
// @Summary List backups
// @Description List all available database backups
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /backup [get]
func listBackups(c *gin.Context) {
	backupService := services.NewBackupService()

	backups, err := backupService.ListBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
		"count":   len(backups),
	})
}

// getBackupConfig returns backup configuration
// @Summary Get backup config
// @Description Get backup configuration settings
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /backup/config [get]
func getBackupConfig(c *gin.Context) {
	config := services.DefaultBackupConfig

	c.JSON(http.StatusOK, gin.H{
		"config": gin.H{
			"max_backups":     config.MaxBackups,
			"backup_interval": config.BackupInterval.String(),
			"enabled":         config.Enabled,
			"message":         "Manual backup creation is deprecated. Use automatic scheduled backups.",
		},
	})
}

// getBackupInfo returns information about a specific backup
// @Summary Get backup info
// @Description Get information about a specific backup file
// @Tags backup
// @Accept json
// @Produce json
// @Param filename path string true "Backup filename"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Security BearerAuth
// @Router /backup/{filename} [get]
func getBackupInfo(c *gin.Context) {
	filename := c.Param("filename")

	backupService := services.NewBackupService()

	backupInfo, err := backupService.GetBackupInfo(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// VULNERABLE: This endpoint can leak sensitive wallet information
	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"info":     backupInfo,
	})
}

// getBackupStats returns backup statistics
// @Summary Get backup stats
// @Description Get backup statistics and metadata
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /backup/stats [get]
func getBackupStats(c *gin.Context) {
	backupService := services.NewBackupService()

	stats := backupService.GetBackupStats()

	c.JSON(http.StatusOK, stats)
}
