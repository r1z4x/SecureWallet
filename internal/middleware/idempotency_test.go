package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestIdempotencyMiddleware_MissingKey_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.IdempotencyMiddleware("transfer"))
	r.POST("/transfer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/transfer", nil)
	c.Request = req
	c.Set("user", user)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing idempotency key, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := parseJSON(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := resp["error"]; !ok {
		t.Error("Expected error field in missing idempotency key response")
	}
}

func TestIdempotencyMiddleware_InvalidUUIDFormat_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.IdempotencyMiddleware("transfer"))
	r.POST("/transfer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/transfer", nil)
	req.Header.Set("Idempotency-Key", "not-a-valid-uuid")
	c.Request = req
	c.Set("user", user)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid UUID format, got %d", w.Code)
	}
}

func TestIdempotencyMiddleware_NoUserInContext_RejectedBeforeDatabaseUse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validUUID := uuid.New().String()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.IdempotencyMiddleware("transfer"))
	r.POST("/transfer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/transfer", nil)
	req.Header.Set("Idempotency-Key", validUUID)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500 when user is missing from context, got %d", w.Code)
	}
}

func TestIdempotencyKey_UUIDFormatValidation(t *testing.T) {
	validKeys := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		uuid.New().String(),
	}

	for _, key := range validKeys {
		_, err := uuid.Parse(key)
		if err != nil {
			t.Errorf("Key %q should be valid UUID, got error: %v", key, err)
		}
	}

	invalidKeys := []string{
		"",
		"not-a-uuid",
		"12345",
		"550e8400-e29b-41d4-a716",
		"xyz-abc-def",
	}

	for _, key := range invalidKeys {
		_, err := uuid.Parse(key)
		if err == nil {
			t.Errorf("Key %q should be invalid UUID", key)
		}
	}
}

func TestIdempotencyMiddleware_EmptyKey_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.IdempotencyMiddleware("transfer"))
	r.POST("/transfer", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/transfer", nil)
	req.Header.Set("Idempotency-Key", "")
	c.Request = req
	c.Set("user", user)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty idempotency key, got %d", w.Code)
	}
}
