package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

// TwoFactorService handles 2FA operations
type TwoFactorService struct{}

// NewTwoFactorService creates a new 2FA service
func NewTwoFactorService() *TwoFactorService {
	return &TwoFactorService{}
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TwoFactorService) GenerateSecret(username, email string) (string, string, error) {
	// Generate a random secret
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", "", fmt.Errorf("failed to generate secret: %v", err)
	}

	// Encode as base32
	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// Create TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SecureWallet",
		AccountName: fmt.Sprintf("%s (%s)", username, email),
		Secret:      secret,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP key: %v", err)
	}

	return secretBase32, key.URL(), nil
}

// ValidateCode validates a TOTP code against a secret
func (s *TwoFactorService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateQRCodeURL generates a QR code URL for the secret
func (s *TwoFactorService) GenerateQRCodeURL(secret, username, email string) (string, error) {
	// Decode base32 secret to bytes
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %v", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SecureWallet",
		AccountName: fmt.Sprintf("%s (%s)", username, email),
		Secret:      secretBytes,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code URL: %v", err)
	}

	return key.URL(), nil
}

// GetCurrentCode returns the current TOTP code for testing purposes
func (s *TwoFactorService) GetCurrentCode(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate current code: %v", err)
	}
	return code, nil
}
