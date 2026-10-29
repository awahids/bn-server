package authrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/authhandler"

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
		auth.POST("/google/oauth", rateLimiters.Google, handler.GoogleOAuthLogin)
		auth.POST("/refresh", rateLimiters.Refresh, handler.RefreshToken)
		auth.POST("/logout", rateLimiters.Logout, handler.Logout)
		auth.GET("/me", authMiddleware, handler.Me)
	}
}
