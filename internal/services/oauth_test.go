package services_test

import (
	"testing"

	"securewallet/internal/services"
)

func TestErrEmailCollision_IsExported(t *testing.T) {
	if services.ErrEmailCollision == nil {
		t.Fatal("ErrEmailCollision should be exported and non-nil")
	}
}

func TestLoginMethodsResponse_StructFields(t *testing.T) {
	lm := services.LoginMethodsResponse{
		HasPassword:      true,
		HasTwoFactor:     false,
		OAuthProviders:   []string{"google", "github"},
		RequiresPassword: false,
	}

	if !lm.HasPassword {
		t.Error("HasPassword should be true")
	}
	if lm.HasTwoFactor {
		t.Error("HasTwoFactor should be false")
	}
	if len(lm.OAuthProviders) != 2 {
		t.Errorf("Expected 2 OAuth providers, got %d", len(lm.OAuthProviders))
	}
	if lm.RequiresPassword {
		t.Error("RequiresPassword should be false")
	}
}

func TestLoginMethodsResponse_OAuthOnlyUser(t *testing.T) {
	lm := services.LoginMethodsResponse{
		HasPassword:      false,
		HasTwoFactor:     false,
		OAuthProviders:   []string{"google"},
		RequiresPassword: false,
	}

	if lm.HasPassword {
		t.Error("OAuth-only user should not have password")
	}
	if len(lm.OAuthProviders) != 1 {
		t.Errorf("Expected 1 OAuth provider, got %d", len(lm.OAuthProviders))
	}
}

func TestLoginMethodsResponse_NativeOnlyUser(t *testing.T) {
	lm := services.LoginMethodsResponse{
		HasPassword:      true,
		HasTwoFactor:     true,
		OAuthProviders:   []string{},
		RequiresPassword: false,
	}

	if !lm.HasPassword {
		t.Error("Native user should have password")
	}
	if !lm.HasTwoFactor {
		t.Error("Native user should have 2FA")
	}
	if len(lm.OAuthProviders) != 0 {
		t.Errorf("Expected 0 OAuth providers, got %d", len(lm.OAuthProviders))
	}
}
