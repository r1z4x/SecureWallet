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
	"github.com/golang-jwt/jwt/v5"
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
		auth.POST("/login/2fa", middleware.RateLimitMiddleware(), login2FA)
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

// Token represents JWT token response
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
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

	// VULNERABILITY: Advanced authentication with multiple bypass techniques
	var user models.User
	if err := db.Where("username = ?", userCredentials.Username).First(&user).Error; err != nil {
		// Record failed login attempt
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(0, "failed", c.Request)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	// SECURE: Use bcrypt for password verification
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userCredentials.Password)); err != nil {
		// Record failed login attempt
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	// Authentication successful
	// Record successful login attempt
	loginHistoryService := services.NewLoginHistoryService()
	loginHistoryService.RecordLoginAttempt(user.ID, "success", c.Request)

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is disabled"})
		return
	}

	// Check if 2FA is enabled
	if user.TwoFactorEnabled {
		c.JSON(http.StatusOK, gin.H{
			"requires_2fa": true,
			"message":      "2FA code required",
			"user_id":      user.ID,
		})
		return
	}

	// SECURE: Use environment variable for JWT secret and reasonable expiration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}
	expireMinutes := 60 // 1 hour - reasonable expiration time

	// Create access token with vulnerability-specific settings
	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, Token{
		AccessToken: accessToken,
		TokenType:   "bearer",
	})
}

// Login2FARequest represents 2FA login request
type Login2FARequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

// login2FA handles 2FA verification during login
func login2FA(c *gin.Context) {
	var req Login2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	// Get user
	var user models.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if 2FA is enabled
	if !user.TwoFactorEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled for this user"})
		return
	}

	// Validate 2FA code
	twoFactorService := services.NewTwoFactorService()
	if !twoFactorService.ValidateCode(user.TwoFactorSecret, req.Code) {
		// Record failed 2FA attempt
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Record successful 2FA login
	loginHistoryService := services.NewLoginHistoryService()
	loginHistoryService.RecordLoginAttempt(user.ID, "success", c.Request)

	// SECURE: Create JWT token with environment variable and reasonable expiration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}
	expireMinutes := 60 // 1 hour - reasonable expiration time

	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, Token{
		AccessToken: accessToken,
		TokenType:   "bearer",
	})
}

// @Summary Logout user
// @Description Logout current user session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /auth/logout [post]
func logout(c *gin.Context) {
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

// @Summary Refresh token
// @Description Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /auth/refresh [post]
func refreshToken(c *gin.Context) {
	// Implementation for token refresh
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// @Summary Reset password
// @Description Reset password for a user
// @Tags auth
// @Accept json
// @Produce json
// @Param email body PasswordResetRequest true "Email for password reset"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /auth/password-reset [post]
// passwordReset handles password reset
func passwordReset(c *gin.Context) {
	email := c.PostForm("email")
	// In passwordReset: prefer JSON body if form empty
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

	// SECURE: Validate email format
	if !strings.Contains(email, "@") || len(email) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	db := config.GetDB()

	// SECURE: Check if user exists and generate secure reset token
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		// Generate cryptographically secure random token
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}
		resetToken := fmt.Sprintf("%x", tokenBytes)

		// Store token in Redis with 24h expiration; fallback to memory
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

// passwordVerify handles password reset verification and password change
func passwordVerify(c *gin.Context) {
	var req PasswordVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	// SECURE: Validate token securely
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or token"})
		return
	}

	// Validate token with Redis first, then fallback
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Enforce strong password policy in reset verify
	if !isStrongPassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 12 chars, include upper, lower, number, special"})
		return
	}

	// SECURE: Update user password with bcrypt hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update user password in database
	if err := db.Model(&user).Update("password_hash", passwordHash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Clear reset token from storage
	delete(secondOrderStorage[req.Email], "reset_token")
	delete(secondOrderStorage[req.Email], "timestamp")

	// SECURE: Return minimal information
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
