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

func TestCreateAccessToken_ContainsJTI(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-unit-tests-only")
	defer os.Unsetenv("JWT_SECRET_KEY")

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	tokenStr, err := services.CreateAccessToken(user)
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Token claims are not MapClaims")
	}

	jti, exists := claims["jti"]
	if !exists {
		t.Fatal("Token missing jti claim")
	}

	jtiStr, ok := jti.(string)
	if !ok {
		t.Fatal("jti claim is not a string")
	}

	if _, err := uuid.Parse(jtiStr); err != nil {
		t.Fatalf("jti claim is not a valid UUID: %v", err)
	}
}

func TestCreateAccessToken_HasIssuerAndAudience(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-unit-tests-only")
	defer os.Unsetenv("JWT_SECRET_KEY")

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	tokenStr, err := services.CreateAccessToken(user)
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["iss"] != "SecureWallet" {
		t.Errorf("Expected iss=SecureWallet, got %v", claims["iss"])
	}

	if claims["aud"] != "SecureWallet-Users" {
		t.Errorf("Expected aud=SecureWallet-Users, got %v", claims["aud"])
	}
}

func TestCreateAccessToken_ExpirationWithinBounds(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-unit-tests-only")
	os.Setenv("ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	defer os.Unsetenv("JWT_SECRET_KEY")
	defer os.Unsetenv("ACCESS_TOKEN_EXPIRE_MINUTES")

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	tokenStr, err := services.CreateAccessToken(user)
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("exp claim is not a number")
	}

	expTime := time.Unix(int64(exp), 0)
	now := time.Now()

	if expTime.Before(now.Add(14 * time.Minute)) {
		t.Errorf("Token expires too soon: %v", expTime.Sub(now))
	}

	if expTime.After(now.Add(16 * time.Minute)) {
		t.Errorf("Token expires too late: %v", expTime.Sub(now))
	}
}

func TestGetJWTSecret_MissingEnv(t *testing.T) {
	os.Unsetenv("JWT_SECRET_KEY")

	_, err := services.GetJWTSecret()
	if err == nil {
		t.Fatal("Expected error when JWT_SECRET_KEY is not set, got nil")
	}
}

func TestGenerateRecoveryCodes_CountAndUniqueness(t *testing.T) {
	svc := services.NewTwoFactorService()

	codes, err := svc.GenerateRecoveryCodes()
	if err != nil {
		t.Fatalf("GenerateRecoveryCodes() error = %v", err)
	}

	if len(codes) != 10 {
		t.Errorf("Expected 10 recovery codes, got %d", len(codes))
	}

	seen := make(map[string]bool)
	for _, code := range codes {
		if seen[code] {
			t.Errorf("Duplicate recovery code: %s", code)
		}
		seen[code] = true

		if len(code) != 8 {
			t.Errorf("Expected recovery code length 8, got %d", len(code))
		}
	}
}

func TestGenerateRecoveryCodes_Alphanumeric(t *testing.T) {
	svc := services.NewTwoFactorService()

	codes, err := svc.GenerateRecoveryCodes()
	if err != nil {
		t.Fatalf("GenerateRecoveryCodes() error = %v", err)
	}

	for _, code := range codes {
		for _, ch := range code {
			if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z')) {
				t.Errorf("Recovery code contains invalid character: %c in code %s", ch, code)
			}
		}
	}
}

func TestValidateCode_ValidTOTP(t *testing.T) {
	svc := services.NewTwoFactorService()

	currentCode, err := svc.GetCurrentCode("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("GetCurrentCode() error = %v", err)
	}

	if !svc.ValidateCode("JBSWY3DPEHPK3PXP", currentCode) {
		t.Error("Expected valid TOTP code to pass validation")
	}
}

func TestValidateCode_InvalidCode(t *testing.T) {
	svc := services.NewTwoFactorService()

	if svc.ValidateCode("JBSWY3DPEHPK3PXP", "000000") {
		t.Error("Expected invalid TOTP code to fail validation")
	}
}
