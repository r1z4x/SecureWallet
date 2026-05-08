package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"securewallet/internal/config"

	"github.com/gin-gonic/gin"
)

const (
	twoFARateLimitWindow   = 15 * time.Minute
	twoFARateLimitMax      = 5
	twoFALockoutDuration   = 30 * time.Minute
	twoFARateLimitRedisKey = "2fa_rate:%s"
)

// TwoFARateLimitMiddleware provides per-user rate limiting for 2FA attempts using Redis
func TwoFARateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("GIN_MODE") == "debug" {
			c.Next()
			return
		}

		rdb := config.GetRedis()
		if rdb == nil {
			fallbackInMemory2FARateLimit(c)
			return
		}

		userID := c.Param("user_id")
		if userID == "" {
			var body struct {
				UserID string `json:"user_id"`
			}
			if err := c.ShouldBindBodyWithJSON(&body); err == nil && body.UserID != "" {
				userID = body.UserID
			}
		}
		if userID == "" {
			c.Next()
			return
		}

		lockoutKey := fmt.Sprintf("2fa_lockout:%s", userID)
		rateKey := fmt.Sprintf(twoFARateLimitRedisKey, userID)

		locked, err := rdb.Exists(c.Request.Context(), lockoutKey).Result()
		if err == nil && locked > 0 {
			ttl, _ := rdb.TTL(c.Request.Context(), lockoutKey).Result()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":          "Too many failed 2FA attempts. Account temporarily locked.",
				"retry_after":    int(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		attempts, err := rdb.Incr(c.Request.Context(), rateKey).Result()
		if err != nil {
			c.Next()
			return
		}

		if attempts == 1 {
			rdb.Expire(c.Request.Context(), rateKey, twoFARateLimitWindow)
		}

		maxAttempts := twoFARateLimitMax
		if envMax := os.Getenv("TWO_FA_MAX_ATTEMPTS"); envMax != "" {
			if parsed, err := strconv.Atoi(envMax); err == nil && parsed > 0 {
				maxAttempts = parsed
			}
		}

		if attempts > int64(maxAttempts) {
			lockoutDur := twoFALockoutDuration
			if envLockout := os.Getenv("TWO_FA_LOCKOUT_MINUTES"); envLockout != "" {
				if parsed, err := strconv.Atoi(envLockout); err == nil && parsed > 0 {
					lockoutDur = time.Duration(parsed) * time.Minute
				}
			}

			rdb.Set(c.Request.Context(), lockoutKey, "1", lockoutDur)
			rdb.Del(c.Request.Context(), rateKey)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many failed 2FA attempts. Account temporarily locked.",
				"retry_after": int(lockoutDur.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func fallbackInMemory2FARateLimit(c *gin.Context) {
	c.Next()
}
