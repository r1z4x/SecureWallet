package middleware

import (
	"net/http"

	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// RequirePermission returns a middleware that checks if the user has the specified permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
			c.Abort()
			return
		}

		// Legacy admin check - admins have all permissions
		if currentUser.IsAdmin {
			c.Next()
			return
		}

		// Check RBAC permissions
		rbacService := services.NewRBACService()
		hasPerm, err := rbacService.HasPermission(currentUser.ID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !hasPerm {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "Insufficient permissions",
				"required": permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission returns a middleware that checks if the user has ANY of the specified permissions
func RequireAnyPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
			c.Abort()
			return
		}

		// Legacy admin check
		if currentUser.IsAdmin {
			c.Next()
			return
		}

		// Check if user has any of the required permissions
		rbacService := services.NewRBACService()
		for _, perm := range permissions {
			hasPerm, err := rbacService.HasPermission(currentUser.ID, perm)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}
			if hasPerm {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "Insufficient permissions",
			"required": permissions,
		})
		c.Abort()
	}
}

// RequireRole returns a middleware that checks if the user has the specified role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
			c.Abort()
			return
		}

		// Legacy admin check
		if currentUser.IsAdmin && roleName == models.RoleAdmin {
			c.Next()
			return
		}

		// Check RBAC roles
		rbacService := services.NewRBACService()
		roles, err := rbacService.GetUserRoles(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check roles"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if role.Name == roleName {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":    "Insufficient role",
			"required": roleName,
		})
		c.Abort()
	}
}
