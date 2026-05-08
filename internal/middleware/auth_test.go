package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for missing auth header, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid header format, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic some-token")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong scheme, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidJWT(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-middleware-tests")
	defer os.Unsetenv("JWT_SECRET_KEY")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-jwt-token-here")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid JWT, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongSignatureRejectedBeforeDatabaseUse(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-middleware-tests")
	defer os.Unsetenv("JWT_SECRET_KEY")

	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	tokenStr, err := generateTestTokenWithSecret(user, "wrong-secret-for-middleware-tests")
	if err != nil {
		t.Fatalf("Failed to generate wrong-signature test token: %v", err)
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 for wrong-signature JWT, got %d", w.Code)
	}
}

func TestOptionalAuthMiddleware_NoHeader_Passes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.OptionalAuthMiddleware())
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 when no auth header, got %d", w.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidToken_Passes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.OptionalAuthMiddleware())
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for optional auth with invalid token, got %d", w.Code)
	}
}

func TestAdminMiddleware_NonAdmin_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "regularuser",
		IsAdmin:  false,
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	})
	r.Use(middleware.AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 for non-admin, got %d", w.Code)
	}
}

func TestAdminMiddleware_Admin_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "adminuser",
		IsAdmin:  true,
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	})
	r.Use(middleware.AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for admin, got %d", w.Code)
	}
}

func TestAdminMiddleware_NoUser_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 when no user in context, got %d", w.Code)
	}
}

func TestAdminOnlyMiddleware_NonAdmin_Rejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "regularuser",
		IsAdmin:  false,
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	})
	r.Use(middleware.AdminOnlyMiddleware())
	r.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 for non-admin, got %d", w.Code)
	}
}

func TestAdminOnlyMiddleware_ErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := &models.User{
		ID:       uuid.New(),
		Username: "regularuser",
		IsAdmin:  false,
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.AdminOnlyMiddleware())
	r.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	c.Request = req
	c.Set("user", user)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	if err := parseJSON(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := resp["error"]; !ok {
		t.Error("Expected error field in admin rejection response")
	}
}

func TestSecurityHeadersMiddleware_SetsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expected := range expectedHeaders {
		got := w.Header().Get(header)
		if got != expected {
			t.Errorf("Header %q: expected %q, got %q", header, expected, got)
		}
	}
}
