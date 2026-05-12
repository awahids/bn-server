package hijaiyahrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/hijaiyahhandler"

	"github.com/gin-gonic/gin"
)

func RegisterHijaiyahRoutes(group *gin.RouterGroup, handler *hijaiyahhandler.HijaiyahHandler) {
	group.GET("/hijaiyah/letters", handler.GetHijaiyahLetters)
	group.GET("/hijaiyah/letters/:id", handler.GetHijaiyahLetterByID)
}
