package services

import (
	"fmt"
	"log"
	"os"
	"securewallet/internal/config"
	"securewallet/internal/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
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

// InitializeDatabase initializes the database tables
func (dm *SampleDataManager) InitializeDatabase() error {
	log.Println("Initializing database tables...")

	// Drop all existing tables first to ensure clean slate
	log.Println("Dropping all existing tables...")
	if err := dm.db.Migrator().DropTable(
		&models.Transaction{},
		&models.LoginHistory{},
		&models.SupportTicket{},
		&models.AuditLog{},
		&models.Session{},
		&models.Wallet{},
		&models.User{},
	); err != nil {
		log.Printf("Error dropping tables: %v", err)
		return err
	}

	// Auto migrate to create new tables with correct schema
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

	// Step 1: Check if we need to completely recreate the database
	log.Println("Step 1: Checking database schema compatibility...")
	if needsCompleteReset, err := dm.checkSchemaCompatibility(); err != nil {
		log.Printf("Warning: Could not check schema compatibility: %v", err)
	} else if needsCompleteReset {
		log.Println("Schema incompatibility detected, performing complete database recreation...")
		if err := dm.CompleteDatabaseRecreation(); err != nil {
			return fmt.Errorf("failed to recreate database: %v", err)
		}
		return nil
	}

	// Step 2: Clear all existing sample data
	log.Println("Step 2: Clearing existing sample data...")
	if err := dm.ClearSampleData(); err != nil {
		return fmt.Errorf("failed to clear sample data: %v", err)
	}

	// Step 3: Initialize database tables
	log.Println("Step 3: Initializing database tables...")
	if err := dm.InitializeDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Step 4: Create sample data
	log.Println("Step 4: Creating sample data...")
	if err := dm.CreateSampleUsers(); err != nil {
		return fmt.Errorf("failed to create sample users: %v", err)
	}

	if err := dm.CreateSampleWallets(); err != nil {
		return fmt.Errorf("failed to create sample wallets: %v", err)
	}

	if err := dm.CreateSampleTransactions(); err != nil {
		return fmt.Errorf("failed to create sample transactions: %v", err)
	}

	if err := dm.CreateSampleLoginHistory(); err != nil {
		return fmt.Errorf("failed to create sample login history: %v", err)
	}

	log.Println("Database reset completed successfully")
	return nil
}

// checkSchemaCompatibility checks if the existing database schema is compatible with the new UUID-based schema
func (dm *SampleDataManager) checkSchemaCompatibility() (bool, error) {
	// Check if users table exists and has the correct ID column type
	var result struct {
		ColumnType string `gorm:"column:COLUMN_TYPE"`
	}

	err := dm.db.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'id'").Scan(&result).Error
	if err != nil {
		// Table doesn't exist, which is fine
		return false, nil
	}

	// If the ID column is not CHAR(36), we need a complete reset
	if result.ColumnType != "char(36)" {
		log.Printf("Schema incompatibility detected: users.id is %s, expected char(36)", result.ColumnType)
		return true, nil
	}

	return false, nil
}

// CompleteDatabaseRecreation completely drops and recreates the database
func (dm *SampleDataManager) CompleteDatabaseRecreation() error {
	log.Println("Performing complete database recreation...")

	// Get database name
	dbName := dm.db.Migrator().CurrentDatabase()

	// Close current connection
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.Close()

	// Connect to MySQL without specifying a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	tempDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Drop and recreate the database
	log.Printf("Dropping database: %s", dbName)
	if err := tempDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)).Error; err != nil {
		return fmt.Errorf("failed to drop database: %v", err)
	}

	log.Printf("Creating database: %s", dbName)
	if err := tempDB.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error; err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// Close temporary connection
	tempSqlDB, _ := tempDB.DB()
	tempSqlDB.Close()

	// Reconnect to the recreated database
	if err := dm.reconnectToDatabase(); err != nil {
		return fmt.Errorf("failed to reconnect to database: %v", err)
	}

	// Update the global database connection as well
	if err := dm.updateGlobalDatabaseConnection(); err != nil {
		return fmt.Errorf("failed to update global database connection: %v", err)
	}

	log.Println("Database recreation completed successfully")
	return nil
}

// reconnectToDatabase reconnects to the database after recreation
func (dm *SampleDataManager) reconnectToDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	dm.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to reconnect to database: %v", err)
	}

	return nil
}

// CreateSampleUsers creates sample users
func (dm *SampleDataManager) CreateSampleUsers() error {
	// Create users with properly hashed passwords using bcrypt
	hashPassword := func(password string) (string, error) {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(hash), nil
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

	// Now set a secure bcrypt password
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Admin#2025" // Default fallback
	}
	pw, err := hashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %v", err)
	}
	adminUser.PasswordHash = pw
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

	// Now set a secure bcrypt password
	userPassword := os.Getenv("USER_PASSWORD")
	if userPassword == "" {
		userPassword = "User#2025" // Default fallback
	}
	pw2, err := hashPassword(userPassword)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %v", err)
	}
	standardUser.PasswordHash = pw2
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

		// Set a secure bcrypt password
		randomUserPassword := os.Getenv("RANDOM_USER_PASSWORD")
		if randomUserPassword == "" {
			randomUserPassword = "ChangeMe_User#2025" // Default fallback
		}
		if pw3, err := hashPassword(randomUserPassword); err == nil {
			user.PasswordHash = pw3
		} else {
			log.Printf("Error hashing password for user %s: %v", username, err)
		}
		if err := dm.db.Save(&user).Error; err != nil {
			log.Printf("Error updating password hash for user %s: %v", username, err)
		}
	}

	log.Printf("Created admin user, standard user, and 50 random users with properly hashed passwords")

	return nil
}

// CreateSampleWallets creates sample wallets
func (dm *SampleDataManager) CreateSampleWallets() error {
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
		// Use a hash of the user ID to generate consistent but varied balances
		userIDHash := fmt.Sprintf("%x", user.ID)
		balance := float64(100 + (i * 50) + (len(userIDHash) * 10))
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
			log.Printf("Error creating wallet for user %s: %v", user.ID, err)
			continue // Continue with next wallet instead of failing completely
		}
	}

	log.Printf("Created 500 random wallets")
	return nil
}

// CreateSampleTransactions creates sample transactions
func (dm *SampleDataManager) CreateSampleTransactions() error {
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
		// Use a hash of the wallet ID to generate consistent but varied amounts
		walletIDHash := fmt.Sprintf("%x", wallet.ID)
		amount := float64(1 + (i * 10) + (len(walletIDHash) * 5))
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

// CreateSampleLoginHistory creates sample login history
func (dm *SampleDataManager) CreateSampleLoginHistory() error {
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

// updateGlobalDatabaseConnection updates the global database connection after recreation
func (dm *SampleDataManager) updateGlobalDatabaseConnection() error {
	// Create a new database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create new global database connection: %v", err)
	}

	// Update the global database connection using the config package function
	config.UpdateDB(newDB)

	// Also update our local connection to ensure consistency
	dm.db = newDB

	log.Println("Global database connection updated successfully")
	return nil
}
