package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookmarkType string

const (
	BookmarkTypeQuran BookmarkType = "quran"
	BookmarkTypeDhikr BookmarkType = "dhikr"
)

type Bookmark struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;index;not null" json:"userId"`
	Type      string         `gorm:"size:20;index;not null" json:"type"`
	ContentID string         `gorm:"column:content_id;size:191;index;not null" json:"contentId"`
	Note      *string        `gorm:"size:500" json:"note,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;index;not null" json:"createdAt"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Bookmark) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	return nil
}
