package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DhikrSession string

const (
	DhikrSessionMorning DhikrSession = "morning"
	DhikrSessionEvening DhikrSession = "evening"
)

type DhikrCounter struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;index;not null" json:"userId"`
	DhikrID   string         `gorm:"column:dhikr_id;size:191;index;not null" json:"dhikrId"`
	Count     int            `gorm:"not null;default:0" json:"count"`
	Target    int            `gorm:"not null;default:33" json:"target"`
	Date      string         `gorm:"size:10;index;not null" json:"date"`
	Session   string         `gorm:"size:20;index;not null" json:"session"`
	Completed bool           `gorm:"not null;default:false" json:"completed"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DhikrCounter) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
