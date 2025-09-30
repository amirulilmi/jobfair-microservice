package main

import (
	"log"
	"os"

	"jobfair-company-service/internal/config"
	"jobfair-company-service/internal/handlers"
	"jobfair-company-service/internal/middleware"
	"jobfair-company-service/internal/repository"
	"jobfair-company-service/internal/services"
	"jobfair-company-service/pkg/database"

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

	companyRepo := repository.NewCompanyRepository(db)
	jobRepo := repository.NewJobRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)

	companyService := services.NewCompanyService(companyRepo, jobRepo, applicationRepo)
	jobService := services.NewJobService(jobRepo, companyRepo, applicationRepo)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo)

	companyHandler := handlers.NewCompanyHandler(companyService)
	jobHandler := handlers.NewJobHandler(companyService, jobService)
	applicationHandler := handlers.NewApplicationHandler(companyService, applicationService)

	router := gin.Default()
	router.MaxMultipartMemory = 10 << 20

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "company-service",
			"version": "1.0.0",
		})
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret"
	}

	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	api := router.Group("/api/v1")
	{
		public := api.Group("")
		{
			public.GET("/companies", companyHandler.ListCompanies)
			public.GET("/companies/:id", companyHandler.GetCompany)
			public.GET("/jobs/:id", jobHandler.GetJob)
		}

		protected := api.Group("")
		protected.Use(jwtMiddleware)
		{
			protected.GET("/my-company", companyHandler.GetMyCompany)
			protected.POST("/companies", companyHandler.CreateCompany)
			protected.PUT("/companies/:id", companyHandler.UpdateCompany)

			protected.POST("/companies/:id/logo", companyHandler.UploadLogo)
			protected.POST("/companies/:id/banner", companyHandler.UploadBanner)
			protected.POST("/companies/:id/videos", companyHandler.UploadVideo)
			protected.POST("/companies/:id/gallery", companyHandler.UploadGallery)

			protected.GET("/companies/:id/analytics", companyHandler.GetAnalytics)
			protected.GET("/dashboard", companyHandler.GetDashboard)

			protected.POST("/jobs", jobHandler.CreateJob)
			protected.GET("/jobs", jobHandler.ListJobs)
			protected.PUT("/jobs/:id", jobHandler.UpdateJob)
			protected.DELETE("/jobs/:id", jobHandler.DeleteJob)
			protected.POST("/jobs/:id/publish", jobHandler.PublishJob)
			protected.POST("/jobs/:id/close", jobHandler.CloseJob)

			protected.GET("/applications", applicationHandler.ListApplications)
			protected.GET("/applications/:id", applicationHandler.GetApplication)
			protected.GET("/jobs/:job_id/applications", applicationHandler.GetApplicationsByJobID)
			protected.PUT("/applications/:id/status", applicationHandler.UpdateApplicationStatus)
			protected.GET("/applications/stats", applicationHandler.GetApplicationStats)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Company service starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("🔑 API endpoint: http://localhost:%s/api/v1", port)
	router.Run(":" + port)
}
