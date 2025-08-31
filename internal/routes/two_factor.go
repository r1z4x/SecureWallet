package routes

import (
	"net/http"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupTwoFactorRoutes sets up 2FA routes
func SetupTwoFactorRoutes(router *gin.RouterGroup) {
	twoFactor := router.Group("/2fa")
	{
		twoFactor.POST("/enable", middleware.AuthMiddleware(), enable2FA)
		twoFactor.POST("/disable", middleware.AuthMiddleware(), disable2FA)
		twoFactor.POST("/verify", middleware.AuthMiddleware(), verify2FA)
		twoFactor.GET("/status", middleware.AuthMiddleware(), get2FAStatus)
	}
}

// Enable2FARequest represents enable 2FA request
type Enable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// enable2FA enables 2FA for a user
func enable2FA(c *gin.Context) {
	var req Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Check if 2FA is already enabled
	if userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is already enabled"})
		return
	}

	// Validate the code
	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Enable 2FA
	if err := db.Model(&userData).Update("two_factor_enabled", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA enabled successfully",
		"two_factor_enabled": true,
	})
}

// Disable2FARequest represents disable 2FA request
type Disable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// disable2FA disables 2FA for a user
func disable2FA(c *gin.Context) {
	var req Disable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Check if 2FA is enabled
	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled"})
		return
	}

	// Validate the code
	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Disable 2FA and clear secret
	if err := db.Model(&userData).Updates(map[string]interface{}{
		"two_factor_enabled": false,
		"two_factor_secret":  "",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA disabled successfully",
		"two_factor_enabled": false,
	})
}

// Verify2FARequest represents verify 2FA request
type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// verify2FA verifies a 2FA code
func verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Check if 2FA is enabled
	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled"})
		return
	}

	// Validate the code
	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA code verified successfully",
		"valid": true,
	})
}

// get2FAStatus returns the 2FA status for a user
func get2FAStatus(c *gin.Context) {
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

	// If 2FA is not enabled, generate a new secret and QR code
	if !userData.TwoFactorEnabled {
		twoFactorService := services.NewTwoFactorService()
		secret, qrURL, err := twoFactorService.GenerateSecret(userData.Username, userData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA secret"})
			return
		}

		// Save the secret temporarily
		if err := db.Model(&userData).Update("two_factor_secret", secret).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save 2FA secret"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"two_factor_enabled": false,
			"qr_code_url": qrURL,
			"secret": secret,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"two_factor_enabled": true,
	})
}
