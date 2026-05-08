package routes

import (
	"log"
	"net/http"
	"time"

	"securewallet/internal/middleware"
	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupOAuthRoutes sets up OAuth2 authentication routes
func SetupOAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/oauth/providers", getOAuthProviders)
		auth.GET("/oauth/:provider/authorize", startOAuthFlow)
		auth.GET("/oauth/:provider/callback", handleOAuthCallback)
		auth.POST("/oauth/link", middleware.AuthMiddleware(), middleware.RequirePermission(models.PermSessionWrite), linkOAuthAccount)
		auth.POST("/oauth/link-existing", middleware.AuthMiddleware(), middleware.RequirePermission(models.PermSessionWrite), linkExistingOAuthAccount)
		auth.POST("/oauth/unlink", middleware.AuthMiddleware(), middleware.RequirePermission(models.PermSessionWrite), unlinkOAuthAccount)
		auth.GET("/login-methods/:user_id", getLoginMethods)
		auth.GET("/login-methods", middleware.AuthMiddleware(), middleware.RequirePermission(models.PermSessionRead), getMyLoginMethods)
		auth.POST("/oauth/set-password", middleware.AuthMiddleware(), middleware.RequirePermission(models.PermSessionWrite), setOAuthUserPassword)
	}
}

// OAuthProviderResponse represents an OAuth provider for API responses
type OAuthProviderResponse struct {
	Name string `json:"name"`
}

// getOAuthProviders returns all active OAuth providers
func getOAuthProviders(c *gin.Context) {
	oauthService := services.NewOAuthService()
	providers, err := oauthService.GetOAuthProviders()
	if err != nil {
		log.Printf("Failed to get OAuth providers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load OAuth providers"})
		return
	}

	result := make([]OAuthProviderResponse, 0, len(providers))
	for _, p := range providers {
		result = append(result, OAuthProviderResponse{Name: p.Name})
	}

	c.JSON(http.StatusOK, gin.H{"providers": result})
}

// startOAuthFlow initiates the OAuth flow by redirecting to the provider
func startOAuthFlow(c *gin.Context) {
	providerName := c.Param("provider")
	if providerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider name is required"})
		return
	}

	oauthService := services.NewOAuthService()
	state, err := oauthService.GenerateState(providerName)
	if err != nil {
		log.Printf("Failed to generate OAuth state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in cookie for CSRF validation
	c.SetCookie("oauth_state", state, 300, "/api/auth/oauth/"+providerName, "", false, true)

	authURL, err := oauthService.GetAuthURL(providerName, state)
	if err != nil {
		log.Printf("Failed to get OAuth auth URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// handleOAuthCallback handles the OAuth callback from the provider
func handleOAuthCallback(c *gin.Context) {
	providerName := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		errorDesc := c.Query("error_description")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            errorParam,
			"error_description": errorDesc,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})
		return
	}

	// Verify state from cookie
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/api/auth/oauth/"+providerName, "", false, true)

	oauthService := services.NewOAuthService()
	userInfo, err := oauthService.HandleCallback(providerName, code, state)
	if err != nil {
		log.Printf("OAuth callback failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to complete OAuth flow"})
		return
	}

	// Find or create user
	user, err := oauthService.FindOrCreateUserByOAuth(providerName, userInfo)
	if err != nil {
		if err == services.ErrEmailCollision {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "email_collision",
				"message": "An account with this email already exists. Login with your password first, then link this OAuth provider from your profile settings.",
			})
			return
		}
		log.Printf("Failed to find/create user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Complete login with 2FA check and audit logging
	tokenPair, requires2FA, err := oauthService.CompleteOAuthLogin(user, c)
	if err != nil {
		log.Printf("OAuth login completion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if requires2FA {
		c.JSON(http.StatusOK, gin.H{
			"requires_2fa": true,
			"message":      "2FA code required",
			"user_id":      user.ID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OAuth login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
		},
		"tokens": tokenPair,
	})
}

// LinkOAuthAccountRequest represents the request to link an OAuth account
type LinkOAuthAccountRequest struct {
	ProviderName string `json:"provider_name" binding:"required"`
	Code         string `json:"code" binding:"required"`
	State        string `json:"state" binding:"required"`
}

// linkOAuthAccount links an OAuth account to the current user
func linkOAuthAccount(c *gin.Context) {
	var req LinkOAuthAccountRequest
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

	oauthService := services.NewOAuthService()
	userInfo, err := oauthService.HandleCallback(req.ProviderName, req.Code, req.State)
	if err != nil {
		log.Printf("OAuth callback failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to complete OAuth flow"})
		return
	}

	db := middleware.DefaultDB(c)
	var provider models.OAuthProvider
	if err := db.Where("name = ?", req.ProviderName).First(&provider).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not found"})
		return
	}

	tokenResp := &services.TokenResponse{}
	if err := oauthService.LinkOAuthAccount(currentUser.ID, req.ProviderName, userInfo, tokenResp); err != nil {
		log.Printf("Failed to link OAuth account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OAuth account linked successfully"})
}

// linkExistingOAuthAccount links an OAuth account to the current user after email collision
func linkExistingOAuthAccount(c *gin.Context) {
	var req LinkOAuthAccountRequest
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

	oauthService := services.NewOAuthService()
	userInfo, err := oauthService.HandleCallback(req.ProviderName, req.Code, req.State)
	if err != nil {
		log.Printf("OAuth callback failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to complete OAuth flow"})
		return
	}

	if err := oauthService.LinkExistingUserToOAuth(currentUser.ID, req.ProviderName, userInfo); err != nil {
		log.Printf("Failed to link OAuth account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audit := services.NewAuditLogger().WithGinContext(c)
	audit.Log(currentUser.ID, "OAUTH_ACCOUNT_LINKED", "auth", "OAuth account linked: "+req.ProviderName, services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{"message": "OAuth account linked successfully"})
}

// UnlinkOAuthAccountRequest represents the request to unlink an OAuth account
type UnlinkOAuthAccountRequest struct {
	ProviderName string `json:"provider_name" binding:"required"`
}

// unlinkOAuthAccount unlinks an OAuth account from the current user
func unlinkOAuthAccount(c *gin.Context) {
	var req UnlinkOAuthAccountRequest
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

	// Check if user has a password set (must have at least one login method)
	if currentUser.PasswordHash == "" {
		// Check if user has other OAuth accounts linked
		var oauthCount int64
		db.Model(&models.OAuthAccount{}).Where("user_id = ?", currentUser.ID).Count(&oauthCount)
		if oauthCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unlink last login method"})
			return
		}
	}

	result := db.Where("user_id = ? AND provider_name = ?", currentUser.ID, req.ProviderName).Delete(&models.OAuthAccount{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink account"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "OAuth account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OAuth account unlinked successfully"})
}

// getLoginMethods returns the available login methods for a user (public, used during login flow)
func getLoginMethods(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	oauthService := services.NewOAuthService()
	methods, err := oauthService.GetUserLoginMethods(userID)
	if err != nil {
		log.Printf("Failed to get login methods: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// getMyLoginMethods returns the login methods for the authenticated user
func getMyLoginMethods(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	oauthService := services.NewOAuthService()
	methods, err := oauthService.GetUserLoginMethods(currentUser.ID)
	if err != nil {
		log.Printf("Failed to get login methods: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load login methods"})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// SetPasswordRequest represents the request to set a password for an OAuth-only user
type SetPasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password" binding:"required,min=12"`
}

// setOAuthUserPassword allows an OAuth-only user to set a native password
func setOAuthUserPassword(c *gin.Context) {
	var req SetPasswordRequest
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

	oauthService := services.NewOAuthService()
	if err := oauthService.SetPasswordForOAuthUser(currentUser.ID, req.CurrentPassword, req.NewPassword); err != nil {
		log.Printf("Failed to set password: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audit := services.NewAuditLogger().WithGinContext(c)
	audit.Log(currentUser.ID, "PASSWORD_SET", "auth", "Password set for OAuth user", services.AuditResultSuccess)

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully"})
}

// SetupSessionRoutes sets up session management routes
func SetupSessionRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/sessions", middleware.RequirePermission(models.PermSessionRead), listSessions)
		auth.DELETE("/sessions/:id", middleware.RequirePermission(models.PermSessionWrite), revokeSession)
		auth.POST("/sessions/revoke-all", middleware.RequirePermission(models.PermSessionWrite), revokeAllSessions)
		auth.POST("/sessions/revoke-others", middleware.RequirePermission(models.PermSessionWrite), revokeOtherSessions)
	}
}

// SessionResponse represents a session for API responses
type SessionResponse struct {
	ID           string    `json:"id"`
	DeviceName   string    `json:"device_name"`
	IPAddress    string    `json:"ip_address"`
	LastAccessed time.Time `json:"last_accessed"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsCurrent    bool      `json:"is_current"`
}

// listSessions returns all active sessions for the current user
func listSessions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)

	// Get current session ID from context (set by auth middleware)
	currentSessionID := uuid.Nil
	if sid, ok := c.Get("session_id"); ok {
		if parsed, err := uuid.Parse(sid.(string)); err == nil {
			currentSessionID = parsed
		}
	}

	sessionService := services.NewSessionService()
	sessions, err := sessionService.GetUserSessions(currentUser.ID, currentSessionID)
	if err != nil {
		log.Printf("Failed to get sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load sessions"})
		return
	}

	result := make([]SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, SessionResponse{
			ID:           s.ID,
			DeviceName:   s.DeviceName,
			IPAddress:    s.IPAddress,
			LastAccessed: s.LastAccessed,
			CreatedAt:    s.CreatedAt,
			ExpiresAt:    s.ExpiresAt,
			IsCurrent:    s.IsCurrent,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": result})
}

// revokeSession revokes a specific session
func revokeSession(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	sessionService := services.NewSessionService()
	if err := sessionService.RevokeSession(currentUser.ID, sessionID); err != nil {
		log.Printf("Failed to revoke session: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// revokeAllSessions revokes all sessions including current
func revokeAllSessions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	sessionService := services.NewSessionService()
	if err := sessionService.RevokeAllUserSessions(currentUser.ID); err != nil {
		log.Printf("Failed to revoke all sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke sessions"})
		return
	}

	services.BlacklistAllAccessTokensForUser(currentUser.ID)

	c.JSON(http.StatusOK, gin.H{"message": "All sessions revoked"})
}

// revokeOtherSessions revokes all sessions except the current one
func revokeOtherSessions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(*models.User)
	currentSessionID := uuid.Nil
	if sid, ok := c.Get("session_id"); ok {
		if parsed, err := uuid.Parse(sid.(string)); err == nil {
			currentSessionID = parsed
		}
	}

	sessionService := services.NewSessionService()
	if err := sessionService.RevokeAllOtherSessions(currentUser.ID, currentSessionID); err != nil {
		log.Printf("Failed to revoke other sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All other sessions revoked"})
}
