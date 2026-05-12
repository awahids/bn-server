package tajwidrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/tajwidhandler"

	"github.com/gin-gonic/gin"
)

func RegisterTajwidRoutes(group *gin.RouterGroup, handler *tajwidhandler.TajwidHandler) {
	group.GET("/tajwid/rules", handler.GetTajwidRules)
	group.GET("/tajwid/rules/:id", handler.GetTajwidRuleByID)
}
