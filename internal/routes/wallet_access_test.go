package routes

import (
	"os"
	"strings"
	"testing"
)

func TestWalletRoutesRequireAuthenticatedOwnershipFilters(t *testing.T) {
	source := readRouteSource(t, "wallets.go")

	required := []string{
		`Where("user_id = ?", currentUser.ID).Find(&wallets)`,
		`Where("id = ? AND user_id = ?", id, currentUser.ID).First(&wallet)`,
		`Where("user_id = ?", currentUser.ID).First(&userWallet)`,
		`Where("wallet_id = ?", userWallet.ID).Count(&transactionCount)`,
	}

	for _, fragment := range required {
		if !strings.Contains(source, fragment) {
			t.Fatalf("wallet access route is missing ownership guard: %s", fragment)
		}
	}
}

func TestTransactionRoutesResolveWalletBeforeListingTransactions(t *testing.T) {
	source := readRouteSource(t, "transactions.go")

	required := []string{
		`Where("user_id = ?", currentUser.ID).First(&userWallet)`,
		`Where("wallet_id = ?", userWallet.ID)`,
	}

	for _, fragment := range required {
		if !strings.Contains(source, fragment) {
			t.Fatalf("transaction route is missing authenticated wallet guard: %s", fragment)
		}
	}
}

func readRouteSource(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read route source %s: %v", name, err)
	}

	return string(data)
}
