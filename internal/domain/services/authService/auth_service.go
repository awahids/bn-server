package authservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bn-mobile/server/configs"
	"bn-mobile/server/internal/domain/models"
	"bn-mobile/server/internal/domain/repositories/repoInterface"
	"bn-mobile/server/internal/domain/services/serviceInterface"
	"bn-mobile/server/pkg/utils"
)

var (
	ErrInvalidGoogleToken  = errors.New("invalid google token")
	ErrEmailNotVerified    = errors.New("google email is not verified")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type authService struct {
	repo repointerface.AuthRepository
	cfg  *configs.Config
}

func NewAuthService(repo repointerface.AuthRepository, cfg *configs.Config) serviceinterface.AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) LoginWithGoogle(ctx context.Context, idTokenRaw string) (*models.User, *serviceinterface.TokenPair, error) {
	if strings.TrimSpace(idTokenRaw) == "" {
		return nil, nil, ErrInvalidGoogleToken
	}
	if strings.TrimSpace(s.cfg.Google.ClientID) == "" {
		return nil, nil, errors.New("GOOGLE_CLIENT_ID is not configured")
	}

	tokenInfo, err := s.validateGoogleIDToken(ctx, idTokenRaw)
	if err != nil {
		return nil, nil, ErrInvalidGoogleToken
	}

	email := strings.ToLower(tokenInfo.Email)
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
		s.cfg.JWT.Secret,
		s.cfg.JWT.Issuer,
		storedToken.User.ID,
		storedToken.User.Email,
		string(storedToken.User.Role),
		s.cfg.JWT.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	newRefreshRaw, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshHash := utils.HashToken(newRefreshRaw)
	newExpiresAt := time.Now().Add(s.cfg.JWT.RefreshTokenTTL)

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
		s.cfg.JWT.Secret,
		s.cfg.JWT.Issuer,
		user.ID,
		user.Email,
		string(user.Role),
		s.cfg.JWT.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	rawRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.cfg.JWT.RefreshTokenTTL)
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

type googleTokenInfoResponse struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified,string"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	ExpiresIn     string `json:"expires_in"`
}

func (s *authService) validateGoogleIDToken(ctx context.Context, rawIDToken string) (*googleTokenInfoResponse, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(rawIDToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google tokeninfo responded with status %d", res.StatusCode)
	}

	var payload googleTokenInfoResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	audience := strings.TrimSpace(payload.Audience)
	if audience == "" {
		audience = strings.TrimSpace(payload.IssuedTo)
	}
	if audience != strings.TrimSpace(s.cfg.Google.ClientID) {
		return nil, errors.New("google token audience mismatch")
	}
	if strings.TrimSpace(payload.Subject) == "" {
		return nil, errors.New("google token subject is empty")
	}

	return &payload, nil
}
