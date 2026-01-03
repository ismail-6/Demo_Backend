package handlers

import (
	"learning-app-backend/database"
	"learning-app-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllChapters - Get all chapters
func GetAllChapters(c *gin.Context) {
	var chapters []models.Chapter

	result := database.DB.Order("order_index ASC").Find(&chapters)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch chapters",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"chapters": chapters,
	})
}

// GetChapterByID - Get chapter details by ID
func GetChapterByID(c *gin.Context) {
	chapterID := c.Param("id")

	var chapter models.Chapter
	result := database.DB.First(&chapter, chapterID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"chapter": chapter,
	})
}

// GetChapterVideo - Get video for a specific chapter
func GetChapterVideo(c *gin.Context) {
	chapterID := c.Param("id")

	var video models.Video
	result := database.DB.Where("chapter_id = ?", chapterID).First(&video)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Video not found for this chapter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"video":   video,
	})
}

// GetChapterQuiz - Get quiz questions for a specific chapter
func GetChapterQuiz(c *gin.Context) {
	chapterID := c.Param("id")

	var questions []models.QuizQuestion
	result := database.DB.Where("chapter_id = ?", chapterID).Order("order_index ASC").Find(&questions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz questions",
		})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No quiz questions found for this chapter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"questions": questions,
	})
}

// GetChapterContent - Get both video and quiz for a chapter
func GetChapterContent(c *gin.Context) {
	chapterID := c.Param("id")

	var chapter models.Chapter
	result := database.DB.Preload("Video").Preload("QuizQuestions", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_index ASC")
	}).First(&chapter, chapterID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"chapter": chapter,
	})
}
