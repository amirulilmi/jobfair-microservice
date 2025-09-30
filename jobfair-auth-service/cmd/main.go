// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	// "github.com/jobfair/jobfair-auth-service/internal/config"
// 	// "github.com/jobfair/jobfair-auth-service/internal/handlers"
// 	// "github.com/jobfair/jobfair-auth-service/internal/repository"
// 	// "github.com/jobfair/jobfair-auth-service/internal/services"
// 	// "github.com/jobfair/jobfair-auth-service/pkg/database"
// 	"jobfair-auth-service/internal/config"
//     "jobfair-auth-service/internal/handlers"
//     "jobfair-auth-service/internal/repository"
//     "jobfair-auth-service/internal/services"
//     "jobfair-auth-service/pkg/database"
// )

// func main() {
// 	cfg := config.Load()

// 	db, err := database.Connect(cfg.DatabaseURL)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Correct way to close connection in GORM v2
// 		defer func() {
// 			sqlDB, err := db.DB()
// 			if err != nil {
// 				log.Printf("Error getting underlying SQL DB: %v", err)
// 				return
// 			}
// 			if err := sqlDB.Close(); err != nil {
// 				log.Printf("Error closing database connection: %v", err)
// 			}
// 		}()

// 	userRepo := repository.NewUserRepository(db)
// 	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
// 	authHandler := handlers.NewAuthHandler(authService)

// 	router := gin.Default()

// 	// Health check endpoint
// 	router.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":  "healthy",
// 			"service": "auth-service",
// 			"version": "1.0.0",
// 		})
// 	})

// 	// CORS middleware
// 	router.Use(func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	api := router.Group("/api/v1")
// 	{
// 		api.POST("/register", authHandler.Register)
// 		api.POST("/login", authHandler.Login)
// 		api.POST("/refresh", authHandler.RefreshToken)
// 	}

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	log.Printf("🚀 Auth service starting on port %s", port)
// 	log.Printf("📊 Health check: http://localhost:%s/health", port)
// 	log.Printf("🔑 API endpoint: http://localhost:%s/api/v1", port)

// 	if err := router.Run(":" + port); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }

package main

import (
	"log"
	"os"

	"jobfair-auth-service/internal/config"
	"jobfair-auth-service/internal/handlers"
	"jobfair-auth-service/internal/middleware"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/services"
	"jobfair-auth-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// sql := `ALTER TABLE users ALTER COLUMN phone_number DROP NOT NULL;`
	// if err := db.Exec(sql).Error; err != nil {
	// 	log.Fatal("Failed to alter column:", err)
	// }

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewJobSeekerProfileRepository(db)
	otpRepo := repository.NewOTPRepository(db)

	// Initialize services
	registrationService := services.NewRegistrationService(userRepo, profileRepo, otpRepo, cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	// Initialize handlers
	registrationHandler := handlers.NewRegistrationHandler(registrationService)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	api := router.Group("/api/v1")
	{
		// Registration flow (multi-step)
		register := api.Group("/register")
		{
			register.POST("/step1", registrationHandler.RegisterStep1)                            // Email & Password
			register.PUT("/profile", authMiddleware, registrationHandler.CompleteBasicProfile)    // Basic info
			register.POST("/send-otp", authMiddleware, registrationHandler.SendPhoneOTP)          // Send OTP
			register.POST("/verify-otp", registrationHandler.VerifyPhoneOTP)                      // Verify OTP
			register.POST("/employment", authMiddleware, registrationHandler.SetEmploymentStatus) // Job seeker only
			register.POST("/preferences", authMiddleware, registrationHandler.SetJobPreferences)  // Job seeker only
			// register.POST("/permissions", authMiddleware, registrationHandler.SetPermissions)     // Notifications & Location
			register.POST("/photo", authMiddleware, registrationHandler.UploadProfilePhoto) // Profile photo
			register.GET("/users", authHandler.GetAllUsers)
		}

		// Authentication
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		// api.GET("/profile", authMiddleware, authHandler.GetProfile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Auth service starting on port %s", port)
	router.Run(":" + port)
}
