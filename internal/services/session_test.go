package services

import (
	"testing"

	"securewallet/internal/models"

	"github.com/google/uuid"
)

func TestSessionTableName(t *testing.T) {
	session := models.Session{}
	if session.TableName() != "sessions" {
		t.Errorf("Expected table name 'sessions', got '%s'", session.TableName())
	}
}

func TestSessionBeforeCreate(t *testing.T) {
	session := models.Session{
		UserID:  uuid.New(),
		Token:   "test_token_12345678",
	}
	err := session.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if session.ID == uuid.Nil {
		t.Errorf("BeforeCreate should set ID")
	}
	if session.TokenLast8 != "12345678" {
		t.Errorf("Expected TokenLast8 '12345678', got '%s'", session.TokenLast8)
	}
	if session.LastAccessed.IsZero() {
		t.Errorf("BeforeCreate should set LastAccessed")
	}
}

func TestSessionBeforeCreateShortToken(t *testing.T) {
	session := models.Session{
		UserID:  uuid.New(),
		Token:   "short",
	}
	err := session.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
	if session.TokenLast8 != "" {
		t.Errorf("Expected empty TokenLast8 for short token, got '%s'", session.TokenLast8)
	}
}
