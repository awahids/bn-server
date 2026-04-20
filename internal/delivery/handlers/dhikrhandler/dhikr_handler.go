package dhikrhandler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/request/dhikrreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type DhikrHandler struct {
	appService serviceinterface.AppService
}

func NewDhikrHandler(appService serviceinterface.AppService) *DhikrHandler {
	return &DhikrHandler{appService: appService}
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
func (h *DhikrHandler) GetDhikrCounters(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	dateValue := strings.TrimSpace(c.Query("date"))
	if dateValue == "" {
		dateValue = time.Now().UTC().Format("2006-01-02")
	}
	if !handlerutil.IsValidDate(dateValue) {
		response.Failed(c, http.StatusBadRequest, "invalid date format", "date must be in YYYY-MM-DD format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
// @Param payload body dhikrreq.UpsertDhikrCounterRequest true "Upsert dhikr counter payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /dhikr/counters [post]
func (h *DhikrHandler) PostDhikrCounter(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req dhikrreq.UpsertDhikrCounterRequest
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

	target := handlerutil.IntOrDefault(req.Target, 33)
	if target < 1 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "target must be >= 1")
		return
	}

	req.Date = strings.TrimSpace(req.Date)
	if !handlerutil.IsValidDate(req.Date) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "date must be in YYYY-MM-DD format")
		return
	}

	req.Session = strings.TrimSpace(req.Session)
	if !handlerutil.IsValidDhikrSession(req.Session) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "session must be one of: morning, evening")
		return
	}

	completed := handlerutil.BoolOrDefault(req.Completed, req.Count >= target)

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
