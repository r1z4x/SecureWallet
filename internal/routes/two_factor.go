package routes

import (
	"net/http"

	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupTwoFactorRoutes sets up 2FA routes
func SetupTwoFactorRoutes(router *gin.RouterGroup) {
	twoFactor := router.Group("/2fa")
	twoFactor.Use(middleware.AuthMiddleware())
	{
		twoFactor.POST("/enable", middleware.RequirePermission(models.PermTwoFactorWrite), enable2FA)
		twoFactor.POST("/disable", middleware.RequirePermission(models.PermTwoFactorWrite), disable2FA)
		twoFactor.POST("/verify", middleware.RequirePermission(models.PermTwoFactorWrite), verify2FA)
		twoFactor.GET("/status", middleware.RequirePermission(models.PermTwoFactorRead), get2FAStatus)
		twoFactor.POST("/recovery-codes/generate", middleware.RequirePermission(models.PermTwoFactorWrite), generateRecoveryCodes)
		twoFactor.GET("/recovery-codes", middleware.RequirePermission(models.PermTwoFactorRead), getRecoveryCodes)
	}
}

// Enable2FARequest represents enable 2FA request
type Enable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// enable2FA enables 2FA for a user
// @Summary Enable 2FA
// @Description Enable TOTP two-factor authentication for the authenticated user
// @Tags 2fa
// @Accept json
// @Produce json
// @Param body body Enable2FARequest true "2FA enable data with TOTP code"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /2fa/enable [post]
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
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is already enabled"})
		return
	}

	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	if err := db.Model(&userData).Update("two_factor_enabled", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	recoveryCodes, err := twoFactorService.GenerateRecoveryCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recovery codes"})
		return
	}

	if err := twoFactorService.StoreRecoveryCodes(userData.ID.String(), recoveryCodes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store recovery codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "2FA enabled successfully",
		"two_factor_enabled": true,
		"recovery_codes":     recoveryCodes,
	})
}

// Disable2FARequest represents disable 2FA request
type Disable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// disable2FA disables 2FA for a user
// @Summary Disable 2FA
// @Description Disable TOTP two-factor authentication for the authenticated user
// @Tags 2fa
// @Accept json
// @Produce json
// @Param body body Disable2FARequest true "2FA disable data with TOTP code"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /2fa/disable [post]
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
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled"})
		return
	}

	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	if err := db.Model(&userData).Updates(map[string]interface{}{
		"two_factor_enabled": false,
		"two_factor_secret":  "",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

	_ = twoFactorService.ClearRecoveryCodes(userData.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"message":            "2FA disabled successfully",
		"two_factor_enabled": false,
	})
}

// Verify2FARequest represents verify 2FA request
type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// verify2FA verifies a 2FA code
// @Summary Verify 2FA code
// @Description Verify a TOTP code for the authenticated user
// @Tags 2fa
// @Accept json
// @Produce json
// @Param body body Verify2FARequest true "TOTP code"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /2fa/verify [post]
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
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled"})
		return
	}

	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(userData.TwoFactorSecret, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA code verified successfully",
		"valid":   true,
	})
}

// get2FAStatus returns the 2FA status for a user
// @Summary Get 2FA status
// @Description Get the 2FA enrollment status for the authenticated user. Returns QR code URL if not yet enabled.
// @Tags 2fa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Security BearerAuth
// @Router /2fa/status [get]
func get2FAStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !userData.TwoFactorEnabled {
		twoFactorService := services.NewTwoFactorService()
		secret, qrURL, err := twoFactorService.GenerateSecret(userData.Username, userData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA secret"})
			return
		}

		if err := db.Model(&userData).Update("two_factor_secret", secret).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save 2FA secret"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"two_factor_enabled": false,
			"qr_code_url":        qrURL,
			"secret":             secret,
		})
		return
	}

	twoFactorService := services.NewTwoFactorService()
	remainingCodes, _ := twoFactorService.GetRemainingRecoveryCodes(userData.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"two_factor_enabled":      true,
		"recovery_codes_remaining": remainingCodes,
	})
}

// generateRecoveryCodes generates new recovery codes for a user
// @Summary Generate recovery codes
// @Description Generate new TOTP recovery codes for the authenticated user (requires 2FA enabled)
// @Tags 2fa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /2fa/recovery-codes/generate [post]
func generateRecoveryCodes(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA must be enabled to generate recovery codes"})
		return
	}

	twoFactorService := services.NewTwoFactorService()

	_ = twoFactorService.ClearRecoveryCodes(userData.ID.String())

	recoveryCodes, err := twoFactorService.GenerateRecoveryCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recovery codes"})
		return
	}

	if err := twoFactorService.StoreRecoveryCodes(userData.ID.String(), recoveryCodes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store recovery codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Recovery codes regenerated successfully",
		"recovery_codes": recoveryCodes,
	})
}

// getRecoveryCodes returns the count of remaining recovery codes
// @Summary Get recovery codes count
// @Description Get the count of remaining unused recovery codes
// @Tags 2fa
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Security BearerAuth
// @Router /2fa/recovery-codes [get]
func getRecoveryCodes(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	db := middleware.DefaultDB(c)

	var userData models.User
	if err := db.First(&userData, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !userData.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA must be enabled"})
		return
	}

	twoFactorService := services.NewTwoFactorService()
	remaining, err := twoFactorService.GetRemainingRecoveryCodes(userData.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recovery codes count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recovery_codes_remaining": remaining,
	})
}
