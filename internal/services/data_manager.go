package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"securewallet/internal/config"
	"securewallet/internal/models"
	"strings"

	"gorm.io/gorm"
)

// SampleDataManager manages sample data insertion
type SampleDataManager struct {
	db *gorm.DB
}

// NewSampleDataManager creates a new sample data manager
func NewSampleDataManager() *SampleDataManager {
	return &SampleDataManager{
		db: config.GetDB(),
	}
}

// InitializeDatabase initializes database tables
func (dm *SampleDataManager) InitializeDatabase() error {
	log.Println("Initializing database tables...")

	// Auto-migrate all models to create tables
	if err := dm.db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
		&models.Session{},
		&models.AuditLog{},
		&models.SupportTicket{},
		&models.LoginHistory{},
	); err != nil {
		log.Printf("Error migrating database: %v", err)
		return err
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// ClearSampleData clears all sample data from the database
func (dm *SampleDataManager) ClearSampleData() error {
	log.Println("Clearing all sample data...")

	// Use a transaction to ensure all operations are atomic
	tx := dm.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear data in the correct order to avoid foreign key constraint issues
	// Use Where("1=1") to satisfy GORM's requirement for WHERE conditions
	if err := tx.Unscoped().Where("1=1").Delete(&models.Transaction{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear transactions: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.LoginHistory{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear login history: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.SupportTicket{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear support tickets: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.AuditLog{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear audit logs: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.Session{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear sessions: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.Wallet{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear wallets: %v", err)
	}

	if err := tx.Unscoped().Where("1=1").Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear users: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("All sample data cleared successfully")
	return nil
}

// ResetDatabase performs a complete database reset
func (dm *SampleDataManager) ResetDatabase() error {
	log.Println("Starting complete database reset...")

	// Step 1: Clear all existing sample data
	log.Println("Step 1: Clearing existing sample data...")
	if err := dm.ClearSampleData(); err != nil {
		return fmt.Errorf("failed to clear sample data: %v", err)
	}

	// Step 2: Initialize database tables
	log.Println("Step 2: Initializing database tables...")
	if err := dm.InitializeDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Step 3: Create sample data
	log.Println("Step 3: Creating sample data...")
	if err := dm.createSampleUsers(); err != nil {
		return fmt.Errorf("failed to create sample users: %v", err)
	}

	if err := dm.createSampleWallets(); err != nil {
		return fmt.Errorf("failed to create sample wallets: %v", err)
	}

	if err := dm.createSampleTransactions(); err != nil {
		return fmt.Errorf("failed to create sample transactions: %v", err)
	}

	if err := dm.createSampleLoginHistory(); err != nil {
		return fmt.Errorf("failed to create sample login history: %v", err)
	}

	log.Println("Database reset completed successfully")
	return nil
}

// createSampleUsers creates sample users
func (dm *SampleDataManager) createSampleUsers() error {
	// Create users with properly hashed passwords
	// Using SHA1 with salt as expected by the login function

	// Helper function to hash password with salt
	hashPassword := func(password string, userID uint) string {
		salt := fmt.Sprintf("user_%d_salt", userID)
		sha1Hash := sha1.Sum([]byte(password + salt))
		return hex.EncodeToString(sha1Hash[:])
	}

	// Ensure no existing users before creating new ones
	log.Println("Clearing any existing users...")
	if err := dm.db.Unscoped().Where("1=1").Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to clear existing users: %v", err)
	}

	// Verify the table is empty
	var count int64
	if err := dm.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("users table is not empty after clearing, count: %d", count)
	}

	// Create admin user first
	adminUser := models.User{
		Username:     "admin",
		Email:        "admin@securewallet.com",
		PasswordHash: "", // Will be set after creation
		IsActive:     true,
		IsAdmin:      true,
	}

	if err := dm.db.Create(&adminUser).Error; err != nil {
		log.Printf("Error creating admin user: %v", err)
		return err
	}

	// Now hash the password with the actual user ID
	adminUser.PasswordHash = hashPassword("admin123", adminUser.ID)
	if err := dm.db.Save(&adminUser).Error; err != nil {
		log.Printf("Error updating admin password hash: %v", err)
		return err
	}

	// Create standard user
	standardUser := models.User{
		Username:     "user",
		Email:        "user@securewallet.com",
		PasswordHash: "", // Will be set after creation
		IsActive:     true,
		IsAdmin:      false,
	}

	if err := dm.db.Create(&standardUser).Error; err != nil {
		log.Printf("Error creating standard user: %v", err)
		return err
	}

	// Now hash the password with the actual user ID
	standardUser.PasswordHash = hashPassword("password123", standardUser.ID)
	if err := dm.db.Save(&standardUser).Error; err != nil {
		log.Printf("Error updating standard user password hash: %v", err)
		return err
	}

	// Create random users
	firstNames := []string{"John", "Jane", "Bob", "Alice", "Charlie", "Diana", "Edward", "Fiona", "George", "Helen", "Ian", "Julia", "Kevin", "Laura", "Michael", "Nancy", "Oliver", "Patricia", "Quinn", "Rachel", "Steven", "Tina", "Ulysses", "Victoria", "William", "Xena", "Yuki", "Zachary"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson"}
	domains := []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "example.com", "test.com", "demo.com", "sample.com"}

	// Create 50 random users
	for i := 1; i <= 50; i++ {
		firstName := firstNames[i%len(firstNames)]
		lastName := lastNames[i%len(lastNames)]
		domain := domains[i%len(domains)]

		username := fmt.Sprintf("%s_%s_%d", strings.ToLower(firstName), strings.ToLower(lastName), i)
		email := fmt.Sprintf("%s.%s@%s", strings.ToLower(firstName), strings.ToLower(lastName), domain)

		// Make some users admin (about 10%)
		isAdmin := i%10 == 0

		user := models.User{
			Username:     username,
			Email:        email,
			PasswordHash: "", // Will be set after creation
			IsActive:     true,
			IsAdmin:      isAdmin,
		}

		if err := dm.db.Create(&user).Error; err != nil {
			log.Printf("Error creating user %s: %v", username, err)
			continue // Continue with next user instead of failing completely
		}

		// Hash the password with the actual user ID
		user.PasswordHash = hashPassword("password123", user.ID)
		if err := dm.db.Save(&user).Error; err != nil {
			log.Printf("Error updating password hash for user %s: %v", username, err)
		}
	}

	log.Printf("Created admin user, standard user, and 50 random users with properly hashed passwords")

	return nil
}

// createSampleWallets creates sample wallets
func (dm *SampleDataManager) createSampleWallets() error {
	// Get all user IDs
	var users []models.User
	if err := dm.db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY"}

	// Create 500 wallets with random data
	for i := 0; i < 500; i++ {
		// Randomly select a user
		userIndex := i % len(users)
		user := users[userIndex]

		// Random balance between 100 and 50000
		balance := float64(100 + (i * 50) + (int(user.ID) * 10))
		if balance > 50000 {
			balance = 50000 - float64(i*10)
		}

		// Random currency
		currency := currencies[i%len(currencies)]

		wallet := models.Wallet{
			UserID:   user.ID,
			Balance:  balance,
			Currency: currency,
		}

		if err := dm.db.Create(&wallet).Error; err != nil {
			log.Printf("Error creating wallet for user %d: %v", user.ID, err)
			continue // Continue with next wallet instead of failing completely
		}
	}

	log.Printf("Created 500 random wallets")
	return nil
}

// createSampleTransactions creates sample transactions
func (dm *SampleDataManager) createSampleTransactions() error {
	// Get all wallet IDs
	var wallets []models.Wallet
	if err := dm.db.Find(&wallets).Error; err != nil {
		return fmt.Errorf("failed to get wallets: %v", err)
	}

	transactionTypes := []string{"TRANSFER", "DEPOSIT", "WITHDRAWAL", "PAYMENT", "REFUND", "FEE", "BONUS", "INTEREST"}
	statuses := []string{"completed", "pending", "failed", "cancelled", "processing"}
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY"}

	descriptions := []string{
		"Salary payment", "Grocery shopping", "Online purchase", "Restaurant payment", "Gas station",
		"Movie tickets", "Coffee shop", "Book store", "Electronics store", "Clothing store",
		"Pharmacy", "Hardware store", "Subscription payment", "Utility bill", "Insurance payment",
		"Tax payment", "Investment deposit", "Loan payment", "Credit card payment", "ATM withdrawal",
		"Hotel booking", "Flight ticket", "Car rental", "Medical expense", "Dental care",
		"Gym membership", "Netflix subscription", "Spotify premium", "Amazon purchase", "Uber ride",
		"Food delivery", "Home improvement", "Pet care", "Education fee", "Charity donation",
		"Gift purchase", "Travel expense", "Entertainment", "Sports equipment", "Art supplies",
		"Music lesson", "Language course", "Fitness class", "Spa treatment", "Hair salon",
		"Car maintenance", "Home insurance", "Life insurance", "Health insurance", "Property tax",
	}

	// Create 500 transaction records
	for i := 0; i < 500; i++ {
		// Randomly select a wallet
		walletIndex := i % len(wallets)
		wallet := wallets[walletIndex]

		// Random amount between 1 and 5000
		amount := float64(1 + (i * 10) + (int(wallet.ID) * 5))
		if amount > 5000 {
			amount = 5000 - float64(i*5)
		}

		transaction := models.Transaction{
			WalletID:    wallet.ID,
			Type:        transactionTypes[i%len(transactionTypes)],
			Amount:      amount,
			Currency:    currencies[i%len(currencies)],
			Description: descriptions[i%len(descriptions)],
			Status:      statuses[i%len(statuses)],
		}

		if err := dm.db.Create(&transaction).Error; err != nil {
			log.Printf("Error creating transaction: %v", err)
			continue // Continue with next transaction instead of failing completely
		}
	}

	log.Printf("Created 500 random transactions")
	return nil
}

// createSampleLoginHistory creates sample login history
func (dm *SampleDataManager) createSampleLoginHistory() error {
	// Get all user IDs
	var users []models.User
	if err := dm.db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}
	ipAddresses := []string{
		"192.168.1.100", "203.0.113.45", "198.51.100.123", "172.16.0.50",
		"10.0.0.15", "192.168.0.25", "172.20.0.10", "10.1.1.5",
		"192.168.2.30", "203.0.113.67", "198.51.100.89", "172.16.0.75",
		"8.8.8.8", "1.1.1.1", "208.67.222.222", "9.9.9.9",
		"185.228.168.9", "76.76.19.19", "94.140.14.14", "176.103.130.130",
	}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15",
		"Mozilla/5.0 (Android 10; Mobile; rv:68.0) Gecko/68.0 Firefox/68.0",
		"Mozilla/5.0 (iPad; CPU OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/91.0.864.59",
	}
	statuses := []string{"success", "failed", "blocked", "suspicious", "timeout"}
	locations := []string{"Local", "Unknown", "New York", "London", "Tokyo", "Berlin", "Paris", "Sydney", "Toronto", "Singapore", "Mumbai", "São Paulo", "Moscow", "Seoul", "Mexico City", "Cairo", "Lagos", "Bangkok", "Jakarta", "Manila"}

	// Create 500 login history records
	for i := 0; i < 500; i++ {
		// Randomly select a user
		userIndex := i % len(users)
		user := users[userIndex]

		loginHistory := models.LoginHistory{
			UserID:    user.ID,
			IPAddress: ipAddresses[i%len(ipAddresses)],
			UserAgent: userAgents[i%len(userAgents)],
			Status:    statuses[i%len(statuses)],
			Location:  locations[i%len(locations)],
		}

		if err := dm.db.Create(&loginHistory).Error; err != nil {
			log.Printf("Error creating login history: %v", err)
			continue // Continue with next record instead of failing completely
		}
	}

	log.Printf("Created 500 random login history records")
	return nil
}

// GetSampleDataStats returns statistics about sample data
func (dm *SampleDataManager) GetSampleDataStats() map[string]interface{} {
	var userCount, walletCount, transactionCount, loginHistoryCount int64

	dm.db.Model(&models.User{}).Count(&userCount)
	dm.db.Model(&models.Wallet{}).Count(&walletCount)
	dm.db.Model(&models.Transaction{}).Count(&transactionCount)
	dm.db.Model(&models.LoginHistory{}).Count(&loginHistoryCount)

	return map[string]interface{}{
		"users":         userCount,
		"wallets":       walletCount,
		"transactions":  transactionCount,
		"login_history": loginHistoryCount,
	}
}
