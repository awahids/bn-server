package pushhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/request/pushreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	appService     serviceinterface.AppService
	vapidPublicKey string
	enabled        bool
}

func NewPushHandler(appService serviceinterface.AppService, vapidPublicKey string, enabled bool) *PushHandler {
	return &PushHandler{
		appService:     appService,
		vapidPublicKey: strings.TrimSpace(vapidPublicKey),
		enabled:        enabled,
	}
}

// GetPushPublicKey godoc
// @Summary Get push public key
// @Description Get VAPID public key for web push subscription.
// @Tags Push
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 503 {object} response.APIResponse
// @Router /push/public-key [get]
func (h *PushHandler) GetPushPublicKey(c *gin.Context) {
	if !h.enabled || h.vapidPublicKey == "" {
		response.Failed(c, http.StatusServiceUnavailable, "push notification is not configured", "missing VAPID public key")
		return
	}

	response.Success(c, http.StatusOK, "ok", gin.H{
		"publicKey": h.vapidPublicKey,
	})
}

// PostPushSubscription godoc
// @Summary Upsert push subscription
// @Description Register or update web push subscription for authenticated user.
// @Tags Push
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body pushreq.UpsertPushSubscriptionRequest true "Push subscription payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /push/subscriptions [post]
func (h *PushHandler) PostPushSubscription(c *gin.Context) {
	if !h.enabled {
		response.Failed(c, http.StatusServiceUnavailable, "push notification is not configured", "push disabled")
		return
	}

	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req pushreq.UpsertPushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	err = h.appService.UpsertPushSubscription(ctx, userID, serviceinterface.UpsertPushSubscriptionInput{
		Endpoint:       req.Endpoint,
		ExpirationTime: req.ExpirationTime,
		Timezone:       req.Timezone,
		Keys: serviceinterface.PushSubscriptionKeysInput{
			P256DH: req.Keys.P256DH,
			Auth:   req.Keys.Auth,
		},
	})
	if err != nil {
		if errors.Is(err, appservice.ErrPushInvalidData) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to save push subscription", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "created", nil)
}

// DeletePushSubscription godoc
// @Summary Delete push subscription
// @Description Delete web push subscription for authenticated user.
// @Tags Push
// @Produce json
// @Security BearerAuth
// @Param endpoint query string false "Push endpoint"
// @Param payload body pushreq.DeletePushSubscriptionRequest false "Delete payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /push/subscriptions [delete]
func (h *PushHandler) DeletePushSubscription(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	endpoint := strings.TrimSpace(c.Query("endpoint"))
	if endpoint == "" {
		var req pushreq.DeletePushSubscriptionRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			endpoint = strings.TrimSpace(req.Endpoint)
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	err = h.appService.DeletePushSubscription(ctx, userID, endpoint)
	if err != nil {
		if errors.Is(err, appservice.ErrPushInvalidData) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to delete push subscription", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "deleted", nil)
}
