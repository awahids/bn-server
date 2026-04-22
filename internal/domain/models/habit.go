package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Habit struct {
	ID              string            `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          string            `gorm:"type:uuid;index;not null" json:"userId"`
	Name            string            `gorm:"size:191;not null" json:"name"`
	Category        string            `gorm:"size:50;not null;default:'Other'" json:"category"`
	ReminderTime    string            `gorm:"column:reminder_time;size:5;not null;default:''" json:"reminderTime"`
	ReminderEnabled bool              `gorm:"column:reminder_enabled;not null;default:false" json:"reminderEnabled"`
	User            User              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Completions     []HabitCompletion `gorm:"foreignKey:HabitID" json:"-"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (h *Habit) BeforeCreate(_ *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
	return nil
}

type HabitCompletion struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;index;not null" json:"userId"`
	HabitID   string         `gorm:"type:uuid;index;not null" json:"habitId"`
	Date      string         `gorm:"size:10;index;not null" json:"date"`
	Completed bool           `gorm:"not null;default:true" json:"completed"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Habit     Habit          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HabitCompletion) TableName() string {
	return "habit_completions"
}

func (h *HabitCompletion) BeforeCreate(_ *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
	return nil
}
