package repointerface

import (
	"context"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
)

type AuthRepository interface {
	FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error

	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error
	RotateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash string, newExpiresAt time.Time) (bool, error)
	RevokeRefreshTokensByUserID(ctx context.Context, userID string) error
}
