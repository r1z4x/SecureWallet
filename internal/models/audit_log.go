package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	Action    string         `json:"action" gorm:"size:100;not null"`
	Resource  string         `json:"resource" gorm:"size:100"`
	Details   string         `json:"details" gorm:"type:text"`
	IPAddress string         `json:"ip_address" gorm:"size:45"`
	UserAgent string         `json:"user_agent" gorm:"size:500"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
