package authhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"bn-mobile/server/configs"
	authreq "bn-mobile/server/internal/delivery/data/request/authReq"
	"bn-mobile/server/internal/delivery/data/response"
	authres "bn-mobile/server/internal/delivery/data/response/authRes"
	"bn-mobile/server/internal/delivery/middleware"
	"bn-mobile/server/internal/domain/models"
	authservice "bn-mobile/server/internal/domain/services/authService"
	"bn-mobile/server/internal/domain/services/serviceInterface"

	"github.com/gin-gonic/gin"
)

const requestTimeout = 8 * time.Second

type AuthHandler struct {
	authService serviceinterface.AuthService
	cookieCfg   configs.AuthCookieConfig
}

func NewAuthHandler(authService serviceinterface.AuthService, cookieCfg configs.AuthCookieConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cookieCfg:   cookieCfg,
	}
}

// GoogleLogin godoc
// @Summary Login with Google ID token
// @Description Validate Google ID token and return access token with refresh-token cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body authreq.GoogleLoginRequest true "Google login payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req authreq.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	user, tokens, err := h.authService.LoginWithGoogle(ctx, req.IDToken)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	response.Success(c, http.StatusOK, "login success", toAuthResponse(user, tokens))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Issue a new access token and rotate refresh-token cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := h.getRefreshTokenCookie(c)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, "missing refresh token cookie", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	tokens, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidRefreshToken) {
			response.Failed(c, http.StatusUnauthorized, "invalid refresh token", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to refresh token", err.Error())
		return
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	response.Success(c, http.StatusOK, "token refreshed", authres.TokenResponse{
		AccessToken: tokens.AccessToken,
		TokenType:   tokens.TokenType,
		ExpiresIn:   tokens.ExpiresIn,
	})
}

// Logout godoc
// @Summary Logout current session
// @Description Revoke refresh token and clear refresh-token cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := h.getRefreshTokenCookie(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	if err := h.authService.Logout(ctx, refreshToken); err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to logout", err.Error())
		return
	}

	h.clearRefreshTokenCookie(c)
	response.Success(c, http.StatusOK, "logout success", nil)
}

// Me godoc
// @Summary Get current authenticated user
// @Description Return current user profile from Bearer token.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		response.Failed(c, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		response.Failed(c, http.StatusUnauthorized, "unauthorized", "invalid user context")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	user, err := h.authService.GetCurrentUser(ctx, userID)
	if err != nil {
		if errors.Is(err, authservice.ErrUserNotFound) {
			response.Failed(c, http.StatusNotFound, "user not found", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to get profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", authres.UserProfile{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		LastLoginAt: user.LastLoginAt,
	})
}

func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, authservice.ErrInvalidGoogleToken):
		response.Failed(c, http.StatusUnauthorized, "google token invalid", err.Error())
	case errors.Is(err, authservice.ErrEmailNotVerified):
		response.Failed(c, http.StatusForbidden, "google email not verified", err.Error())
	default:
		response.Failed(c, http.StatusInternalServerError, "authentication failed", err.Error())
	}
}

func toAuthResponse(user *models.User, tokens *serviceinterface.TokenPair) authres.AuthResponse {
	return authres.AuthResponse{
		User: authres.UserProfile{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			AvatarURL:   user.AvatarURL,
			Role:        user.Role,
			LastLoginAt: user.LastLoginAt,
		},
		Tokens: authres.TokenResponse{
			AccessToken: tokens.AccessToken,
			TokenType:   tokens.TokenType,
			ExpiresIn:   tokens.ExpiresIn,
		},
	}
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetSameSite(h.cookieCfg.RefreshTokenSameSite)
	c.SetCookie(
		h.cookieCfg.RefreshTokenName,
		token,
		h.cookieCfg.RefreshTokenMaxAge,
		h.cookieCfg.RefreshTokenPath,
		h.cookieCfg.RefreshTokenDomain,
		h.cookieCfg.RefreshTokenSecure,
		h.cookieCfg.RefreshTokenHTTPOnly,
	)
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetSameSite(h.cookieCfg.RefreshTokenSameSite)
	c.SetCookie(
		h.cookieCfg.RefreshTokenName,
		"",
		-1,
		h.cookieCfg.RefreshTokenPath,
		h.cookieCfg.RefreshTokenDomain,
		h.cookieCfg.RefreshTokenSecure,
		h.cookieCfg.RefreshTokenHTTPOnly,
	)
}

func (h *AuthHandler) getRefreshTokenCookie(c *gin.Context) (string, error) {
	return c.Cookie(h.cookieCfg.RefreshTokenName)
}
