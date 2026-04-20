package router

import (
	"net/http"
	"time"

	"github.com/awahids/bn-server/configs"
	"github.com/awahids/bn-server/internal/delivery/handlers/authhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/bookmarkhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/dhikrhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/progresshandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/publichandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/quizhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/userhandler"
	"github.com/awahids/bn-server/internal/delivery/middleware"
	authrouter "github.com/awahids/bn-server/internal/delivery/router/authrouter"
	bookmarkrouter "github.com/awahids/bn-server/internal/delivery/router/bookmarkrouter"
	dhikrrouter "github.com/awahids/bn-server/internal/delivery/router/dhikrrouter"
	progressrouter "github.com/awahids/bn-server/internal/delivery/router/progressrouter"
	publicrouter "github.com/awahids/bn-server/internal/delivery/router/publicrouter"
	quizrouter "github.com/awahids/bn-server/internal/delivery/router/quizrouter"
	userrouter "github.com/awahids/bn-server/internal/delivery/router/userrouter"

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
	userHandler *userhandler.UserHandler,
	progressHandler *progresshandler.ProgressHandler,
	bookmarkHandler *bookmarkhandler.BookmarkHandler,
	dhikrHandler *dhikrhandler.DhikrHandler,
	quizHandler *quizhandler.QuizHandler,
	publicHandler *publichandler.PublicHandler,
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

	userrouter.RegisterUserRoutes(apiV1, userHandler, authMiddleware)
	progressrouter.RegisterProgressRoutes(apiV1, progressHandler, authMiddleware)
	bookmarkrouter.RegisterBookmarkRoutes(apiV1, bookmarkHandler, authMiddleware)
	dhikrrouter.RegisterDhikrRoutes(apiV1, dhikrHandler, authMiddleware)
	quizrouter.RegisterQuizRoutes(apiV1, quizHandler, authMiddleware)

	userrouter.RegisterUserRoutes(apiLegacy, userHandler, authMiddleware)
	progressrouter.RegisterProgressRoutes(apiLegacy, progressHandler, authMiddleware)
	bookmarkrouter.RegisterBookmarkRoutes(apiLegacy, bookmarkHandler, authMiddleware)
	dhikrrouter.RegisterDhikrRoutes(apiLegacy, dhikrHandler, authMiddleware)
	quizrouter.RegisterQuizRoutes(apiLegacy, quizHandler, authMiddleware)

	return engine
}
