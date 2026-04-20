package quizhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/request/quizreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	appService serviceinterface.AppService
}

func NewQuizHandler(appService serviceinterface.AppService) *QuizHandler {
	return &QuizHandler{appService: appService}
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
func (h *QuizHandler) GetQuizAttempts(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	categoryQuery := strings.TrimSpace(c.Query("category"))
	var category *string
	if categoryQuery != "" {
		category = &categoryQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
// @Param payload body quizreq.CreateQuizAttemptRequest true "Create quiz attempt payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /quiz/attempts [post]
func (h *QuizHandler) PostQuizAttempt(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req quizreq.CreateQuizAttemptRequest
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
func (h *QuizHandler) GetQuizStats(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	category := strings.TrimSpace(c.Query("category"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
