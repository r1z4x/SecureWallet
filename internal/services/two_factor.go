package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
	"time"

	"securewallet/internal/config"

	"github.com/pquerna/otp/totp"
)

const (
	recoveryCodeCount  = 10
	recoveryCodeLength = 8
)

// TwoFactorService handles 2FA operations
type TwoFactorService struct{}

// NewTwoFactorService creates a new 2FA service
func NewTwoFactorService() *TwoFactorService {
	return &TwoFactorService{}
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TwoFactorService) GenerateSecret(username, email string) (string, string, error) {
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", "", fmt.Errorf("failed to generate secret: %v", err)
	}

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

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

// GenerateRecoveryCodes generates a set of cryptographically secure recovery codes
func (s *TwoFactorService) GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, recoveryCodeCount)
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := 0; i < recoveryCodeCount; i++ {
		code := make([]byte, recoveryCodeLength)
		for j := range code {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return nil, fmt.Errorf("failed to generate recovery code: %v", err)
			}
			code[j] = charset[n.Int64()]
		}
		codes[i] = string(code)
	}

	return codes, nil
}

// StoreRecoveryCodes stores recovery codes in Redis with 90-day expiration
func (s *TwoFactorService) StoreRecoveryCodes(userID string, codes []string) error {
	rdb := config.GetRedis()
	if rdb == nil {
		return fmt.Errorf("Redis not available for recovery code storage")
	}

	ctx := rdb.Context()
	key := fmt.Sprintf("2fa_recovery:%s", userID)

	for _, code := range codes {
		if err := rdb.SAdd(ctx, key, code).Err(); err != nil {
			return fmt.Errorf("failed to store recovery code: %v", err)
		}
	}

	rdb.Expire(ctx, key, 90*24*time.Hour)
	return nil
}

// ValidateRecoveryCode checks and consumes a recovery code
func (s *TwoFactorService) ValidateRecoveryCode(userID, code string) (bool, error) {
	rdb := config.GetRedis()
	if rdb == nil {
		return false, fmt.Errorf("Redis not available for recovery code validation")
	}

	ctx := rdb.Context()
	key := fmt.Sprintf("2fa_recovery:%s", userID)

	exists, err := rdb.SIsMember(ctx, key, code).Result()
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	if err := rdb.SRem(ctx, key, code).Err(); err != nil {
		return false, fmt.Errorf("failed to consume recovery code: %v", err)
	}

	return true, nil
}

// GetRemainingRecoveryCodes returns the count of unused recovery codes
func (s *TwoFactorService) GetRemainingRecoveryCodes(userID string) (int64, error) {
	rdb := config.GetRedis()
	if rdb == nil {
		return 0, fmt.Errorf("Redis not available")
	}

	ctx := rdb.Context()
	key := fmt.Sprintf("2fa_recovery:%s", userID)

	return rdb.SCard(ctx, key).Result()
}

// ClearRecoveryCodes removes all recovery codes for a user
func (s *TwoFactorService) ClearRecoveryCodes(userID string) error {
	rdb := config.GetRedis()
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}

	ctx := rdb.Context()
	key := fmt.Sprintf("2fa_recovery:%s", userID)

	return rdb.Del(ctx, key).Err()
}
