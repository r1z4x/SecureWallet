package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditResult string

const (
	AuditResultSuccess AuditResult = "success"
	AuditResultFailure AuditResult = "failure"
	AuditResultDenied  AuditResult = "denied"
	AuditResultUnknown AuditResult = "unknown"
)

type AuditLogger struct {
	db            *gorm.DB
	correlationID string
	ipAddress     string
	userAgent     string
}

func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		db: config.GetDB(),
	}
}

func (l *AuditLogger) WithCorrelationID(id string) *AuditLogger {
	l.correlationID = id
	return l
}

func (l *AuditLogger) WithRequest(r *http.Request) *AuditLogger {
	if r != nil {
		l.ipAddress = extractClientIP(r)
		l.userAgent = r.UserAgent()
	}
	return l
}

func (l *AuditLogger) WithGinContext(c *gin.Context) *AuditLogger {
	if c == nil {
		return l
	}
	if c.Request != nil {
		l.ipAddress = extractClientIP(c.Request)
		l.userAgent = c.Request.UserAgent()
	}
	if corrID, exists := c.Get("correlation_id"); exists {
		if id, ok := corrID.(string); ok {
			l.correlationID = id
		}
	}
	return l
}

func (l *AuditLogger) Log(userID uuid.UUID, action, resource, details string, result AuditResult) {
	entry := models.AuditLog{
		ID:            uuid.New(),
		UserID:        userID,
		Action:        action,
		Resource:      resource,
		Details:       sanitizeDetails(details),
		Result:        string(result),
		CorrelationID: l.correlationID,
		IPAddress:     l.ipAddress,
		UserAgent:     l.userAgent,
	}

	if err := l.db.Create(&entry).Error; err != nil {
		log.Printf("WARNING: failed to create audit log entry: %v", err)
	}
}

func (l *AuditLogger) LogJSON(userID uuid.UUID, action, resource string, metadata map[string]interface{}, result AuditResult) {
	sanitized := sanitizeMetadata(metadata)
	detailsBytes, err := json.Marshal(sanitized)
	details := ""
	if err == nil {
		details = string(detailsBytes)
	}

	l.Log(userID, action, resource, details, result)
}

func sanitizeDetails(details string) string {
	sanitized := details
	sensitivePatterns := []string{
		"password",
		"token",
		"secret",
		"authorization",
		"cookie",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(strings.ToLower(sanitized), pattern) {
			sanitized = "[REDACTED]"
		}
	}
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000] + "..."
	}
	return sanitized
}

func sanitizeMetadata(metadata map[string]interface{}) map[string]interface{} {
	if metadata == nil {
		return nil
	}

	sensitiveKeys := map[string]bool{
		"password":          true,
		"password_hash":     true,
		"token":             true,
		"access_token":      true,
		"refresh_token":     true,
		"secret":            true,
		"two_factor_secret": true,
		"recovery_code":     true,
		"authorization":     true,
		"cookie":            true,
		"jwt_secret":        true,
	}

	sanitized := make(map[string]interface{})
	for k, v := range metadata {
		if sensitiveKeys[strings.ToLower(k)] {
			sanitized[k] = "[REDACTED]"
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

func extractClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if idx := strings.Index(ip, ","); idx != -1 {
			return strings.TrimSpace(ip[:idx])
		}
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
