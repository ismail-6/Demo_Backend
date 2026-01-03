package main

import (
	"learning-app-backend/config"
	"learning-app-backend/database"
	"learning-app-backend/handlers"
	"learning-app-backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Starting in %s mode", cfg.Environment)

	// Initialize database
	database.InitDatabase(cfg)

	// Create Gin router
	router := gin.Default()

	// Setup CORS
	router.Use(middleware.SetupCORS())

	// API Routes
	api := router.Group("/api")
	{
		// Auth routes (Raw SQL)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.LoginRaw)
			auth.POST("/logout", handlers.LogoutRaw)
			auth.GET("/user/:userId", handlers.GetUserRaw)
		}

		// Chapter routes (Raw SQL)
		chapters := api.Group("/chapters")
		{
			chapters.GET("", handlers.GetAllChaptersRaw)
			chapters.GET("/:id", handlers.GetChapterByIDRaw)
			chapters.GET("/:id/video", handlers.GetChapterVideoRaw)
			chapters.GET("/:id/quiz", handlers.GetChapterQuizRaw)
			chapters.GET("/:id/content", handlers.GetChapterContentRaw)
		}

		// Progress routes (Raw SQL)
		progress := api.Group("/progress")
		{
			progress.POST("", handlers.SaveProgressRaw)
			progress.GET("/user/:userId", handlers.GetUserProgressRaw)
			progress.GET("/user/:userId/all", handlers.GetAllUserProgressRaw)
			progress.GET("/user/:userId/chapter/:chapterId", handlers.GetChapterProgressRaw)
			progress.DELETE("/user/:userId/reset", handlers.ResetProgressRaw)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Learning App API is running",
			"environment": cfg.Environment,
			"database": cfg.DatabaseType,
		})
	})

	// Start server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s...", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
