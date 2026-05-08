package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
)

// IdempotencyMiddleware handles idempotency checks for financial operations
// Ensures that identical requests are processed only once, preventing duplicate operations
// Uses Idempotency-Key header and stores records in idempotency_records table
func IdempotencyMiddleware(operation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract idempotency key from header
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required"})
			c.Abort()
			return
		}

		// Validate UUID format
		_, err := uuid.Parse(idempotencyKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid idempotency key format. Must be a valid UUID."})
			c.Abort()
			return
		}
	// Get context without user
	_, exists := c.Get("user")
	if !exists {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		db := DB(c)

		// Check if idempotency record exists
		var existingRecord models.IdempotencyRecord
		if err := db.Where("key = ?", idempotencyKey).First(&existingRecord).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				// Database error
				log.Printf("Database error when checking idempotency record: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				c.Abort()
				return
			}

			// Record doesn't exist - proceed with request
			c.Next()
			return

		}

		// Record exists - return cached response
		if time.Now().After(existingRecord.ExpiresAt) {
			// Record expired - allow processing
			c.Next()
			return
		}

		// Return cached response
		c.Data(existingRecord.HttpStatus, "application/json", []byte(existingRecord.ResponseBody))
		c.Abort()
	}
}

// CreateIdempotencyRecord creates a new idempotency record after successful processing
// Should be called after the operation completes successfully
func CreateIdempotencyRecord(db *gorm.DB, key, operation string, user models.User, payloadHash string, status string, httpStatus int, responseBody string, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	record := models.IdempotencyRecord{
		Key:          key,
		UserID:       user.ID,
		Operation:    operation,
		PayloadHash:  payloadHash,
		Status:       status,
		HttpStatus:   httpStatus,
		ResponseBody: responseBody,
		ExpiresAt:    expiresAt,
	}

	// Use db instead of direct c.MustGet("db")
	result := db.Create(&record)
	if result.Error != nil {
		log.Printf("Failed to create idempotency record: %v", result.Error)
		return result.Error
	}

	return nil
}
