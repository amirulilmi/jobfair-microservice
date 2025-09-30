package main

import (
	"log"
	"os"

	"jobfair-company-service/internal/config"
	"jobfair-company-service/internal/handlers"
	"jobfair-company-service/internal/repository"
	"jobfair-company-service/internal/services"
	"jobfair-company-service/pkg/database"

	"jobfair-company-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
	// defer db.Close()

	companyRepo := repository.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Initialize jwtSecret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret" // fallback untuk dev
	}

	// Setup middleware
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	api := router.Group("/api/v1")
	api.Use(jwtMiddleware) // apply JWT middleware
	{
		api.POST("/companies", companyHandler.CreateCompany)
		api.GET("/companies/:id", companyHandler.GetCompany)
		api.PUT("/companies/:id", companyHandler.UpdateCompany)
		api.POST("/companies/:id/logo", companyHandler.UploadLogo)
		api.POST("/companies/:id/banner", companyHandler.UploadBanner)
		api.POST("/companies/:id/videos", companyHandler.UploadVideo)
		api.GET("/companies/:id/analytics", companyHandler.GetAnalytics)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Company service starting on port %s", port)
	router.Run(":" + port)
}
