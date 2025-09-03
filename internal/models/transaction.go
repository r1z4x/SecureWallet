package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	WalletID    uuid.UUID      `json:"wallet_id" gorm:"type:char(36);not null"`
	Type        string         `json:"type" gorm:"size:20;not null"` // deposit, withdrawal, transfer
	Amount      float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string         `json:"currency" gorm:"size:3;default:'USD'"`
	Description string         `json:"description" gorm:"size:255"`
	Status      string         `json:"status" gorm:"size:20;default:'pending'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Wallet Wallet `json:"wallet,omitempty" gorm:"foreignKey:WalletID"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
