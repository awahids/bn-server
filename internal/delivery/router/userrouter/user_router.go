package userrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/userhandler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(group *gin.RouterGroup, handler *userhandler.UserHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/user", handler.GetUser)
		protected.PATCH("/user", handler.PatchUser)
	}
}
