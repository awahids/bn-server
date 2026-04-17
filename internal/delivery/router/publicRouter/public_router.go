package publicrouter

import (
	apphandler "bn-mobile/server/internal/delivery/handlers/appHandler"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(group *gin.RouterGroup, handler *apphandler.PublicHandler) {
	group.GET("/audio-proxy", handler.GetAudioProxy)
	group.GET("/prayer-times", handler.GetPrayerTimes)
}
