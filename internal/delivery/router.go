package delivery

import (
	"fmt"
	"strings"
	"time"

	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/audit"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/middleware"
	"git.trovefin.com/poc/ctrlc-service-go/internal/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authHandler *AuthHandler, auditRepo audit.Repository) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		tcID := param.Request.Header.Get("X-Ctrlc-TC-ID")
		flowID := param.Request.Header.Get("X-Flow-Id")
		return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %#v | tc_id=%s flow_id=%s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			tcID,
			flowID,
		)
	}), gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = strings.Split(cfg.CORSOrigins, ",")
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Request-ID", "X-Idempotency-Key", "X-Ctrlc-TC-ID", "X-Flow-Id"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestID())

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(middleware.AuditWrites(auditRepo))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", middleware.RequireAuth(cfg), authHandler.Logout)
			auth.GET("/me", middleware.RequireAuth(cfg), authHandler.GetMe)
		}
	}

	return r
}
