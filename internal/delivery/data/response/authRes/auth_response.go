package authres

import (
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
)

type UserProfile struct {
	ID          string          `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	AvatarURL   string          `json:"avatarUrl"`
	Role        models.UserRole `json:"role"`
	LastLoginAt *time.Time      `json:"lastLoginAt,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type AuthResponse struct {
	User   UserProfile   `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}
