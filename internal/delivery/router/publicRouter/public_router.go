package publicrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/publichandler"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(group *gin.RouterGroup, handler *publichandler.PublicHandler) {
	group.GET("/audio-proxy", handler.GetAudioProxy)
	group.GET("/prayer-times", handler.GetPrayerTimes)
}
