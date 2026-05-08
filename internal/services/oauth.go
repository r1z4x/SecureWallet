package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OAuthService handles OAuth2 authentication flows
type OAuthService struct {
	db *gorm.DB
}

// NewOAuthService creates a new OAuth service
func NewOAuthService() *OAuthService {
	return &OAuthService{
		db: config.GetDB(),
	}
}

// GenerateState generates a cryptographically secure state token for OAuth CSRF protection
func (s *OAuthService) GenerateState(provider string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %v", err)
	}
	state := hex.EncodeToString(b)
	return state, nil
}

// GetAuthURL returns the OAuth authorization URL for a provider
func (s *OAuthService) GetAuthURL(providerName, state string) (string, error) {
	var provider models.OAuthProvider
	if err := s.db.Where("name = ? AND is_active = ?", providerName, true).First(&provider).Error; err != nil {
		return "", fmt.Errorf("oauth provider not found or inactive: %s", providerName)
	}

	params := url.Values{}
	params.Set("client_id", provider.ClientID)
	params.Set("redirect_uri", s.getRedirectURL(providerName))
	params.Set("response_type", "code")
	params.Set("state", state)
	if provider.Scopes != "" {
		params.Set("scope", provider.Scopes)
	}

	authURL := provider.AuthURL + "?" + params.Encode()
	return authURL, nil
}

// HandleCallback handles the OAuth callback and returns user info
func (s *OAuthService) HandleCallback(providerName, code, state string) (*OAuthUserInfo, error) {
	var provider models.OAuthProvider
	if err := s.db.Where("name = ? AND is_active = ?", providerName, true).First(&provider).Error; err != nil {
		return nil, fmt.Errorf("oauth provider not found or inactive: %s", providerName)
	}

	tokenResp, err := s.exchangeCodeForToken(provider, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	userInfo, err := s.getUserInfo(provider, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return userInfo, nil
}

// OAuthUserInfo holds user information from OAuth provider
type OAuthUserInfo struct {
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

// TokenResponse holds the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *OAuthService) exchangeCodeForToken(provider models.OAuthProvider, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", provider.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.getRedirectURL(provider.Name))
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(provider.TokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %s", string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokenResp, nil
}

func (s *OAuthService) getUserInfo(provider models.OAuthProvider, accessToken string) (*OAuthUserInfo, error) {
	if provider.UserInfoURL == "" {
		return nil, fmt.Errorf("provider does not support user info endpoint")
	}

	req, err := http.NewRequest("GET", provider.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info request failed: %s", string(body))
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	result := &OAuthUserInfo{}
	if id, ok := userInfo["id"].(string); ok {
		result.ProviderUserID = id
	} else if id, ok := userInfo["sub"].(string); ok {
		result.ProviderUserID = id
	}
	if email, ok := userInfo["email"].(string); ok {
		result.Email = email
	}
	if name, ok := userInfo["name"].(string); ok {
		result.Name = name
	}
	if avatar, ok := userInfo["avatar_url"].(string); ok {
		result.AvatarURL = avatar
	} else if avatar, ok := userInfo["picture"].(string); ok {
		result.AvatarURL = avatar
	}

	return result, nil
}

// LinkOAuthAccount links an OAuth account to an existing user
func (s *OAuthService) LinkOAuthAccount(userID uuid.UUID, providerName string, userInfo *OAuthUserInfo, tokenResp *TokenResponse) error {
	var provider models.OAuthProvider
	if err := s.db.Where("name = ?", providerName).First(&provider).Error; err != nil {
		return fmt.Errorf("provider not found: %v", err)
	}

	var existing models.OAuthAccount
	if err := s.db.Where("provider_id = ? AND provider_user_id = ?", provider.ID, userInfo.ProviderUserID).
		First(&existing).Error; err == nil {
		return fmt.Errorf("oauth account already linked")
	}

	account := models.OAuthAccount{
		UserID:         userID,
		ProviderID:     provider.ID,
		ProviderName:   providerName,
		ProviderUserID: userInfo.ProviderUserID,
		Email:          userInfo.Email,
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   tokenResp.RefreshToken,
	}
	if tokenResp.ExpiresIn > 0 {
		account.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return s.db.Create(&account).Error
}

// FindOrCreateUserByOAuth finds an existing user by OAuth account or creates a new one
func (s *OAuthService) FindOrCreateUserByOAuth(providerName string, userInfo *OAuthUserInfo) (*models.User, error) {
	var provider models.OAuthProvider
	if err := s.db.Where("name = ?", providerName).First(&provider).Error; err != nil {
		return nil, fmt.Errorf("provider not found: %v", err)
	}

	var oauthAccount models.OAuthAccount
	err := s.db.Where("provider_id = ? AND provider_user_id = ?", provider.ID, userInfo.ProviderUserID).
		Preload("User").First(&oauthAccount).Error

	if err == nil {
		return &oauthAccount.User, nil
	}

	if userInfo.Email != "" {
		var user models.User
		if err := s.db.Where("email = ?", userInfo.Email).First(&user).Error; err == nil {
			return nil, fmt.Errorf("user with email already exists, please login with password first")
		}
	}

	user := models.User{
		ID:           uuid.New(),
		Username:     generateOAuthUsername(providerName, userInfo),
		Name:         userInfo.Name,
		Email:        userInfo.Email,
		PasswordHash: "",
		IsActive:     true,
		Avatar:       userInfo.AvatarURL,
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	var userRole models.Role
	if err := tx.Where("name = ?", models.RoleUser).First(&userRole).Error; err == nil {
		tx.Model(&user).Association("Roles").Append(&userRole)
	}

	account := models.OAuthAccount{
		UserID:         user.ID,
		ProviderID:     provider.ID,
		ProviderName:   providerName,
		ProviderUserID: userInfo.ProviderUserID,
		Email:          userInfo.Email,
	}
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create oauth account: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit: %v", err)
	}

	return &user, nil
}

// GetOAuthProviders returns all active OAuth providers
func (s *OAuthService) GetOAuthProviders() ([]models.OAuthProvider, error) {
	var providers []models.OAuthProvider
	if err := s.db.Where("is_active = ?", true).Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch providers: %v", err)
	}
	return providers, nil
}

// RegisterProvider registers a new OAuth provider
func (s *OAuthService) RegisterProvider(name, clientID, clientSecret, authURL, tokenURL, userInfoURL, scopes string) error {
	provider := models.OAuthProvider{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		UserInfoURL:  userInfoURL,
		Scopes:       scopes,
		IsActive:     true,
	}
	return s.db.Create(&provider).Error
}

func (s *OAuthService) getRedirectURL(providerName string) string {
	baseURL := config.GetAPIBaseURL()
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/api/auth/oauth/%s/callback", baseURL, providerName)
}

func generateOAuthUsername(providerName string, userInfo *OAuthUserInfo) string {
	if userInfo.Email != "" {
		parts := strings.Split(userInfo.Email, "@")
		return fmt.Sprintf("%s_%s", providerName, parts[0])
	}
	return fmt.Sprintf("%s_%s", providerName, userInfo.ProviderUserID[:8])
}
