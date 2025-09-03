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

	db := config.GetDB()

	// Check if username already exists
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	newUser := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsAdmin:      req.IsAdmin,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create wallet for new user
	wallet := models.Wallet{
		UserID:    newUser.ID,
		Balance:   0.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&wallet).Error; err != nil {
		// If wallet creation fails, delete the user
		db.Delete(&newUser)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet for user"})
		return
	}

	// Return created user (without password hash)
	result := gin.H{
		"id":         newUser.ID,
		"username":   newUser.Username,
		"email":      newUser.Email,
		"is_admin":   newUser.IsAdmin,
		"is_active":  newUser.IsActive,
		"created_at": newUser.CreatedAt,
	}

	c.JSON(http.StatusCreated, result)
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
