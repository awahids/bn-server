package habithandler

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/request/habitreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/models"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

var reminderTimeRegex = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

type HabitHandler struct {
	appService serviceinterface.AppService
}

func NewHabitHandler(appService serviceinterface.AppService) *HabitHandler {
	return &HabitHandler{appService: appService}
}

type HabitsResponse struct {
	Habits      []models.Habit           `json:"habits"`
	Completions []models.HabitCompletion `json:"completions"`
}

// GetHabits godoc
// @Summary Get habits
// @Description Get authenticated user habits and completion logs.
// @Tags Habit
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /habits [get]
func (h *HabitHandler) GetHabits(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	habits, err := h.appService.GetHabits(ctx, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get habits", err.Error())
		return
	}

	completions, err := h.appService.GetHabitCompletions(ctx, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get habit completions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", HabitsResponse{
		Habits:      habits,
		Completions: completions,
	})
}

// PostHabit godoc
// @Summary Create habit
// @Description Create a habit for authenticated user.
// @Tags Habit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body habitreq.CreateHabitRequest true "Create habit payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /habits [post]
func (h *HabitHandler) PostHabit(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req habitreq.CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 191 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "name must be 1-191 characters")
		return
	}

	category := strings.TrimSpace(req.Category)
	if category == "" {
		category = "Other"
	}
	if len(category) > 50 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "category must be at most 50 characters")
		return
	}

	reminderTime := strings.TrimSpace(req.ReminderTime)
	if reminderTime != "" && !reminderTimeRegex.MatchString(reminderTime) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "reminderTime must be in HH:MM format")
		return
	}

	reminderEnabled := handlerutil.BoolOrDefault(req.ReminderEnabled, false)
	if reminderEnabled && reminderTime == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "reminderTime is required when reminderEnabled is true")
		return
	}
	if !reminderEnabled {
		reminderTime = ""
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	habit, err := h.appService.CreateHabit(ctx, userID, serviceinterface.CreateHabitInput{
		Name:            name,
		Category:        category,
		ReminderTime:    reminderTime,
		ReminderEnabled: reminderEnabled,
	})
	if err != nil {
		if errors.Is(err, appservice.ErrHabitInvalidData) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to create habit", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "created", habit)
}

// PatchHabit godoc
// @Summary Update habit
// @Description Update habit for authenticated user.
// @Tags Habit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Habit ID"
// @Param payload body habitreq.UpdateHabitRequest true "Update habit payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /habits/{id} [patch]
func (h *HabitHandler) PatchHabit(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	habitID := strings.TrimSpace(c.Param("id"))
	if habitID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid habit id", "habit ID is required")
		return
	}

	var req habitreq.UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" || len(trimmed) > 191 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "name must be 1-191 characters")
			return
		}
		req.Name = &trimmed
	}

	if req.Category != nil {
		trimmed := strings.TrimSpace(*req.Category)
		if trimmed == "" || len(trimmed) > 50 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "category must be 1-50 characters")
			return
		}
		req.Category = &trimmed
	}

	if req.ReminderTime != nil {
		trimmed := strings.TrimSpace(*req.ReminderTime)
		if trimmed != "" && !reminderTimeRegex.MatchString(trimmed) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "reminderTime must be in HH:MM format")
			return
		}
		req.ReminderTime = &trimmed
	}

	if req.ReminderEnabled != nil && *req.ReminderEnabled {
		if req.ReminderTime != nil && strings.TrimSpace(*req.ReminderTime) == "" {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "reminderTime is required when reminderEnabled is true")
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	habit, err := h.appService.UpdateHabit(ctx, userID, habitID, serviceinterface.UpdateHabitInput{
		Name:            req.Name,
		Category:        req.Category,
		ReminderTime:    req.ReminderTime,
		ReminderEnabled: req.ReminderEnabled,
	})
	if err != nil {
		switch {
		case errors.Is(err, appservice.ErrHabitNotFound):
			response.Failed(c, http.StatusNotFound, "habit not found", err.Error())
		case errors.Is(err, appservice.ErrHabitInvalidData):
			response.Failed(c, http.StatusBadRequest, "invalid request data", err.Error())
		default:
			response.Failed(c, http.StatusInternalServerError, "failed to update habit", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "ok", habit)
}

// DeleteHabit godoc
// @Summary Delete habit
// @Description Delete habit for authenticated user.
// @Tags Habit
// @Produce json
// @Security BearerAuth
// @Param id path string true "Habit ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /habits/{id} [delete]
func (h *HabitHandler) DeleteHabit(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	habitID := strings.TrimSpace(c.Param("id"))
	if habitID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid habit id", "habit ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	if err := h.appService.DeleteHabit(ctx, userID, habitID); err != nil {
		switch {
		case errors.Is(err, appservice.ErrHabitNotFound):
			response.Failed(c, http.StatusNotFound, "habit not found", err.Error())
		default:
			response.Failed(c, http.StatusInternalServerError, "failed to delete habit", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Habit deleted successfully", nil)
}

// PostHabitCompletion godoc
// @Summary Set habit completion
// @Description Set completion status for a habit on specific date.
// @Tags Habit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body habitreq.SetHabitCompletionRequest true "Set habit completion payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /habits/completions [post]
func (h *HabitHandler) PostHabitCompletion(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req habitreq.SetHabitCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.HabitID = strings.TrimSpace(req.HabitID)
	if req.HabitID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "habitId is required")
		return
	}

	req.Date = strings.TrimSpace(req.Date)
	if !handlerutil.IsValidDate(req.Date) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "date must be in YYYY-MM-DD format")
		return
	}

	completed := handlerutil.BoolOrDefault(req.Completed, true)

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	completion, err := h.appService.SetHabitCompletion(ctx, userID, serviceinterface.SetHabitCompletionInput{
		HabitID:   req.HabitID,
		Date:      req.Date,
		Completed: completed,
	})
	if err != nil {
		switch {
		case errors.Is(err, appservice.ErrHabitNotFound):
			response.Failed(c, http.StatusNotFound, "habit not found", err.Error())
		default:
			response.Failed(c, http.StatusInternalServerError, "failed to set habit completion", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "ok", completion)
}
