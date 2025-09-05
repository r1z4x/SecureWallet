package services

import (
	"log"
	"securewallet/internal/config"
	"time"

	"gorm.io/gorm"
)

// CommentService handles comment operations
type CommentService struct {
	db *gorm.DB
}

// CommentConfig holds comment configuration
type CommentConfig struct {
	AutoApproveDelay time.Duration // How long to wait before auto-approving comments
	Enabled          bool          // Whether auto-approval is enabled
}

// Default comment configuration
var DefaultCommentConfig = CommentConfig{
	AutoApproveDelay: 1 * time.Minute, // Auto-approve after 1 minute (for testing)
	Enabled:          true,              // Enable auto-approval
}

// NewCommentService creates a new comment service
func NewCommentService() *CommentService {
	cs := &CommentService{
		db: config.GetDB(),
	}

	// Start automatic comment approval scheduler
	go cs.startCommentApprovalScheduler()

	return cs
}

// startCommentApprovalScheduler starts the automatic comment approval scheduler
func (cs *CommentService) startCommentApprovalScheduler() {
	if !DefaultCommentConfig.Enabled {
		log.Println("Automatic comment approval is disabled")
		return
	}

	log.Printf("Starting automatic comment approval scheduler (delay: %v)",
		DefaultCommentConfig.AutoApproveDelay)

	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for range ticker.C {
		cs.approvePendingComments()
	}
}

// approvePendingComments approves comments that are older than the configured delay
func (cs *CommentService) approvePendingComments() {
	cutoffTime := time.Now().Add(-DefaultCommentConfig.AutoApproveDelay)
	
	// Update comments that are pending and older than cutoff time
	result := cs.db.Table("blog_comments").
		Where("status = ? AND created_at <= ?", "pending", cutoffTime).
		Update("status", "approved")

	if result.Error != nil {
		log.Printf("Failed to approve pending comments: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Auto-approved %d pending comments", result.RowsAffected)
	}
}

// GetCommentStats returns statistics about comments
func (cs *CommentService) GetCommentStats() map[string]interface{} {
	var totalComments int64
	var pendingComments int64
	var approvedComments int64

	cs.db.Table("blog_comments").Count(&totalComments)
	cs.db.Table("blog_comments").Where("status = ?", "pending").Count(&pendingComments)
	cs.db.Table("blog_comments").Where("status = ?", "approved").Count(&approvedComments)

	return map[string]interface{}{
		"total_comments":    totalComments,
		"pending_comments":  pendingComments,
		"approved_comments": approvedComments,
		"auto_approval": map[string]interface{}{
			"enabled": DefaultCommentConfig.Enabled,
			"delay":   DefaultCommentConfig.AutoApproveDelay.String(),
		},
	}
}
