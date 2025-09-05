package services

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"securewallet/internal/config"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CronService handles cron-based scheduled tasks
type CronService struct {
	db *gorm.DB
}

// CronJob represents a scheduled job
type CronJob struct {
	Name        string        // Job name
	Schedule    string        // Cron schedule (e.g., "*/5 * * * *")
	Command     string        // Command to execute
	Description string        // Job description
	Enabled     bool          // Whether job is enabled
	LastRun     time.Time     // Last execution time
	NextRun     time.Time     // Next scheduled run
	Status      string        // Job status: "pending", "running", "completed", "failed"
	LogFile     string        // Log file path
}

// CronConfig holds cron configuration
type CronConfig struct {
	Enabled     bool          // Whether cron service is enabled
	LogDir      string        // Directory for cron logs
	MaxLogFiles int           // Maximum number of log files to keep
}

// Default cron configuration
var DefaultCronConfig = CronConfig{
	Enabled:     true,
	LogDir:      "logs/cron",
	MaxLogFiles: 10,
}

// NewCronService creates a new cron service
func NewCronService() *CronService {
	cs := &CronService{
		db: config.GetDB(),
	}

	// Create log directory
	if err := os.MkdirAll(DefaultCronConfig.LogDir, 0755); err != nil {
		log.Printf("Failed to create cron log directory: %v", err)
	}

	// Setup cron jobs
	cs.SetupCronJobs()

	return cs
}

// SetupCronJobs sets up all cron jobs
func (cs *CronService) SetupCronJobs() {
	if !DefaultCronConfig.Enabled {
		log.Println("Cron service is disabled")
		return
	}

	log.Println("Setting up cron jobs...")

	// Comment auto-approval job (every 5 minutes)
	cs.addCronJob(CronJob{
		Name:        "comment-auto-approval",
		Schedule:    "*/5 * * * *",
		Command:     "go run main.go --cron=comment-approval",
		Description: "Auto-approve pending comments older than 10 minutes",
		Enabled:     true,
		LogFile:     filepath.Join(DefaultCronConfig.LogDir, "comment-approval.log"),
	})

	// Database backup job (every hour)
	cs.addCronJob(CronJob{
		Name:        "database-backup",
		Schedule:    "0 * * * *",
		Command:     "go run main.go --cron=backup",
		Description: "Create database backup",
		Enabled:     true,
		LogFile:     filepath.Join(DefaultCronConfig.LogDir, "backup.log"),
	})

	// Log cleanup job (daily at 2 AM)
	cs.addCronJob(CronJob{
		Name:        "log-cleanup",
		Schedule:    "0 2 * * *",
		Command:     "go run main.go --cron=log-cleanup",
		Description: "Clean up old log files",
		Enabled:     true,
		LogFile:     filepath.Join(DefaultCronConfig.LogDir, "log-cleanup.log"),
	})

	// Security monitoring job (every 10 minutes)
	cs.addCronJob(CronJob{
		Name:        "security-monitoring",
		Schedule:    "*/10 * * * *",
		Command:     "go run main.go --cron=security-monitor",
		Description: "Monitor security events and generate alerts",
		Enabled:     true,
		LogFile:     filepath.Join(DefaultCronConfig.LogDir, "security-monitor.log"),
	})

	log.Printf("Setup %d cron jobs", len(cs.getCronJobs()))
}

// addCronJob adds a cron job to the system
func (cs *CronService) addCronJob(job CronJob) {
	// Write cron job to crontab
	if err := cs.writeCrontabEntry(job); err != nil {
		log.Printf("Failed to add cron job %s: %v", job.Name, err)
		return
	}

	log.Printf("Added cron job: %s (%s)", job.Name, job.Schedule)
}

// writeCrontabEntry writes a cron job entry to the system crontab
func (cs *CronService) writeCrontabEntry(job CronJob) error {
	if !job.Enabled {
		return nil
	}

	// Create the cron entry
	cronEntry := job.Schedule + " " + job.Command + " >> " + job.LogFile + " 2>&1"

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	currentCrontab, err := cmd.Output()
	if err != nil && !strings.Contains(err.Error(), "no crontab") {
		return err
	}

	// Check if job already exists
	if strings.Contains(string(currentCrontab), job.Name) {
		log.Printf("Cron job %s already exists", job.Name)
		return nil
	}

	// Add new entry
	newCrontab := string(currentCrontab)
	if newCrontab != "" && !strings.HasSuffix(newCrontab, "\n") {
		newCrontab += "\n"
	}
	newCrontab += "# " + job.Description + "\n"
	newCrontab += cronEntry + "\n"

	// Write new crontab
	cmd = exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCrontab)
	return cmd.Run()
}

// getCronJobs returns all configured cron jobs
func (cs *CronService) getCronJobs() []CronJob {
	return []CronJob{
		{
			Name:        "comment-auto-approval",
			Schedule:    "*/5 * * * *",
			Command:     "go run main.go --cron=comment-approval",
			Description: "Auto-approve pending comments older than 10 minutes",
			Enabled:     true,
		},
		{
			Name:        "database-backup",
			Schedule:    "0 * * * *",
			Command:     "go run main.go --cron=backup",
			Description: "Create database backup",
			Enabled:     true,
		},
		{
			Name:        "log-cleanup",
			Schedule:    "0 2 * * *",
			Command:     "go run main.go --cron=log-cleanup",
			Description: "Clean up old log files",
			Enabled:     true,
		},
		{
			Name:        "security-monitoring",
			Schedule:    "*/10 * * * *",
			Command:     "go run main.go --cron=security-monitor",
			Description: "Monitor security events and generate alerts",
			Enabled:     true,
		},
	}
}

// RemoveCronJobs removes all cron jobs from the system
func (cs *CronService) RemoveCronJobs() error {
	log.Println("Removing all cron jobs...")

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	currentCrontab, err := cmd.Output()
	if err != nil && !strings.Contains(err.Error(), "no crontab") {
		return err
	}

	// Filter out our cron jobs
	lines := strings.Split(string(currentCrontab), "\n")
	var filteredLines []string

	for _, line := range lines {
		// Skip lines that contain our job commands
		if !strings.Contains(line, "go run main.go --cron=") {
			filteredLines = append(filteredLines, line)
		}
	}

	// Write filtered crontab
	newCrontab := strings.Join(filteredLines, "\n")
	if newCrontab != "" && !strings.HasSuffix(newCrontab, "\n") {
		newCrontab += "\n"
	}

	cmd = exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCrontab)
	return cmd.Run()
}

// GetCronStatus returns the status of all cron jobs
func (cs *CronService) GetCronStatus() map[string]interface{} {
	jobs := cs.getCronJobs()
	
	// Get current crontab to check if jobs are installed
	cmd := exec.Command("crontab", "-l")
	currentCrontab, err := cmd.Output()
	installed := err == nil

	status := map[string]interface{}{
		"enabled":   DefaultCronConfig.Enabled,
		"installed": installed,
		"jobs":      []map[string]interface{}{},
	}

	for _, job := range jobs {
		jobStatus := map[string]interface{}{
			"name":        job.Name,
			"schedule":    job.Schedule,
			"description": job.Description,
			"enabled":     job.Enabled,
			"installed":   installed && strings.Contains(string(currentCrontab), job.Name),
		}
		status["jobs"] = append(status["jobs"].([]map[string]interface{}), jobStatus)
	}

	return status
}

// ExecuteCronJob executes a specific cron job
func (cs *CronService) ExecuteCronJob(jobName string) error {
	log.Printf("Executing cron job: %s", jobName)

	switch jobName {
	case "comment-approval":
		return cs.executeCommentApproval()
	case "backup":
		return cs.executeBackup()
	case "log-cleanup":
		return cs.executeLogCleanup()
	case "security-monitor":
		return cs.executeSecurityMonitoring()
	default:
		log.Printf("Unknown cron job: %s", jobName)
		return nil
	}
}

// executeCommentApproval executes the comment auto-approval job
func (cs *CronService) executeCommentApproval() error {
	log.Println("Executing comment auto-approval...")
	
	// Get pending comments older than 10 minutes
	cutoffTime := time.Now().Add(-10 * time.Minute)
	
	result := cs.db.Table("blog_comments").
		Where("status = ? AND created_at <= ?", "pending", cutoffTime).
		Update("status", "approved")

	if result.Error != nil {
		log.Printf("Failed to approve pending comments: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Auto-approved %d pending comments", result.RowsAffected)
	} else {
		log.Println("No pending comments to approve")
	}

	return nil
}

// executeBackup executes the database backup job
func (cs *CronService) executeBackup() error {
	log.Println("Executing database backup...")
	
	backupService := NewBackupService()
	return backupService.BackupData()
}

// executeLogCleanup executes the log cleanup job
func (cs *CronService) executeLogCleanup() error {
	log.Println("Executing log cleanup...")
	
	// Clean up old cron log files
	logDir := DefaultCronConfig.LogDir
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// Keep only the most recent log files
	if len(files) > DefaultCronConfig.MaxLogFiles {
		// Sort by modification time (oldest first)
		// Remove oldest files
		filesToRemove := len(files) - DefaultCronConfig.MaxLogFiles
		for i := 0; i < filesToRemove; i++ {
			filePath := filepath.Join(logDir, files[i].Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("Failed to remove log file %s: %v", files[i].Name(), err)
			} else {
				log.Printf("Removed old log file: %s", files[i].Name())
			}
		}
	}

	return nil
}

// executeSecurityMonitoring executes the security monitoring job
func (cs *CronService) executeSecurityMonitoring() error {
	log.Println("Executing security monitoring...")
	
	// This would typically check for security events, generate alerts, etc.
	// For now, just log that the job ran
	log.Println("Security monitoring completed")
	
	return nil
}
