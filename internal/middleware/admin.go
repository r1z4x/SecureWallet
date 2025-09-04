package middleware

import (
	"net/http"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
)

// AdminOnlyMiddleware ensures only admin users can access the endpoint
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			c.Abort()
			return
		}

		if !currentUser.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
				"message": "This endpoint is restricted to administrators only",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
