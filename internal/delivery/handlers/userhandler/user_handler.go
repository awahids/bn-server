package userhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/request/userreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/models"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	appService serviceinterface.AppService
}

type UserProfileResponse struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Email         string                `json:"email"`
	Image         string                `json:"image,omitempty"`
	Username      *string               `json:"username,omitempty"`
	Streak        int                   `json:"streak"`
	DailyProgress int                   `json:"dailyProgress"`
	LastActive    time.Time             `json:"lastActive"`
	Preferences   map[string]any        `json:"preferences"`
	Progress      []models.UserProgress `json:"progress,omitempty"`
	Bookmarks     []models.Bookmark     `json:"bookmarks,omitempty"`
	DhikrCounters []models.DhikrCounter `json:"dhikrCounters,omitempty"`
	QuizAttempts  []models.QuizAttempt  `json:"quizAttempts,omitempty"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
}

func NewUserHandler(appService serviceinterface.AppService) *UserHandler {
	return &UserHandler{appService: appService}
}

// GetUser godoc
// @Summary Get user profile
// @Description Return authenticated user profile and related app data.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	user, err := h.appService.GetUserProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, appservice.ErrUserNotFound) {
			response.Failed(c, http.StatusNotFound, "user not found", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to get user profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", toUserProfileResponse(user))
}

// PatchUser godoc
// @Summary Update user profile
// @Description Update authenticated user profile fields.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body userreq.UpdateUserRequest true "Update user payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user [patch]
func (h *UserHandler) PatchUser(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req userreq.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" || len(trimmed) > 100 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "name must be between 1 and 100 characters")
			return
		}
		req.Name = &trimmed
	}

	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		if len(trimmed) < 3 || len(trimmed) > 50 || !handlerutil.IsValidUsername(trimmed) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "username must be 3-50 chars with letters, numbers, and underscore only")
			return
		}
		req.Username = &trimmed
	}

	if req.Streak != nil && *req.Streak < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "streak must be >= 0")
		return
	}

	if req.DailyProgress != nil && *req.DailyProgress < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "dailyProgress must be >= 0")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	updated, err := h.appService.UpdateUserProfile(ctx, userID, serviceinterface.UpdateUserInput{
		Name:          req.Name,
		Username:      req.Username,
		Streak:        req.Streak,
		DailyProgress: req.DailyProgress,
		Preferences:   req.Preferences,
	})
	if err != nil {
		switch {
		case errors.Is(err, appservice.ErrUserNotFound):
			response.Failed(c, http.StatusNotFound, "user not found", err.Error())
		case errors.Is(err, appservice.ErrUsernameTaken):
			response.Failed(c, http.StatusConflict, "username already taken", err.Error())
		default:
			response.Failed(c, http.StatusInternalServerError, "failed to update user profile", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "ok", toUserProfileResponse(updated))
}

func toUserProfileResponse(user *models.User) UserProfileResponse {
	return UserProfileResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Image:         user.AvatarURL,
		Username:      user.Username,
		Streak:        user.Streak,
		DailyProgress: user.DailyProgress,
		LastActive:    user.LastActive,
		Preferences:   parsePreferences(user.Preferences),
		Progress:      user.Progress,
		Bookmarks:     user.Bookmarks,
		DhikrCounters: user.DhikrCounters,
		QuizAttempts:  user.QuizAttempts,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func parsePreferences(raw []byte) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}

	var preferences map[string]any
	if err := json.Unmarshal(raw, &preferences); err != nil || preferences == nil {
		return map[string]any{}
	}
	return preferences
}
