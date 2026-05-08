package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestSanitizeDetails_RemovesSensitiveData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"password in details", "User login with password secret123", "[REDACTED]"},
		{"token in details", "Bearer token abc123 used", "[REDACTED]"},
		{"secret in details", "Secret key rotated", "[REDACTED]"},
		{"safe details", "User transferred 50.00 USD", "User transferred 50.00 USD"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeDetails(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeDetails(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSanitizeDetails_TruncatesLongStrings(t *testing.T) {
	long := ""
	for i := 0; i < 1500; i++ {
		long += "a"
	}

	got := sanitizeDetails(long)
	if len(got) > 1003 {
		t.Errorf("sanitizeDetails should truncate to ~1000 chars, got %d", len(got))
	}
}

func TestSanitizeMetadata_RemovesSensitiveKeys(t *testing.T) {
	metadata := map[string]interface{}{
		"amount":            100.0,
		"password":          "secret123",
		"token":             "abc123",
		"access_token":      "jwt-token-here",
		"refresh_token":     "refresh-token-here",
		"two_factor_secret": "totp-secret",
		"recovery_code":     "X7K9M2P4",
		"recipient":         "user@example.com",
		"currency":          "USD",
	}

	sanitized := sanitizeMetadata(metadata)

	sensitiveKeys := []string{"password", "token", "access_token", "refresh_token", "two_factor_secret", "recovery_code"}
	for _, key := range sensitiveKeys {
		if sanitized[key] != "[REDACTED]" {
			t.Errorf("sanitizeMetadata should redact %q, got %v", key, sanitized[key])
		}
	}

	if sanitized["amount"] != 100.0 {
		t.Errorf("sanitizeMetadata should preserve amount, got %v", sanitized["amount"])
	}
	if sanitized["recipient"] != "user@example.com" {
		t.Errorf("sanitizeMetadata should preserve recipient, got %v", sanitized["recipient"])
	}
	if sanitized["currency"] != "USD" {
		t.Errorf("sanitizeMetadata should preserve currency, got %v", sanitized["currency"])
	}
}

func TestSanitizeMetadata_NilInput(t *testing.T) {
	got := sanitizeMetadata(nil)
	if got != nil {
		t.Errorf("sanitizeMetadata(nil) should return nil, got %v", got)
	}
}

func TestSanitizeMetadata_EmptyInput(t *testing.T) {
	got := sanitizeMetadata(map[string]interface{}{})
	if len(got) != 0 {
		t.Errorf("sanitizeMetadata(empty) should return empty map, got %v", got)
	}
}

func TestExtractClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")

	got := extractClientIP(req)
	if got != "203.0.113.195" {
		t.Errorf("extractClientIP should return first IP from X-Forwarded-For, got %q", got)
	}
}

func TestExtractClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.100")

	got := extractClientIP(req)
	if got != "192.168.1.100" {
		t.Errorf("extractClientIP should return X-Real-IP, got %q", got)
	}
}

func TestExtractClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	got := extractClientIP(req)
	if got != "192.0.2.1:1234" {
		t.Errorf("extractClientIP should return RemoteAddr, got %q", got)
	}
}

func TestExtractClientIP_XForwardedFor_SingleIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")

	got := extractClientIP(req)
	if got != "10.0.0.1" {
		t.Errorf("extractClientIP should return single X-Forwarded-For IP, got %q", got)
	}
}

func TestAuditLogger_Chaining(t *testing.T) {
	logger := NewAuditLogger()
	if logger == nil {
		t.Fatal("NewAuditLogger should return non-nil logger")
	}

	correlationID := uuid.New().String()
	logger2 := logger.WithCorrelationID(correlationID)
	if logger2.correlationID != correlationID {
		t.Errorf("WithCorrelationID should set correlation ID, got %q", logger2.correlationID)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.Header.Set("User-Agent", "test-agent/1.0")
	logger3 := logger2.WithRequest(req)
	if logger3.ipAddress != "10.0.0.1" {
		t.Errorf("WithRequest should set IP address, got %q", logger3.ipAddress)
	}
	if logger3.userAgent != "test-agent/1.0" {
		t.Errorf("WithRequest should set user agent, got %q", logger3.userAgent)
	}
}

func TestAuditResult_Values(t *testing.T) {
	results := []AuditResult{
		AuditResultSuccess,
		AuditResultFailure,
		AuditResultDenied,
		AuditResultUnknown,
	}

	for _, r := range results {
		if string(r) == "" {
			t.Errorf("AuditResult %v should have non-empty string value", r)
		}
	}
}

func TestSanitizeDetails_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PASSWORD in details", "[REDACTED]"},
		{"Token value here", "[REDACTED]"},
		{"SECRET key", "[REDACTED]"},
		{"Authorization header", "[REDACTED]"},
		{"Cookie data", "[REDACTED]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeDetails(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeDetails(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSanitizeMetadata_CaseInsensitiveKeys(t *testing.T) {
	metadata := map[string]interface{}{
		"PASSWORD": "secret",
		"Token":    "abc123",
		"Amount":   100.0,
	}

	sanitized := sanitizeMetadata(metadata)

	if sanitized["PASSWORD"] != "[REDACTED]" {
		t.Errorf("sanitizeMetadata should redact PASSWORD (uppercase), got %v", sanitized["PASSWORD"])
	}
	if sanitized["Token"] != "[REDACTED]" {
		t.Errorf("sanitizeMetadata should redact Token (mixed case), got %v", sanitized["Token"])
	}
	if sanitized["Amount"] != 100.0 {
		t.Errorf("sanitizeMetadata should preserve Amount, got %v", sanitized["Amount"])
	}
}
