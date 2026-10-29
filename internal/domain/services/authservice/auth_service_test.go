package authservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/pkg/utils"
)

type mockAuthRepo struct {
	findRefreshTokenResult *models.RefreshToken
	findRefreshTokenErr    error
	rotateResult           bool
	rotateErr              error
	revokedTokenHash       string

	lastOldTokenHash string
	lastNewTokenHash string
	lastNewExpiresAt time.Time

	findUserResult *models.User
	findUserErr    error
}

func (m *mockAuthRepo) FindUserByGoogleID(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (m *mockAuthRepo) FindUserByEmail(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (m *mockAuthRepo) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	_ = ctx
	_ = userID
	return m.findUserResult, m.findUserErr
}

func (m *mockAuthRepo) CreateUser(context.Context, *models.User) error {
	return nil
}

func (m *mockAuthRepo) UpdateUser(context.Context, *models.User) error {
	return nil
}

func (m *mockAuthRepo) SaveRefreshToken(context.Context, *models.RefreshToken) error {
	return nil
}

func (m *mockAuthRepo) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	_ = ctx
	m.lastOldTokenHash = tokenHash
	return m.findRefreshTokenResult, m.findRefreshTokenErr
}

func (m *mockAuthRepo) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	_ = ctx
	m.revokedTokenHash = tokenHash
	return nil
}

func (m *mockAuthRepo) RotateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash string, newExpiresAt time.Time) (bool, error) {
	_ = ctx
	m.lastOldTokenHash = oldTokenHash
	m.lastNewTokenHash = newTokenHash
	m.lastNewExpiresAt = newExpiresAt
	return m.rotateResult, m.rotateErr
}

func (m *mockAuthRepo) RevokeRefreshTokensByUserID(context.Context, string) error {
	return nil
}

func TestRefreshToken_Success(t *testing.T) {
	repo := &mockAuthRepo{
		findRefreshTokenResult: &models.RefreshToken{
			User: models.User{
				ID:    "user-123",
				Email: "user@example.com",
				Role:  models.RoleUser,
			},
		},
		rotateResult: true,
	}

	svc := NewAuthService(repo, TokenConfig{
		Issuer:          "test-issuer",
		Secret:          "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}, nil)

	pair, err := svc.RefreshToken(context.Background(), "old-refresh-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pair == nil {
		t.Fatal("expected token pair, got nil")
	}
	if pair.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if pair.RefreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
	if pair.TokenType != "Bearer" {
		t.Fatalf("expected token type Bearer, got %s", pair.TokenType)
	}
	if pair.ExpiresIn <= 0 {
		t.Fatalf("expected positive expiresIn, got %d", pair.ExpiresIn)
	}

	expectedOldHash := utils.HashToken("old-refresh-token")
	if repo.lastOldTokenHash != expectedOldHash {
		t.Fatalf("expected old token hash %s, got %s", expectedOldHash, repo.lastOldTokenHash)
	}
	if repo.lastNewTokenHash == "" {
		t.Fatal("expected new token hash to be set")
	}
	if repo.lastNewExpiresAt.Before(time.Now()) {
		t.Fatal("expected new expires at to be in the future")
	}
}

func TestRefreshToken_InvalidToken_NotFound(t *testing.T) {
	repo := &mockAuthRepo{}

	svc := NewAuthService(repo, TokenConfig{Issuer: "test", Secret: "secret", AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour}, nil)

	pair, err := svc.RefreshToken(context.Background(), "missing")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
	if pair != nil {
		t.Fatal("expected nil pair for invalid token")
	}
}

func TestRefreshToken_InvalidToken_RotationFailed(t *testing.T) {
	repo := &mockAuthRepo{
		findRefreshTokenResult: &models.RefreshToken{
			User: models.User{ID: "user-1", Email: "user@example.com", Role: models.RoleUser},
		},
		rotateResult: false,
	}

	svc := NewAuthService(repo, TokenConfig{Issuer: "test", Secret: "secret", AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour}, nil)

	pair, err := svc.RefreshToken(context.Background(), "invalid-after-check")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
	if pair != nil {
		t.Fatal("expected nil pair when rotation fails")
	}
}

func TestLogout_WithRefreshToken_RevokesToken(t *testing.T) {
	repo := &mockAuthRepo{}

	svc := NewAuthService(repo, TokenConfig{Issuer: "test", Secret: "secret", AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour}, nil)

	err := svc.Logout(context.Background(), "raw-refresh-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedHash := utils.HashToken("raw-refresh-token")
	if repo.revokedTokenHash != expectedHash {
		t.Fatalf("expected revoked hash %s, got %s", expectedHash, repo.revokedTokenHash)
	}
}

func TestGetCurrentUser_NotFound(t *testing.T) {
	repo := &mockAuthRepo{findUserResult: nil}

	svc := NewAuthService(repo, TokenConfig{Issuer: "test", Secret: "secret", AccessTokenTTL: time.Minute, RefreshTokenTTL: time.Hour}, nil)

	user, err := svc.GetCurrentUser(context.Background(), "missing")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if user != nil {
		t.Fatal("expected nil user")
	}
}
