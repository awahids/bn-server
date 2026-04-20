package progressrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/progresshandler"

	"github.com/gin-gonic/gin"
)

func RegisterProgressRoutes(group *gin.RouterGroup, handler *progresshandler.ProgressHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/progress", handler.GetProgress)
		protected.POST("/progress", handler.PostProgress)
		protected.GET("/progress/activity", handler.GetWeeklyActivity)
		protected.GET("/progress/achievements", handler.GetAchievements)
		protected.GET("/progress/:module/:itemId", handler.GetProgressItem)
	}
}
