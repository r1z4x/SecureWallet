package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	Subject     string         `json:"subject" gorm:"size:200;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      string         `json:"status" gorm:"size:20;default:'open'"`
	Priority    string         `json:"priority" gorm:"size:20;default:'medium'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for SupportTicket
func (SupportTicket) TableName() string {
	return "support_tickets"
}

// BeforeCreate will set a UUID rather than numeric ID
func (s *SupportTicket) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
