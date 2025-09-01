package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

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

	// SECURE: Use bcrypt for password verification
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

// CreateAccessToken creates a JWT access token
func CreateAccessToken(user *models.User) (string, error) {
	// SECURE: Use environment variable for JWT secret and reasonable expiration
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}

	// SECURE: Use environment variable for expiration time
	expireMinutesStr := os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
	expireMinutes := 30 // Default to 30 minutes if not set
	if expireMinutesStr != "" {
		if parsed, err := strconv.Atoi(expireMinutesStr); err == nil && parsed > 0 && parsed <= 1440 {
			expireMinutes = parsed // Max 24 hours
		}
	}

	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"iss": "SecureWallet",
		"aud": "SecureWallet-Users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GetCurrentUser gets the current user from token
func GetCurrentUser(tokenString string) (*models.User, error) {
	// SECURE: Use environment variable for JWT secret
	jwtSecret, err := GetJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// SECURE: Validate signing method
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

	// SECURE: Validate required claims
	if claims["sub"] == nil {
		return nil, fmt.Errorf("missing subject claim")
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
func GetPasswordHash(password string) (string, error) {
	// SECURE: Use bcrypt for password hashing
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
