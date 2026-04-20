package authservice

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/repositories/repointerface"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
	"github.com/awahids/bn-server/pkg/utils"
)

var (
	ErrInvalidGoogleToken       = errors.New("invalid google token")
	ErrInvalidGoogleOAuthCode   = errors.New("invalid google oauth code")
	ErrGoogleOAuthNotConfigured = errors.New("google oauth is not configured")
	ErrEmailNotVerified         = errors.New("google email is not verified")
	ErrInvalidRefreshToken      = errors.New("invalid refresh token")
	ErrUserNotFound             = errors.New("user not found")
)

type TokenConfig struct {
	Issuer          string
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type authService struct {
	repo               repointerface.AuthRepository
	tokenCfg           TokenConfig
	googleAuthProvider serviceinterface.GoogleAuthProvider
}

func NewAuthService(
	repo repointerface.AuthRepository,
	tokenCfg TokenConfig,
	googleAuthProvider serviceinterface.GoogleAuthProvider,
) serviceinterface.AuthService {
	return &authService{
		repo:               repo,
		tokenCfg:           tokenCfg,
		googleAuthProvider: googleAuthProvider,
	}
}

func (s *authService) LoginWithGoogle(ctx context.Context, idTokenRaw string) (*models.User, *serviceinterface.TokenPair, error) {
	if strings.TrimSpace(idTokenRaw) == "" {
		return nil, nil, ErrInvalidGoogleToken
	}
	if s.googleAuthProvider == nil {
		return nil, nil, ErrGoogleOAuthNotConfigured
	}

	tokenInfo, err := s.googleAuthProvider.GetTokenInfoByIDToken(ctx, idTokenRaw)
	if err != nil {
		return nil, nil, mapGoogleAuthProviderError(err)
	}

	return s.loginWithGoogleTokenInfo(ctx, tokenInfo)
}

func (s *authService) LoginWithGoogleOAuthCode(
	ctx context.Context,
	code, redirectURI string,
) (*models.User, *serviceinterface.TokenPair, error) {
	if strings.TrimSpace(code) == "" {
		return nil, nil, ErrInvalidGoogleOAuthCode
	}
	if s.googleAuthProvider == nil {
		return nil, nil, ErrGoogleOAuthNotConfigured
	}

	tokenInfo, err := s.googleAuthProvider.GetTokenInfoByOAuthCode(ctx, code, redirectURI)
	if err != nil {
		return nil, nil, mapGoogleAuthProviderError(err)
	}

	return s.loginWithGoogleTokenInfo(ctx, tokenInfo)
}

func (s *authService) loginWithGoogleTokenInfo(
	ctx context.Context,
	tokenInfo *serviceinterface.GoogleTokenInfo,
) (*models.User, *serviceinterface.TokenPair, error) {
	email := strings.ToLower(strings.TrimSpace(tokenInfo.Email))
	if email == "" {
		return nil, nil, ErrInvalidGoogleToken
	}

	if !tokenInfo.EmailVerified {
		return nil, nil, ErrEmailNotVerified
	}

	googleID := strings.TrimSpace(tokenInfo.Subject)
	if googleID == "" {
		return nil, nil, ErrInvalidGoogleToken
	}

	name := strings.TrimSpace(tokenInfo.Name)
	avatar := strings.TrimSpace(tokenInfo.Picture)

	user, err := s.repo.FindUserByGoogleID(ctx, googleID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		user, err = s.repo.FindUserByEmail(ctx, email)
		if err != nil {
			return nil, nil, err
		}
	}

	now := time.Now()
	if user == nil {
		user = &models.User{
			GoogleID:    googleID,
			Email:       email,
			Name:        fallbackName(name, email),
			AvatarURL:   avatar,
			Role:        models.RoleUser,
			LastLoginAt: &now,
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, nil, err
		}
	} else {
		if user.GoogleID == "" {
			user.GoogleID = googleID
		}
		if name != "" {
			user.Name = name
		}
		if avatar != "" {
			user.AvatarURL = avatar
		}
		user.LastLoginAt = &now

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return nil, nil, err
		}
	}

	tokenPair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*serviceinterface.TokenPair, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrInvalidRefreshToken
	}

	hash := utils.HashToken(refreshToken)
	storedToken, err := s.repo.FindRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil || storedToken.User.ID == "" {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, expiresIn, err := utils.GenerateAccessToken(
		s.tokenCfg.Secret,
		s.tokenCfg.Issuer,
		storedToken.User.ID,
		storedToken.User.Email,
		string(storedToken.User.Role),
		s.tokenCfg.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	newRefreshRaw, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshHash := utils.HashToken(newRefreshRaw)
	newExpiresAt := time.Now().Add(s.tokenCfg.RefreshTokenTTL)

	rotated, err := s.repo.RotateRefreshToken(ctx, hash, newRefreshHash, newExpiresAt)
	if err != nil {
		return nil, err
	}
	if !rotated {
		return nil, ErrInvalidRefreshToken
	}

	return &serviceinterface.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshRaw,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}

	hash := utils.HashToken(refreshToken)
	return s.repo.RevokeRefreshTokenByHash(ctx, hash)
}

func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *authService) issueTokenPair(ctx context.Context, user *models.User) (*serviceinterface.TokenPair, error) {
	accessToken, expiresIn, err := utils.GenerateAccessToken(
		s.tokenCfg.Secret,
		s.tokenCfg.Issuer,
		user.ID,
		user.Email,
		string(user.Role),
		s.tokenCfg.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	rawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.tokenCfg.RefreshTokenTTL)
	refreshEntity := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: utils.HashToken(rawRefreshToken),
		ExpiresAt: expiresAt,
	}
	if err := s.repo.SaveRefreshToken(ctx, refreshEntity); err != nil {
		return nil, err
	}

	return &serviceinterface.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func fallbackName(name, email string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	localPart := strings.Split(email, "@")
	if len(localPart) == 0 || localPart[0] == "" {
		return "User"
	}
	return localPart[0]
}

func mapGoogleAuthProviderError(err error) error {
	switch {
	case errors.Is(err, serviceinterface.ErrGoogleAuthInvalidIDToken):
		return ErrInvalidGoogleToken
	case errors.Is(err, serviceinterface.ErrGoogleAuthInvalidOAuthCode):
		return ErrInvalidGoogleOAuthCode
	case errors.Is(err, serviceinterface.ErrGoogleAuthNotConfigured):
		return ErrGoogleOAuthNotConfigured
	default:
		return err
	}
}
