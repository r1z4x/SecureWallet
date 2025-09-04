package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"securewallet/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get underlying sql.DB object
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Check if we need to handle schema compatibility issues
	if err := handleSchemaCompatibility(); err != nil {
		return fmt.Errorf("failed to handle schema compatibility: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// handleSchemaCompatibility checks and handles schema compatibility issues
func handleSchemaCompatibility() error {
	// Check if we're in a startup reset mode
	if os.Getenv("RESET_DATABASE_ON_STARTUP") == "true" {
		log.Println("Database reset requested on startup, skipping auto-migration")
		return nil
	}

	// Check if we need to force database recreation
	if os.Getenv("FORCE_DATABASE_RECREATION") == "true" {
		log.Println("Force database recreation requested on startup, skipping auto-migration")
		return nil
	}

	// Check if the database schema is compatible
	if isSchemaCompatible() {
		// Schema is compatible, proceed with auto-migration
		log.Println("Schema is compatible, proceeding with auto-migration")
		return autoMigrate()
	}

	// Schema is incompatible, log warning and skip auto-migration
	log.Println("WARNING: Database schema is incompatible with current models")
	log.Println("Please use one of the following options to fix this:")
	log.Println("1. Set RESET_DATABASE_ON_STARTUP=true and restart")
	log.Println("2. Set FORCE_DATABASE_RECREATION=true and restart")
	log.Println("3. Call POST /api/data/force-recreate endpoint")
	log.Println("4. Manually drop and recreate the database")

	// Skip auto-migration to prevent errors
	return nil
}

// isSchemaCompatible checks if the existing database schema is compatible with the new UUID-based schema
func isSchemaCompatible() bool {
	// First check if the database is empty (no tables exist)
	var tableCount int64
	err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tableCount).Error
	if err != nil {
		log.Printf("Warning: Could not check table count: %v", err)
		return false
	}

	// If no tables exist, schema is compatible (will be created)
	if tableCount == 0 {
		return true
	}

	// Check if users table exists and has the correct ID column type
	var result struct {
		ColumnType string `gorm:"column:COLUMN_TYPE"`
	}

	err = DB.Raw("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'id'").Scan(&result).Error
	if err != nil {
		// Table doesn't exist, which means schema is compatible (will be created)
		return true
	}

	// If the ID column is CHAR(36), schema is compatible
	if result.ColumnType == "char(36)" {
		return true
	}

	// Schema is incompatible
	log.Printf("Schema incompatibility detected: users.id is %s, expected char(36)", result.ColumnType)
	return false
}

// autoMigrate runs database migrations
func autoMigrate() error {
	// Import models here to avoid circular imports
	// This will be implemented when we create the models
	return DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
		&models.Session{},
		&models.AuditLog{},
		&models.SupportTicket{},
		&models.LoginHistory{},
		&models.BlogPost{},
		&models.BlogComment{},
		&models.BlogCategory{},
		&models.BlogTag{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// UpdateDB updates the global database connection
// This is useful after database recreation to ensure all parts of the app use the new connection
func UpdateDB(newDB *gorm.DB) {
	DB = newDB
}
