package models

import "time"

type QuranSurah struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	ArabicName     string    `gorm:"size:200;not null" json:"arabicName"`
	EnglishName    string    `gorm:"size:200;not null" json:"englishName"`
	NumberOfAyahs  int       `gorm:"not null" json:"numberOfAyahs"`
	RevelationType string    `gorm:"size:10;not null" json:"revelationType"`
	AudioUrl       *string   `gorm:"size:500" json:"audioUrl,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
