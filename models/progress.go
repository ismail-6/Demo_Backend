package models

import (
	"time"
	"gorm.io/gorm"
)

type ContentType string

const (
	ContentTypeVideo ContentType = "video"
	ContentTypeQuiz  ContentType = "quiz"
)

type Progress struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            string         `gorm:"not null;index:idx_user_chapter" json:"user_id"`
	ChapterID         uint           `gorm:"not null;index:idx_user_chapter" json:"chapter_id"`
	ContentType       ContentType    `gorm:"type:varchar(10);not null" json:"content_type"`
	VideoTimestamp    *int           `json:"video_timestamp,omitempty"`     // nullable, in seconds
	QuizQuestionIndex *int           `json:"quiz_question_index,omitempty"` // nullable, 0-based index
	IsCompleted       bool           `gorm:"default:false" json:"is_completed"`
	LastUpdated       time.Time      `gorm:"autoUpdateTime" json:"last_updated"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type SaveProgressRequest struct {
	UserID            string      `json:"user_id" binding:"required"`
	ChapterID         uint        `json:"chapter_id" binding:"required"`
	ContentType       ContentType `json:"content_type" binding:"required"`
	VideoTimestamp    *int        `json:"video_timestamp"`
	QuizQuestionIndex *int        `json:"quiz_question_index"`
	IsCompleted       bool        `json:"is_completed"`
}

type ProgressResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Progress *Progress `json:"progress,omitempty"`
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
