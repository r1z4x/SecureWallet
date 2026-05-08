package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrRecipientNotFound   = errors.New("recipient not found")
	ErrSelfTransfer        = errors.New("cannot transfer to yourself")
	ErrDuplicateRequest    = errors.New("duplicate request")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrCurrencyMismatch    = errors.New("currency mismatch")
)

const (
	TransferFeePercentage = 0.01
	MinTransferFee        = 1.0
	MaxTransferFee        = 50.0
	MinTransferAmount     = 1.0
	MaxTransferAmount     = 1000.0
	IdempotencyTTL        = 24 * time.Hour
)

type TransferRequest struct {
	SenderUserID    uuid.UUID
	RecipientEmail  string
	Amount          float64
	Description     string
	IdempotencyKey  string
}

type TransferResult struct {
	TransactionID      uuid.UUID
	SenderWalletID     uuid.UUID
	RecipientWalletID  uuid.UUID
	Amount             float64
	Fee                float64
	TotalDeducted      float64
	SenderBalance      float64
	RecipientBalance   float64
}

type DepositRequest struct {
	UserID      uuid.UUID
	Amount      float64
	Description string
}

type DepositResult struct {
	TransactionID uuid.UUID
	WalletID      uuid.UUID
	Amount        float64
	NewBalance    float64
}

type TransferService interface {
	Transfer(req TransferRequest) (*TransferResult, error)
	Deposit(req DepositRequest) (*DepositResult, error)
	CalculateFee(amount float64) float64
}

type transferService struct {
	db *gorm.DB
}

func NewTransferService(db *gorm.DB) TransferService {
	return &transferService{db: db}
}

func (s *transferService) CalculateFee(amount float64) float64 {
	fee := amount * TransferFeePercentage
	if fee < MinTransferFee {
		return MinTransferFee
	}
	if fee > MaxTransferFee {
		return MaxTransferFee
	}
	return fee
}

func (s *transferService) Transfer(req TransferRequest) (*TransferResult, error) {
	if req.Amount < MinTransferAmount || req.Amount > MaxTransferAmount {
		return nil, fmt.Errorf("%w: must be between %.2f and %.2f", ErrInvalidAmount, MinTransferAmount, MaxTransferAmount)
	}

	fee := s.CalculateFee(req.Amount)
	totalAmount := req.Amount + fee

	var result *TransferResult

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if req.IdempotencyKey != "" {
			if _, err := checkIdempotency(tx, req.IdempotencyKey, req.SenderUserID); err != nil {
				if errors.Is(err, ErrDuplicateRequest) {
					return err
				}
			}
		}

		senderWallet, err := lockWalletByUserID(tx, req.SenderUserID)
		if err != nil {
			return fmt.Errorf("%w: sender", ErrWalletNotFound)
		}

		if senderWallet.Balance < totalAmount {
			return ErrInsufficientBalance
		}

		recipientUser, err := getUserByEmail(tx, req.RecipientEmail)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrRecipientNotFound, req.RecipientEmail)
		}

		if recipientUser.ID == req.SenderUserID {
			return ErrSelfTransfer
		}

		recipientWallet, err := lockWalletByUserID(tx, recipientUser.ID)
		if err != nil {
			return fmt.Errorf("%w: recipient", ErrWalletNotFound)
		}

		if senderWallet.Currency != recipientWallet.Currency {
			return ErrCurrencyMismatch
		}

		if err := tx.Model(&senderWallet).Update("balance", senderWallet.Balance-totalAmount).Error; err != nil {
			return fmt.Errorf("failed to deduct from sender: %w", err)
		}

		if err := tx.Model(&recipientWallet).Update("balance", recipientWallet.Balance+req.Amount).Error; err != nil {
			return fmt.Errorf("failed to credit recipient: %w", err)
		}

		outgoingTx := models.Transaction{
			ID:          uuid.New(),
			WalletID:    senderWallet.ID,
			Type:        "transfer_out",
			Amount:      totalAmount,
			Currency:    senderWallet.Currency,
			Description: req.Description,
			Status:      "completed",
		}
		if err := tx.Create(&outgoingTx).Error; err != nil {
			return fmt.Errorf("failed to create outgoing transaction: %w", err)
		}

		incomingTx := models.Transaction{
			ID:          uuid.New(),
			WalletID:    recipientWallet.ID,
			Type:        "transfer_in",
			Amount:      req.Amount,
			Currency:    recipientWallet.Currency,
			Description: req.Description,
			Status:      "completed",
		}
		if err := tx.Create(&incomingTx).Error; err != nil {
			return fmt.Errorf("failed to create incoming transaction: %w", err)
		}

		if req.IdempotencyKey != "" {
			payloadHash := hashPayload(req.RecipientEmail, req.Amount, req.Description)
			idemRecord := models.IdempotencyRecord{
				ID:         uuid.New(),
				Key:        req.IdempotencyKey,
				UserID:     req.SenderUserID,
				Operation:  models.IdempotencyOperationTransfer,
				PayloadHash: payloadHash,
				Status:     models.IdempotencyStatusCompleted,
				HttpStatus: 200,
				ExpiresAt:  time.Now().Add(IdempotencyTTL),
			}
			if err := tx.Create(&idemRecord).Error; err != nil {
				log.Printf("WARNING: failed to create idempotency record: %v", err)
			}
		}

		result = &TransferResult{
			TransactionID:     outgoingTx.ID,
			SenderWalletID:    senderWallet.ID,
			RecipientWalletID: recipientWallet.ID,
			Amount:            req.Amount,
			Fee:               fee,
			TotalDeducted:     totalAmount,
			SenderBalance:     senderWallet.Balance - totalAmount,
			RecipientBalance:  recipientWallet.Balance + req.Amount,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *transferService) Deposit(req DepositRequest) (*DepositResult, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("%w: must be positive", ErrInvalidAmount)
	}

	var result *DepositResult

	err := s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := lockWalletByUserID(tx, req.UserID)
		if err != nil {
			return fmt.Errorf("%w", ErrWalletNotFound)
		}

		if err := tx.Model(&wallet).Update("balance", wallet.Balance+req.Amount).Error; err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		transaction := models.Transaction{
			ID:          uuid.New(),
			WalletID:    wallet.ID,
			Type:        "deposit",
			Amount:      req.Amount,
			Currency:    wallet.Currency,
			Description: req.Description,
			Status:      "completed",
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		result = &DepositResult{
			TransactionID: transaction.ID,
			WalletID:      wallet.ID,
			Amount:        req.Amount,
			NewBalance:    wallet.Balance + req.Amount,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func lockWalletByUserID(tx *gorm.DB, userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_id = ?", userID).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func getUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	defaultDB := config.GetShardManager().GetDefaultDB()
	var user models.User
	if err := defaultDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func checkIdempotency(tx *gorm.DB, key string, userID uuid.UUID) (*models.IdempotencyRecord, error) {
	var record models.IdempotencyRecord
	if err := tx.Where("key = ? AND user_id = ?", key, userID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if time.Now().After(record.ExpiresAt) {
		return nil, nil
	}

	return &record, ErrDuplicateRequest
}

func hashPayload(recipient string, amount float64, description string) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%.2f:%s", recipient, amount, description)))
	return hex.EncodeToString(h.Sum(nil))
}
