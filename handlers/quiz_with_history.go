package handlers

import (
	"database/sql"
	"learning-app-backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QuizQuestionWithUserAnswer struct {
	ID            uint       `json:"id"`
	ChapterID     uint       `json:"chapter_id"`
	QuestionText  string     `json:"question_text"`
	OptionA       string     `json:"option_a"`
	OptionB       string     `json:"option_b"`
	OptionC       string     `json:"option_c"`
	OptionD       string     `json:"option_d"`
	OrderIndex    int        `json:"order_index"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// User's answer history (if exists)
	HasAnswered      bool       `json:"has_answered"`
	UserAnswer       *string    `json:"user_answer,omitempty"`
	IsCorrect        *bool      `json:"is_correct,omitempty"`
	AnsweredAt       *time.Time `json:"answered_at,omitempty"`
	CorrectAnswer    *string    `json:"correct_answer,omitempty"` // Only shown if user has answered
	TimesAttempted   int        `json:"times_attempted"`
}

type ChapterQuizWithProgress struct {
	ChapterID        uint                          `json:"chapter_id"`
	ChapterTitle     string                        `json:"chapter_title"`
	TotalQuestions   int                           `json:"total_questions"`
	QuestionsAnswered int                          `json:"questions_answered"`
	CorrectAnswers   int                           `json:"correct_answers"`
	Questions        []QuizQuestionWithUserAnswer  `json:"questions"`
}

// GetChapterQuizWithHistoryRaw - Get quiz questions with user's answer history
func GetChapterQuizWithHistory(c *gin.Context) {
	chapterID := c.Param("id")
	userID := c.Query("user_id") // Required query parameter

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "user_id query parameter is required",
		})
		return
	}

	sqlDB, _ := database.DB.DB()

	// Get chapter info
	var chapterTitle string
	chapterQuery := `SELECT title FROM chapters WHERE id = $1 AND deleted_at IS NULL`
	err := sqlDB.QueryRow(chapterQuery, chapterID).Scan(&chapterTitle)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Chapter not found",
		})
		return
	}

	// Get all quiz questions with user's latest answer (if exists)
	query := `
		SELECT
			qq.id, qq.chapter_id, qq.question_text,
			qq.option_a, qq.option_b, qq.option_c, qq.option_d,
			qq.order_index, qq.created_at, qq.updated_at,
			qa.user_answer, qa.is_correct, qa.answered_at, qq.correct_answer,
			COALESCE((SELECT COUNT(*) FROM quiz_answers
					  WHERE quiz_question_id = qq.id AND user_id = $2 AND deleted_at IS NULL), 0) as times_attempted
		FROM quiz_questions qq
		LEFT JOIN LATERAL (
			SELECT user_answer, is_correct, answered_at
			FROM quiz_answers
			WHERE quiz_question_id = qq.id
			AND user_id = $2
			AND deleted_at IS NULL
			ORDER BY answered_at DESC
			LIMIT 1
		) qa ON true
		WHERE qq.chapter_id = $1 AND qq.deleted_at IS NULL
		ORDER BY qq.order_index ASC
	`

	rows, err := sqlDB.Query(query, chapterID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz questions",
		})
		return
	}
	defer rows.Close()

	var questions []QuizQuestionWithUserAnswer
	var questionsAnswered, correctAnswers int

	for rows.Next() {
		var q QuizQuestionWithUserAnswer
		var userAnswer, correctAnswer sql.NullString
		var isCorrect sql.NullBool
		var answeredAt sql.NullTime

		err := rows.Scan(
			&q.ID, &q.ChapterID, &q.QuestionText,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
			&q.OrderIndex, &q.CreatedAt, &q.UpdatedAt,
			&userAnswer, &isCorrect, &answeredAt, &correctAnswer,
			&q.TimesAttempted,
		)
		if err != nil {
			continue
		}

		// Check if user has answered this question
		if userAnswer.Valid {
			q.HasAnswered = true
			q.UserAnswer = &userAnswer.String
			q.IsCorrect = &isCorrect.Bool
			q.AnsweredAt = &answeredAt.Time
			q.CorrectAnswer = &correctAnswer.String

			questionsAnswered++
			if isCorrect.Bool {
				correctAnswers++
			}
		} else {
			q.HasAnswered = false
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

	var chapterIDUint uint
	sqlDB.QueryRow("SELECT $1::integer", chapterID).Scan(&chapterIDUint)

	result := ChapterQuizWithProgress{
		ChapterID:         chapterIDUint,
		ChapterTitle:      chapterTitle,
		TotalQuestions:    len(questions),
		QuestionsAnswered: questionsAnswered,
		CorrectAnswers:    correctAnswers,
		Questions:         questions,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"quiz":    result,
	})
}

// GetQuizResumePointRaw - Get where user should resume in a quiz
func GetQuizResumePoint(c *gin.Context) {
	userID := c.Param("userId")
	chapterID := c.Param("chapterId")

	sqlDB, _ := database.DB.DB()

	// Find the first unanswered question
	query := `
		SELECT qq.id, qq.order_index, qq.question_text
		FROM quiz_questions qq
		WHERE qq.chapter_id = $1
		AND qq.deleted_at IS NULL
		AND NOT EXISTS (
			SELECT 1 FROM quiz_answers qa
			WHERE qa.quiz_question_id = qq.id
			AND qa.user_id = $2
			AND qa.deleted_at IS NULL
		)
		ORDER BY qq.order_index ASC
		LIMIT 1
	`

	var questionID uint
	var orderIndex int
	var questionText string

	err := sqlDB.QueryRow(query, chapterID, userID).Scan(&questionID, &orderIndex, &questionText)

	if err == sql.ErrNoRows {
		// All questions answered
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"completed": true,
			"message": "All quiz questions completed",
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
		"completed": false,
		"resume_point": gin.H{
			"question_id": questionID,
			"order_index": orderIndex,
			"question_text": questionText,
		},
	})
}

