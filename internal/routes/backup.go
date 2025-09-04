package routes

import (
	"net/http"
	"securewallet/internal/middleware"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupBackupRoutes sets up backup routes
func SetupBackupRoutes(router *gin.RouterGroup) {
	backup := router.Group("/backup")
	{
		// VULNERABLE: These endpoints can leak sensitive information through backups
		// Manual backup creation is deprecated - use automatic scheduled backups
		backup.GET("/", middleware.AuthMiddleware(), listBackups)
		backup.GET("/:filename", middleware.AuthMiddleware(), getBackupInfo)
		backup.GET("/stats", middleware.AuthMiddleware(), getBackupStats)
		backup.GET("/config", middleware.AuthMiddleware(), getBackupConfig)
	}
}

// listBackups lists all available backups
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
func getBackupConfig(c *gin.Context) {
	config := services.DefaultBackupConfig
	
	c.JSON(http.StatusOK, gin.H{
		"config": gin.H{
			"max_backups":      config.MaxBackups,
			"backup_interval":  config.BackupInterval.String(),
			"enabled":          config.Enabled,
			"message":          "Manual backup creation is deprecated. Use automatic scheduled backups.",
		},
	})
}

// getBackupInfo returns information about a specific backup
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
func getBackupStats(c *gin.Context) {
	backupService := services.NewBackupService()
	
	stats := backupService.GetBackupStats()
	
	c.JSON(http.StatusOK, stats)
}
