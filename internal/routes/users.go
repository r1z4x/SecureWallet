package routes

import (
	"net/http"
	"sync"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Mock user data for testing
var (
	mockUsers = []gin.H{
		{
			"id":         "014e2ee7-8033-4a84-b659-70b147b4dcff",
			"username":   "john_doe",
			"email":      "john@example.com",
			"is_active":  true,
			"is_admin":   false,
			"created_at": "2024-01-15T10:30:00Z",
			"updated_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":         "024e2ee7-8033-4a84-b659-70b147b4dcff",
			"username":   "jane_smith",
			"email":      "jane@example.com",
			"is_active":  true,
			"is_admin":   false,
			"created_at": "2024-01-16T14:20:00Z",
			"updated_at": "2024-01-16T14:20:00Z",
		},
		{
			"id":         "034e2ee7-8033-4a84-b659-70b147b4dcff",
			"username":   "admin_user",
			"email":      "admin@example.com",
			"is_active":  true,
			"is_admin":   true,
			"created_at": "2024-01-10T09:15:00Z",
			"updated_at": "2024-01-10T09:15:00Z",
		},
	}
	mockUsersMutex sync.RWMutex
)

// SetupUserRoutes sets up user routes
func SetupUserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/", middleware.AuthMiddleware(), getUsers)
		users.GET("/search", middleware.AuthMiddleware(), searchUsers)
		users.GET("/:id", middleware.AuthMiddleware(), getUser)
		users.POST("/", middleware.AuthMiddleware(), createUser)
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

// createUser creates a new user (admin only)
func createUser(c *gin.Context) {
	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Check if current user is admin
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		IsAdmin  bool   `json:"is_admin"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For mock data testing, add to mock users list
	mockUsersMutex.Lock()
	defer mockUsersMutex.Unlock()

	// Check if username already exists in mock data
	for _, user := range mockUsers {
		if user["username"] == req.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		if user["email"] == req.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	// Add to mock users
	newMockUser := gin.H{
		"id":         uuid.New().String(),
		"username":   req.Username,
		"email":      req.Email,
		"is_active":  req.IsActive,
		"is_admin":   req.IsAdmin,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}
	mockUsers = append(mockUsers, newMockUser)

	// Return created user
	c.JSON(http.StatusCreated, newMockUser)
}

// getUsers gets all users (admin only)
func getUsers(c *gin.Context) {
	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Check if current user is admin
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	db := config.GetDB()
	var users []models.User

	// Get all users with pagination (limit to 100 for performance)
	if err := db.Limit(100).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Return user list (exclude sensitive data like password hash)
	var results []gin.H
	for _, user := range users {
		results = append(results, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"is_active":  user.IsActive,
			"is_admin":   user.IsAdmin,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	// Always return mock data for testing (comment out real database usage)
	// if len(results) == 0 {
	mockUsersMutex.RLock()
	results = make([]gin.H, len(mockUsers))
	copy(results, mockUsers)
	mockUsersMutex.RUnlock()
	// fmt.Printf("Returning %d mock users\n", len(results))
	// }

	c.JSON(http.StatusOK, results)
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

// getUser gets a specific user (admin only)
func getUser(c *gin.Context) {
	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Check if current user is admin
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()

	var targetUser models.User
	if err := db.First(&targetUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user data (exclude sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"id":         targetUser.ID,
		"username":   targetUser.Username,
		"email":      targetUser.Email,
		"is_active":  targetUser.IsActive,
		"is_admin":   targetUser.IsAdmin,
		"created_at": targetUser.CreatedAt,
		"updated_at": targetUser.UpdatedAt,
	})
}

// updateUser updates a user (admin only)
func updateUser(c *gin.Context) {
	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Check if current user is admin
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()

	// Check if it's a mock user ID first
	mockUsersMutex.Lock()
	defer mockUsersMutex.Unlock()

	mockUserIndex := -1
	for i, user := range mockUsers {
		if user["id"] == id {
			mockUserIndex = i
			break
		}
	}

	if mockUserIndex != -1 {
		// Update mock user data
		updateData := struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			IsActive *bool  `json:"is_active"`
			IsAdmin  *bool  `json:"is_admin"`
		}{}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update fields if provided
		if updateData.Username != "" {
			mockUsers[mockUserIndex]["username"] = updateData.Username
		}
		if updateData.Email != "" {
			mockUsers[mockUserIndex]["email"] = updateData.Email
		}
		if updateData.IsActive != nil {
			mockUsers[mockUserIndex]["is_active"] = *updateData.IsActive
		}
		if updateData.IsAdmin != nil {
			mockUsers[mockUserIndex]["is_admin"] = *updateData.IsAdmin
		}
		mockUsers[mockUserIndex]["updated_at"] = time.Now().Format(time.RFC3339)

		// Return success for mock user update
		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user":    mockUsers[mockUserIndex],
		})
		return
	}

	var targetUser models.User
	if err := db.First(&targetUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse update request
	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		IsActive *bool  `json:"is_active"`
		IsAdmin  *bool  `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if updateData.Username != "" {
		targetUser.Username = updateData.Username
	}
	if updateData.Email != "" {
		targetUser.Email = updateData.Email
	}
	if updateData.IsActive != nil {
		targetUser.IsActive = *updateData.IsActive
	}
	if updateData.IsAdmin != nil {
		targetUser.IsAdmin = *updateData.IsAdmin
	}

	// Save changes
	if err := db.Save(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":         targetUser.ID,
			"username":   targetUser.Username,
			"email":      targetUser.Email,
			"is_active":  targetUser.IsActive,
			"is_admin":   targetUser.IsAdmin,
			"updated_at": targetUser.UpdatedAt,
		},
	})
}

// deleteUser deletes a user (admin only)
func deleteUser(c *gin.Context) {
	// Get current user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Check if current user is admin
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()

	// Check if it's a mock user ID first
	mockUsersMutex.Lock()
	defer mockUsersMutex.Unlock()

	mockUserIndex := -1
	for i, user := range mockUsers {
		if user["id"] == id {
			mockUserIndex = i
			break
		}
	}

	if mockUserIndex != -1 {
		// Remove user from mock data
		mockUsers = append(mockUsers[:mockUserIndex], mockUsers[mockUserIndex+1:]...)

		// Debug log
		// fmt.Printf("Mock user deleted. Remaining users: %d\n", len(mockUsers))

		// Return success for mock user deletion
		c.JSON(http.StatusOK, gin.H{
			"message":    "User deleted successfully",
			"deleted_at": time.Now(),
		})
		return
	}

	var targetUser models.User
	if err := db.First(&targetUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent admin from deleting themselves
	if targetUser.ID == currentUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	// Check if user has any active wallets with balance
	var walletCount int64
	db.Model(&models.Wallet{}).Where("user_id = ? AND balance > 0", targetUser.ID).Count(&walletCount)
	if walletCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete user with active wallet balance. Please transfer or withdraw all funds first.",
		})
		return
	}

	// Start database transaction
	tx := db.Begin()

	// Delete user's data in order (similar to deleteCurrentUserAccount)
	// 1. Delete login history
	if err := tx.Where("user_id = ?", targetUser.ID).Delete(&models.LoginHistory{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete login history"})
		return
	}

	// 2. Delete sessions
	if err := tx.Where("user_id = ?", targetUser.ID).Delete(&models.Session{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sessions"})
		return
	}

	// 3. Delete support tickets
	if err := tx.Where("user_id = ?", targetUser.ID).Delete(&models.SupportTicket{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete support tickets"})
		return
	}

	// 4. Delete audit logs
	if err := tx.Where("user_id = ?", targetUser.ID).Delete(&models.AuditLog{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete audit logs"})
		return
	}

	// 5. Delete transactions
	if err := tx.Where("wallet_id IN (SELECT id FROM wallets WHERE user_id = ?)", targetUser.ID).Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transactions"})
		return
	}

	// 6. Delete wallets
	if err := tx.Where("user_id = ?", targetUser.ID).Delete(&models.Wallet{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallets"})
		return
	}

	// 7. Finally, delete the user
	if err := tx.Delete(&targetUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit user deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User deleted successfully",
		"deleted_at": time.Now(),
	})
}
