package handlers

import (
	"database/sql"
	"learning-app-backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Chapter struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Video struct {
	ID              uint      `json:"id"`
	ChapterID       uint      `json:"chapter_id"`
	Title           string    `json:"title"`
	VideoURL        string    `json:"video_url"`
	DurationSeconds int       `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type QuizQuestion struct {
	ID            uint      `json:"id"`
	ChapterID     uint      `json:"chapter_id"`
	QuestionText  string    `json:"question_text"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectAnswer string    `json:"correct_answer"`
	OrderIndex    int       `json:"order_index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChapterWithContent struct {
	Chapter
	Video         *Video         `json:"video,omitempty"`
	QuizQuestions []QuizQuestion `json:"quiz_questions,omitempty"`
}

// GetAllChaptersRaw - Get all chapters using raw SQL
func GetAllChapters(c *gin.Context) {
	sqlDB, _ := database.DB.DB()

	query := `SELECT id, title, description, order_index, created_at, updated_at
			  FROM chapters WHERE deleted_at IS NULL ORDER BY order_index ASC`

	rows, err := sqlDB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch chapters",
		})
		return
	}
	defer rows.Close()

	var chapters []Chapter
	for rows.Next() {
		var ch Chapter
		err := rows.Scan(&ch.ID, &ch.Title, &ch.Description, &ch.OrderIndex, &ch.CreatedAt, &ch.UpdatedAt)
		if err != nil {
			continue
		}
		chapters = append(chapters, ch)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"chapters": chapters,
	})
}

// GetChapterByIDRaw - Get chapter by ID
func GetChapterByID(c *gin.Context) {
	chapterID := c.Param("id")
	sqlDB, _ := database.DB.DB()

	var chapter Chapter
	query := `SELECT id, title, description, order_index, created_at, updated_at
			  FROM chapters WHERE id = $1 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(query, chapterID).Scan(
		&chapter.ID, &chapter.Title, &chapter.Description, &chapter.OrderIndex,
		&chapter.CreatedAt, &chapter.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"chapter": chapter,
	})
}

// GetChapterVideoRaw - Get video for a chapter
func GetChapterVideo(c *gin.Context) {
	chapterID := c.Param("id")
	sqlDB, _ := database.DB.DB()

	var video Video
	query := `SELECT id, chapter_id, title, video_url, duration_seconds, created_at, updated_at
			  FROM videos WHERE chapter_id = $1 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(query, chapterID).Scan(
		&video.ID, &video.ChapterID, &video.Title, &video.VideoURL,
		&video.DurationSeconds, &video.CreatedAt, &video.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Video not found for this chapter",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"video":   video,
	})
}

// GetChapterQuizRaw - Get quiz questions for a chapter
func GetChapterQuiz(c *gin.Context) {
	chapterID := c.Param("id")
	sqlDB, _ := database.DB.DB()

	query := `SELECT id, chapter_id, question_text, option_a, option_b, option_c, option_d,
			  correct_answer, order_index, created_at, updated_at
			  FROM quiz_questions WHERE chapter_id = $1 AND deleted_at IS NULL
			  ORDER BY order_index ASC`

	rows, err := sqlDB.Query(query, chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz questions",
		})
		return
	}
	defer rows.Close()

	var questions []QuizQuestion
	for rows.Next() {
		var q QuizQuestion
		err := rows.Scan(&q.ID, &q.ChapterID, &q.QuestionText, &q.OptionA, &q.OptionB,
			&q.OptionC, &q.OptionD, &q.CorrectAnswer, &q.OrderIndex,
			&q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			continue
		}
		questions = append(questions, q)
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

// GetChapterContentRaw - Get chapter with video and quiz questions
func GetChapterContent(c *gin.Context) {
	chapterID := c.Param("id")
	sqlDB, _ := database.DB.DB()

	// Get chapter
	var chapter ChapterWithContent
	chapterQuery := `SELECT id, title, description, order_index, created_at, updated_at
					 FROM chapters WHERE id = $1 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(chapterQuery, chapterID).Scan(
		&chapter.ID, &chapter.Title, &chapter.Description, &chapter.OrderIndex,
		&chapter.CreatedAt, &chapter.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	// Get video
	var video Video
	videoQuery := `SELECT id, chapter_id, title, video_url, duration_seconds, created_at, updated_at
				   FROM videos WHERE chapter_id = $1 AND deleted_at IS NULL`

	err = sqlDB.QueryRow(videoQuery, chapterID).Scan(
		&video.ID, &video.ChapterID, &video.Title, &video.VideoURL,
		&video.DurationSeconds, &video.CreatedAt, &video.UpdatedAt,
	)

	if err == nil {
		chapter.Video = &video
	}

	// Get quiz questions
	quizQuery := `SELECT id, chapter_id, question_text, option_a, option_b, option_c, option_d,
				  correct_answer, order_index, created_at, updated_at
				  FROM quiz_questions WHERE chapter_id = $1 AND deleted_at IS NULL
				  ORDER BY order_index ASC`

	rows, err := sqlDB.Query(quizQuery, chapterID)
	if err == nil {
		defer rows.Close()
		var questions []QuizQuestion
		for rows.Next() {
			var q QuizQuestion
			err := rows.Scan(&q.ID, &q.ChapterID, &q.QuestionText, &q.OptionA, &q.OptionB,
				&q.OptionC, &q.OptionD, &q.CorrectAnswer, &q.OrderIndex,
				&q.CreatedAt, &q.UpdatedAt)
			if err == nil {
				questions = append(questions, q)
			}
		}
		chapter.QuizQuestions = questions
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"chapter": chapter,
	})
}
