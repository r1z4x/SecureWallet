package services_test

import (
	"os"
	"testing"
	"time"

	"securewallet/internal/models"
	"securewallet/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCreateAccessTokenIncludesExpectedAuthClaims(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-auth-flow-tests")
	os.Setenv("ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	defer os.Unsetenv("JWT_SECRET_KEY")
	defer os.Unsetenv("ACCESS_TOKEN_EXPIRE_MINUTES")

	user := &models.User{ID: uuid.New(), Username: "authflowuser"}

	tokenStr, err := services.CreateAccessToken(user)
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["sub"] != user.Username {
		t.Fatalf("sub claim = %v, want %s", claims["sub"], user.Username)
	}
	if claims["iss"] != "SecureWallet" {
		t.Fatalf("iss claim = %v, want SecureWallet", claims["iss"])
	}
	if claims["aud"] != "SecureWallet-Users" {
		t.Fatalf("aud claim = %v, want SecureWallet-Users", claims["aud"])
	}
	if _, err := uuid.Parse(claims["jti"].(string)); err != nil {
		t.Fatalf("jti claim is not a UUID: %v", err)
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	if exp.Before(time.Now().Add(14*time.Minute)) || exp.After(time.Now().Add(16*time.Minute)) {
		t.Fatalf("expiration outside configured 15 minute window: %s", exp)
	}
}

func TestGetCurrentUserRejectsMalformedTokensBeforeDatabaseLookup(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-auth-flow-tests")
	defer os.Unsetenv("JWT_SECRET_KEY")

	tests := []struct {
		name  string
		token string
	}{
		{"not a jwt", "not-a-jwt"},
		{"empty", ""},
		{"wrong signature", signedTestToken(t, "different-secret", jwt.MapClaims{
			"sub": "authflowuser",
			"exp": time.Now().Add(time.Hour).Unix(),
		})},
		{"expired", signedTestToken(t, "test-secret-for-auth-flow-tests", jwt.MapClaims{
			"sub": "authflowuser",
			"exp": time.Now().Add(-time.Hour).Unix(),
		})},
		{"missing subject", signedTestToken(t, "test-secret-for-auth-flow-tests", jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := services.GetCurrentUser(tt.token)
			if err == nil {
				t.Fatalf("GetCurrentUser returned nil error and user %#v for invalid token", user)
			}
			if user != nil {
				t.Fatalf("GetCurrentUser returned user for invalid token: %#v", user)
			}
		})
	}
}

func TestPasswordHashFlow(t *testing.T) {
	password := "TestPassword123!"

	hash1, err := services.GetPasswordHash(password)
	if err != nil {
		t.Fatalf("GetPasswordHash() error = %v", err)
	}
	hash2, err := services.GetPasswordHash(password)
	if err != nil {
		t.Fatalf("GetPasswordHash() second call error = %v", err)
	}

	if hash1 == password || hash2 == password {
		t.Fatal("password hash must not equal plaintext")
	}
	if hash1 == hash2 {
		t.Fatal("hashes for the same password should use different salts")
	}
}

func signedTestToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	return signed
}
