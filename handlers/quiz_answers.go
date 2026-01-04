package handlers

import (
	"database/sql"
	"learning-app-backend/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitQuizAnswerRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	ChapterID      uint   `json:"chapter_id" binding:"required"`
	QuizQuestionID uint   `json:"quiz_question_id" binding:"required"`
	UserAnswer     string `json:"user_answer" binding:"required"`
}

type QuizAnswer struct {
	ID             uint      `json:"id"`
	UserID         string    `json:"user_id"`
	ChapterID      uint      `json:"chapter_id"`
	QuizQuestionID uint      `json:"quiz_question_id"`
	UserAnswer     string    `json:"user_answer"`
	IsCorrect      bool      `json:"is_correct"`
	AnsweredAt     time.Time `json:"answered_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type QuizAnswerWithDetails struct {
	QuizAnswer
	QuestionText  string `json:"question_text"`
	CorrectAnswer string `json:"correct_answer"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
}

type QuizHistorySummary struct {
	ChapterID     uint   `json:"chapter_id"`
	ChapterTitle  string `json:"chapter_title"`
	TotalAnswered int    `json:"total_answered"`
	TotalCorrect  int    `json:"total_correct"`
	TotalWrong    int    `json:"total_wrong"`
	Percentage    float64 `json:"percentage"`
}

// SubmitQuizAnswerRaw - Submit a quiz answer and save to history
func SubmitQuizAnswer(c *gin.Context) {
	var req SubmitQuizAnswerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate answer
	if req.UserAnswer != "A" && req.UserAnswer != "B" && req.UserAnswer != "C" && req.UserAnswer != "D" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid answer. Must be A, B, C, or D",
		})
		return
	}

	sqlDB, _ := database.DB.DB()

	// Get the correct answer from the question
	var correctAnswer string
	questionQuery := `SELECT correct_answer FROM quiz_questions
					  WHERE id = $1 AND chapter_id = $2 AND deleted_at IS NULL`

	err := sqlDB.QueryRow(questionQuery, req.QuizQuestionID, req.ChapterID).Scan(&correctAnswer)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Quiz question not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	// Check if answer is correct
	isCorrect := req.UserAnswer == correctAnswer

	// Save the answer to history
	var answer QuizAnswer
	insertQuery := `INSERT INTO quiz_answers (user_id, chapter_id, quiz_question_id, user_answer,
				   is_correct, answered_at, created_at, updated_at)
				   VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), NOW())
				   RETURNING id, user_id, chapter_id, quiz_question_id, user_answer,
				   is_correct, answered_at, created_at, updated_at`

	err = sqlDB.QueryRow(insertQuery, req.UserID, req.ChapterID, req.QuizQuestionID,
		req.UserAnswer, isCorrect).Scan(
		&answer.ID, &answer.UserID, &answer.ChapterID, &answer.QuizQuestionID,
		&answer.UserAnswer, &answer.IsCorrect, &answer.AnsweredAt,
		&answer.CreatedAt, &answer.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save answer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "Answer submitted successfully",
		"is_correct":     isCorrect,
		"correct_answer": correctAnswer,
		"answer":         answer,
	})
}

// GetQuizHistoryRaw - Get all quiz answers for a user in a specific chapter
func GetQuizHistory(c *gin.Context) {
	userID := c.Param("userId")
	chapterID := c.Param("chapterId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT qa.id, qa.user_id, qa.chapter_id, qa.quiz_question_id, qa.user_answer,
			  qa.is_correct, qa.answered_at, qa.created_at, qa.updated_at,
			  qq.question_text, qq.correct_answer, qq.option_a, qq.option_b, qq.option_c, qq.option_d
			  FROM quiz_answers qa
			  JOIN quiz_questions qq ON qa.quiz_question_id = qq.id
			  WHERE qa.user_id = $1 AND qa.chapter_id = $2 AND qa.deleted_at IS NULL
			  ORDER BY qa.answered_at DESC`

	rows, err := sqlDB.Query(query, userID, chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz history",
		})
		return
	}
	defer rows.Close()

	var answers []QuizAnswerWithDetails
	for rows.Next() {
		var a QuizAnswerWithDetails
		err := rows.Scan(&a.ID, &a.UserID, &a.ChapterID, &a.QuizQuestionID, &a.UserAnswer,
			&a.IsCorrect, &a.AnsweredAt, &a.CreatedAt, &a.UpdatedAt,
			&a.QuestionText, &a.CorrectAnswer, &a.OptionA, &a.OptionB, &a.OptionC, &a.OptionD)
		if err == nil {
			answers = append(answers, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"answers": answers,
	})
}

// GetAllQuizHistoryRaw - Get all quiz answers for a user across all chapters
func GetAllQuizHistory(c *gin.Context) {
	userID := c.Param("userId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT qa.id, qa.user_id, qa.chapter_id, qa.quiz_question_id, qa.user_answer,
			  qa.is_correct, qa.answered_at, qa.created_at, qa.updated_at,
			  qq.question_text, qq.correct_answer, qq.option_a, qq.option_b, qq.option_c, qq.option_d
			  FROM quiz_answers qa
			  JOIN quiz_questions qq ON qa.quiz_question_id = qq.id
			  WHERE qa.user_id = $1 AND qa.deleted_at IS NULL
			  ORDER BY qa.chapter_id ASC, qa.answered_at DESC`

	rows, err := sqlDB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz history",
		})
		return
	}
	defer rows.Close()

	var answers []QuizAnswerWithDetails
	for rows.Next() {
		var a QuizAnswerWithDetails
		err := rows.Scan(&a.ID, &a.UserID, &a.ChapterID, &a.QuizQuestionID, &a.UserAnswer,
			&a.IsCorrect, &a.AnsweredAt, &a.CreatedAt, &a.UpdatedAt,
			&a.QuestionText, &a.CorrectAnswer, &a.OptionA, &a.OptionB, &a.OptionC, &a.OptionD)
		if err == nil {
			answers = append(answers, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"answers": answers,
	})
}

// GetQuizScoreRaw - Get quiz score summary for a user
func GetQuizScore(c *gin.Context) {
	userID := c.Param("userId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT qa.chapter_id, ch.title,
			  COUNT(*) as total_answered,
			  SUM(CASE WHEN qa.is_correct = true THEN 1 ELSE 0 END) as total_correct,
			  SUM(CASE WHEN qa.is_correct = false THEN 1 ELSE 0 END) as total_wrong
			  FROM quiz_answers qa
			  JOIN chapters ch ON qa.chapter_id = ch.id
			  WHERE qa.user_id = $1 AND qa.deleted_at IS NULL
			  GROUP BY qa.chapter_id, ch.title
			  ORDER BY qa.chapter_id ASC`

	rows, err := sqlDB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch quiz scores",
		})
		return
	}
	defer rows.Close()

	var summaries []QuizHistorySummary
	for rows.Next() {
		var s QuizHistorySummary
		err := rows.Scan(&s.ChapterID, &s.ChapterTitle, &s.TotalAnswered, &s.TotalCorrect, &s.TotalWrong)
		if err == nil {
			if s.TotalAnswered > 0 {
				s.Percentage = (float64(s.TotalCorrect) / float64(s.TotalAnswered)) * 100
			}
			summaries = append(summaries, s)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"scores":  summaries,
	})
}

// GetQuestionAnswerHistoryRaw - Get all attempts for a specific question by a user
func GetQuestionAnswerHistory(c *gin.Context) {
	userID := c.Param("userId")
	questionID := c.Param("questionId")
	sqlDB, _ := database.DB.DB()

	query := `SELECT qa.id, qa.user_id, qa.chapter_id, qa.quiz_question_id, qa.user_answer,
			  qa.is_correct, qa.answered_at, qa.created_at, qa.updated_at,
			  qq.question_text, qq.correct_answer, qq.option_a, qq.option_b, qq.option_c, qq.option_d
			  FROM quiz_answers qa
			  JOIN quiz_questions qq ON qa.quiz_question_id = qq.id
			  WHERE qa.user_id = $1 AND qa.quiz_question_id = $2 AND qa.deleted_at IS NULL
			  ORDER BY qa.answered_at DESC`

	rows, err := sqlDB.Query(query, userID, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch answer history",
		})
		return
	}
	defer rows.Close()

	var answers []QuizAnswerWithDetails
	for rows.Next() {
		var a QuizAnswerWithDetails
		err := rows.Scan(&a.ID, &a.UserID, &a.ChapterID, &a.QuizQuestionID, &a.UserAnswer,
			&a.IsCorrect, &a.AnsweredAt, &a.CreatedAt, &a.UpdatedAt,
			&a.QuestionText, &a.CorrectAnswer, &a.OptionA, &a.OptionB, &a.OptionC, &a.OptionD)
		if err == nil {
			answers = append(answers, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"answers": answers,
	})
}

// ClearQuizHistoryRaw - Clear quiz history for a user (optional)
func ClearQuizHistory(c *gin.Context) {
	userID := c.Param("userId")
	chapterID := c.Query("chapter_id") // Optional: clear only for specific chapter
	sqlDB, _ := database.DB.DB()

	var query string
	var result sql.Result
	var err error

	if chapterID != "" {
		// Clear for specific chapter
		query = `UPDATE quiz_answers SET deleted_at = NOW()
				 WHERE user_id = $1 AND chapter_id = $2 AND deleted_at IS NULL`
		result, err = sqlDB.Exec(query, userID, chapterID)
	} else {
		// Clear all history
		query = `UPDATE quiz_answers SET deleted_at = NOW()
				 WHERE user_id = $1 AND deleted_at IS NULL`
		result, err = sqlDB.Exec(query, userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to clear quiz history",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quiz history cleared successfully",
		"deleted": rowsAffected,
	})
}
