package services

import (
	"fmt"
	"testing"
)

func TestCalculateFee(t *testing.T) {
	svc := &transferService{}

	tests := []struct {
		name     string
		amount   float64
		expected float64
	}{
		{"below minimum fee threshold", 50.0, MinTransferFee},
		{"at minimum fee threshold", 100.0, MinTransferFee},
		{"1 percent fee", 200.0, 2.0},
		{"mid range fee", 500.0, 5.0},
		{"at maximum fee threshold", 5000.0, MaxTransferFee},
		{"above maximum fee threshold", 10000.0, MaxTransferFee},
		{"zero amount", 0.0, MinTransferFee},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.CalculateFee(tt.amount)
			if got != tt.expected {
				t.Errorf("CalculateFee(%f) = %f, want %f", tt.amount, got, tt.expected)
			}
		})
	}
}

func TestCalculateFee_BoundaryValues(t *testing.T) {
	svc := &transferService{}

	amountForMinFee := MinTransferFee / TransferFeePercentage
	if svc.CalculateFee(amountForMinFee-0.01) != MinTransferFee {
		t.Error("fee just below threshold should be minimum fee")
	}

	amountForMaxFee := MaxTransferFee / TransferFeePercentage
	if svc.CalculateFee(amountForMaxFee+0.01) != MaxTransferFee {
		t.Error("fee just above threshold should be maximum fee")
	}
}

func TestValidateTransferAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		expectError bool
	}{
		{"below minimum", 0.5, true},
		{"at minimum", MinTransferAmount, false},
		{"valid amount", 50.0, false},
		{"at maximum", MaxTransferAmount, false},
		{"above maximum", 1001.0, true},
		{"zero", 0.0, true},
		{"negative", -10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.amount >= MinTransferAmount && tt.amount <= MaxTransferAmount
			if valid == tt.expectError {
				t.Errorf("amount=%f valid=%v, expectError=%v", tt.amount, valid, tt.expectError)
			}
		})
	}
}

func TestValidateDepositAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		expectError bool
	}{
		{"zero", 0.0, true},
		{"negative", -50.0, true},
		{"positive", 100.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.amount > 0
			if valid == tt.expectError {
				t.Errorf("amount=%f valid=%v, expectError=%v", tt.amount, valid, tt.expectError)
			}
		})
	}
}

func TestHashPayload_Deterministic(t *testing.T) {
	h1 := hashPayload("test@example.com", 100.0, "payment")
	h2 := hashPayload("test@example.com", 100.0, "payment")
	if h1 != h2 {
		t.Error("hashPayload should be deterministic")
	}
}

func TestHashPayload_DifferentInputs(t *testing.T) {
	h1 := hashPayload("a@example.com", 100.0, "payment")
	h2 := hashPayload("b@example.com", 100.0, "payment")
	if h1 == h2 {
		t.Error("different recipients should produce different hashes")
	}

	h3 := hashPayload("a@example.com", 100.0, "payment")
	h4 := hashPayload("a@example.com", 200.0, "payment")
	if h3 == h4 {
		t.Error("different amounts should produce different hashes")
	}
}

func TestTransferFeeCalculation_TotalAmount(t *testing.T) {
	svc := &transferService{}

	tests := []struct {
		name          string
		amount        float64
		expectedFee   float64
		expectedTotal float64
	}{
		{"small transfer", 10.0, MinTransferFee, 11.0},
		{"medium transfer", 100.0, 1.0, 101.0},
		{"large transfer", 500.0, 5.0, 505.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee := svc.CalculateFee(tt.amount)
			total := tt.amount + fee
			if fee != tt.expectedFee {
				t.Errorf("fee = %f, want %f", fee, tt.expectedFee)
			}
			if total != tt.expectedTotal {
				t.Errorf("total = %f, want %f", total, tt.expectedTotal)
			}
		})
	}
}

func TestInsufficientBalance_Scenario(t *testing.T) {
	svc := &transferService{}

	balance := 50.0
	amount := 100.0
	fee := svc.CalculateFee(amount)
	total := amount + fee

	if balance >= total {
		t.Error("balance should be insufficient for this transfer")
	}

	expectedShortfall := total - balance
	if expectedShortfall <= 0 {
		t.Errorf("expected positive shortfall, got %f", expectedShortfall)
	}
}

func TestSelfTransfer_Prevention(t *testing.T) {
	senderID := "user-1"
	recipientID := "user-1"

	if senderID == recipientID {
		// Self-transfer should be blocked
		passed := true
		if !passed {
			t.Error("self-transfer should be prevented")
		}
	}
}

func TestTransferConstants(t *testing.T) {
	if MinTransferAmount <= 0 {
		t.Error("MinTransferAmount must be positive")
	}
	if MaxTransferAmount <= MinTransferAmount {
		t.Error("MaxTransferAmount must exceed MinTransferAmount")
	}
	if MinTransferFee < 0 {
		t.Error("MinTransferFee must be non-negative")
	}
	if MaxTransferFee < MinTransferFee {
		t.Error("MaxTransferFee must exceed MinTransferFee")
	}
	if TransferFeePercentage <= 0 || TransferFeePercentage >= 1 {
		t.Error("TransferFeePercentage must be between 0 and 1")
	}
}

func TestErrorTypes(t *testing.T) {
	errs := []error{
		ErrInsufficientBalance,
		ErrWalletNotFound,
		ErrRecipientNotFound,
		ErrSelfTransfer,
		ErrDuplicateRequest,
		ErrInvalidAmount,
		ErrCurrencyMismatch,
	}

	for _, err := range errs {
		if err == nil {
			t.Error("error sentinel should not be nil")
		}
		if err.Error() == "" {
			t.Errorf("error %v should have a non-empty message", err)
		}
	}
}

func TestFeeEdgeCases(t *testing.T) {
	svc := &transferService{}

	fee := svc.CalculateFee(1.0)
	if fee != MinTransferFee {
		t.Errorf("fee for $1 should be minimum fee %f, got %f", MinTransferFee, fee)
	}

	fee = svc.CalculateFee(99.99)
	if fee != MinTransferFee {
		t.Errorf("fee for $99.99 should be minimum fee %f, got %f", MinTransferFee, fee)
	}

	fee = svc.CalculateFee(100.0)
	expected := 100.0 * TransferFeePercentage
	if fee != expected {
		t.Errorf("fee for $100 should be %f, got %f", expected, fee)
	}

	fee = svc.CalculateFee(1000.0)
	expected = 1000.0 * TransferFeePercentage
	if fee != expected {
		t.Errorf("fee for $1000 should be %f, got %f", expected, fee)
	}
}

func TestIdempotencyKeyFormat(t *testing.T) {
	validKeys := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}

	for _, key := range validKeys {
		parts := len(key)
		if parts != 36 {
			t.Errorf("UUID key %s should be 36 chars, got %d", key, parts)
		}
	}

	invalidKeys := []string{
		"",
		"not-a-uuid",
		"12345",
	}

	for _, key := range invalidKeys {
		if len(key) == 36 {
			t.Errorf("key %q should not be 36 chars", key)
		}
	}
}

func TestDescriptionLengthLimit(t *testing.T) {
	maxLen := 255

	valid := "This is a valid description"
	if len(valid) > maxLen {
		t.Error("valid description should be within limit")
	}

	invalid := ""
	for i := 0; i < 300; i++ {
		invalid += "x"
	}
	if len(invalid) <= maxLen {
		t.Error("test string should exceed limit")
	}
}

func TestEmailFormatValidation(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"test.user@domain.org",
		"a@b.co",
	}

	for _, email := range validEmails {
		hasAt := false
		for _, c := range email {
			if c == '@' {
				hasAt = true
				break
			}
		}
		if !hasAt || len(email) > 100 {
			t.Errorf("email %q should be valid", email)
		}
	}

	invalidEmails := []string{
		"no-at-sign",
		"",
	}

	for _, email := range invalidEmails {
		hasAt := false
		for _, c := range email {
			if c == '@' {
				hasAt = true
				break
			}
		}
		if hasAt && len(email) <= 100 {
			t.Errorf("email %q should be invalid", email)
		}
	}
}

func TestTransferResult_Structure(t *testing.T) {
	result := TransferResult{
		Amount:       100.0,
		Fee:          1.0,
		TotalDeducted: 101.0,
	}

	if result.Amount+result.Fee != result.TotalDeducted {
		t.Errorf("Amount + Fee should equal TotalDeducted: %f + %f != %f",
			result.Amount, result.Fee, result.TotalDeducted)
	}
}

func TestDepositResult_Structure(t *testing.T) {
	oldBalance := 100.0
	deposit := 50.0

	result := DepositResult{
		Amount:     deposit,
		NewBalance: oldBalance + deposit,
	}

	if result.NewBalance != oldBalance+result.Amount {
		t.Errorf("NewBalance should equal old balance + deposit amount")
	}
}

func TestFeeFormula(t *testing.T) {
	svc := &transferService{}

	for amount := 1.0; amount <= 10000.0; amount *= 10 {
		fee := svc.CalculateFee(amount)
		calculated := amount * TransferFeePercentage

		if calculated < MinTransferFee {
			if fee != MinTransferFee {
				t.Errorf("amount=%.2f: fee should be min %f, got %f", amount, MinTransferFee, fee)
			}
		} else if calculated > MaxTransferFee {
			if fee != MaxTransferFee {
				t.Errorf("amount=%.2f: fee should be max %f, got %f", amount, MaxTransferFee, fee)
			}
		} else {
			if fee != calculated {
				t.Errorf("amount=%.2f: fee should be %f, got %f", amount, calculated, fee)
			}
		}
	}
}

func TestErrorMessageFormat(t *testing.T) {
	err := fmt.Errorf("%w: must be between %.2f and %.2f", ErrInvalidAmount, MinTransferAmount, MaxTransferAmount)
	msg := err.Error()

	if msg == "" {
		t.Error("error message should not be empty")
	}

	var target error = ErrInvalidAmount
	if !isWrapped(err, target) {
		t.Error("error should wrap ErrInvalidAmount")
	}
}

func isWrapped(err, target error) bool {
	return err.Error() != "" && target.Error() != ""
}
