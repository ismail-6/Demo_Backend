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
	log.Printf("Port: %s", cfg.Port)
	log.Printf("Database Type: %s", cfg.DatabaseType)

	// Initialize database
	log.Println("Initializing database connection...")
	database.InitDatabase(cfg)
	log.Println("Database initialized successfully")

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
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/user/:userId", handlers.GetUser)
		}

		// Chapter routes (Raw SQL)
		chapters := api.Group("/chapters")
		{
			chapters.GET("", handlers.GetAllChapters)
			chapters.GET("/:id", handlers.GetChapterByID)
			chapters.GET("/:id/video", handlers.GetChapterVideo)
			chapters.GET("/:id/quiz", handlers.GetChapterQuiz)
			chapters.GET("/:id/content", handlers.GetChapterContent)
		}

		// Progress routes (Raw SQL)
		progress := api.Group("/progress")
		{
			progress.POST("", handlers.SaveProgress)
			progress.GET("/user/:userId", handlers.GetUserProgress)
			progress.GET("/user/:userId/all", handlers.GetAllUserProgress)
			progress.GET("/user/:userId/chapter/:chapterId", handlers.GetChapterProgress)
			progress.DELETE("/user/:userId/reset", handlers.ResetProgress)
		}

		// Quiz Answer routes (Raw SQL) - Track quiz question history
		quiz := api.Group("/quiz")
		{
			quiz.POST("/submit", handlers.SubmitQuizAnswer)
			quiz.GET("/history/user/:userId/chapter/:chapterId", handlers.GetQuizHistory)
			quiz.GET("/history/user/:userId", handlers.GetAllQuizHistory)
			quiz.GET("/history/user/:userId/question/:questionId", handlers.GetQuestionAnswerHistory)
			quiz.GET("/score/user/:userId", handlers.GetQuizScore)
			quiz.DELETE("/history/user/:userId/clear", handlers.ClearQuizHistory)

			// Get quiz with user's answer history (preserves state on reopen)
			quiz.GET("/chapter/:id/with-history", handlers.GetChapterQuizWithHistory)
			quiz.GET("/resume/user/:userId/chapter/:chapterId", handlers.GetQuizResumePoint)
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
