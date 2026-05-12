package models

import (
	"encoding/json"
	"time"
)

type HijaiyahLetter struct {
	ID              string          `gorm:"primaryKey;size:50" json:"id"`
	Arabic          string          `gorm:"size:10;not null" json:"arabic"`
	Name            string          `gorm:"size:50;not null" json:"name"`
	Transliteration string          `gorm:"size:10;not null" json:"transliteration"`
	Pronunciation   string          `gorm:"size:50;not null" json:"pronunciation"`
	AudioUrl        *string         `gorm:"size:500" json:"audioUrl,omitempty"`
	SortOrder       int             `gorm:"not null;default:0" json:"order"`
	Description     string          `gorm:"type:text" json:"description"`
	WritingSteps    json.RawMessage `gorm:"type:jsonb" json:"writingSteps"`
	StrokePoints    json.RawMessage `gorm:"type:jsonb" json:"strokePoints"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}
