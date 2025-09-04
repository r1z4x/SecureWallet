package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"securewallet/internal/config"
	"securewallet/internal/models"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BackupService handles database backups
type BackupService struct {
	db *gorm.DB
}

// BackupConfig holds backup configuration
type BackupConfig struct {
	MaxBackups     int           // Maximum number of backups to keep
	BackupInterval time.Duration // How often to create backups
	Enabled        bool          // Whether automatic backups are enabled
}

// Default backup configuration
var DefaultBackupConfig = BackupConfig{
	MaxBackups:     7,             // Keep max 7 backups
	BackupInterval: 1 * time.Hour, // Create backup every hour
	Enabled:        true,          // Enable automatic backups
}

// NewBackupService creates a new backup service
func NewBackupService() *BackupService {
	bs := &BackupService{
		db: config.GetDB(),
	}

	// Start automatic backup scheduler
	go bs.startBackupScheduler()

	return bs
}

// BackupData creates a backup of the database
func (bs *BackupService) BackupData() error {
	// Create backup directory if it doesn't exist
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.json", timestamp))

	// Create backup data structure
	backupData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"data":      map[string]interface{}{},
	}

	// Backup users (with masked sensitive data)
	var users []models.User
	if err := bs.db.Select("id, username, email, is_active, is_admin, created_at, updated_at").Find(&users).Error; err != nil {
		log.Printf("Warning: Failed to backup users: %v", err)
	} else {
		// Mask sensitive email addresses
		for i := range users {
			if users[i].Email != "" {
				parts := strings.Split(users[i].Email, "@")
				if len(parts) == 2 {
					username := parts[0]
					domain := parts[1]
					if len(username) > 2 {
						users[i].Email = username[:2] + "***@" + domain
					} else {
						users[i].Email = "***@" + domain
					}
				}
			}
		}
		backupData["data"].(map[string]interface{})["users"] = users
	}

	// Backup wallets (VULNERABLE: This leaks wallet IDs and balances)
	// Limit to first 50 wallets to reduce backup size
	var wallets []models.Wallet
	if err := bs.db.Limit(50).Find(&wallets).Error; err != nil {
		log.Printf("Warning: Failed to backup wallets: %v", err)
	} else {
		// VULNERABLE: Include sensitive wallet information in backup
		log.Printf("Found %d wallets to backup (limited to 50)", len(wallets))

		// Get user information for each wallet and mask sensitive data
		for i := range wallets {
			var user models.User
			if err := bs.db.First(&user, wallets[i].UserID).Error; err == nil {
				// Mask user email
				if user.Email != "" {
					parts := strings.Split(user.Email, "@")
					if len(parts) == 2 {
						username := parts[0]
						domain := parts[1]
						if len(username) > 2 {
							user.Email = username[:2] + "***@" + domain
						} else {
							user.Email = "***@" + domain
						}
					}
				}
				wallets[i].User = user
			}
		}

		backupData["data"].(map[string]interface{})["wallets"] = wallets
	}

	// Backup transactions (VULNERABLE: This leaks transaction details)
	// Limit to first 100 transactions to reduce backup size
	var transactions []models.Transaction
	if err := bs.db.Limit(100).Find(&transactions).Error; err != nil {
		log.Printf("Warning: Failed to backup transactions: %v", err)
	} else {
		log.Printf("Found %d transactions to backup (limited to 100)", len(transactions))

		// Mask sensitive transaction descriptions
		for i := range transactions {
			if len(transactions[i].Description) > 20 {
				transactions[i].Description = transactions[i].Description[:20] + "..."
			}
		}

		backupData["data"].(map[string]interface{})["transactions"] = transactions
	}

	// Backup audit logs (VULNERABLE: This leaks access patterns)
	// Limit to first 50 audit logs to reduce backup size
	var auditLogs []models.AuditLog
	if err := bs.db.Limit(50).Find(&auditLogs).Error; err != nil {
		log.Printf("Warning: Failed to backup audit logs: %v", err)
	} else {
		log.Printf("Found %d audit logs to backup (limited to 50)", len(auditLogs))

		// Mask IP addresses
		for i := range auditLogs {
			if auditLogs[i].IPAddress != "" {
				parts := strings.Split(auditLogs[i].IPAddress, ".")
				if len(parts) == 4 {
					auditLogs[i].IPAddress = parts[0] + "." + parts[1] + ".***.***"
				}
			}
		}

		backupData["data"].(map[string]interface{})["audit_logs"] = auditLogs
	}

	// Backup security alerts (VULNERABLE: This leaks security information)
	var securityAlerts []SecurityAlert
	if err := bs.db.Limit(20).Find(&securityAlerts).Error; err != nil {
		log.Printf("Warning: Failed to backup security alerts: %v", err)
	} else {
		log.Printf("Found %d security alerts to backup (limited to 20)", len(securityAlerts))

		// Mask sensitive details
		for i := range securityAlerts {
			if securityAlerts[i].Details != nil {
				// Mask IP addresses in details
				if ip, exists := securityAlerts[i].Details["ip_address"]; exists {
					if ipStr, ok := ip.(string); ok && ipStr != "" {
						parts := strings.Split(ipStr, ".")
						if len(parts) == 4 {
							securityAlerts[i].Details["ip_address"] = parts[0] + "." + parts[1] + ".***.***"
						}
					}
				}
			}
		}

		backupData["data"].(map[string]interface{})["security_alerts"] = securityAlerts
	}

	// Note: IDOR attempts are now stored in Redis, not in database
	log.Printf("IDOR attempts are stored in Redis for real-time monitoring")

	// Write backup to file
	backupJSON, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %v", err)
	}

	if err := os.WriteFile(backupFile, backupJSON, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	log.Printf("Backup created successfully: %s", backupFile)
	return nil
}

// ListBackups returns a list of available backups
func (bs *BackupService) ListBackups() ([]string, error) {
	backupDir := "backups"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

// GetBackupInfo returns information about a specific backup
func (bs *BackupService) GetBackupInfo(backupName string) (map[string]interface{}, error) {
	backupFile := filepath.Join("backups", backupName)

	data, err := os.ReadFile(backupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %v", err)
	}

	var backupData map[string]interface{}
	if err := json.Unmarshal(data, &backupData); err != nil {
		return nil, fmt.Errorf("failed to parse backup file: %v", err)
	}

	return backupData, nil
}

// startBackupScheduler starts the automatic backup scheduler
func (bs *BackupService) startBackupScheduler() {
	if !DefaultBackupConfig.Enabled {
		log.Println("Automatic backups are disabled")
		return
	}

	log.Printf("Starting automatic backup scheduler (interval: %v, max backups: %d)",
		DefaultBackupConfig.BackupInterval, DefaultBackupConfig.MaxBackups)

	ticker := time.NewTicker(DefaultBackupConfig.BackupInterval)
	defer ticker.Stop()

	// Create initial backup
	bs.createScheduledBackup()

	for range ticker.C {
		bs.createScheduledBackup()
	}
}

// createScheduledBackup creates a scheduled backup and manages cleanup
func (bs *BackupService) createScheduledBackup() {
	log.Println("Creating scheduled backup...")

	if err := bs.BackupData(); err != nil {
		log.Printf("Failed to create scheduled backup: %v", err)
		return
	}

	// Cleanup old backups to maintain max count
	bs.cleanupOldBackups()
}

// cleanupOldBackups removes old backups to maintain maximum count
func (bs *BackupService) cleanupOldBackups() {
	backups, err := bs.ListBackups()
	if err != nil {
		log.Printf("Failed to list backups for cleanup: %v", err)
		return
	}

	if len(backups) <= DefaultBackupConfig.MaxBackups {
		return
	}

	// Sort backups by creation time (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i] < backups[j]
	})

	// Remove oldest backups
	backupsToRemove := len(backups) - DefaultBackupConfig.MaxBackups
	for i := 0; i < backupsToRemove; i++ {
		backupFile := filepath.Join("backups", backups[i])
		if err := os.Remove(backupFile); err != nil {
			log.Printf("Failed to remove old backup %s: %v", backups[i], err)
		} else {
			log.Printf("Removed old backup: %s", backups[i])
		}
	}
}

// CreateBackupEndpoint creates a backup via API endpoint (DEPRECATED)
// This endpoint is deprecated and will be removed in future versions
// Use automatic scheduled backups instead
func (bs *BackupService) CreateBackupEndpoint() error {
	log.Println("WARNING: Manual backup creation is deprecated. Use automatic scheduled backups.")
	return bs.BackupData()
}

// GetBackupStats returns statistics about backups
func (bs *BackupService) GetBackupStats() map[string]interface{} {
	backups, err := bs.ListBackups()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Calculate total size of backups
	var totalSize int64
	for _, backup := range backups {
		backupFile := filepath.Join("backups", backup)
		if info, err := os.Stat(backupFile); err == nil {
			totalSize += info.Size()
		}
	}

	return map[string]interface{}{
		"total_backups": len(backups),
		"total_size_mb": float64(totalSize) / (1024 * 1024),
		"backup_files":  backups,
	}
}
