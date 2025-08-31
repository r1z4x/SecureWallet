package models

import (
	"time"

	"gorm.io/gorm"
)

// SupportTicket represents a support ticket
type SupportTicket struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
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
