package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type School struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        string         `gorm:"type:uuid;index;not null" json:"userId"`
	Name          string         `gorm:"size:191;not null" json:"name"`
	Location      string         `gorm:"size:255;not null" json:"location"`
	Jenjang       string         `gorm:"size:32;not null;default:'Lainnya'" json:"jenjang"`
	StatusSekolah string         `gorm:"column:status_sekolah;size:20;not null;default:'swasta'" json:"statusSekolah"`
	MonthlyFee    int            `gorm:"column:monthly_fee;not null;default:0" json:"monthlyFee"`
	MapURL        string         `gorm:"column:map_url;size:1024;not null" json:"mapUrl"`
	Contact       string         `gorm:"size:100;not null;default:''" json:"contact"`
	Description   string         `gorm:"size:1000;not null;default:''" json:"description"`
	User          User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *School) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}
