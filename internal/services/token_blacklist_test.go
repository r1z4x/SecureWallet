package services

import (
	"os"
	"testing"
	"time"

	"securewallet/internal/models"

	"github.com/google/uuid"
)

func TestTokenBlacklistService_IsTokenBlacklisted_WhenRedisUnavailable(t *testing.T) {
	svc := &TokenBlacklistService{rdb: nil}

	result, err := svc.IsTokenBlacklisted("some-jti")
	if err != nil {
		t.Fatalf("expected no error when Redis unavailable, got: %v", err)
	}
	if result {
		t.Fatal("expected false when Redis unavailable")
	}
}

func TestTokenBlacklistService_IsUserBlacklisted_WhenRedisUnavailable(t *testing.T) {
	svc := &TokenBlacklistService{rdb: nil}
	userID := uuid.New()

	result, err := svc.IsUserBlacklisted(userID)
	if err != nil {
		t.Fatalf("expected no error when Redis unavailable, got: %v", err)
	}
	if result {
		t.Fatal("expected false when Redis unavailable")
	}
}

func TestTokenBlacklistService_IsBlacklisted_WhenRedisUnavailable(t *testing.T) {
	svc := &TokenBlacklistService{rdb: nil}
	userID := uuid.New()

	result, err := svc.IsBlacklisted("some-jti", userID)
	if err != nil {
		t.Fatalf("expected no error when Redis unavailable, got: %v", err)
	}
	if result {
		t.Fatal("expected false when Redis unavailable")
	}
}

func TestTokenBlacklistService_BlacklistToken_WhenRedisUnavailable(t *testing.T) {
	svc := &TokenBlacklistService{rdb: nil}

	err := svc.BlacklistToken("some-jti", 30*time.Minute)
	if err == nil {
		t.Fatal("expected error when Redis unavailable")
	}
}

func TestTokenBlacklistService_BlacklistAllUserTokens_WhenRedisUnavailable(t *testing.T) {
	svc := &TokenBlacklistService{rdb: nil}
	userID := uuid.New()

	err := svc.BlacklistAllUserTokens(userID, 30*time.Minute)
	if err == nil {
		t.Fatal("expected error when Redis unavailable")
	}
}

func TestBlacklistAccessToken_BlacklistsTokenInRedis(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "test-secret-blacklist")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET_KEY", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET_KEY")
		}
	}()

	userID := uuid.New()
	user := &models.User{ID: userID, Username: "blacklist-test-user"}
	token, err := CreateAccessToken(user)
	if err != nil {
		t.Skipf("cannot create test token: %v", err)
	}

	BlacklistAccessToken(token, userID)

	svc := NewTokenBlacklistService()
	if svc.rdb == nil {
		t.Skip("Redis not available, skipping blacklist verification")
	}

	isBlacklisted, err := svc.IsTokenBlacklisted(extractJTI(token))
	if err != nil {
		t.Fatalf("failed to check blacklist: %v", err)
	}
	if !isBlacklisted {
		t.Fatal("expected token to be blacklisted after BlacklistAccessToken call")
	}
}

func TestBlacklistAllAccessTokensForUser_BlacklistsUserInRedis(t *testing.T) {
	userID := uuid.New()
	BlacklistAllAccessTokensForUser(userID)

	svc := NewTokenBlacklistService()
	if svc.rdb == nil {
		t.Skip("Redis not available, skipping blacklist verification")
	}

	isBlacklisted, err := svc.IsUserBlacklisted(userID)
	if err != nil {
		t.Fatalf("failed to check user blacklist: %v", err)
	}
	if !isBlacklisted {
		t.Fatal("expected user to be blacklisted after BlacklistAllAccessTokensForUser call")
	}
}

func TestGetCurrentUser_RejectsBlacklistedToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "test-secret-blacklist-reject")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET_KEY", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET_KEY")
		}
	}()

	svc := NewTokenBlacklistService()
	if svc.rdb == nil {
		t.Skip("Redis not available, skipping blacklist rejection test")
	}

	userID := uuid.New()
	user := &models.User{ID: userID, Username: "blacklist-reject-user"}
	token, err := CreateAccessToken(user)
	if err != nil {
		t.Skipf("cannot create test token: %v", err)
	}

	jti := extractJTI(token)
	err = svc.BlacklistToken(jti, 30*time.Minute)
	if err != nil {
		t.Fatalf("failed to blacklist token: %v", err)
	}

	_, err = GetCurrentUser(token)
	if err == nil {
		t.Fatal("expected error for blacklisted token, got nil")
	}
}

func TestGetCurrentUser_RejectsBlacklistedUser(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "test-secret-blacklist-user-reject")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET_KEY", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET_KEY")
		}
	}()

	svc := NewTokenBlacklistService()
	if svc.rdb == nil {
		t.Skip("Redis not available, skipping user blacklist rejection test")
	}

	userID := uuid.New()
	user := &models.User{ID: userID, Username: "blacklist-user-reject"}
	token, err := CreateAccessToken(user)
	if err != nil {
		t.Skipf("cannot create test token: %v", err)
	}

	err = svc.BlacklistAllUserTokens(userID, 30*time.Minute)
	if err != nil {
		t.Fatalf("failed to blacklist user: %v", err)
	}

	_, err = GetCurrentUser(token)
	if err == nil {
		t.Fatal("expected error for blacklisted user token, got nil")
	}
}
