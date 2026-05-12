package quizcontenthandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type QuizContentHandler struct {
	appService serviceinterface.AppService
}

func NewQuizContentHandler(appService serviceinterface.AppService) *QuizContentHandler {
	return &QuizContentHandler{appService: appService}
}

func (h *QuizContentHandler) GetQuizCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	categories, err := h.appService.GetQuizCategories(ctx)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get quiz categories", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", categories)
}

func (h *QuizContentHandler) GetQuizQuestions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	var categoryID *string
	if cat := strings.TrimSpace(c.Query("category")); cat != "" {
		categoryID = &cat
	}

	var difficulty *string
	if d := strings.TrimSpace(c.Query("difficulty")); d != "" {
		if d != "easy" && d != "medium" && d != "hard" {
			response.Failed(c, http.StatusBadRequest, "invalid difficulty", "difficulty must be easy, medium, or hard")
			return
		}
		difficulty = &d
	}

	questions, err := h.appService.GetQuizQuestions(ctx, categoryID, difficulty)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get quiz questions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", questions)
}
