package main

import (
	"log"
	"log/slog"
	"os"

	authusecase "git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/delivery"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/cache"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/database"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/errcache"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/logging"
	"git.trovefin.com/poc/ctrlc-service-go/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Config load itself isn't logged through slog yet (LOG_LEVEL comes
	// from cfg), so keep stdlib log for this one fatal path only.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Config] %v", err)
	}

	logging.Init(cfg.LogLevel)

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("[DB] connect failed", "error", err)
		os.Exit(1)
	}

	cacheInstance, err := cache.NewCache(cfg)
	if err != nil {
		slog.Error("[Cache] durable cache unavailable", "error", err)
		os.Exit(1)
	}
	slog.Info("[Cache] ready", "durable", cacheInstance.Durable())

	// Repositories
	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewRefreshSessionRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	errmsgRepo := repository.NewErrMsgRepo(db)

	// Error cache
	resolver, err := errcache.New(errmsgRepo, cfg.DefaultLocale)
	if err != nil {
		log.Fatalf("[ErrCache] %v", err)
	}

	// Usecases
	loginUC := authusecase.NewLoginUseCase(userRepo, sessionRepo, auditRepo, cfg, resolver)
	refreshUC := authusecase.NewRefreshUseCase(sessionRepo, userRepo, cfg, resolver)
	logoutUC := authusecase.NewLogoutUseCase(sessionRepo, auditRepo, resolver)
	getMeUC := authusecase.NewGetMeUseCase(userRepo, resolver)

	// Handlers + router
	authHandler := delivery.NewAuthHandler(loginUC, refreshUC, logoutUC, getMeUC, cfg)
	r := delivery.SetupRouter(cfg, authHandler, auditRepo)

	slog.Info("[API] starting", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		slog.Error("[API] server stopped", "error", err)
		os.Exit(1)
	}
}
