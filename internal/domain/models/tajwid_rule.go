package models

import (
	"encoding/json"
	"time"
)

type TajwidRule struct {
	ID             string          `gorm:"primaryKey;size:50"  json:"id"`
	Name           string          `gorm:"size:100;not null"   json:"name"`
	ArabicName     string          `gorm:"size:200"            json:"arabicName"`
	Category       string          `gorm:"size:50;not null"    json:"category"`
	Description    string          `gorm:"type:text;not null"  json:"description"`
	TriggerLetters string          `gorm:"size:200"            json:"triggerLetters"`
	Examples       json.RawMessage `gorm:"type:jsonb"          json:"examples"`
	AudioUrl       *string         `gorm:"size:500"            json:"audioUrl,omitempty"`
	SortOrder      int             `gorm:"not null;default:0"  json:"sortOrder"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}
