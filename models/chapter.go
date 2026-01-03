package models

import (
	"time"
	"gorm.io/gorm"
)

type Chapter struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	OrderIndex  int            `gorm:"not null" json:"order_index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Video         Video          `gorm:"foreignKey:ChapterID" json:"video,omitempty"`
	QuizQuestions []QuizQuestion `gorm:"foreignKey:ChapterID" json:"quiz_questions,omitempty"`
}

type Video struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ChapterID       uint           `gorm:"not null" json:"chapter_id"`
	Title           string         `gorm:"not null" json:"title"`
	VideoURL        string         `gorm:"not null" json:"video_url"`
	DurationSeconds int            `gorm:"not null" json:"duration_seconds"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type QuizQuestion struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ChapterID     uint           `gorm:"not null" json:"chapter_id"`
	QuestionText  string         `gorm:"not null" json:"question_text"`
	OptionA       string         `gorm:"not null" json:"option_a"`
	OptionB       string         `gorm:"not null" json:"option_b"`
	OptionC       string         `gorm:"not null" json:"option_c"`
	OptionD       string         `gorm:"not null" json:"option_d"`
	CorrectAnswer string         `gorm:"not null" json:"correct_answer"` // A, B, C, or D
	OrderIndex    int            `gorm:"not null" json:"order_index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
