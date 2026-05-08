package middleware_test

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"securewallet/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateTestToken(user *models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "test-secret-for-middleware-tests"
	}
	return generateTestTokenWithSecret(user, jwtSecret)
}

func generateTestTokenWithSecret(user *models.User, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.Username,
		"jti": uuid.New().String(),
		"exp": time.Now().Add(30 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"iss": "SecureWallet",
		"aud": "SecureWallet-Users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func generateExpiredTestToken(user *models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "test-secret-for-middleware-tests"
	}

	claims := jwt.MapClaims{
		"sub": user.Username,
		"jti": uuid.New().String(),
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
		"iss": "SecureWallet",
		"aud": "SecureWallet-Users",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func formatErrorResponse(msg string, details ...interface{}) string {
	return fmt.Sprintf(msg, details...)
}
