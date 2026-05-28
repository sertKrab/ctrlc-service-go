package main

import (
	"log"

	"github.com/joho/godotenv"
	authusecase "git.trovefin.com/poc/ctrlc-service-go/internal/application/usecase/auth"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/delivery"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/cache"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/database"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/errcache"
	"git.trovefin.com/poc/ctrlc-service-go/internal/repository"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Config] %v", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("[DB] %v", err)
	}

	_ = cache.NewCache(cfg)

	// Repositories
	userRepo    := repository.NewUserRepo(db)
	sessionRepo := repository.NewRefreshSessionRepo(db)
	auditRepo   := repository.NewAuditRepo(db)
	errmsgRepo  := repository.NewErrMsgRepo(db)

	// Error cache
	resolver, err := errcache.New(errmsgRepo, cfg.DefaultLocale)
	if err != nil {
		log.Fatalf("[ErrCache] %v", err)
	}

	// Usecases
	loginUC   := authusecase.NewLoginUseCase(userRepo, sessionRepo, auditRepo, cfg, resolver)
	refreshUC := authusecase.NewRefreshUseCase(sessionRepo, userRepo, cfg, resolver)
	logoutUC  := authusecase.NewLogoutUseCase(sessionRepo, auditRepo, resolver)
	getMeUC   := authusecase.NewGetMeUseCase(userRepo, resolver)

	// Handlers + router
	authHandler := delivery.NewAuthHandler(loginUC, refreshUC, logoutUC, getMeUC, cfg)
	r := delivery.SetupRouter(cfg, authHandler, auditRepo)

	log.Printf("[API] starting on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("[API] %v", err)
	}
}
