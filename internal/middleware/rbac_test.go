package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"securewallet/internal/middleware"
	"securewallet/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequirePermission_RejectsUnauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/test", middleware.RequirePermission(models.PermWalletRead), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["error"] != "Unauthorized" {
		t.Errorf("Expected 'Unauthorized' error, got %v", resp)
	}
}

func TestRequirePermission_AllowsLegacyAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	adminUser := &models.User{
		ID:       uuid.New(),
		Username: "admin",
		IsAdmin:  true,
	}

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user", adminUser)
		c.Next()
	})
	router.GET("/test", middleware.RequirePermission(models.PermAdminAll), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestRequireAnyPermission_RejectsUnauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/test", middleware.RequireAnyPermission([]string{models.PermWalletRead, models.PermWalletWrite}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestRequireAnyPermission_AllowsLegacyAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	adminUser := &models.User{
		ID:       uuid.New(),
		Username: "admin",
		IsAdmin:  true,
	}

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user", adminUser)
		c.Next()
	})
	router.GET("/test", middleware.RequireAnyPermission([]string{models.PermAdminAll, models.PermDataManage}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestRequireRole_RejectsUnauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.GET("/test", middleware.RequireRole(models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestRequireRole_AllowsLegacyAdminForAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	adminUser := &models.User{
		ID:       uuid.New(),
		Username: "admin",
		IsAdmin:  true,
	}

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user", adminUser)
		c.Next()
	})
	router.GET("/test", middleware.RequireRole(models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestPermissionConstants_CoverAllResources(t *testing.T) {
	expectedPerms := map[string]string{
		"wallet:read":        "wallet:read",
		"wallet:write":       "wallet:write",
		"wallet:delete":      "wallet:delete",
		"transfer:read":      "transfer:read",
		"transfer:write":     "transfer:write",
		"transfer:delete":    "transfer:delete",
		"user:read":          "user:read",
		"user:write":         "user:write",
		"user:delete":        "user:delete",
		"session:read":       "session:read",
		"session:write":      "session:write",
		"2fa:read":           "2fa:read",
		"2fa:write":          "2fa:write",
		"login_history:read": "login_history:read",
		"support:read":       "support:read",
		"support:write":      "support:write",
		"audit:read":         "audit:read",
		"backup:read":        "backup:read",
		"security:read":      "security:read",
		"security:write":     "security:write",
		"data:manage":        "data:manage",
		"cron:manage":        "cron:manage",
		"blog:read":          "blog:read",
		"blog:comment":       "blog:comment",
		"admin:*":            "admin:*",
	}

	allPerms := []string{
		models.PermWalletRead, models.PermWalletWrite, models.PermWalletDelete,
		models.PermTransferRead, models.PermTransferWrite, models.PermTransferDelete,
		models.PermUserRead, models.PermUserWrite, models.PermUserDelete,
		models.PermSessionRead, models.PermSessionWrite,
		models.PermTwoFactorRead, models.PermTwoFactorWrite,
		models.PermLoginHistoryRead,
		models.PermSupportRead, models.PermSupportWrite,
		models.PermAuditRead,
		models.PermBackupRead,
		models.PermSecurityRead, models.PermSecurityWrite,
		models.PermDataManage,
		models.PermCronManage,
		models.PermBlogRead, models.PermBlogComment,
		models.PermAdminAll,
	}

	for _, perm := range allPerms {
		if _, ok := expectedPerms[perm]; !ok {
			t.Errorf("Unexpected permission constant: %s", perm)
		}
	}

	if len(allPerms) != len(expectedPerms) {
		t.Errorf("Expected %d permissions, got %d", len(expectedPerms), len(allPerms))
	}
}

func TestRoleConstants_Exist(t *testing.T) {
	roles := map[string]string{
		"admin":     models.RoleAdmin,
		"user":      models.RoleUser,
		"auditor":   models.RoleAuditor,
		"support":   models.RoleSupport,
	}

	for name, constant := range roles {
		if constant != name {
			t.Errorf("Role %s: expected %q, got %q", name, name, constant)
		}
	}
}
