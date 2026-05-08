package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IdempotencyRecord represents an idempotency record for API requests
// Ensures that identical requests are processed only once, preventing duplicate operations
// Especially critical for financial transactions
type IdempotencyRecord struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`

	Key          string         `json:"key" gorm:"type:char(36);uniqueIndex;not null"` // UUID string value for idempotency key
	UserID       uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	Operation    string         `json:"operation" gorm:"size:100;not null"`      // e.g., "transfer", "deposit", etc.
	PayloadHash  string         `json:"payload_hash" gorm:"size:64;not null"`    // SHA256 of request body to detect tampering
	Status       string         `json:"status" gorm:"size:20;default:'pending'"` // pending, completed, failed
	HttpStatus   int            `json:"http_status"`                             // HTTP status to return on replay
	ResponseBody string         `json:"response_body" gorm:"type:text"`          // Response body to return on replay
	ExpiresAt    time.Time      `json:"expires_at"`                              // Time after which record is deleted
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`


	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for IdempotencyRecord
func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}

// BeforeCreate will set a UUID rather than numeric ID
func (i *IdempotencyRecord) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// Status constants
const (
	IdempotencyStatusPending   = "pending"
	IdempotencyStatusCompleted = "completed"
	IdempotencyStatusFailed    = "failed"
)

// Operation constants
const (
	IdempotencyOperationTransfer = "transfer"
	IdempotencyOperationDeposit  = "deposit"
)
