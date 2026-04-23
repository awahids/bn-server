package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushSubscription struct {
	ID             string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         string         `gorm:"type:uuid;index;not null" json:"userId"`
	Endpoint       string         `gorm:"size:2048;uniqueIndex;not null" json:"endpoint"`
	P256DH         string         `gorm:"column:p256dh;size:512;not null" json:"-"`
	Auth           string         `gorm:"column:auth;size:255;not null" json:"-"`
	ExpirationTime *int64         `gorm:"column:expiration_time" json:"expirationTime,omitempty"`
	Timezone       string         `gorm:"size:100;not null;default:'UTC'" json:"timezone"`
	Enabled        bool           `gorm:"not null;default:true" json:"enabled"`
	LastSeenAt     time.Time      `gorm:"column:last_seen_at;not null;default:CURRENT_TIMESTAMP" json:"lastSeenAt"`
	User           User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type PushReminderTarget struct {
	SubscriptionID string `json:"-"`
	Endpoint       string `json:"endpoint"`
	P256DH         string `json:"-"`
	Auth           string `json:"-"`
	Timezone       string `json:"timezone"`
	UserID         string `json:"userId"`
	HabitID        string `json:"habitId"`
	HabitName      string `json:"habitName"`
	ReminderTime   string `json:"reminderTime"`
}

func (p *PushSubscription) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.LastSeenAt.IsZero() {
		p.LastSeenAt = time.Now()
	}
	return nil
}
