package services

import (
	"testing"

	"securewallet/internal/models"

	"github.com/google/uuid"
)

func TestRBACConstants(t *testing.T) {
	// Test that role constants are defined correctly
	roles := []string{
		models.RoleAdmin,
		models.RoleUser,
		models.RoleAuditor,
		models.RoleSupport,
	}

	for _, role := range roles {
		if role == "" {
			t.Errorf("Role constant should not be empty")
		}
	}

	// Test that permission constants are defined correctly
	perms := []string{
		models.PermWalletRead,
		models.PermWalletWrite,
		models.PermWalletDelete,
		models.PermTransferRead,
		models.PermTransferWrite,
		models.PermUserRead,
		models.PermUserWrite,
		models.PermUserDelete,
		models.PermAdminAll,
		models.PermAuditRead,
		models.PermSupportRead,
		models.PermSupportWrite,
	}

	for _, perm := range perms {
		if perm == "" {
			t.Errorf("Permission constant should not be empty")
		}
	}
}

func TestRoleTableName(t *testing.T) {
	role := models.Role{}
	if role.TableName() != "roles" {
		t.Errorf("Expected table name 'roles', got '%s'", role.TableName())
	}
}

func TestPermissionTableName(t *testing.T) {
	perm := models.Permission{}
	if perm.TableName() != "permissions" {
		t.Errorf("Expected table name 'permissions', got '%s'", perm.TableName())
	}
}

func TestOAuthProviderTableName(t *testing.T)	{
	provider := models.OAuthProvider{}
	if provider.TableName() != "oauth_providers" {
		t.Errorf("Expected table name 'oauth_providers', got '%s'", provider.TableName())
	}
}

func TestOAuthAccountTableName(t *testing.T) {
	account := models.OAuthAccount{}
	if account.TableName() != "oauth_accounts" {
		t.Errorf("Expected table name 'oauth_accounts', got '%s'", account.TableName())
	}
}

func TestRoleBeforeCreate(t *testing.T) {
	role := models.Role{Name: "test"}
	err := role.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if role.ID == uuid.Nil {
		t.Errorf("BeforeCreate should set ID")
	}
}

func TestPermissionBeforeCreate(t *testing.T) {
	perm := models.Permission{Name: "test", Resource: "test", Action: "test"}
	err := perm.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if perm.ID == uuid.Nil {
		t.Errorf("BeforeCreate should set ID")
	}
}

func TestOAuthProviderBeforeCreate(t *testing.T) {
	provider := models.OAuthProvider{Name: "test", ClientID: "test", ClientSecret: "test", AuthURL: "test", TokenURL: "test"}
	err := provider.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if provider.ID == uuid.Nil {
		t.Errorf("BeforeCreate should set ID")
	}
}

func TestOAuthAccountBeforeCreate(t *testing.T) {
	account := models.OAuthAccount{UserID: uuid.New(), ProviderID: uuid.New(), ProviderName: "test", ProviderUserID: "test"}
	err := account.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if account.ID == uuid.Nil {
		t.Errorf("BeforeCreate should set ID")
	}
}
