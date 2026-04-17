package approuter

import (
	apphandler "bn-mobile/server/internal/delivery/handlers/appHandler"

	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(
	group *gin.RouterGroup,
	handler *apphandler.AppHandler,
	authMiddleware gin.HandlerFunc,
) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/user", handler.GetUser)
		protected.PATCH("/user", handler.PatchUser)

		protected.GET("/progress", handler.GetProgress)
		protected.POST("/progress", handler.PostProgress)
		protected.GET("/progress/:module/:itemId", handler.GetProgressItem)

		protected.GET("/bookmarks", handler.GetBookmarks)
		protected.POST("/bookmarks", handler.PostBookmark)
		protected.DELETE("/bookmarks/:id", handler.DeleteBookmark)

		protected.GET("/dhikr/counters", handler.GetDhikrCounters)
		protected.POST("/dhikr/counters", handler.PostDhikrCounter)

		protected.GET("/quiz/attempts", handler.GetQuizAttempts)
		protected.POST("/quiz/attempts", handler.PostQuizAttempt)
		protected.GET("/quiz/stats", handler.GetQuizStats)
	}
}
