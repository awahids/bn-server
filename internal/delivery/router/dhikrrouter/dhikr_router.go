package dhikrrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/dhikrhandler"

	"github.com/gin-gonic/gin"
)

func RegisterDhikrRoutes(group *gin.RouterGroup, handler *dhikrhandler.DhikrHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/dhikr/counters", handler.GetDhikrCounters)
		protected.POST("/dhikr/counters", handler.PostDhikrCounter)
	}
}
