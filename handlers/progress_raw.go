package handlers

import (
	"database/sql"
	"learning-app-backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SaveProgressRequest struct {
	UserID            string `json:"user_id" binding:"required"`
	ChapterID         uint   `json:"chapter_id" binding:"required"`
	ContentType       string `json:"content_type" binding:"required"`
	VideoTimestamp    *int   `json:"video_timestamp"`
	QuizQuestionIndex *int   `json:"quiz_question_index"`
	IsCompleted       bool   `json:"is_completed"`
}

type Progress struct {
	ID                uint       `json:"id"`
	UserID            string     `json:"user_id"`
	ChapterID         uint       `json:"chapter_id"`
	ContentType       string     `json:"content_type"`
	VideoTimestamp    *int       `json:"video_timestamp,omitempty"`
	QuizQuestionIndex *int       `json:"quiz_question_index,omitempty"`
	IsCompleted       bool       `json:"is_completed"`
	LastUpdated       time.Time  `json:"last_updated"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type UserProgressSummary struct {
	HasProgress      bool       `json:"has_progress"`
	LastChapterID    *uint      `json:"last_chapter_id,omitempty"`
	LastContentType  *string    `json:"last_content_type,omitempty"`
	LastVideoTime    *int       `json:"last_video_time,omitempty"`
	LastQuizQuestion *int       `json:"last_quiz_question,omitempty"`
	ChapterTitle     *string    `json:"chapter_title,omitempty"`
	LastUpdated      *time.Time `json:"last_updated,omitempty"`
}

// SaveProgressRaw - Save or update progress using raw SQL
func SaveProgressRaw(c *gin.Context) {
	var req SaveProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate content type
	if req.ContentType != "video" && req.ContentType != "quiz" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid content type. Must be 'video' or 'quiz'",
		})
		return
	}

	// Validate required fields
	if req.ContentType == "video" && req.VideoTimestamp == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "video_timestamp is required for video content type",
		})
		return
	}

	if req.ContentType == "quiz" && req.QuizQuestionIndex == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "quiz_question_index is required for quiz content type",
		})
		return
	}

	sqlDB, _ := database.DB.DB()

	// Check if user exists
	var userExists bool
	userQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1 AND deleted_at IS NULL)`
	err := sqlDB.QueryRow(userQuery, req.UserID).Scan(&userExists)
	if err != nil || !userExists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Check if chapter exists
	var chapterExists bool
	chapterQuery := `SELECT EXISTS(SELECT 1 FROM chapters WHERE id = $1 AND deleted_at IS NULL)`
	err = sqlDB.QueryRow(chapterQuery, req.ChapterID).Scan(&chapterExists)
	if err != nil || !chapterExists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	}

	// Check if progress exists
	var progressID sql.NullInt64
	checkQuery := `SELECT id FROM progresses
				   WHERE user_id = $1 AND chapter_id = $2 AND content_type = $3 AND deleted_at IS NULL`
	err = sqlDB.QueryRow(checkQuery, req.UserID, req.ChapterID, req.ContentType).Scan(&progressID)

	var progress Progress

	if err == sql.ErrNoRows {
		// Create new progress
		insertQuery := `INSERT INTO progresses (user_id, chapter_id, content_type, video_timestamp,
				   quiz_question_index, is_completed, last_updated, created_at, updated_at)
				   VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW())
				   RETURNING id, user_id, chapter_id, content_type, video_timestamp,
				   quiz_question_index, is_completed, last_updated, created_at, updated_at`

		err = sqlDB.QueryRow(insertQuery, req.UserID, req.ChapterID, req.ContentType,
			req.VideoTimestamp, req.QuizQuestionIndex, req.IsCompleted).Scan(
			&progress.ID, &progress.UserID, &progress.ChapterID, &progress.ContentType,
			&progress.VideoTimestamp, &progress.QuizQuestionIndex, &progress.IsCompleted,
			&progress.LastUpdated, &progress.CreatedAt, &progress.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save progress",
			})
			return
		}
	} else if err == nil {
		// Update existing progress
		updateQuery := `UPDATE progresses SET video_timestamp = $1, quiz_question_index = $2,
						is_completed = $3, last_updated = NOW(), updated_at = NOW()
						WHERE id = $4
						RETURNING id, user_id, chapter_id, content_type, video_timestamp,
						quiz_question_index, is_completed, last_updated, created_at, updated_at`

		err = sqlDB.QueryRow(updateQuery, req.VideoTimestamp, req.QuizQuestionIndex,
			req.IsCompleted, progressID.Int64).Scan(
			&progress.ID, &progress.UserID, &progress.ChapterID, &progress.ContentType,
			&progress.VideoTimestamp, &progress.QuizQuestionIndex, &progress.IsCompleted,
			&progress.LastUpdated, &progress.CreatedAt, &progress.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update progress",
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Progress saved successfully",
		"progress": progress,
	})
}

// GetUserProgressRaw - Get latest progress for a user
func GetUserProgressRaw(c *gin.Context) {
	userID := c.Param("userId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT p.id, p.user_id, p.chapter_id, p.content_type, p.video_timestamp,
			  p.quiz_question_index, p.is_completed, p.last_updated, ch.title
			  FROM progresses p
			  JOIN chapters ch ON p.chapter_id = ch.id
			  WHERE p.user_id = $1 AND p.deleted_at IS NULL
			  ORDER BY p.last_updated DESC LIMIT 1`

	var progress Progress
	var chapterTitle string

	err := sqlDB.QueryRow(query, userID).Scan(
		&progress.ID, &progress.UserID, &progress.ChapterID, &progress.ContentType,
		&progress.VideoTimestamp, &progress.QuizQuestionIndex, &progress.IsCompleted,
		&progress.LastUpdated, &chapterTitle,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"has_progress": false,
			"message":      "No progress found for this user",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	contentType := progress.ContentType
	summary := UserProgressSummary{
		HasProgress:      true,
		LastChapterID:    &progress.ChapterID,
		LastContentType:  &contentType,
		LastVideoTime:    progress.VideoTimestamp,
		LastQuizQuestion: progress.QuizQuestionIndex,
		ChapterTitle:     &chapterTitle,
		LastUpdated:      &progress.LastUpdated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": summary,
	})
}

// GetChapterProgressRaw - Get progress for a specific chapter
func GetChapterProgressRaw(c *gin.Context) {
	userID := c.Param("userId")
	chapterID := c.Param("chapterId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT id, user_id, chapter_id, content_type, video_timestamp,
			  quiz_question_index, is_completed, last_updated, created_at, updated_at
			  FROM progresses
			  WHERE user_id = $1 AND chapter_id = $2 AND deleted_at IS NULL
			  ORDER BY last_updated DESC`

	rows, err := sqlDB.Query(query, userID, chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch progress",
		})
		return
	}
	defer rows.Close()

	var progressRecords []Progress
	for rows.Next() {
		var p Progress
		err := rows.Scan(&p.ID, &p.UserID, &p.ChapterID, &p.ContentType,
			&p.VideoTimestamp, &p.QuizQuestionIndex, &p.IsCompleted,
			&p.LastUpdated, &p.CreatedAt, &p.UpdatedAt)
		if err == nil {
			progressRecords = append(progressRecords, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": progressRecords,
	})
}

// GetAllUserProgressRaw - Get all progress for a user
func GetAllUserProgressRaw(c *gin.Context) {
	userID := c.Param("userId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT id, user_id, chapter_id, content_type, video_timestamp,
			  quiz_question_index, is_completed, last_updated, created_at, updated_at
			  FROM progresses
			  WHERE user_id = $1 AND deleted_at IS NULL
			  ORDER BY chapter_id ASC, content_type ASC`

	rows, err := sqlDB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch progress",
		})
		return
	}
	defer rows.Close()

	var progressRecords []Progress
	for rows.Next() {
		var p Progress
		err := rows.Scan(&p.ID, &p.UserID, &p.ChapterID, &p.ContentType,
			&p.VideoTimestamp, &p.QuizQuestionIndex, &p.IsCompleted,
			&p.LastUpdated, &p.CreatedAt, &p.UpdatedAt)
		if err == nil {
			progressRecords = append(progressRecords, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"progress": progressRecords,
	})
}

// ResetProgressRaw - Reset all progress for a user
func ResetProgressRaw(c *gin.Context) {
	userID := c.Param("userId")
	sqlDB, _ := database.DB.DB()

	// Soft delete
	query := `UPDATE progresses SET deleted_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`

	result, err := sqlDB.Exec(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to reset progress",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Progress reset successfully",
		"deleted": rowsAffected,
	})
}
