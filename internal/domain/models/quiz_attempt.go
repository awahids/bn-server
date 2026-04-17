package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizAttempt struct {
	ID             string          `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         string          `gorm:"type:uuid;index;not null" json:"userId"`
	Category       string          `gorm:"size:100;index;not null" json:"category"`
	Score          int             `gorm:"not null" json:"score"`
	TotalQuestions int             `gorm:"column:total_questions;not null" json:"totalQuestions"`
	TimeSpent      int             `gorm:"column:time_spent;not null" json:"timeSpent"`
	Answers        json.RawMessage `gorm:"type:jsonb;not null" json:"answers"`
	CompletedAt    time.Time       `gorm:"column:completed_at;index;not null" json:"completedAt"`
	User           User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (q *QuizAttempt) BeforeCreate(_ *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	if q.CompletedAt.IsZero() {
		q.CompletedAt = time.Now()
	}
	if len(q.Answers) == 0 {
		q.Answers = json.RawMessage("[]")
	}
	return nil
}
