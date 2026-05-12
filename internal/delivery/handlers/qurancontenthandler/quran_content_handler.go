package qurancontenthandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type QuranContentHandler struct {
	appService serviceinterface.AppService
}

func NewQuranContentHandler(appService serviceinterface.AppService) *QuranContentHandler {
	return &QuranContentHandler{appService: appService}
}

func (h *QuranContentHandler) GetQuranSurahs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	surahs, err := h.appService.GetQuranSurahs(ctx)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get quran surahs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", surahs)
}

func (h *QuranContentHandler) GetQuranSurahByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > 114 {
		response.Failed(c, http.StatusBadRequest, "invalid surah id", "id must be a number between 1 and 114")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	surah, err := h.appService.GetQuranSurahByID(ctx, id)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get surah", err.Error())
		return
	}
	if surah == nil {
		response.Failed(c, http.StatusNotFound, "surah not found", "no surah with the given id")
		return
	}

	response.Success(c, http.StatusOK, "ok", surah)
}
