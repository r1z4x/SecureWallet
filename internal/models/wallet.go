package models

import (
	"time"

	"gorm.io/gorm"
)

// Wallet represents a user's wallet
type Wallet struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Balance   float64        `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Currency  string         `json:"currency" gorm:"size:3;default:'USD'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:WalletID"`
}

// TableName specifies the table name for Wallet
func (Wallet) TableName() string {
	return "wallets"
}
