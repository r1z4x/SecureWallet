package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a system role with associated permissions
type Role struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	IsSystem    bool           `json:"is_system" gorm:"default:false"` // System roles cannot be deleted
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Users       []User         `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Permissions []Permission   `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate will set a UUID rather than numeric ID
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Permission represents a granular permission
type Permission struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:100;not null"` // e.g., "wallet:read", "wallet:write", "admin:users"
	Description string         `json:"description" gorm:"size:255"`
	Resource    string         `json:"resource" gorm:"size:50;not null"`   // e.g., "wallet", "transaction", "user"
	Action      string         `json:"action" gorm:"size:50;not null"`     // e.g., "read", "write", "delete", "admin"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName specifies the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// System role names
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleAuditor   = "auditor"
	RoleSupport   = "support"
)

// System permission names
const (
	PermWalletRead    = "wallet:read"
	PermWalletWrite   = "wallet:write"
	PermWalletDelete  = "wallet:delete"
	PermTransferRead  = "transfer:read"
	PermTransferWrite = "transfer:write"
	PermUserRead      = "user:read"
	PermUserWrite     = "user:write"
	PermUserDelete    = "user:delete"
	PermAdminAll      = "admin:*"
	PermAuditRead     = "audit:read"
	PermSupportRead   = "support:read"
	PermSupportWrite  = "support:write"
)
