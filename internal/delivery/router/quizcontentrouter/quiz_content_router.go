package quizcontentrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/quizcontenthandler"

	"github.com/gin-gonic/gin"
)

func RegisterQuizContentRoutes(group *gin.RouterGroup, handler *quizcontenthandler.QuizContentHandler) {
	group.GET("/quiz/categories", handler.GetQuizCategories)
	group.GET("/quiz/questions", handler.GetQuizQuestions)
}
