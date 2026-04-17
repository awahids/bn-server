package authrouter

import (
	authhandler "bn-mobile/server/internal/delivery/handlers/authHandler"

	"github.com/gin-gonic/gin"
)

type RateLimitMiddlewares struct {
	Google  gin.HandlerFunc
	Refresh gin.HandlerFunc
	Logout  gin.HandlerFunc
}

func RegisterAuthRoutes(
	group *gin.RouterGroup,
	handler *authhandler.AuthHandler,
	authMiddleware gin.HandlerFunc,
	rateLimiters RateLimitMiddlewares,
) {
	auth := group.Group("/auth")
	{
		auth.POST("/google", rateLimiters.Google, handler.GoogleLogin)
		auth.POST("/refresh", rateLimiters.Refresh, handler.RefreshToken)
		auth.POST("/logout", rateLimiters.Logout, handler.Logout)
		auth.GET("/me", authMiddleware, handler.Me)
	}
}
