package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginHistory represents a user login attempt
type LoginHistory struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);not null;index"`
	IPAddress string         `json:"ip_address" gorm:"size:45"`
	UserAgent string         `json:"user_agent" gorm:"size:500"`
	Status    string         `json:"status" gorm:"size:20;not null"` // success, failed, blocked
	Location  string         `json:"location" gorm:"size:100"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for LoginHistory
func (LoginHistory) TableName() string {
	return "login_history"
}

// BeforeCreate will set a UUID rather than numeric ID
func (l *LoginHistory) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
