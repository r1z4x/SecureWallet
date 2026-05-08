package routes

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Global storage for second-order attacks
var secondOrderStorage = make(map[string]map[string]string)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		// SECURE: Add rate limiting to sensitive endpoints
		auth.POST("/register", middleware.RateLimitMiddleware(), register)
		auth.POST("/login", middleware.RateLimitMiddleware(), login)
		auth.POST("/login/2fa", middleware.RateLimitMiddleware(), middleware.TwoFARateLimitMiddleware(), login2FA)
		auth.POST("/logout", logout)
		auth.GET("/me", middleware.AuthMiddleware(), getCurrentUser)
		auth.POST("/refresh", refreshToken)
		auth.POST("/password-reset", middleware.RateLimitMiddleware(), passwordReset)
		auth.POST("/password-verify", middleware.RateLimitMiddleware(), passwordVerify)
	}
}

// UserCreate represents user registration data
type UserCreate struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLogin represents user login data
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Token represents JWT token response (legacy, kept for backward compatibility)
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// TokenPairResponse represents the full token response with refresh token
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// @Summary Register a new user
// @Description Register a new user account with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body UserCreate true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} gin.H
// @Router /auth/register [post]
func register(c *gin.Context) {
	var userData UserCreate
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", userData.Username, userData.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already registered"})
		return
	}

	// SECURE: Use bcrypt for password hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	user := models.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: string(passwordHash),
		IsActive:     true,
		IsAdmin:      false,
	}

	// Enforce strong password policy at registration
	if !isStrongPassword(userData.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 12 chars, include upper, lower, number, special"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// login handles user login
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body UserLogin true "User login credentials"
// @Success 200 {object} Token
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /auth/login [post]
func login(c *gin.Context) {
	var userCredentials UserLogin
	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	audit := services.NewAuditLogger().WithGinContext(c)

	var user models.User
	if err := db.Where("username = ?", userCredentials.Username).First(&user).Error; err != nil {
		audit.Log(uuid.Nil, "LOGIN_ATTEMPT", "auth", "User not found: "+userCredentials.Username, services.AuditResultFailure)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userCredentials.Password)); err != nil {
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

		audit.Log(user.ID, "LOGIN_ATTEMPT", "auth", "Invalid password", services.AuditResultFailure)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	loginHistoryService := services.NewLoginHistoryService()
	loginHistoryService.RecordLoginAttempt(user.ID, "success", c.Request)

	audit.Log(user.ID, "LOGIN_SUCCESS", "auth", "Authentication successful", services.AuditResultSuccess)

	if !user.IsActive {
		audit.Log(user.ID, "LOGIN_DENIED", "auth", "Account disabled", services.AuditResultDenied)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is disabled"})
		return
	}

	if user.TwoFactorEnabled {
		c.JSON(http.StatusOK, gin.H{
			"requires_2fa": true,
			"message":      "2FA code required",
			"user_id":      user.ID,
		})
		return
	}

	tokenPair, err := services.CreateTokenPair(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

// Login2FARequest represents 2FA login request
type Login2FARequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	Code         string    `json:"code"`
	RecoveryCode string    `json:"recovery_code"`
}

// login2FA handles 2FA verification during login
// @Summary Verify 2FA code during login
// @Description Verify TOTP code or recovery code to complete login after initial authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param body body Login2FARequest true "2FA verification data"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /auth/login/2fa [post]
func login2FA(c *gin.Context) {
	var req Login2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Code == "" && req.RecoveryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either code or recovery_code is required"})
		return
	}

	db := config.GetDB()
	audit := services.NewAuditLogger().WithGinContext(c)

	var user models.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !user.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled for this user"})
		return
	}

	twoFactorService := services.NewTwoFactorService()

	valid := false
	if req.RecoveryCode != "" {
		recoveryValid, err := twoFactorService.ValidateRecoveryCode(user.ID.String(), req.RecoveryCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate recovery code"})
			return
		}
		if !recoveryValid {
			loginHistoryService := services.NewLoginHistoryService()
			loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

			audit.Log(user.ID, "2FA_VERIFY", "auth", "Invalid recovery code", services.AuditResultFailure)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery code"})
			return
		}
		valid = true
	} else {
		if !twoFactorService.ValidateCode(user.TwoFactorSecret, req.Code) {
			loginHistoryService := services.NewLoginHistoryService()
			loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

			audit.Log(user.ID, "2FA_VERIFY", "auth", "Invalid TOTP code", services.AuditResultFailure)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
			return
		}
		valid = true
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication code"})
		return
	}

	loginHistoryService := services.NewLoginHistoryService()
	loginHistoryService.RecordLoginAttempt(user.ID, "success", c.Request)

	audit.Log(user.ID, "2FA_VERIFY", "auth", "2FA verification successful", services.AuditResultSuccess)

	tokenPair, err := services.CreateTokenPair(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

// LogoutRequest represents a logout request with optional refresh token revocation
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// @Summary Logout user
// @Description Logout current user session. If a refresh_token is provided, it will be revoked. The current access token is blacklisted.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LogoutRequest false "Optional refresh token to revoke"
// @Success 200 {object} gin.H
// @Router /auth/logout [post]
func logout(c *gin.Context) {
	var req LogoutRequest
	_ = c.ShouldBindJSON(&req)

	user, exists := c.Get("user")
	if exists {
		if u, ok := user.(*models.User); ok {
			audit := services.NewAuditLogger().WithGinContext(c)
			audit.Log(u.ID, "LOGOUT", "auth", "User logged out", services.AuditResultSuccess)

			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenParts := strings.Split(authHeader, " ")
				if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
					services.BlacklistAccessToken(tokenParts[1], u.ID)
				}
			}
		}
	}

	if req.RefreshToken != "" {
		_ = services.RevokeToken(req.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 401 {object} gin.H
// @Router /auth/me [get]
func getCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Refresh token
// @Description Refresh access token using a valid refresh token. Implements token rotation: the old refresh token is invalidated and a new one is issued.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /auth/refresh [post]
func refreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	db := config.GetDB()
	audit := services.NewAuditLogger().WithGinContext(c)

	var session models.Session
	if err := db.Where("token = ?", req.RefreshToken).First(&session).Error; err != nil {
		audit.Log(uuid.Nil, "TOKEN_REFRESH", "auth", "Invalid refresh token", services.AuditResultFailure)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newRefreshToken, err := services.RotateRefreshToken(session.UserID, req.RefreshToken)
	if err != nil {
		audit.Log(session.UserID, "TOKEN_REFRESH", "auth", "Token rotation failed", services.AuditResultFailure)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
		return
	}

	var user models.User
	if err := db.First(&user, session.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := services.CreateAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	audit.Log(user.ID, "TOKEN_REFRESH", "auth", "Token rotated successfully", services.AuditResultSuccess)

	expireMinutesStr := os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
	expireMinutes := 30
	if expireMinutesStr != "" {
		if parsed, err := strconv.Atoi(expireMinutesStr); err == nil && parsed > 0 {
			expireMinutes = parsed
		}
	}

	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    expireMinutes * 60,
	})
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// @Summary Request password reset
// @Description Request a password reset link for the given email address
// @Tags auth
// @Accept json
// @Produce json
// @Param body body PasswordResetRequest true "Email for password reset"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /auth/password-reset [post]
// passwordReset handles password reset
func passwordReset(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		var body struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&body); err == nil && body.Email != "" {
			email = body.Email
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}
	}

	if !strings.Contains(email, "@") || len(email) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	db := config.GetDB()
	audit := services.NewAuditLogger().WithGinContext(c)

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}
		resetToken := fmt.Sprintf("%x", tokenBytes)

		if rdb := config.GetRedis(); rdb != nil {
			key := fmt.Sprintf("password_reset:%s", resetToken)
			if err := rdb.Set(c.Request.Context(), key, fmt.Sprintf("%d", user.ID), 24*time.Hour).Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store reset token"})
				return
			}
		} else {
			secondOrderStorage[email] = map[string]string{
				"reset_token": resetToken,
				"timestamp":   fmt.Sprintf("%d", time.Now().Unix()),
			}
		}

		audit.Log(user.ID, "PASSWORD_RESET_REQUEST", "auth", "Password reset requested", services.AuditResultSuccess)

		c.JSON(http.StatusOK, gin.H{"message": "Reset link sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link will be sent"})
}

// PasswordVerifyRequest represents password verification request data
type PasswordVerifyRequest struct {
	Email       string `json:"email" binding:"required"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// @Summary Complete password reset
// @Description Reset password using the token received via email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body PasswordVerifyRequest true "Password reset verification data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /auth/password-verify [post]
// passwordVerify handles password reset verification and password change
func passwordVerify(c *gin.Context) {
	var req PasswordVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	audit := services.NewAuditLogger().WithGinContext(c)

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		audit.Log(uuid.Nil, "PASSWORD_RESET_COMPLETE", "auth", "Invalid email or token", services.AuditResultFailure)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or token"})
		return
	}

	valid := false
	if rdb := config.GetRedis(); rdb != nil {
		key := fmt.Sprintf("password_reset:%s", req.Token)
		val, err := rdb.Get(c.Request.Context(), key).Result()
		if err == nil && val == fmt.Sprintf("%d", user.ID) {
			valid = true
			rdb.Del(c.Request.Context(), key)
		}
	}
	if !valid {
		if storedData, exists := secondOrderStorage[req.Email]; exists {
			timestampStr := storedData["timestamp"]
			if timestampStr != "" {
				if unixTime, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
					timestamp := time.Unix(unixTime, 0)
					if time.Since(timestamp) <= 24*time.Hour && storedData["reset_token"] == req.Token {
						valid = true
						delete(secondOrderStorage, req.Email)
					}
				}
			}
		}
	}
	if !valid {
		audit.Log(user.ID, "PASSWORD_RESET_COMPLETE", "auth", "Invalid or expired token", services.AuditResultFailure)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	if !isStrongPassword(req.NewPassword) {
		audit.Log(user.ID, "PASSWORD_RESET_COMPLETE", "auth", "Weak password rejected", services.AuditResultFailure)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 12 chars, include upper, lower, number, special"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := db.Model(&user).Update("password_hash", passwordHash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	_ = services.RevokeAllUserTokens(user.ID)
	services.BlacklistAllAccessTokensForUser(user.ID)

	delete(secondOrderStorage[req.Email], "reset_token")
	delete(secondOrderStorage[req.Email], "timestamp")

	audit.Log(user.ID, "PASSWORD_RESET_COMPLETE", "auth", "Password updated, all sessions revoked", services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password successfully updated",
	})
}

// helper: strong password checker
func isStrongPassword(pw string) bool {
	if len(pw) < 12 {
		return false
	}
	hasU, hasL, hasD, hasS := false, false, false, false
	for _, ch := range pw {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasU = true
		case ch >= 'a' && ch <= 'z':
			hasL = true
		case ch >= '0' && ch <= '9':
			hasD = true
		default:
			hasS = true
		}
	}
	return hasU && hasL && hasD && hasS
}
