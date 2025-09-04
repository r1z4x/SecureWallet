package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID               uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Username         string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Name             string         `json:"name" gorm:"size:100;not null"`
	Email            string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Title            string         `json:"title" gorm:"size:100"`
	Avatar           string         `json:"avatar" gorm:"size:500"`
	Bio              string         `json:"bio" gorm:"type:text"`
	PasswordHash     string         `json:"-" gorm:"size:255;not null"`
	TwoFactorSecret  string         `json:"-" gorm:"size:255"`
	TwoFactorEnabled bool           `json:"two_factor_enabled" gorm:"default:false"`
	IsActive         bool           `json:"is_active" gorm:"default:true"`
	IsAdmin          bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Wallets        []Wallet        `json:"wallets,omitempty" gorm:"foreignKey:UserID"`
	Sessions       []Session       `json:"sessions,omitempty" gorm:"foreignKey:UserID"`
	AuditLogs      []AuditLog      `json:"audit_logs,omitempty" gorm:"foreignKey:UserID"`
	SupportTickets []SupportTicket `json:"support_tickets,omitempty" gorm:"foreignKey:UserID"`
	BlogPosts      []BlogPost      `json:"blog_posts,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
