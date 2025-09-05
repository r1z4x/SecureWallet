package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"securewallet/internal/config"
	"securewallet/internal/middleware"
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
	// Parse command line flags
	cronJob := flag.String("cron", "", "Execute a specific cron job (comment-approval, backup, log-cleanup, security-monitor)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Handle cron job execution
	if *cronJob != "" {
		executeCronJob(*cronJob)
		return
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Check if database reset is requested on startup BEFORE initializing the database
	if os.Getenv("RESET_DATABASE_ON_STARTUP") == "true" {
		log.Println("RESET_DATABASE_ON_STARTUP is enabled, will reset database after initialization...")
	}

	// Initialize database connection (but skip auto-migration if reset is requested)
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	services.InitServices()
	
	// Initialize comment service (starts auto-approval scheduler)
	services.NewCommentService()

	// Now perform database reset if requested
	if os.Getenv("RESET_DATABASE_ON_STARTUP") == "true" {
		log.Println("Performing database reset...")
		dataManager := services.NewSampleDataManager()

		// Check if force recreation is requested
		if os.Getenv("FORCE_DATABASE_RECREATION") == "true" {
			log.Println("FORCE_DATABASE_RECREATION is enabled, forcing complete database recreation...")
			if err := dataManager.CompleteDatabaseRecreation(); err != nil {
				log.Printf("Warning: Failed to force database recreation on startup: %v", err)
			} else {
				log.Println("Database recreation completed successfully on startup")
			}
		} else {
			if err := dataManager.ResetDatabase(); err != nil {
				log.Printf("Warning: Failed to reset database on startup: %v", err)
			} else {
				log.Println("Database reset completed successfully on startup")
			}
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
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.InputValidationMiddleware())
	r.Use(func(c *gin.Context) {
		// SECURE: Get allowed origins from environment variable
		allowedOrigins := os.Getenv("CORS_ORIGINS")
		origin := c.Request.Header.Get("Origin")

		// Debug logging
		log.Printf("CORS Debug - Origin: %s, AllowedOrigins: %s", origin, allowedOrigins)

		// SECURE: Validate origin against allowed origins
		if allowedOrigins != "" && origin != "" {
			// Parse the JSON array from environment variable
			// Remove brackets and split by comma
			originsStr := strings.Trim(allowedOrigins, "[]")
			origins := strings.Split(originsStr, ",")

			// Check if origin is in allowed list
			for _, allowedOrigin := range origins {
				allowedOrigin = strings.Trim(allowedOrigin, `" `)
				log.Printf("CORS Debug - Checking origin: %s against allowed: %s", origin, allowedOrigin)
				if origin == allowedOrigin {
					c.Header("Access-Control-Allow-Origin", origin)
					log.Printf("CORS Debug - Origin allowed: %s", origin)
					break
				}
			}
		} else {
			// Fallback to default - include both common ports
			if origin == "http://localhost:3001" || origin == "http://127.0.0.1:3001" {
				c.Header("Access-Control-Allow-Origin", origin)
				log.Printf("CORS Debug - Fallback origin allowed: %s", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
				log.Printf("CORS Debug - Default origin set: http://localhost:3000")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

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
		routes.SetupTwoFactorRoutes(api)
		routes.SetupLoginHistoryRoutes(api)
		routes.SetupBackupRoutes(api)
		routes.SetupSecurityRoutes(api)
	}

	// Blog routes (public access)
	routes.BlogRoutes(r, config.GetDB())

	// Cron routes (admin access)
	routes.CronRoutes(r)

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SecureWallet - Digital Banking Platform (Vulnerable) API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
			"health":  "/health",
		})
	})

	// Test CORS endpoint
	r.GET("/test-cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":      "CORS test successful",
			"origin":       c.Request.Header.Get("Origin"),
			"cors_working": true,
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
			env = "production"
		}

		// SECURE: Only return minimal, non-sensitive information
		c.JSON(200, gin.H{
			"message":     "SecureWallet API",
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

// executeCronJob executes a specific cron job
func executeCronJob(jobName string) {
	log.Printf("Executing cron job: %s", jobName)

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize database connection
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	services.InitServices()

	// Create cron service and execute the job
	cronService := services.NewCronService()
	if err := cronService.ExecuteCronJob(jobName); err != nil {
		log.Fatal("Failed to execute cron job:", err)
	}

	log.Printf("Cron job %s completed successfully", jobName)
}
