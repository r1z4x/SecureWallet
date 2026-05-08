package services

import (
	"context"
	"fmt"
	"time"

	"securewallet/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	blacklistTokenPrefix = "jwt:blacklist:token:"
	blacklistUserPrefix  = "jwt:blacklist:user:"
)

// TokenBlacklistService manages JWT invalidation via Redis denylist
type TokenBlacklistService struct {
	rdb *redis.Client
}

// NewTokenBlacklistService creates a new token blacklist service
func NewTokenBlacklistService() *TokenBlacklistService {
	return &TokenBlacklistService{
		rdb: config.GetRedis(),
	}
}

// BlacklistToken adds a specific JWT (by jti) to the denylist with a TTL
func (s *TokenBlacklistService) BlacklistToken(jti string, ttl time.Duration) error {
	if s.rdb == nil {
		return fmt.Errorf("Redis not available")
	}
	key := fmt.Sprintf("%s%s", blacklistTokenPrefix, jti)
	return s.rdb.Set(context.Background(), key, "1", ttl).Err()
}

// IsTokenBlacklisted checks if a specific JWT (by jti) is in the denylist
func (s *TokenBlacklistService) IsTokenBlacklisted(jti string) (bool, error) {
	if s.rdb == nil {
		return false, nil
	}
	key := fmt.Sprintf("%s%s", blacklistTokenPrefix, jti)
	val, err := s.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// BlacklistAllUserTokens adds a user-level denylist entry that invalidates all JWTs for that user
func (s *TokenBlacklistService) BlacklistAllUserTokens(userID uuid.UUID, ttl time.Duration) error {
	if s.rdb == nil {
		return fmt.Errorf("Redis not available")
	}
	key := fmt.Sprintf("%s%s", blacklistUserPrefix, userID.String())
	return s.rdb.Set(context.Background(), key, "1", ttl).Err()
}

// IsUserBlacklisted checks if a user has a blanket token denylist entry
func (s *TokenBlacklistService) IsUserBlacklisted(userID uuid.UUID) (bool, error) {
	if s.rdb == nil {
		return false, nil
	}
	key := fmt.Sprintf("%s%s", blacklistUserPrefix, userID.String())
	val, err := s.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// IsBlacklisted checks both token-level and user-level denylist entries
func (s *TokenBlacklistService) IsBlacklisted(jti string, userID uuid.UUID) (bool, error) {
	tokenBlacklisted, err := s.IsTokenBlacklisted(jti)
	if err != nil {
		return false, err
	}
	if tokenBlacklisted {
		return true, nil
	}
	return s.IsUserBlacklisted(userID)
}
