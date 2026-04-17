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
	"syscall"
	"time"

	"bn-mobile/server/configs"
	_ "bn-mobile/server/docs"
	apphandler "bn-mobile/server/internal/delivery/handlers/appHandler"
	authhandler "bn-mobile/server/internal/delivery/handlers/authHandler"
	"bn-mobile/server/internal/delivery/router"
	apprepo "bn-mobile/server/internal/domain/repositories/appRepo"
	authrepo "bn-mobile/server/internal/domain/repositories/authRepo"
	appservice "bn-mobile/server/internal/domain/services/appService"
	authservice "bn-mobile/server/internal/domain/services/authService"
	"bn-mobile/server/internal/infrastructure/database"
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
	authService := authservice.NewAuthService(authRepository, cfg)
	authHandler := authhandler.NewAuthHandler(authService, cfg.AuthCookie)
	appRepository := apprepo.NewAppRepository(db)
	appService := appservice.NewAppService(appRepository)
	appHandler := apphandler.NewAppHandler(appService)
	publicHandler := apphandler.NewPublicHandler()

	engine := router.NewRouter(cfg, authHandler, appHandler, publicHandler)

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
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

	log.Println("server stopped")
}
