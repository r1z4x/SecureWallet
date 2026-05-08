package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session
type Session struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:char(36);not null;index"`
	Token        string         `json:"-" gorm:"size:255;not null;uniqueIndex"` // Hide token from JSON responses
	TokenLast8   string         `json:"token_last_8" gorm:"size:8"`             // Last 8 chars for identification
	DeviceName   string         `json:"device_name" gorm:"size:100"`            // e.g., "Chrome on macOS"
	IPAddress    string         `json:"ip_address" gorm:"size:45"`              // IPv4 or IPv6
	UserAgent    string         `json:"user_agent" gorm:"size:500"`             // Full user agent string
	IsActive     bool           `json:"is_active" gorm:"default:true;index"`    // Can be revoked without deletion
	ExpiresAt    time.Time      `json:"expires_at"`
	LastAccessed time.Time      `json:"last_accessed"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

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
	if s.TokenLast8 == "" && len(s.Token) >= 8 {
		s.TokenLast8 = s.Token[len(s.Token)-8:]
	}
	if s.LastAccessed.IsZero() {
		s.LastAccessed = time.Now()
	}
	return nil
}

// OAuthProvider represents an OAuth2/OIDC provider configuration
type OAuthProvider struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name         string         `json:"name" gorm:"uniqueIndex;size:50;not null"` // e.g., "google", "github"
	ClientID     string         `json:"client_id" gorm:"size:255;not null"`
	ClientSecret string         `json:"-" gorm:"size:255;not null"`
	AuthURL      string         `json:"auth_url" gorm:"size:500;not null"`
	TokenURL     string         `json:"token_url" gorm:"size:500;not null"`
	UserInfoURL  string         `json:"user_info_url" gorm:"size:500"`
	Scopes       string         `json:"scopes" gorm:"size:500"` // Comma-separated scopes
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for OAuthProvider
func (OAuthProvider) TableName() string {
	return "oauth_providers"
}

// BeforeCreate will set a UUID rather than numeric ID
func (o *OAuthProvider) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// OAuthAccount links a user to an OAuth provider account
type OAuthAccount struct {
	ID            uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:char(36);not null;index"`
	ProviderID    uuid.UUID      `json:"provider_id" gorm:"type:char(36);not null;index"`
	ProviderName  string         `json:"provider_name" gorm:"size:50;not null"`
	ProviderUserID string        `json:"provider_user_id" gorm:"size:255;not null"` // External user ID from provider
	Email         string         `json:"email" gorm:"size:255"`
	AccessToken   string         `json:"-" gorm:"size:1000"`
	RefreshToken  string         `json:"-" gorm:"size:1000"`
	TokenExpiry   time.Time      `json:"token_expiry"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Provider OAuthProvider `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
}

// TableName specifies the table name for OAuthAccount
func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}

// BeforeCreate will set a UUID rather than numeric ID
func (o *OAuthAccount) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
