package qurancontentrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/qurancontenthandler"

	"github.com/gin-gonic/gin"
)

func RegisterQuranContentRoutes(group *gin.RouterGroup, handler *qurancontenthandler.QuranContentHandler) {
	group.GET("/quran/surahs", handler.GetQuranSurahs)
	group.GET("/quran/surahs/:id", handler.GetQuranSurahByID)
}
