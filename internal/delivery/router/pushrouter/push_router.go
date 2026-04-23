package pushrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/pushhandler"

	"github.com/gin-gonic/gin"
)

func RegisterPushRoutes(group *gin.RouterGroup, handler *pushhandler.PushHandler, authMiddleware gin.HandlerFunc) {
	group.GET("/push/public-key", handler.GetPushPublicKey)

	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("/push/subscriptions", handler.PostPushSubscription)
		protected.DELETE("/push/subscriptions", handler.DeletePushSubscription)
	}
}
