package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"securewallet/internal/config"
	"securewallet/internal/models"
)

// LoginHistoryService handles login history operations
type LoginHistoryService struct{}

// NewLoginHistoryService creates a new login history service
func NewLoginHistoryService() *LoginHistoryService {
	return &LoginHistoryService{}
}

// RecordLoginAttempt records a login attempt
func (s *LoginHistoryService) RecordLoginAttempt(userID uuid.UUID, status string, r *http.Request) error {
	db := config.GetDB()

	ipAddress := s.getClientIP(r)
	userAgent := r.UserAgent()
	location := s.getLocationFromIP(ipAddress) // Simplified - in production use a geolocation service

	loginHistory := models.LoginHistory{
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    status,
		Location:  location,
	}

	return db.Create(&loginHistory).Error
}

// GetLoginHistory gets login history for a user
func (s *LoginHistoryService) GetLoginHistory(userID uuid.UUID, limit int) ([]models.LoginHistory, error) {
	db := config.GetDB()

	var history []models.LoginHistory
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error

	return history, err
}

// GetFailedLoginAttempts gets recent failed login attempts for a user
func (s *LoginHistoryService) GetFailedLoginAttempts(userID uuid.UUID, since time.Time) (int64, error) {
	db := config.GetDB()

	var count int64
	err := db.Model(&models.LoginHistory{}).
		Where("user_id = ? AND status = ? AND created_at > ?", userID, "failed", since).
		Count(&count).Error

	return count, err
}

// getClientIP gets the real client IP address
func (s *LoginHistoryService) getClientIP(r *http.Request) string {
	// Check for forwarded headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
			return strings.TrimSpace(ip[:commaIndex])
		}
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}

	// Fallback to remote address
	return r.RemoteAddr
}

// getLocationFromIP gets location from IP (simplified)
func (s *LoginHistoryService) getLocationFromIP(ip string) string {
	// In a real application, you would use a geolocation service
	// For now, return a placeholder
	if strings.HasPrefix(ip, "127.0.0.1") || strings.HasPrefix(ip, "::1") {
		return "Local"
	}
	return "Unknown"
}
