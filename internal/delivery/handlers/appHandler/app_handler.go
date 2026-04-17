package apphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	appreq "bn-mobile/server/internal/delivery/data/request/appReq"
	"bn-mobile/server/internal/delivery/data/response"
	"bn-mobile/server/internal/domain/models"
	appservice "bn-mobile/server/internal/domain/services/appService"
	"bn-mobile/server/internal/domain/services/serviceInterface"

	"github.com/gin-gonic/gin"
)

const appRequestTimeout = 8 * time.Second

type AppHandler struct {
	appService serviceinterface.AppService
}

type userProfileResponse struct {
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

func NewAppHandler(appService serviceinterface.AppService) *AppHandler {
	return &AppHandler{appService: appService}
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
func (h *AppHandler) GetUser(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
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
// @Param payload body appreq.UpdateUserRequest true "Update user payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user [patch]
func (h *AppHandler) PatchUser(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	var req appreq.UpdateUserRequest
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
		if len(trimmed) < 3 || len(trimmed) > 50 || !usernameRegex.MatchString(trimmed) {
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
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

// GetProgress godoc
// @Summary Get user progress
// @Description Get authenticated user progress, optionally filtered by module.
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Param module query string false "Filter by module (hijaiyah|quran|dhikr|quiz)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress [get]
func (h *AppHandler) GetProgress(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	moduleQuery := strings.TrimSpace(c.Query("module"))
	var module *string
	if moduleQuery != "" {
		if !isValidProgressModule(moduleQuery) {
			response.Failed(c, http.StatusBadRequest, "invalid module", "module must be one of: hijaiyah, quran, dhikr, quiz")
			return
		}
		module = &moduleQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	progress, err := h.appService.GetProgress(ctx, userID, module)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get progress", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", progress)
}

// PostProgress godoc
// @Summary Create or update progress
// @Description Upsert user progress item by module and item id.
// @Tags Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body appreq.UpsertProgressRequest true "Upsert progress payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress [post]
func (h *AppHandler) PostProgress(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	var req appreq.UpsertProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.Module = strings.TrimSpace(req.Module)
	if !isValidProgressModule(req.Module) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "module must be one of: hijaiyah, quran, dhikr, quiz")
		return
	}

	req.ItemID = strings.TrimSpace(req.ItemID)
	if req.ItemID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "itemId is required")
		return
	}

	if req.Progress < 0 || req.Progress > 100 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "progress must be between 0 and 100")
		return
	}

	if req.Score != nil && *req.Score < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "score must be >= 0")
		return
	}

	if req.TimeSpent != nil && *req.TimeSpent < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "timeSpent must be >= 0")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	progress, err := h.appService.UpsertProgress(ctx, userID, serviceinterface.UpsertProgressInput{
		Module:    req.Module,
		ItemID:    req.ItemID,
		Progress:  req.Progress,
		Completed: boolOrDefault(req.Completed, false),
		Score:     intOrDefault(req.Score, 0),
		TimeSpent: intOrDefault(req.TimeSpent, 0),
	})
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to update progress", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", progress)
}

// GetProgressItem godoc
// @Summary Get progress item
// @Description Get progress item for a specific module and item id.
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Param module path string true "Module name"
// @Param itemId path string true "Item id"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress/{module}/{itemId} [get]
func (h *AppHandler) GetProgressItem(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	module := strings.TrimSpace(c.Param("module"))
	if !isValidProgressModule(module) {
		response.Failed(c, http.StatusBadRequest, "invalid module", "module must be one of: hijaiyah, quran, dhikr, quiz")
		return
	}

	itemID := strings.TrimSpace(decodePath(c.Param("itemId")))
	if itemID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid itemId", "itemId is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	progress, err := h.appService.GetProgressItem(ctx, userID, module, itemID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get progress item", err.Error())
		return
	}

	if progress == nil {
		response.Success(c, http.StatusOK, "ok", models.UserProgress{
			UserID:       userID,
			Module:       module,
			ItemID:       itemID,
			Progress:     0,
			Completed:    false,
			Score:        0,
			TimeSpent:    0,
			LastAccessed: time.Now(),
		})
		return
	}

	response.Success(c, http.StatusOK, "ok", progress)
}

// GetBookmarks godoc
// @Summary Get bookmarks
// @Description Get authenticated user bookmarks, optionally filtered by type.
// @Tags Bookmark
// @Produce json
// @Security BearerAuth
// @Param type query string false "Bookmark type (quran|dhikr)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /bookmarks [get]
func (h *AppHandler) GetBookmarks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	typeQuery := strings.TrimSpace(c.Query("type"))
	var bookmarkType *string
	if typeQuery != "" {
		if !isValidBookmarkType(typeQuery) {
			response.Failed(c, http.StatusBadRequest, "invalid type", "type must be one of: quran, dhikr")
			return
		}
		bookmarkType = &typeQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	bookmarks, err := h.appService.GetBookmarks(ctx, userID, bookmarkType)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get bookmarks", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", bookmarks)
}

// PostBookmark godoc
// @Summary Create bookmark
// @Description Create bookmark for authenticated user.
// @Tags Bookmark
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body appreq.CreateBookmarkRequest true "Create bookmark payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /bookmarks [post]
func (h *AppHandler) PostBookmark(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	var req appreq.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.Type = strings.TrimSpace(req.Type)
	if !isValidBookmarkType(req.Type) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "type must be one of: quran, dhikr")
		return
	}

	req.ContentID = strings.TrimSpace(req.ContentID)
	if req.ContentID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "contentId is required")
		return
	}

	if req.Note != nil {
		trimmed := strings.TrimSpace(*req.Note)
		if len(trimmed) > 500 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "note must be at most 500 characters")
			return
		}
		req.Note = &trimmed
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	bookmark, err := h.appService.CreateBookmark(ctx, userID, serviceinterface.CreateBookmarkInput{
		Type:      req.Type,
		ContentID: req.ContentID,
		Note:      req.Note,
	})
	if err != nil {
		if errors.Is(err, appservice.ErrBookmarkExists) {
			response.Failed(c, http.StatusConflict, "bookmark already exists", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to create bookmark", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "created", bookmark)
}

// DeleteBookmark godoc
// @Summary Delete bookmark
// @Description Delete bookmark by id for authenticated user.
// @Tags Bookmark
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bookmark ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /bookmarks/{id} [delete]
func (h *AppHandler) DeleteBookmark(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	bookmarkID := strings.TrimSpace(c.Param("id"))
	if bookmarkID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid bookmark id", "bookmark ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	err = h.appService.DeleteBookmark(ctx, userID, bookmarkID)
	if err != nil {
		switch {
		case errors.Is(err, appservice.ErrBookmarkNotFound):
			response.Failed(c, http.StatusNotFound, "bookmark not found", err.Error())
		case errors.Is(err, appservice.ErrBookmarkForbidden):
			response.Failed(c, http.StatusForbidden, "forbidden", err.Error())
		default:
			response.Failed(c, http.StatusInternalServerError, "failed to delete bookmark", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Bookmark deleted successfully", nil)
}

// GetDhikrCounters godoc
// @Summary Get dhikr counters
// @Description Get dhikr counters for authenticated user by date.
// @Tags Dhikr
// @Produce json
// @Security BearerAuth
// @Param date query string false "Date in YYYY-MM-DD"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /dhikr/counters [get]
func (h *AppHandler) GetDhikrCounters(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	dateValue := strings.TrimSpace(c.Query("date"))
	if dateValue == "" {
		dateValue = time.Now().UTC().Format("2006-01-02")
	}
	if !dateRegex.MatchString(dateValue) {
		response.Failed(c, http.StatusBadRequest, "invalid date format", "date must be in YYYY-MM-DD format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	counters, err := h.appService.GetDhikrCounters(ctx, userID, dateValue)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get dhikr counters", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", counters)
}

// PostDhikrCounter godoc
// @Summary Create or update dhikr counter
// @Description Upsert dhikr counter entry for authenticated user.
// @Tags Dhikr
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body appreq.UpsertDhikrCounterRequest true "Upsert dhikr counter payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /dhikr/counters [post]
func (h *AppHandler) PostDhikrCounter(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	var req appreq.UpsertDhikrCounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.DhikrID = strings.TrimSpace(req.DhikrID)
	if req.DhikrID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "dhikrId is required")
		return
	}

	if req.Count < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "count must be >= 0")
		return
	}

	target := intOrDefault(req.Target, 33)
	if target < 1 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "target must be >= 1")
		return
	}

	req.Date = strings.TrimSpace(req.Date)
	if !dateRegex.MatchString(req.Date) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "date must be in YYYY-MM-DD format")
		return
	}

	req.Session = strings.TrimSpace(req.Session)
	if !isValidDhikrSession(req.Session) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "session must be one of: morning, evening")
		return
	}

	completed := boolOrDefault(req.Completed, req.Count >= target)

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	counter, err := h.appService.UpsertDhikrCounter(ctx, userID, serviceinterface.UpsertDhikrCounterInput{
		DhikrID:   req.DhikrID,
		Count:     req.Count,
		Target:    target,
		Date:      req.Date,
		Session:   req.Session,
		Completed: completed,
	})
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to upsert dhikr counter", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", counter)
}

// GetQuizAttempts godoc
// @Summary Get quiz attempts
// @Description Get quiz attempts for authenticated user, optionally filtered by category.
// @Tags Quiz
// @Produce json
// @Security BearerAuth
// @Param category query string false "Quiz category"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /quiz/attempts [get]
func (h *AppHandler) GetQuizAttempts(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	categoryQuery := strings.TrimSpace(c.Query("category"))
	var category *string
	if categoryQuery != "" {
		category = &categoryQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	attempts, err := h.appService.GetQuizAttempts(ctx, userID, category)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get quiz attempts", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", attempts)
}

// PostQuizAttempt godoc
// @Summary Create quiz attempt
// @Description Save a completed quiz attempt for authenticated user.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body appreq.CreateQuizAttemptRequest true "Create quiz attempt payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /quiz/attempts [post]
func (h *AppHandler) PostQuizAttempt(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	var req appreq.CreateQuizAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.Category = strings.TrimSpace(req.Category)
	if req.Category == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "category is required")
		return
	}
	if req.Score < 0 || req.Score > 100 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "score must be between 0 and 100")
		return
	}
	if req.TotalQuestions < 1 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "totalQuestions must be >= 1")
		return
	}
	if req.TimeSpent < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "timeSpent must be >= 0")
		return
	}

	answers := make([]serviceinterface.QuizAnswerInput, 0, len(req.Answers))
	for _, answer := range req.Answers {
		if answer.TimeSpent != nil && *answer.TimeSpent < 0 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "answer.timeSpent must be >= 0")
			return
		}
		answers = append(answers, serviceinterface.QuizAnswerInput{
			QuestionID:    answer.QuestionID,
			UserAnswer:    answer.UserAnswer,
			CorrectAnswer: answer.CorrectAnswer,
			IsCorrect:     answer.IsCorrect,
			TimeSpent:     answer.TimeSpent,
		})
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	attempt, err := h.appService.CreateQuizAttempt(ctx, userID, serviceinterface.CreateQuizAttemptInput{
		Category:       req.Category,
		Score:          req.Score,
		TotalQuestions: req.TotalQuestions,
		TimeSpent:      req.TimeSpent,
		Answers:        answers,
	})
	if err != nil {
		if errors.Is(err, appservice.ErrScoreMismatch) {
			response.Failed(c, http.StatusBadRequest, "score does not match answers", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to create quiz attempt", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "created", attempt)
}

// GetQuizStats godoc
// @Summary Get quiz stats
// @Description Get overall quiz stats or stats by category for authenticated user.
// @Tags Quiz
// @Produce json
// @Security BearerAuth
// @Param category query string false "Quiz category"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /quiz/stats [get]
func (h *AppHandler) GetQuizStats(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		failUnauthorized(c, err)
		return
	}

	category := strings.TrimSpace(c.Query("category"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), appRequestTimeout)
	defer cancel()

	if category != "" {
		stats, err := h.appService.GetQuizCategoryStats(ctx, userID, category)
		if err != nil {
			response.Failed(c, http.StatusInternalServerError, "failed to get quiz category stats", err.Error())
			return
		}
		response.Success(c, http.StatusOK, "ok", gin.H{
			"category":       category,
			"attempts":       stats.Attempts,
			"averageScore":   stats.AverageScore,
			"bestScore":      stats.BestScore,
			"totalTimeSpent": stats.TotalTimeSpent,
			"lastAttempt":    stats.LastAttempt,
		})
		return
	}

	stats, err := h.appService.GetQuizStats(ctx, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get quiz stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", stats)
}

func toUserProfileResponse(user *models.User) userProfileResponse {
	return userProfileResponse{
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

func isValidProgressModule(module string) bool {
	switch module {
	case string(models.ModuleHijaiyah), string(models.ModuleQuran), string(models.ModuleDhikr), string(models.ModuleQuiz):
		return true
	default:
		return false
	}
}

func isValidBookmarkType(bookmarkType string) bool {
	switch bookmarkType {
	case string(models.BookmarkTypeQuran), string(models.BookmarkTypeDhikr):
		return true
	default:
		return false
	}
}

func isValidDhikrSession(session string) bool {
	switch session {
	case string(models.DhikrSessionMorning), string(models.DhikrSessionEvening):
		return true
	default:
		return false
	}
}
