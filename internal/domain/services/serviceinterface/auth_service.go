package serviceinterface

import (
	"context"

	"github.com/awahids/bn-server/internal/domain/models"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type AuthService interface {
	LoginWithGoogle(ctx context.Context, idToken string) (*models.User, *TokenPair, error)
	LoginWithGoogleOAuthCode(ctx context.Context, code, redirectURI string) (*models.User, *TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetCurrentUser(ctx context.Context, userID string) (*models.User, error)
}
