package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgressModule string

const (
	ModuleHijaiyah ProgressModule = "hijaiyah"
	ModuleQuran    ProgressModule = "quran"
	ModuleDhikr    ProgressModule = "dhikr"
	ModuleQuiz     ProgressModule = "quiz"
)

type UserProgress struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       string         `gorm:"type:uuid;index;not null" json:"userId"`
	Module       string         `gorm:"size:30;index;not null" json:"module"`
	ItemID       string         `gorm:"column:item_id;size:191;index;not null" json:"itemId"`
	Progress     int            `gorm:"not null;default:0" json:"progress"`
	Completed    bool           `gorm:"not null;default:false" json:"completed"`
	Score        int            `gorm:"not null;default:0" json:"score"`
	TimeSpent    int            `gorm:"column:time_spent;not null;default:0" json:"timeSpent"`
	LastAccessed time.Time      `gorm:"column:last_accessed;index;not null" json:"lastAccessed"`
	User         User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserProgress) TableName() string {
	return "user_progress"
}

func (u *UserProgress) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.LastAccessed.IsZero() {
		u.LastAccessed = time.Now()
	}
	return nil
}
