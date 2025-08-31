package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

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

	// VULNERABILITY: Weak authentication
	userID := user.ID
	salt := fmt.Sprintf("user_%d_salt", userID)
	sha1Hash := sha1.Sum([]byte(password + salt))
	inputHash := hex.EncodeToString(sha1Hash[:])

	// Normal authentication with plain text support for testing
	if user.PasswordHash == inputHash || user.PasswordHash == password {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

// CreateAccessToken creates a JWT access token
func CreateAccessToken(user *models.User) (string, error) {
	// VULNERABILITY: Weak JWT
	jwtSecret := "expert_secret_key_789"
	expireMinutes := 525600 // 1 year

	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GetCurrentUser gets the current user from token
func GetCurrentUser(tokenString string) (*models.User, error) {
	// VULNERABILITY: Weak token validation
	jwtSecret := "expert_secret_key_789"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username")
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetPasswordHash creates a password hash
func GetPasswordHash(password string) string {
	// VULNERABILITY: Weak password hashing
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
