package services

import (
	"fmt"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SessionService handles session management operations
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new session service
func NewSessionService() *SessionService {
	return &SessionService{
		db: config.GetDB(),
	}
}

// SessionInfo represents a session summary for API responses
type SessionInfo struct {
	ID           string    `json:"id"`
	DeviceName   string    `json:"device_name"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsActive     bool      `json:"is_active"`
	LastAccessed time.Time `json:"last_accessed"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsCurrent    bool      `json:"is_current"`
}

// GetUserSessions returns all active sessions for a user
func (s *SessionService) GetUserSessions(userID uuid.UUID, currentSessionID uuid.UUID) ([]SessionInfo, error) {
	var sessions []models.Session
	if err := s.db.Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Order("last_accessed DESC").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %v", err)
	}

	result := make([]SessionInfo, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, SessionInfo{
			ID:           sess.ID.String(),
			DeviceName:   sess.DeviceName,
			IPAddress:    sess.IPAddress,
			UserAgent:    sess.UserAgent,
			IsActive:     sess.IsActive,
			LastAccessed: sess.LastAccessed,
			CreatedAt:    sess.CreatedAt,
			ExpiresAt:    sess.ExpiresAt,
			IsCurrent:    sess.ID == currentSessionID,
		})
	}
	return result, nil
}

// RevokeSession revokes a specific session by ID
func (s *SessionService) RevokeSession(userID uuid.UUID, sessionID uuid.UUID) error {
	result := s.db.Model(&models.Session{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Updates(map[string]interface{}{
			"is_active": false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke session: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

// RevokeAllOtherSessions revokes all sessions except the current one
func (s *SessionService) RevokeAllOtherSessions(userID uuid.UUID, currentSessionID uuid.UUID) error {
	result := s.db.Model(&models.Session{}).
		Where("user_id = ? AND id != ?", userID, currentSessionID).
		Updates(map[string]interface{}{
			"is_active": false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke sessions: %v", result.Error)
	}
	return nil
}

// RevokeAllUserSessions revokes all sessions for a user (including current)
func (s *SessionService) RevokeAllUserSessions(userID uuid.UUID) error {
	result := s.db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_active": false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke all sessions: %v", result.Error)
	}
	return nil
}

// UpdateSessionAccess updates the last accessed time for a session
func (s *SessionService) UpdateSessionAccess(sessionID uuid.UUID) error {
	return s.db.Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("last_accessed", time.Now()).Error
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() (int64, error) {
	result := s.db.Where("expires_at < ? OR is_active = ?", time.Now(), false).Delete(&models.Session{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup sessions: %v", result.Error)
	}
	return result.RowsAffected, nil
}

// CreateSessionWithMetadata creates a new session with device and IP metadata
func (s *SessionService) CreateSessionWithMetadata(userID uuid.UUID, token, deviceName, ipAddress, userAgent string, expiresAt time.Time) (*models.Session, error) {
	session := models.Session{
		UserID:     userID,
		Token:      token,
		DeviceName: deviceName,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	return &session, nil
}

// GetSessionByToken retrieves a session by its token
func (s *SessionService) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Where("token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).
		First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired session")
	}
	return &session, nil
}

// GetSessionCount returns the number of active sessions for a user
func (s *SessionService) GetSessionCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count sessions: %v", err)
	}
	return count, nil
}
