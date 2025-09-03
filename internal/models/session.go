package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session
type Session struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	Token     string         `json:"token" gorm:"size:255;not null"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
