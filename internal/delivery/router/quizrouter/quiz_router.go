package quizrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/quizhandler"

	"github.com/gin-gonic/gin"
)

func RegisterQuizRoutes(group *gin.RouterGroup, handler *quizhandler.QuizHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/quiz/attempts", handler.GetQuizAttempts)
		protected.POST("/quiz/attempts", handler.PostQuizAttempt)
		protected.GET("/quiz/stats", handler.GetQuizStats)
	}
}
