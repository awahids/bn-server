package tajwidhandler

import (
	"context"
	"net/http"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type TajwidHandler struct {
	appService serviceinterface.AppService
}

func NewTajwidHandler(appService serviceinterface.AppService) *TajwidHandler {
	return &TajwidHandler{appService: appService}
}

func (h *TajwidHandler) GetTajwidRules(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	rules, err := h.appService.GetTajwidRules(ctx)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get tajwid rules", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", rules)
}

func (h *TajwidHandler) GetTajwidRuleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Failed(c, http.StatusBadRequest, "id is required", "id parameter is missing")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	rule, err := h.appService.GetTajwidRuleByID(ctx, id)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get tajwid rule", err.Error())
		return
	}
	if rule == nil {
		response.Failed(c, http.StatusNotFound, "tajwid rule not found", "no rule with the given id")
		return
	}

	response.Success(c, http.StatusOK, "ok", rule)
}
