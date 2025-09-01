package middleware

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Rate limiting storage
var (
	loginAttempts  = make(map[string][]time.Time)
	rateLimitMutex sync.RWMutex
)

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting in development mode
		if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("GIN_MODE") == "debug" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		rateLimitMutex.Lock()
		now := time.Now()

		// Clean old attempts (older than 15 minutes)
		if attempts, exists := loginAttempts[ip]; exists {
			var validAttempts []time.Time
			for _, attempt := range attempts {
				if now.Sub(attempt) < 15*time.Minute {
					validAttempts = append(validAttempts, attempt)
				}
			}
			loginAttempts[ip] = validAttempts
		}

		// Check rate limit (max 5 attempts per 15 minutes)
		if len(loginAttempts[ip]) >= 5 {
			rateLimitMutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add current attempt
		loginAttempts[ip] = append(loginAttempts[ip], now)
		rateLimitMutex.Unlock()

		c.Next()
	}
}
