package routes

import (
	"net/http"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up user routes
func SetupUserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/", middleware.AuthMiddleware(), getUsers)
		users.GET("/search", middleware.AuthMiddleware(), searchUsers)
		users.GET("/:id", middleware.AuthMiddleware(), getUser)
		users.PUT("/:id", middleware.AuthMiddleware(), updateUser)
		users.DELETE("/:id", middleware.AuthMiddleware(), deleteUser)
		users.DELETE("/account", middleware.AuthMiddleware(), deleteCurrentUserAccount)
	}
}

// searchUsers searches for users by username or email
func searchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// SECURE: Input validation and sanitization
	if len(query) < 2 || len(query) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be between 2 and 50 characters"})
		return
	}

	// SECURE: Check for potentially dangerous characters
	for _, char := range query {
		if char < 32 || char > 126 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query contains invalid characters"})
			return
		}
	}

	db := config.GetDB()
	var users []models.User

	// SECURE: Use parameterized queries to prevent SQL injection
	if err := db.Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").
		Where("is_active = ?", true).
		Limit(10).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// Return only necessary user information (exclude sensitive data)
	var results []gin.H
	for _, user := range users {
		results = append(results, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}

	c.JSON(http.StatusOK, results)
}

// getUsers gets all users
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all users"})
}

// DeleteAccountRequest represents delete account request
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
}

// deleteCurrentUserAccount deletes the current user's account
func deleteCurrentUserAccount(c *gin.Context) {
	var req DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate confirmation
	if req.Confirm != "DELETE MY ACCOUNT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please type 'DELETE MY ACCOUNT' to confirm"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := config.GetDB()

	// Get fresh user data
	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// SECURE: Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Check if user has any active wallets with balance
	var walletCount int64
	db.Model(&models.Wallet{}).Where("user_id = ? AND balance > 0", userData.ID).Count(&walletCount)
	if walletCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete account with active wallet balance. Please transfer or withdraw all funds first.",
		})
		return
	}

	// Start database transaction
	tx := db.Begin()

	// Delete user's data in order
	// 1. Delete login history
	if err := tx.Where("user_id = ?", userData.ID).Delete(&models.LoginHistory{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete login history"})
		return
	}

	// 2. Delete sessions
	if err := tx.Where("user_id = ?", userData.ID).Delete(&models.Session{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sessions"})
		return
	}

	// 3. Delete support tickets
	if err := tx.Where("user_id = ?", userData.ID).Delete(&models.SupportTicket{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete support tickets"})
		return
	}

	// 4. Delete audit logs
	if err := tx.Where("user_id = ?", userData.ID).Delete(&models.AuditLog{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete audit logs"})
		return
	}

	// 5. Delete transactions
	if err := tx.Where("wallet_id IN (SELECT id FROM wallets WHERE user_id = ?)", userData.ID).Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transactions"})
		return
	}

	// 6. Delete wallets
	if err := tx.Where("user_id = ?", userData.ID).Delete(&models.Wallet{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallets"})
		return
	}

	// 7. Finally, delete the user
	if err := tx.Delete(&userData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit account deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Account deleted successfully",
		"deleted_at": time.Now(),
	})
}

// getUser gets a specific user
func getUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get user", "id": id})
}

// updateUser updates a user
func updateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update user", "id": id})
}

// deleteUser deletes a user
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete user", "id": id})
}
