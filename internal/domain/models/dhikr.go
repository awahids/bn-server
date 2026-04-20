package models

import "time"

type Dhikr struct {
	ID              string    `gorm:"primaryKey;size:191" json:"id"`
	Arabic          string    `gorm:"type:text;not null" json:"arabic"`
	Transliteration string    `gorm:"type:text;not null" json:"transliteration"`
	Translation     string    `gorm:"type:text;not null" json:"translation"`
	Meaning         string    `gorm:"type:text;not null" json:"meaning"`
	Count           int       `gorm:"not null;default:1" json:"count"`
	Session         string    `gorm:"size:20;not null" json:"session"`
	Category        string    `gorm:"size:100;not null" json:"category"`
	AudioUrl        *string   `gorm:"size:255" json:"audioUrl,omitempty"`
	Reference       *string   `gorm:"type:text" json:"reference,omitempty"`
	Faedah          *string   `gorm:"type:text" json:"faedah,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
