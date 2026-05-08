package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenPair holds access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
}

// GetJWTSecret returns the JWT secret from environment variable
func GetJWTSecret() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}
	return jwtSecret, nil
}

// InitServices initializes all services
func InitServices() {
	// Initialize services here
}

// AuthenticateUser authenticates a user
func AuthenticateUser(username, password string) (*models.User, error) {
	db := config.GetDB()

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

// CreateAccessToken creates a JWT access token with jti claim for revocation
func CreateAccessToken(user *models.User) (string, error) {
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}

	expireMinutesStr := os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
	expireMinutes := 30
	if expireMinutesStr != "" {
		if parsed, err := strconv.Atoi(expireMinutesStr); err == nil && parsed > 0 && parsed <= 1440 {
			expireMinutes = parsed
		}
	}

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"sub": user.Username,
		"jti": jti,
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"iss": "SecureWallet",
		"aud": "SecureWallet-Users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// CreateRefreshToken generates a cryptographically secure refresh token and persists it
func CreateRefreshToken(user *models.User) (string, error) {
	db := config.GetDB()

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}
	token := hex.EncodeToString(tokenBytes)

	expireHoursStr := os.Getenv("REFRESH_TOKEN_EXPIRE_HOURS")
	expireHours := 168 // default 7 days
	if expireHoursStr != "" {
		if parsed, err := strconv.Atoi(expireHoursStr); err == nil && parsed > 0 && parsed <= 720 {
			expireHours = parsed
		}
	}

	expiresAt := time.Now().Add(time.Duration(expireHours) * time.Hour)

	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := db.Create(&session).Error; err != nil {
		return "", fmt.Errorf("failed to persist refresh token: %v", err)
	}

	return token, nil
}

// ValidateRefreshToken checks if a refresh token exists, is not expired, and belongs to the user
func ValidateRefreshToken(userID uuid.UUID, token string) (*models.Session, error) {
	db := config.GetDB()

	var session models.Session
	if err := db.Where("user_id = ? AND token = ?", userID, token).First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		db.Delete(&session)
		return nil, fmt.Errorf("refresh token expired")
	}

	return &session, nil
}

// RevokeToken deletes a refresh token from the database
func RevokeToken(token string) error {
	db := config.GetDB()

	result := db.Where("token = ?", token).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("token not found")
	}
	return nil
}

// RevokeAllUserTokens deletes all refresh tokens for a user (used on logout-all or password change)
func RevokeAllUserTokens(userID uuid.UUID) error {
	db := config.GetDB()
	return db.Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

// RotateRefreshToken invalidates the old refresh token and issues a new one
func RotateRefreshToken(userID uuid.UUID, oldToken string) (string, error) {
	db := config.GetDB()

	session, err := ValidateRefreshToken(userID, oldToken)
	if err != nil {
		return "", err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to revoke old token: %v", err)
	}

	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("user not found: %v", err)
	}

	newToken, err := CreateRefreshTokenForTx(&user, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit token rotation: %v", err)
	}

	return newToken, nil
}

// CreateRefreshTokenForTx creates a refresh token within an existing transaction
func CreateRefreshTokenForTx(user *models.User, tx *gorm.DB) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}
	token := hex.EncodeToString(tokenBytes)

	expireHoursStr := os.Getenv("REFRESH_TOKEN_EXPIRE_HOURS")
	expireHours := 168
	if expireHoursStr != "" {
		if parsed, err := strconv.Atoi(expireHoursStr); err == nil && parsed > 0 && parsed <= 720 {
			expireHours = parsed
		}
	}

	expiresAt := time.Now().Add(time.Duration(expireHours) * time.Hour)

	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := tx.Create(&session).Error; err != nil {
		return "", fmt.Errorf("failed to persist refresh token: %v", err)
	}

	return token, nil
}

// CreateTokenPair generates both access and refresh tokens for a user
func CreateTokenPair(user *models.User) (*TokenPair, error) {
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := CreateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	expireMinutesStr := os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
	expireMinutes := 30
	if expireMinutesStr != "" {
		if parsed, err := strconv.Atoi(expireMinutesStr); err == nil && parsed > 0 {
			expireMinutes = parsed
		}
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    expireMinutes * 60,
	}, nil
}

// GetCurrentUser gets the current user from token
func GetCurrentUser(tokenString string) (*models.User, error) {
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims["sub"] == nil {
		return nil, fmt.Errorf("missing subject claim")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username")
	}

	jti, _ := claims["jti"].(string)

	db := config.GetDB()
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	blacklistService := NewTokenBlacklistService()
	isBlacklisted, err := blacklistService.IsBlacklisted(jti, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %v", err)
	}
	if isBlacklisted {
		return nil, fmt.Errorf("token has been revoked")
	}

	return &user, nil
}

// GetPasswordHash creates a password hash
func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BlacklistAccessToken extracts the jti from a JWT and adds it to the Redis denylist
func BlacklistAccessToken(tokenString string, userID uuid.UUID) {
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	jti, _ := claims["jti"].(string)
	if jti == "" {
		return
	}

	exp, ok := claims["exp"].(float64)
	ttl := 30 * time.Minute
	if ok {
		remaining := time.Until(time.Unix(int64(exp), 0))
		if remaining > 0 {
			ttl = remaining
		}
	}

	blacklistService := NewTokenBlacklistService()
	_ = blacklistService.BlacklistToken(jti, ttl)
}

// BlacklistAllAccessTokensForUser adds a user-level denylist entry to invalidate all JWTs for that user
func BlacklistAllAccessTokensForUser(userID uuid.UUID) {
	expireHoursStr := os.Getenv("REFRESH_TOKEN_EXPIRE_HOURS")
	expireHours := 168
	if expireHoursStr != "" {
		if parsed, err := strconv.Atoi(expireHoursStr); err == nil && parsed > 0 && parsed <= 720 {
			expireHours = parsed
		}
	}
	ttl := time.Duration(expireHours) * time.Hour

	blacklistService := NewTokenBlacklistService()
	_ = blacklistService.BlacklistAllUserTokens(userID, ttl)
}

func extractJTI(tokenString string) string {
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return ""
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	jti, _ := claims["jti"].(string)
	return jti
}
