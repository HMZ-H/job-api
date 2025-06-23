package main

import (
	"log"
	"os"
	"job-api/config"
	"job-api/handlers"
	"job-api/middleware"
	"job-api/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	config.ConnectDatabase()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Job routes
		jobs := api.Group("/jobs")
		{
			// Company only routes
			jobs.POST("", middleware.RequireRole(models.RoleCompany), handlers.CreateJob)
			jobs.PUT("/:id", middleware.RequireRole(models.RoleCompany), handlers.UpdateJob)
			jobs.DELETE("/:id", middleware.RequireRole(models.RoleCompany), handlers.DeleteJob)
			jobs.GET("/my-jobs", middleware.RequireRole(models.RoleCompany), handlers.GetMyJobs)
			jobs.GET("/:id/applications", middleware.RequireRole(models.RoleCompany), handlers.GetJobApplications)

			// Applicant only routes
			jobs.GET("", middleware.RequireRole(models.RoleApplicant), handlers.BrowseJobs)
			jobs.POST("/:id/apply", middleware.RequireRole(models.RoleApplicant), handlers.ApplyForJob)

			// Both roles can access
			jobs.GET("/:id", handlers.GetJobDetails)
		}

		// Application routes
		applications := api.Group("/applications")
		{
			// Applicant only routes
			applications.GET("/my-applications", middleware.RequireRole(models.RoleApplicant), handlers.GetMyApplications)

			// Company only routes
			applications.PUT("/:id/status", middleware.RequireRole(models.RoleCompany), handlers.UpdateApplicationStatus)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
