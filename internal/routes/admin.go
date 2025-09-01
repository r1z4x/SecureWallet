package routes

import (
	"net/http"
	"securewallet/internal/middleware"

	"github.com/gin-gonic/gin"
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

// getAdminTransactions gets all transactions for admin
func getAdminTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all transactions for admin"})
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
