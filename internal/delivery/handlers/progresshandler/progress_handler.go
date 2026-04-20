package progresshandler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/request/progressreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	appService serviceinterface.AppService
}

func NewProgressHandler(appService serviceinterface.AppService) *ProgressHandler {
	return &ProgressHandler{appService: appService}
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
func (h *ProgressHandler) GetProgress(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	moduleQuery := strings.TrimSpace(c.Query("module"))
	var module *string
	if moduleQuery != "" {
		if !handlerutil.IsValidProgressModule(moduleQuery) {
			response.Failed(c, http.StatusBadRequest, "invalid module", "module must be one of: hijaiyah, quran, dhikr, quiz")
			return
		}
		module = &moduleQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
// @Param payload body progressreq.UpsertProgressRequest true "Upsert progress payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress [post]
func (h *ProgressHandler) PostProgress(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req progressreq.UpsertProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.Module = strings.TrimSpace(req.Module)
	if !handlerutil.IsValidProgressModule(req.Module) {
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	progress, err := h.appService.UpsertProgress(ctx, userID, serviceinterface.UpsertProgressInput{
		Module:    req.Module,
		ItemID:    req.ItemID,
		Progress:  req.Progress,
		Completed: handlerutil.BoolOrDefault(req.Completed, false),
		Score:     handlerutil.IntOrDefault(req.Score, 0),
		TimeSpent: handlerutil.IntOrDefault(req.TimeSpent, 0),
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
func (h *ProgressHandler) GetProgressItem(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	module := strings.TrimSpace(c.Param("module"))
	if !handlerutil.IsValidProgressModule(module) {
		response.Failed(c, http.StatusBadRequest, "invalid module", "module must be one of: hijaiyah, quran, dhikr, quiz")
		return
	}

	itemID := strings.TrimSpace(handlerutil.DecodePath(c.Param("itemId")))
	if itemID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid itemId", "itemId is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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

// GetAchievements godoc
// @Summary Get user achievements
// @Description Get authenticated user achievements with unlocked status.
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress/achievements [get]
func (h *ProgressHandler) GetAchievements(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	achievements, err := h.appService.GetAchievements(ctx, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get achievements", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", achievements)
}

// GetWeeklyActivity godoc
// @Summary Get user weekly activity
// @Description Get user's daily activity completion for the last 7 days.
// @Tags Progress
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /progress/activity [get]
func (h *ProgressHandler) GetWeeklyActivity(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	activity, err := h.appService.GetWeeklyActivity(ctx, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get weekly activity", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", activity)
}
