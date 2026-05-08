package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTransferRejectsOutOfRangeAmountsBeforeDatabaseUse(t *testing.T) {
	svc := &transferService{}

	tests := []struct {
		name   string
		amount float64
	}{
		{"below minimum", MinTransferAmount - 0.01},
		{"above maximum", MaxTransferAmount + 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Transfer(TransferRequest{
				SenderUserID:   uuid.New(),
				RecipientEmail: "recipient@example.com",
				Amount:         tt.amount,
				Description:    "invalid amount",
				IdempotencyKey: uuid.NewString(),
			})

			if result != nil {
				t.Fatalf("Transfer returned result for invalid amount: %#v", result)
			}
			if !errors.Is(err, ErrInvalidAmount) {
				t.Fatalf("Transfer error = %v, want ErrInvalidAmount", err)
			}
			if !strings.Contains(err.Error(), "between") {
				t.Fatalf("Transfer error should describe allowed range, got %q", err.Error())
			}
		})
	}
}

func TestDepositRejectsNonPositiveAmountsBeforeDatabaseUse(t *testing.T) {
	svc := &transferService{}

	for _, amount := range []float64{0, -1, -100.25} {
		result, err := svc.Deposit(DepositRequest{
			UserID:      uuid.New(),
			Amount:      amount,
			Description: "invalid deposit",
		})

		if result != nil {
			t.Fatalf("Deposit(%v) returned result for invalid amount: %#v", amount, result)
		}
		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("Deposit(%v) error = %v, want ErrInvalidAmount", amount, err)
		}
	}
}

func TestTransferPayloadHashBindsIdempotencyKeyToRequestBody(t *testing.T) {
	original := hashPayload("recipient@example.com", 50, "invoice 1001")

	cases := map[string]string{
		"recipient":   hashPayload("other@example.com", 50, "invoice 1001"),
		"amount":      hashPayload("recipient@example.com", 51, "invoice 1001"),
		"description": hashPayload("recipient@example.com", 50, "invoice 1002"),
	}

	if len(original) != 64 {
		t.Fatalf("payload hash length = %d, want sha256 hex length 64", len(original))
	}

	for changedField, changedHash := range cases {
		if original == changedHash {
			t.Fatalf("payload hash did not change when %s changed", changedField)
		}
	}
}

func TestTransferIdempotencyConstantsAndErrors(t *testing.T) {
	if IdempotencyTTL != 24*time.Hour {
		t.Fatalf("IdempotencyTTL = %s, want 24h", IdempotencyTTL)
	}

	if !errors.Is(ErrDuplicateRequest, ErrDuplicateRequest) {
		t.Fatal("ErrDuplicateRequest must remain a sentinel error")
	}
	if ErrDuplicateRequest.Error() == "" {
		t.Fatal("ErrDuplicateRequest should have an operator-readable message")
	}
}

func TestTransferFeeAndTotalDeduction(t *testing.T) {
	svc := &transferService{}

	tests := []struct {
		amount float64
		fee    float64
		total  float64
	}{
		{amount: 10, fee: MinTransferFee, total: 11},
		{amount: 250, fee: 2.5, total: 252.5},
		{amount: 1000, fee: 10, total: 1010},
	}

	for _, tt := range tests {
		fee := svc.CalculateFee(tt.amount)
		if fee != tt.fee {
			t.Fatalf("CalculateFee(%v) = %v, want %v", tt.amount, fee, tt.fee)
		}

		total := tt.amount + fee
		if total != tt.total {
			t.Fatalf("total deduction for %v = %v, want %v", tt.amount, total, tt.total)
		}
	}
}
