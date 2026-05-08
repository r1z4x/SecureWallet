package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"securewallet/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestRateLimitMiddleware_SkipsInDevelopment(t *testing.T) {
	os.Setenv("ENVIRONMENT", "development")
	defer os.Unsetenv("ENVIRONMENT")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.RateLimitMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 in development mode, got %d", w.Code)
	}
}

func TestRateLimitMiddleware_BlocksAfterThreshold(t *testing.T) {
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("GIN_MODE")

	gin.SetMode(gin.ReleaseMode)

	for i := 0; i < 6; i++ {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(middleware.RateLimitMiddleware())
		r.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		c.Request = req
		r.ServeHTTP(w, req)

		if i < 5 {
			if w.Code != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d: expected status 429, got %d", i+1, w.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if _, ok := resp["error"]; !ok {
				t.Error("Expected error field in rate limit response")
			}
		}
	}
}

func TestRateLimitMiddleware_DifferentIPsIndependent(t *testing.T) {
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("GIN_MODE")

	gin.SetMode(gin.ReleaseMode)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(middleware.RateLimitMiddleware())
		r.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("IP 10.0.0.1 request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.RateLimitMiddleware())
	r.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "10.0.0.2:12345"
	c.Request = req
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Different IP should not be rate limited, got %d", w.Code)
	}
}
