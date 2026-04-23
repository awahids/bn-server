package schoolrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/schoolhandler"

	"github.com/gin-gonic/gin"
)

func RegisterSchoolRoutes(group *gin.RouterGroup, handler *schoolhandler.SchoolHandler, authMiddleware gin.HandlerFunc) {
	group.GET("/schools", handler.GetSchools)

	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("/schools", handler.PostSchool)
	}
}
