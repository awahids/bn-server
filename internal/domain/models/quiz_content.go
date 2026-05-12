package models

import (
	"encoding/json"
	"time"
)

type QuizCategory struct {
	ID          string    `gorm:"primaryKey;size:50" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Icon        string    `gorm:"size:50;not null" json:"icon"`
	Color       string    `gorm:"size:50;not null" json:"color"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type QuizQuestion struct {
	ID            string          `gorm:"primaryKey;size:50" json:"id"`
	Question      string          `gorm:"type:text;not null" json:"question"`
	Options       json.RawMessage `gorm:"type:jsonb;not null" json:"options"`
	CorrectAnswer int             `gorm:"not null" json:"correctAnswer"`
	Explanation   string          `gorm:"type:text;not null" json:"explanation"`
	Material      string          `gorm:"type:text;not null" json:"material"`
	CategoryID    string          `gorm:"size:50;not null" json:"category"`
	Difficulty    string          `gorm:"size:10;not null" json:"difficulty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}
