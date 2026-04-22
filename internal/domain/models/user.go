package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID            string            `gorm:"type:uuid;primaryKey" json:"id"`
	GoogleID      string            `gorm:"size:191;uniqueIndex" json:"googleId"`
	Email         string            `gorm:"size:191;uniqueIndex;not null" json:"email"`
	Name          string            `gorm:"size:191;not null" json:"name"`
	AvatarURL     string            `gorm:"size:512" json:"avatarUrl"`
	Username      *string           `gorm:"size:50;uniqueIndex" json:"username,omitempty"`
	Role          UserRole          `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	Streak        int               `gorm:"not null;default:0" json:"streak"`
	DailyProgress int               `gorm:"column:daily_progress;not null;default:0" json:"dailyProgress"`
	LastActive    time.Time         `gorm:"column:last_active;not null;default:CURRENT_TIMESTAMP" json:"lastActive"`
	LastLoginAt   *time.Time        `gorm:"column:last_login_at" json:"lastLoginAt,omitempty"`
	Preferences   json.RawMessage   `gorm:"type:jsonb;not null;default:'{}'" json:"preferences"`
	RefreshTokens []RefreshToken    `gorm:"foreignKey:UserID" json:"-"`
	Progress      []UserProgress    `gorm:"foreignKey:UserID" json:"-"`
	Bookmarks     []Bookmark        `gorm:"foreignKey:UserID" json:"-"`
	Habits        []Habit           `gorm:"foreignKey:UserID" json:"-"`
	HabitLogs     []HabitCompletion `gorm:"foreignKey:UserID" json:"-"`
	DhikrCounters []DhikrCounter    `gorm:"foreignKey:UserID" json:"-"`
	QuizAttempts  []QuizAttempt     `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if len(u.Preferences) == 0 {
		u.Preferences = json.RawMessage("{}")
	}
	if u.LastActive.IsZero() {
		u.LastActive = time.Now()
	}
	return nil
}
