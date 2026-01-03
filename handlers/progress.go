package handlers

import (
	"learning-app-backend/database"
	"learning-app-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SaveProgress - Save or update user progress
func SaveProgress(c *gin.Context) {
	var req models.SaveProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ProgressResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate content type
	if req.ContentType != models.ContentTypeVideo && req.ContentType != models.ContentTypeQuiz {
		c.JSON(http.StatusBadRequest, models.ProgressResponse{
			Success: false,
			Message: "Invalid content type. Must be 'video' or 'quiz'",
		})
		return
	}

	// Validate that appropriate field is provided based on content type
	if req.ContentType == models.ContentTypeVideo && req.VideoTimestamp == nil {
		c.JSON(http.StatusBadRequest, models.ProgressResponse{
			Success: false,
			Message: "video_timestamp is required for video content type",
		})
		return
	}

	if req.ContentType == models.ContentTypeQuiz && req.QuizQuestionIndex == nil {
		c.JSON(http.StatusBadRequest, models.ProgressResponse{
			Success: false,
			Message: "quiz_question_index is required for quiz content type",
		})
		return
	}

	// Check if user exists
	var user models.User
	if err := database.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ProgressResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Check if chapter exists
	var chapter models.Chapter
	if err := database.DB.First(&chapter, req.ChapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ProgressResponse{
			Success: false,
			Message: "Chapter not found",
		})
		return
	}

	// Find existing progress record
	var progress models.Progress
	result := database.DB.Where("user_id = ? AND chapter_id = ? AND content_type = ?",
		req.UserID, req.ChapterID, req.ContentType).First(&progress)

	now := time.Now()

	if result.Error != nil {
		// Create new progress record
		progress = models.Progress{
			UserID:            req.UserID,
			ChapterID:         req.ChapterID,
			ContentType:       req.ContentType,
			VideoTimestamp:    req.VideoTimestamp,
			QuizQuestionIndex: req.QuizQuestionIndex,
			IsCompleted:       req.IsCompleted,
			LastUpdated:       now,
		}

		if err := database.DB.Create(&progress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ProgressResponse{
				Success: false,
				Message: "Failed to save progress",
			})
			return
		}
	} else {
		// Update existing progress
		progress.VideoTimestamp = req.VideoTimestamp
		progress.QuizQuestionIndex = req.QuizQuestionIndex
		progress.IsCompleted = req.IsCompleted
		progress.LastUpdated = now

		if err := database.DB.Save(&progress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ProgressResponse{
				Success: false,
				Message: "Failed to update progress",
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.ProgressResponse{
		Success:  true,
		Message:  "Progress saved successfully",
		Progress: &progress,
	})
}

// GetUserProgress - Get latest progress for a user
func GetUserProgress(c *gin.Context) {
	userID := c.Param("userId")

	// Get the most recent progress entry
	var progress models.Progress
	result := database.DB.Where("user_id = ?", userID).
		Order("last_updated DESC").
		First(&progress)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"has_progress": false,
			"message":     "No progress found for this user",
		})
		return
	}

	// Get chapter details
	var chapter models.Chapter
	database.DB.First(&chapter, progress.ChapterID)

	contentType := string(progress.ContentType)
	summary := models.UserProgressSummary{
		HasProgress:      true,
		LastChapterID:    &progress.ChapterID,
		LastContentType:  &contentType,
		LastVideoTime:    progress.VideoTimestamp,
		LastQuizQuestion: progress.QuizQuestionIndex,
		ChapterTitle:     &chapter.Title,
		LastUpdated:      &progress.LastUpdated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": summary,
	})
}

// GetChapterProgress - Get progress for a specific chapter
func GetChapterProgress(c *gin.Context) {
	userID := c.Param("userId")
	chapterID := c.Param("chapterId")

	var progressRecords []models.Progress
	database.DB.Where("user_id = ? AND chapter_id = ?", userID, chapterID).
		Order("last_updated DESC").
		Find(&progressRecords)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": progressRecords,
	})
}

// GetAllUserProgress - Get all progress for a user (all chapters)
func GetAllUserProgress(c *gin.Context) {
	userID := c.Param("userId")

	var progressRecords []models.Progress
	database.DB.Where("user_id = ?", userID).
		Order("chapter_id ASC, content_type ASC").
		Find(&progressRecords)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": progressRecords,
	})
}

// ResetProgress - Reset all progress for a user (optional feature)
func ResetProgress(c *gin.Context) {
	userID := c.Param("userId")

	result := database.DB.Where("user_id = ?", userID).Delete(&models.Progress{})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Progress reset successfully",
		"deleted": result.RowsAffected,
	})
}
