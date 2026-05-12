package hijaiyahhandler

import (
	"context"
	"net/http"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type HijaiyahHandler struct {
	appService serviceinterface.AppService
}

func NewHijaiyahHandler(appService serviceinterface.AppService) *HijaiyahHandler {
	return &HijaiyahHandler{appService: appService}
}

func (h *HijaiyahHandler) GetHijaiyahLetters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	letters, err := h.appService.GetHijaiyahLetters(ctx)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get hijaiyah letters", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", letters)
}

func (h *HijaiyahHandler) GetHijaiyahLetterByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Failed(c, http.StatusBadRequest, "id is required", "id parameter is missing")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	letter, err := h.appService.GetHijaiyahLetterByID(ctx, id)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get hijaiyah letter", err.Error())
		return
	}
	if letter == nil {
		response.Failed(c, http.StatusNotFound, "hijaiyah letter not found", "no letter with the given id")
		return
	}

	response.Success(c, http.StatusOK, "ok", letter)
}
