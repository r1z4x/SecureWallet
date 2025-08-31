package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string         `json:"-" gorm:"size:255;not null"`
	TwoFactorSecret string      `json:"-" gorm:"size:255"`
	TwoFactorEnabled bool       `json:"two_factor_enabled" gorm:"default:false"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	IsAdmin      bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Wallets        []Wallet        `json:"wallets,omitempty" gorm:"foreignKey:UserID"`
	Sessions       []Session       `json:"sessions,omitempty" gorm:"foreignKey:UserID"`
	AuditLogs      []AuditLog      `json:"audit_logs,omitempty" gorm:"foreignKey:UserID"`
	SupportTickets []SupportTicket `json:"support_tickets,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
