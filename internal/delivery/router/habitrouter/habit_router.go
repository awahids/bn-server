package habitrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/habithandler"

	"github.com/gin-gonic/gin"
)

func RegisterHabitRoutes(group *gin.RouterGroup, handler *habithandler.HabitHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/habits", handler.GetHabits)
		protected.POST("/habits", handler.PostHabit)
		protected.PATCH("/habits/:id", handler.PatchHabit)
		protected.DELETE("/habits/:id", handler.DeleteHabit)
		protected.POST("/habits/completions", handler.PostHabitCompletion)
	}
}
