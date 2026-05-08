package services

import (
	"fmt"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RBACService handles role-based access control operations
type RBACService struct {
	db *gorm.DB
}

// NewRBACService creates a new RBAC service
func NewRBACService() *RBACService {
	return &RBACService{
		db: config.GetDB(),
	}
}

// InitializeSystemRoles creates default system roles and permissions if they don't exist
func (s *RBACService) InitializeSystemRoles() error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create system permissions
	permissions := []models.Permission{
		// Wallet
		{Name: models.PermWalletRead, Description: "Read wallet information", Resource: "wallet", Action: "read"},
		{Name: models.PermWalletWrite, Description: "Modify wallet information", Resource: "wallet", Action: "write"},
		{Name: models.PermWalletDelete, Description: "Delete wallets", Resource: "wallet", Action: "delete"},
		// Transfer/Transaction
		{Name: models.PermTransferRead, Description: "Read transfer history", Resource: "transfer", Action: "read"},
		{Name: models.PermTransferWrite, Description: "Initiate transfers", Resource: "transfer", Action: "write"},
		{Name: models.PermTransferDelete, Description: "Delete transfer records", Resource: "transfer", Action: "delete"},
		// User management (admin)
		{Name: models.PermUserRead, Description: "Read user information (admin)", Resource: "user", Action: "read"},
		{Name: models.PermUserWrite, Description: "Modify user information (admin)", Resource: "user", Action: "write"},
		{Name: models.PermUserDelete, Description: "Delete users (admin)", Resource: "user", Action: "delete"},
		// Sessions (self-service)
		{Name: models.PermSessionRead, Description: "Read own sessions", Resource: "session", Action: "read"},
		{Name: models.PermSessionWrite, Description: "Revoke own sessions", Resource: "session", Action: "write"},
		// 2FA (self-service)
		{Name: models.PermTwoFactorRead, Description: "Read 2FA status and recovery codes", Resource: "2fa", Action: "read"},
		{Name: models.PermTwoFactorWrite, Description: "Enable/disable 2FA and verify codes", Resource: "2fa", Action: "write"},
		// Login history (self-service)
		{Name: models.PermLoginHistoryRead, Description: "Read own login history", Resource: "login_history", Action: "read"},
		// Support
		{Name: models.PermSupportRead, Description: "Read support tickets", Resource: "support", Action: "read"},
		{Name: models.PermSupportWrite, Description: "Create and manage support tickets", Resource: "support", Action: "write"},
		// Audit
		{Name: models.PermAuditRead, Description: "Read audit logs", Resource: "audit", Action: "read"},
		// Backup
		{Name: models.PermBackupRead, Description: "Read backup information", Resource: "backup", Action: "read"},
		// Security
		{Name: models.PermSecurityRead, Description: "Read security alerts and stats", Resource: "security", Action: "read"},
		{Name: models.PermSecurityWrite, Description: "Update security alert status and cleanup", Resource: "security", Action: "write"},
		// Data management (super-admin)
		{Name: models.PermDataManage, Description: "Manage database: reset, recreate, seed", Resource: "data", Action: "manage"},
		// Cron management
		{Name: models.PermCronManage, Description: "Manage cron jobs", Resource: "cron", Action: "manage"},
		// Blog
		{Name: models.PermBlogRead, Description: "Read blog posts and comments", Resource: "blog", Action: "read"},
		{Name: models.PermBlogComment, Description: "Post blog comments", Resource: "blog", Action: "comment"},
		// Admin super-permission
		{Name: models.PermAdminAll, Description: "Full administrative access", Resource: "admin", Action: "*"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		if err := tx.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&perm).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create permission %s: %v", perm.Name, err)
				}
			} else {
				tx.Rollback()
				return fmt.Errorf("failed to check permission %s: %v", perm.Name, err)
			}
		}
	}

	// Create system roles with their permissions
	roles := []struct {
		Name        string
		Description string
		Perms       []string
	}{
		{
			Name:        models.RoleAdmin,
			Description: "Full system administrator access",
			Perms:       []string{models.PermAdminAll},
		},
		{
			Name:        models.RoleUser,
			Description: "Standard user access",
			Perms: []string{
				models.PermWalletRead, models.PermWalletWrite,
				models.PermTransferRead, models.PermTransferWrite,
				models.PermSessionRead, models.PermSessionWrite,
				models.PermTwoFactorRead, models.PermTwoFactorWrite,
				models.PermLoginHistoryRead,
				models.PermSupportRead, models.PermSupportWrite,
				models.PermBlogRead, models.PermBlogComment,
			},
		},
		{
			Name:        models.RoleAuditor,
			Description: "Read-only access to audit logs, security alerts, backups, and system overview",
			Perms: []string{
				models.PermAuditRead,
				models.PermWalletRead, models.PermTransferRead,
				models.PermSecurityRead,
				models.PermBackupRead,
				models.PermUserRead,
			},
		},
		{
			Name:        models.RoleSupport,
			Description: "Customer support access",
			Perms: []string{
				models.PermSupportRead, models.PermSupportWrite,
				models.PermUserRead, models.PermWalletRead,
				models.PermTransferRead,
			},
		},
	}

	for _, roleDef := range roles {
		var role models.Role
		if err := tx.Where("name = ?", roleDef.Name).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				role = models.Role{
					Name:        roleDef.Name,
					Description: roleDef.Description,
					IsSystem:    true,
				}
				if err := tx.Create(&role).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create role %s: %v", roleDef.Name, err)
				}
			} else {
				tx.Rollback()
				return fmt.Errorf("failed to check role %s: %v", roleDef.Name, err)
			}
		}

		// Assign permissions to role
		var perms []models.Permission
		if err := tx.Where("name IN ?", roleDef.Perms).Find(&perms).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find permissions for role %s: %v", roleDef.Name, err)
		}
		if err := tx.Model(&role).Association("Permissions").Replace(&perms); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign permissions to role %s: %v", roleDef.Name, err)
		}
	}

	// Assign admin role to existing admin users
	var adminUsers []models.User
	if err := tx.Where("is_admin = ?", true).Find(&adminUsers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find admin users: %v", err)
	}

	var adminRole models.Role
	if err := tx.Where("name = ?", models.RoleAdmin).First(&adminRole).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find admin role: %v", err)
	}

	for _, user := range adminUsers {
		var existingRoles []models.Role
		if err := tx.Model(&user).Association("Roles").Find(&existingRoles); err == nil {
			hasAdminRole := false
			for _, r := range existingRoles {
				if r.Name == models.RoleAdmin {
					hasAdminRole = true
					break
				}
			}
			if !hasAdminRole {
				if err := tx.Model(&user).Association("Roles").Append(&adminRole); err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to assign admin role to user %s: %v", user.Username, err)
				}
			}
		}
	}

	return tx.Commit().Error
}

// GetUserPermissions returns all permissions for a user through their roles
func (s *RBACService) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	var user models.User
	if err := s.db.Preload("Roles.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	permSet := make(map[string]bool)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permSet[perm.Name] = true
		}
	}

	// Also include legacy admin permissions
	if user.IsAdmin {
		permSet[models.PermAdminAll] = true
	}

	perms := make([]string, 0, len(permSet))
	for perm := range permSet {
		perms = append(perms, perm)
	}
	return perms, nil
}

// HasPermission checks if a user has a specific permission
func (s *RBACService) HasPermission(userID uuid.UUID, permission string) (bool, error) {
	perms, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == permission || p == models.PermAdminAll {
			return true, nil
		}
	}
	return false, nil
}

// HasAnyPermission checks if a user has any of the specified permissions
func (s *RBACService) HasAnyPermission(userID uuid.UUID, permissions []string) (bool, error) {
	for _, perm := range permissions {
		has, err := s.HasPermission(userID, perm)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RBACService) AssignRoleToUser(userID uuid.UUID, roleName string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	var role models.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	if err := s.db.Model(&user).Association("Roles").Append(&role); err != nil {
		return fmt.Errorf("failed to assign role: %v", err)
	}
	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *RBACService) RemoveRoleFromUser(userID uuid.UUID, roleName string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	var role models.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	if err := s.db.Model(&user).Association("Roles").Delete(&role); err != nil {
		return fmt.Errorf("failed to remove role: %v", err)
	}
	return nil
}

// GetUserRoles returns all roles for a user
func (s *RBACService) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	var user models.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user.Roles, nil
}

// CreateRole creates a new custom role
func (s *RBACService) CreateRole(name, description string, permissionNames []string) (*models.Role, error) {
	var existing models.Role
	if err := s.db.Where("name = ?", name).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("role already exists")
	}

	role := models.Role{
		Name:        name,
		Description: description,
		IsSystem:    false,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %v", err)
	}

	if len(permissionNames) > 0 {
		var perms []models.Permission
		if err := s.db.Where("name IN ?", permissionNames).Find(&perms).Error; err != nil {
			return nil, fmt.Errorf("failed to find permissions: %v", err)
		}
		if err := s.db.Model(&role).Association("Permissions").Append(&perms); err != nil {
			return nil, fmt.Errorf("failed to assign permissions: %v", err)
		}
	}

	return &role, nil
}

// DeleteRole deletes a custom role (system roles cannot be deleted)
func (s *RBACService) DeleteRole(roleID uuid.UUID) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	if role.IsSystem {
		return fmt.Errorf("cannot delete system role")
	}

	return s.db.Delete(&role).Error
}

// ListRoles returns all roles
func (s *RBACService) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to list roles: %v", err)
	}
	return roles, nil
}

// ListPermissions returns all permissions
func (s *RBACService) ListPermissions() ([]models.Permission, error) {
	var perms []models.Permission
	if err := s.db.Find(&perms).Error; err != nil {
		return nil, fmt.Errorf("failed to list permissions: %v", err)
	}
	return perms, nil
}
