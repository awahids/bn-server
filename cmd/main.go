// @title BN Mobile API
// @version 1.0
// @description API backend Belajar Ngaji (Go + Gin).
// @BasePath /api/v1
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer access token. Format: "Bearer {token}"
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/awahids/bn-server/configs"
	_ "github.com/awahids/bn-server/docs"
	"github.com/awahids/bn-server/internal/delivery/handlers/aihandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/authhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/bookmarkhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/dhikrhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/habithandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/progresshandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/publichandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/pushhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/quizhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/schoolhandler"
	"github.com/awahids/bn-server/internal/delivery/handlers/userhandler"
	"github.com/awahids/bn-server/internal/delivery/router"
	apprepo "github.com/awahids/bn-server/internal/domain/repositories/apprepo"
	authrepo "github.com/awahids/bn-server/internal/domain/repositories/authrepo"
	aiservice "github.com/awahids/bn-server/internal/domain/services/aiservice"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	authservice "github.com/awahids/bn-server/internal/domain/services/authservice"
	publicservice "github.com/awahids/bn-server/internal/domain/services/publicservice"
	pushservice "github.com/awahids/bn-server/internal/domain/services/pushservice"
	"github.com/awahids/bn-server/internal/infrastructure/database"
	"github.com/awahids/bn-server/internal/infrastructure/googleauth"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}
	defer database.Close(db)

	authRepository := authrepo.NewAuthRepository(db)
	googleAuthProvider := googleauth.NewGoogleAuthProvider(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		nil,
	)
	authService := authservice.NewAuthService(
		authRepository,
		authservice.TokenConfig{
			Issuer:          cfg.JWT.Issuer,
			Secret:          cfg.JWT.Secret,
			AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
			RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		},
		googleAuthProvider,
	)
	authHandler := authhandler.NewAuthHandler(authService, cfg.AuthCookie)

	appRepository := apprepo.NewAppRepository(db)
	appService := appservice.NewAppService(appRepository)
	userHandler := userhandler.NewUserHandler(appService)
	progressHandler := progresshandler.NewProgressHandler(appService)
	bookmarkHandler := bookmarkhandler.NewBookmarkHandler(appService)
	dhikrHandler := dhikrhandler.NewDhikrHandler(appService)
	habitHandler := habithandler.NewHabitHandler(appService)
	pushHandler := pushhandler.NewPushHandler(appService, cfg.Push.VAPIDPublicKey, cfg.Push.Enabled)
	schoolHandler := schoolhandler.NewSchoolHandler(appService)
	quizHandler := quizhandler.NewQuizHandler(appService)

	publicService := publicservice.NewPublicService(nil)
	publicHandler := publichandler.NewPublicHandler(publicService)

	aiService := aiservice.NewAIService(cfg)
	aiHandler := aihandler.NewAIHandler(aiService)

	engine := router.NewRouter(
		cfg,
		authHandler,
		userHandler,
		progressHandler,
		bookmarkHandler,
		dhikrHandler,
		habitHandler,
		pushHandler,
		schoolHandler,
		quizHandler,
		publicHandler,
		aiHandler,
	)

	pushScheduler := pushservice.NewPushService(appRepository, cfg.Push)
	pushCtx, stopPushScheduler := context.WithCancel(context.Background())
	var bgWorkers sync.WaitGroup
	if pushScheduler.IsEnabled() {
		bgWorkers.Add(1)
		go func() {
			defer bgWorkers.Done()
			pushScheduler.Run(pushCtx)
		}()
	}

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		// AI requests can take longer than standard API requests; keep write timeout above AI client timeout.
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server running on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	stopPushScheduler()
	bgWorkers.Wait()

	log.Println("server stopped")
}
