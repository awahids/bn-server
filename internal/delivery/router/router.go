package router

import (
	"net/http"
	"time"

	"bn-mobile/server/configs"
	apphandler "bn-mobile/server/internal/delivery/handlers/appHandler"
	authhandler "bn-mobile/server/internal/delivery/handlers/authHandler"
	"bn-mobile/server/internal/delivery/middleware"
	approuter "bn-mobile/server/internal/delivery/router/appRouter"
	authrouter "bn-mobile/server/internal/delivery/router/authRouter"
	publicrouter "bn-mobile/server/internal/delivery/router/publicRouter"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check server status.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]any
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ok",
	})
}

func NewRouter(
	cfg *configs.Config,
	authHandler *authhandler.AuthHandler,
	appHandler *apphandler.AppHandler,
	publicHandler *apphandler.PublicHandler,
) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/health", healthCheck)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := engine.Group("/api/v1")
	apiLegacy := engine.Group("/api")
	authMiddleware := middleware.AuthMiddleware(cfg.JWT.Secret)
	apiV1.GET("/health", healthCheck)

	publicrouter.RegisterPublicRoutes(apiV1, publicHandler)
	publicrouter.RegisterPublicRoutes(apiLegacy, publicHandler)

	authWindow := time.Duration(cfg.RateLimiter.AuthWindowSeconds) * time.Second
	googleLimiter := middleware.NewIPRateLimiter(cfg.RateLimiter.AuthGoogleLimit, authWindow)
	refreshLimiter := middleware.NewIPRateLimiter(cfg.RateLimiter.AuthRefreshLimit, authWindow)
	logoutLimiter := middleware.NewIPRateLimiter(cfg.RateLimiter.AuthLogoutLimit, authWindow)

	authrouter.RegisterAuthRoutes(
		apiV1,
		authHandler,
		authMiddleware,
		authrouter.RateLimitMiddlewares{
			Google:  googleLimiter.Middleware("auth_google"),
			Refresh: refreshLimiter.Middleware("auth_refresh"),
			Logout:  logoutLimiter.Middleware("auth_logout"),
		},
	)
	approuter.RegisterAppRoutes(apiV1, appHandler, authMiddleware)
	approuter.RegisterAppRoutes(apiLegacy, appHandler, authMiddleware)

	return engine
}
