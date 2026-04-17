package authrepo

import (
	"context"
	"errors"
	"time"

	"bn-mobile/server/internal/domain/models"
	"bn-mobile/server/internal/domain/repositories/repoInterface"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repointerface.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *authRepository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Update("revoked_at", &now).Error
}

func (r *authRepository) RotateRefreshToken(
	ctx context.Context,
	oldTokenHash, newTokenHash string,
	newExpiresAt time.Time,
) (bool, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldToken models.RefreshToken
		findErr := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", oldTokenHash).
			Where("revoked_at IS NULL").
			Where("expires_at > ?", time.Now()).
			First(&oldToken).Error
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return findErr
		}

		now := time.Now()
		if revokeErr := tx.
			Model(&models.RefreshToken{}).
			Where("id = ?", oldToken.ID).
			Where("revoked_at IS NULL").
			Update("revoked_at", &now).Error; revokeErr != nil {
			return revokeErr
		}

		newToken := &models.RefreshToken{
			UserID:    oldToken.UserID,
			TokenHash: newTokenHash,
			ExpiresAt: newExpiresAt,
		}

		return tx.Create(newToken).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *authRepository) RevokeRefreshTokensByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Update("revoked_at", &now).Error
}
