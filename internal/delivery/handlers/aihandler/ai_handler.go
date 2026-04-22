package aihandler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/delivery/middleware"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService             serviceinterface.AIService
	dailyCoachUsage       map[string]string
	dailyCoachUsageLocker sync.Mutex
}

func NewAIHandler(aiService serviceinterface.AIService) *AIHandler {
	return &AIHandler{
		aiService:       aiService,
		dailyCoachUsage: make(map[string]string),
	}
}

type CoachRequest struct {
	System  string `json:"system"`
	Message string `json:"message"`
}

const defaultIslamicHabitCoachSystemPrompt = "You are an Islamic habit coach for a Muslim learning app. Guide users to build istiqamah with practical actions rooted in Islamic values such as ikhlas in intention, disciplined worship, adab, gratitude, and consistency. Keep responses warm, non-judgmental, practical, and concise. Prefer 3-5 actionable steps when giving advice. Use Bahasa Indonesia unless the user asks for another language. Do not use emojis. Do not issue definitive fatwa; for sensitive fiqh issues, advise consulting a trusted ustadz or scholar."
const coachBypassEmail = "awahid.safhadi@gmail.com"
const coachDailyLimitMessage = "AI Coach hanya bisa digunakan 1x per hari. Coba lagi besok."

var coachDayLocation = loadCoachDayLocation()

func loadCoachDayLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.UTC
	}
	return location
}

func normalizedCoachEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func canBypassCoachDailyLimit(email string) bool {
	return normalizedCoachEmail(email) == coachBypassEmail
}

func coachDay(now time.Time) string {
	return now.In(coachDayLocation).Format("2006-01-02")
}

func (h *AIHandler) reserveCoachUsage(userID, userEmail string, now time.Time) bool {
	if canBypassCoachDailyLimit(userEmail) {
		return true
	}

	day := coachDay(now)

	h.dailyCoachUsageLocker.Lock()
	defer h.dailyCoachUsageLocker.Unlock()

	if lastUsageDay, exists := h.dailyCoachUsage[userID]; exists && lastUsageDay == day {
		return false
	}

	h.dailyCoachUsage[userID] = day
	return true
}

func (h *AIHandler) rollbackCoachUsage(userID, userEmail string, now time.Time) {
	if canBypassCoachDailyLimit(userEmail) {
		return
	}

	day := coachDay(now)

	h.dailyCoachUsageLocker.Lock()
	defer h.dailyCoachUsageLocker.Unlock()

	if lastUsageDay, exists := h.dailyCoachUsage[userID]; exists && lastUsageDay == day {
		delete(h.dailyCoachUsage, userID)
	}
}

// GetCoachResponse godoc
// @Summary Get AI coach response
// @Description Get a response from the AI habit coach based on user progress and habits.
// @Tags AI
// @Accept json
// @Produce json
// @Param request body CoachRequest true "AI Coach request payload"
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 429 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /ai/coach [post]
func (h *AIHandler) GetCoachResponse(c *gin.Context) {
	var req CoachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.Message == "" {
		response.Failed(c, http.StatusBadRequest, "message is required", "message is required")
		return
	}

	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	userEmail := c.GetString(middleware.ContextUserEmail)
	now := time.Now()
	if !h.reserveCoachUsage(userID, userEmail, now) {
		response.Failed(c, http.StatusTooManyRequests, "daily AI coach limit reached", coachDailyLimitMessage)
		return
	}

	systemPrompt := defaultIslamicHabitCoachSystemPrompt
	if customSystemPrompt := strings.TrimSpace(req.System); customSystemPrompt != "" {
		systemPrompt = systemPrompt + " Additional app instruction: " + customSystemPrompt
	}

	content, err := h.aiService.GetCoachResponse(c.Request.Context(), systemPrompt, req.Message)
	if err != nil {
		h.rollbackCoachUsage(userID, userEmail, now)
		response.Failed(c, http.StatusInternalServerError, "failed to get AI response", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"content": content,
		},
	})
}
