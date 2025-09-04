package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"securewallet/internal/config"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SecurityEvent represents a general security event
type SecurityEvent struct {
	ID            string                 `json:"id"`
	Category      string                 `json:"category"` // IDOR, BRUTE_FORCE, SQL_INJECTION, etc.
	UserID        string                 `json:"user_id"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	Details       map[string]interface{} `json:"details"`
	Timestamp     time.Time              `json:"timestamp"`
	Severity      string                 `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	Blocked       bool                   `json:"blocked"`
	Resource      string                 `json:"resource,omitempty"`
	ResourceOwner string                 `json:"resource_owner,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`   // For session-based detection
	PatternHash   string                 `json:"pattern_hash,omitempty"` // For pattern detection
}

// SecurityDetector handles security vulnerability detection
type SecurityDetector struct {
	db    *gorm.DB
	redis *redis.Client

	// Detection thresholds
	eventThreshold   int           // Max events before alert
	sessionThreshold int           // Max events in a session before alert
	patternThreshold int           // Max similar pattern events before alert
	rateLimitWindow  time.Duration // Time window for rate limiting
	sessionWindow    time.Duration // Time window for session detection
	alertCooldown    time.Duration // Cooldown between alerts

	// Redis keys
	eventKeyPrefix     string // "security:events:"
	alertKeyPrefix     string // "security:alerts:"
	userEventKeyPrefix string // "security:user:events:"
	sessionKeyPrefix   string // "security:sessions:"
	patternKeyPrefix   string // "security:patterns:"

	mu sync.RWMutex
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`     // IDOR, RATE_LIMIT, SUSPICIOUS_ACTIVITY
	Severity   string                 `json:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	UserID     string                 `json:"user_id"`
	IPAddress  string                 `json:"ip_address"`
	Details    map[string]interface{} `json:"details"`
	Timestamp  time.Time              `json:"timestamp"`
	Status     string                 `json:"status"` // OPEN, INVESTIGATING, RESOLVED, FALSE_POSITIVE
	ResolvedBy string                 `json:"resolved_by"`
	ResolvedAt *time.Time             `json:"resolved_at"`
}

// NewSecurityDetector creates a new security detector
func NewSecurityDetector() *SecurityDetector {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6380" // Default dev port
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "CHANGE_THIS_REDIS_PASSWORD"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	return &SecurityDetector{
		db:    config.GetDB(),
		redis: redisClient,

		// Configuration
		eventThreshold:   10,               // Alert after 10 events
		sessionThreshold: 3,                // Alert after 3 events in a session
		patternThreshold: 2,                // Alert after 2 similar pattern events
		rateLimitWindow:  5 * time.Minute,  // 5 minute window
		sessionWindow:    2 * time.Minute,  // 2 minute session window
		alertCooldown:    10 * time.Minute, // 10 minute cooldown between alerts

		// Redis keys
		eventKeyPrefix:     "security:events:",
		alertKeyPrefix:     "security:alerts:",
		userEventKeyPrefix: "security:user:events:",
		sessionKeyPrefix:   "security:sessions:",
		patternKeyPrefix:   "security:patterns:",
	}
}

// DetectSecurityEvent detects security events and manages alerts with session and pattern detection
func (sd *SecurityDetector) DetectSecurityEvent(event *SecurityEvent) (*SecurityAlert, error) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	ctx := context.Background()

	// Generate session ID if not provided (based on user + time window)
	if event.SessionID == "" {
		event.SessionID = sd.generateSessionID(event.UserID, event.IPAddress)
	}

	// Generate pattern hash for pattern detection
	if event.PatternHash == "" {
		event.PatternHash = sd.generatePatternHash(event)
	}

	// Store event in Redis
	eventKey := fmt.Sprintf("%s%s:%s", sd.eventKeyPrefix, event.Category, event.UserID)
	eventData, _ := json.Marshal(event)

	// Add to Redis sorted set with timestamp as score
	if err := sd.redis.ZAdd(ctx, eventKey, redis.Z{
		Score:  float64(event.Timestamp.Unix()),
		Member: eventData,
	}).Err(); err != nil {
		log.Printf("Failed to store security event in Redis: %v", err)
	}

	// 1. SESSION-BASED DETECTION
	sessionKey := fmt.Sprintf("%s%s:%s", sd.sessionKeyPrefix, event.Category, event.SessionID)
	sessionCount, _ := sd.redis.ZCount(ctx, sessionKey, "-inf", "+inf").Result()

	// Add to session
	if err := sd.redis.ZAdd(ctx, sessionKey, redis.Z{
		Score:  float64(event.Timestamp.Unix()),
		Member: eventData,
	}).Err(); err != nil {
		log.Printf("Failed to store session event: %v", err)
	}

	// 2. PATTERN-BASED DETECTION
	patternKey := fmt.Sprintf("%s%s:%s", sd.patternKeyPrefix, event.Category, event.PatternHash)
	patternCount, _ := sd.redis.ZCount(ctx, patternKey, "-inf", "+inf").Result()

	// Add to pattern
	if err := sd.redis.ZAdd(ctx, patternKey, redis.Z{
		Score:  float64(event.Timestamp.Unix()),
		Member: eventData,
	}).Err(); err != nil {
		log.Printf("Failed to store pattern event: %v", err)
	}

	// 3. TRADITIONAL RATE LIMITING
	windowStart := time.Now().Add(-sd.rateLimitWindow)
	windowStartUnix := float64(windowStart.Unix())
	rateLimitCount, err := sd.redis.ZCount(ctx, eventKey, strconv.FormatFloat(windowStartUnix, 'f', -1, 64), "+inf").Result()
	if err != nil {
		log.Printf("Failed to count events in Redis: %v", err)
		rateLimitCount = 0
	}

	// Log the security event with all detection methods
	log.Printf("🚨 %s EVENT: User %s from IP %s - Session: %d, Pattern: %d, Rate: %d",
		event.Category, event.UserID, event.IPAddress, sessionCount, patternCount, rateLimitCount)

	// Check if we should create an alert (any threshold reached)
	shouldAlert := false
	var alert *SecurityAlert
	alertReason := ""

	if int(sessionCount) >= sd.sessionThreshold {
		shouldAlert = true
		alertReason = fmt.Sprintf("Session threshold reached: %d events", sessionCount)
	} else if int(patternCount) >= sd.patternThreshold {
		shouldAlert = true
		alertReason = fmt.Sprintf("Pattern threshold reached: %d similar events", patternCount)
	} else if int(rateLimitCount) >= sd.eventThreshold {
		shouldAlert = true
		alertReason = fmt.Sprintf("Rate limit threshold reached: %d events", rateLimitCount)
	}

	if shouldAlert {
		// Check cooldown in Redis
		alertKey := fmt.Sprintf("%s%s:%s", sd.alertKeyPrefix, event.Category, event.UserID)
		lastAlertStr, err := sd.redis.Get(ctx, alertKey).Result()

		if err == redis.Nil || lastAlertStr == "" {
			// No previous alert or cooldown expired
		} else {
			// Check if cooldown has passed
			if lastAlert, err := time.Parse(time.RFC3339, lastAlertStr); err == nil {
				if time.Since(lastAlert) <= sd.alertCooldown {
					shouldAlert = false // Still in cooldown
				}
			}
		}
	}

	if shouldAlert {
		// Create security alert
		alert = &SecurityAlert{
			ID:        fmt.Sprintf("%s_%s_%d", event.Category, event.UserID, time.Now().Unix()),
			Type:      event.Category,
			Severity:  event.Severity,
			UserID:    event.UserID,
			IPAddress: event.IPAddress,
			Details: map[string]interface{}{
				"alert_reason":   alertReason,
				"session_count":  sessionCount,
				"pattern_count":  patternCount,
				"rate_count":     rateLimitCount,
				"session_id":     event.SessionID,
				"pattern_hash":   event.PatternHash,
				"resource":       event.Resource,
				"resource_owner": event.ResourceOwner,
				"blocked":        event.Blocked,
			},
			Timestamp: time.Now(),
			Status:    "OPEN",
		}

		// Save alert to database
		if err := sd.db.Create(&alert).Error; err != nil {
			log.Printf("Failed to save security alert: %v", err)
			return nil, err
		}

		// Store alert timestamp in Redis for cooldown
		alertKey := fmt.Sprintf("%s%s:%s", sd.alertKeyPrefix, event.Category, event.UserID)
		if err := sd.redis.Set(ctx, alertKey, time.Now().Format(time.RFC3339), sd.alertCooldown).Err(); err != nil {
			log.Printf("Failed to store alert cooldown: %v", err)
		}

		// Send real-time alert
		sd.sendRealTimeAlert(alert)
	}

	return alert, nil
}

// generateSessionID generates a session ID based on user and time window
func (sd *SecurityDetector) generateSessionID(userID, ipAddress string) string {
	// Create session based on user + IP + time window (2 minutes)
	windowStart := time.Now().Truncate(sd.sessionWindow)
	return fmt.Sprintf("%s_%s_%d", userID, ipAddress, windowStart.Unix())
}

// generatePatternHash generates a hash for pattern detection
func (sd *SecurityDetector) generatePatternHash(event *SecurityEvent) string {
	// Create pattern hash based on event characteristics
	pattern := fmt.Sprintf("%s_%s_%s_%s",
		event.Category,
		event.Resource,
		event.ResourceOwner,
		event.IPAddress)

	// Simple hash (in production, use proper hashing)
	return fmt.Sprintf("%x", len(pattern))
}

// DetectIDOR detects IDOR attempts using the general security event system
func (sd *SecurityDetector) DetectIDOR(userID, attemptedResource, resourceOwner, ipAddress, userAgent string) (*SecurityAlert, error) {
	event := &SecurityEvent{
		ID:            fmt.Sprintf("idor_%s_%d", userID, time.Now().Unix()),
		Category:      "IDOR",
		UserID:        userID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Details:       map[string]interface{}{},
		Timestamp:     time.Now(),
		Severity:      "HIGH",
		Blocked:       false,
		Resource:      attemptedResource,
		ResourceOwner: resourceOwner,
	}

	return sd.DetectSecurityEvent(event)
}

// sendRealTimeAlert sends real-time security alerts
func (sd *SecurityDetector) sendRealTimeAlert(alert *SecurityAlert) {
	// TODO: Implement real-time alerting
	// - Email notifications
	// - Slack/Discord webhooks
	// - SMS alerts
	// Security dashboard updates

	log.Printf("📡 Real-time alert sent for %s: %s", alert.Type, alert.ID)

	// Log to file for real-time monitoring
	alertLog := fmt.Sprintf("[%s] 🚨 %s ALERT: %s | User: %s | IP: %s | Severity: %s | Details: %v\n",
		time.Now().Format("2006-01-02 15:04:05"),
		alert.Type,
		alert.ID,
		alert.UserID,
		alert.IPAddress,
		alert.Severity,
		alert.Details)

	// Write to security alerts log file
	if err := sd.writeToAlertLog(alertLog); err != nil {
		log.Printf("Failed to write to alert log: %v", err)
	}
}

// writeToAlertLog writes alert messages to a log file
func (sd *SecurityDetector) writeToAlertLog(message string) error {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Open or create the security alerts log file
	logFile := filepath.Join(logDir, "security_alerts.log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open alert log file: %v", err)
	}
	defer file.Close()

	// Write the alert message
	if _, err := file.WriteString(message); err != nil {
		return fmt.Errorf("failed to write to alert log: %v", err)
	}

	return nil
}

// GetIDORStats returns IDOR detection statistics from Redis
func (sd *SecurityDetector) GetIDORStats() map[string]interface{} {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	ctx := context.Background()

	// Get total IDOR events from Redis
	var totalAttempts int64
	pattern := fmt.Sprintf("%sIDOR:*", sd.eventKeyPrefix)
	keys, err := sd.redis.Keys(ctx, pattern).Result()
	if err == nil {
		for _, key := range keys {
			count, _ := sd.redis.ZCard(ctx, key).Result()
			totalAttempts += count
		}
	}

	// Get alerts from database
	var totalAlerts int64
	sd.db.Model(&SecurityAlert{}).Where("type = ?", "IDOR").Count(&totalAlerts)

	var openAlerts int64
	sd.db.Model(&SecurityAlert{}).Where("type = ? AND status = ?", "IDOR", "OPEN").Count(&openAlerts)

	// Get active users from Redis
	activeUsers := len(keys)

	return map[string]interface{}{
		"total_attempts":    totalAttempts,
		"total_alerts":      totalAlerts,
		"open_alerts":       openAlerts,
		"threshold":         sd.eventThreshold,
		"rate_limit_window": sd.rateLimitWindow.String(),
		"alert_cooldown":    sd.alertCooldown.String(),
		"active_users":      activeUsers,
	}
}

// GetSecurityAlerts returns security alerts with filtering
func (sd *SecurityDetector) GetSecurityAlerts(alertType, status string, limit int) ([]SecurityAlert, error) {
	var alerts []SecurityAlert
	query := sd.db.Model(&SecurityAlert{})

	if alertType != "" {
		query = query.Where("type = ?", alertType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("timestamp DESC").Find(&alerts).Error
	return alerts, err
}

// UpdateAlertStatus updates the status of a security alert
func (sd *SecurityDetector) UpdateAlertStatus(alertID, status, resolvedBy string) error {
	updateData := map[string]interface{}{
		"status":      status,
		"resolved_by": resolvedBy,
		"resolved_at": time.Now(),
	}

	return sd.db.Model(&SecurityAlert{}).Where("id = ?", alertID).Updates(updateData).Error
}

// ResetUserAttempts resets attempt counters for a user (admin function)
func (sd *SecurityDetector) ResetUserAttempts(userID string) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	ctx := context.Background()

	// Remove all security events for this user from Redis
	pattern := fmt.Sprintf("%s*:%s", sd.eventKeyPrefix, userID)
	keys, err := sd.redis.Keys(ctx, pattern).Result()
	if err == nil {
		for _, key := range keys {
			sd.redis.Del(ctx, key)
		}
	}

	// Remove alert cooldown keys
	alertPattern := fmt.Sprintf("%s*:%s", sd.alertKeyPrefix, userID)
	alertKeys, err := sd.redis.Keys(ctx, alertPattern).Result()
	if err == nil {
		for _, key := range alertKeys {
			sd.redis.Del(ctx, key)
		}
	}

	log.Printf("Reset security attempt counters for user: %s", userID)
}

// CleanupOldData cleans up old security data
func (sd *SecurityDetector) CleanupOldData() {
	// Clean up old IDOR attempts (no database storage needed)
	log.Printf("IDOR attempts cleanup skipped - no database storage")

	// Clean up old alerts (older than 90 days)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)

	if err := sd.db.Where("timestamp < ?", ninetyDaysAgo).Delete(&SecurityAlert{}).Error; err != nil {
		log.Printf("Failed to cleanup old security alerts: %v", err)
	}

	log.Println("Security data cleanup completed")
}
