package routes

import (
	"net/http"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// CronRoutes sets up cron-related routes
func CronRoutes(router *gin.Engine) {
	cron := router.Group("/api/cron")
	{
		// Get cron status
		cron.GET("/status", func(c *gin.Context) {
			cronService := services.NewCronService()
			status := cronService.GetCronStatus()
			c.JSON(http.StatusOK, status)
		})

		// Execute a specific cron job
		cron.POST("/execute/:job", func(c *gin.Context) {
			jobName := c.Param("job")
			
			cronService := services.NewCronService()
			if err := cronService.ExecuteCronJob(jobName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to execute cron job",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Cron job executed successfully",
				"job": jobName,
			})
		})

		// Setup cron jobs
		cron.POST("/setup", func(c *gin.Context) {
			cronService := services.NewCronService()
			cronService.SetupCronJobs()

			c.JSON(http.StatusOK, gin.H{
				"message": "Cron jobs setup completed",
			})
		})

		// Remove cron jobs
		cron.DELETE("/remove", func(c *gin.Context) {
			cronService := services.NewCronService()
			if err := cronService.RemoveCronJobs(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to remove cron jobs",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Cron jobs removed successfully",
			})
		})

		// Get comment stats
		cron.GET("/comments/stats", func(c *gin.Context) {
			commentService := services.NewCommentService()
			stats := commentService.GetCommentStats()
			c.JSON(http.StatusOK, stats)
		})

		// Get backup stats
		cron.GET("/backup/stats", func(c *gin.Context) {
			backupService := services.NewBackupService()
			backups, err := backupService.ListBackups()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to list backups",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"backup_count": len(backups),
				"backups": backups,
			})
		})
	}
}
