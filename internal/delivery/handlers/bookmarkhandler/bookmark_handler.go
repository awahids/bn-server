package bookmarkhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/request/bookmarkreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	appService serviceinterface.AppService
}

func NewBookmarkHandler(appService serviceinterface.AppService) *BookmarkHandler {
	return &BookmarkHandler{appService: appService}
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
func (h *BookmarkHandler) GetBookmarks(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	typeQuery := strings.TrimSpace(c.Query("type"))
	var bookmarkType *string
	if typeQuery != "" {
		if !handlerutil.IsValidBookmarkType(typeQuery) {
			response.Failed(c, http.StatusBadRequest, "invalid type", "type must be one of: quran, dhikr")
			return
		}
		bookmarkType = &typeQuery
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
// @Param payload body bookmarkreq.CreateBookmarkRequest true "Create bookmark payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /bookmarks [post]
func (h *BookmarkHandler) PostBookmark(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req bookmarkreq.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	req.Type = strings.TrimSpace(req.Type)
	if !handlerutil.IsValidBookmarkType(req.Type) {
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
func (h *BookmarkHandler) DeleteBookmark(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	bookmarkID := strings.TrimSpace(c.Param("id"))
	if bookmarkID == "" {
		response.Failed(c, http.StatusBadRequest, "invalid bookmark id", "bookmark ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
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
