package main

import (
	"log"
	"os"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/routes"
	"securewallet/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "securewallet/docs" // Swagger docs
)

// @title SecureWallet - Digital Banking Platform (Vulnerable)
// @description A comprehensive vulnerable application for OWASP Top 10
// @version 1.0.0
// @host localhost:8080
// @BasePath /api
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize services
	services.InitServices()

	// Check if database reset is requested on startup
	if os.Getenv("RESET_DATABASE_ON_STARTUP") == "true" {
		log.Println("RESET_DATABASE_ON_STARTUP is enabled, resetting database...")
		dataManager := services.NewSampleDataManager()
		if err := dataManager.ResetDatabase(); err != nil {
			log.Printf("Warning: Failed to reset database on startup: %v", err)
		} else {
			log.Println("Database reset completed successfully on startup")
		}
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		routes.SetupAuthRoutes(api)
		routes.SetupUserRoutes(api)
		routes.SetupWalletRoutes(api)
		routes.SetupTransactionRoutes(api)
		routes.SetupAdminRoutes(api)
		routes.SetupSupportRoutes(api)
		routes.SetupDataManagementRoutes(api)
		routes.SetupVulnerabilityRoutes(api)
		routes.SetupTwoFactorRoutes(api)
		routes.SetupLoginHistoryRoutes(api)
	}

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SecureWallet - Digital Banking Platform (Vulnerable) API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
			"health":  "/health",
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"app_name":  "SecureWallet - Digital Banking Platform (Vulnerable)",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API info endpoint
	r.GET("/api/info", func(c *gin.Context) {
		env := os.Getenv("ENVIRONMENT")
		if env == "" {
			env = "development"
		}
		c.JSON(200, gin.H{
			"message":     "SecureWallet - Digital Banking Platform (Vulnerable) API",
			"version":     "1.0.0",
			"status":      "running",
			"environment": env,
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
