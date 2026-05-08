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
	// Wallet permissions
	PermWalletRead   = "wallet:read"
	PermWalletWrite  = "wallet:write"
	PermWalletDelete = "wallet:delete"

	// Transfer/Transaction permissions
	PermTransferRead  = "transfer:read"
	PermTransferWrite = "transfer:write"
	PermTransferDelete = "transfer:delete"

	// User management permissions (admin)
	PermUserRead   = "user:read"
	PermUserWrite  = "user:write"
	PermUserDelete = "user:delete"

	// Session permissions (self-service)
	PermSessionRead  = "session:read"
	PermSessionWrite = "session:write"

	// 2FA permissions (self-service)
	PermTwoFactorRead  = "2fa:read"
	PermTwoFactorWrite = "2fa:write"

	// Login history permissions (self-service)
	PermLoginHistoryRead = "login_history:read"

	// Support permissions
	PermSupportRead  = "support:read"
	PermSupportWrite = "support:write"

	// Audit permissions
	PermAuditRead = "audit:read"

	// Backup permissions (admin)
	PermBackupRead = "backup:read"

	// Security monitoring permissions (admin)
	PermSecurityRead  = "security:read"
	PermSecurityWrite = "security:write"

	// Data management permissions (super-admin only)
	PermDataManage = "data:manage"

	// Cron management permissions (admin)
	PermCronManage = "cron:manage"

	// Blog permissions
	PermBlogRead    = "blog:read"
	PermBlogComment = "blog:comment"

	// Admin super-permission
	PermAdminAll = "admin:*"
)
