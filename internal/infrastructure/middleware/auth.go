package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/response"
	"git.trovefin.com/poc/ctrlc-service-go/pkg/jwtutil"
)

const (
	ctxUserID = "user_id"
	ctxRoleID = "role_id"
)

func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "AUTH_UNAUTHORIZED", "missing bearer token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtutil.ValidateAccessToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			response.Unauthorized(c, "AUTH_TOKEN_INVALID", "invalid token")
			c.Abort()
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRoleID, claims.RoleID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(ctxUserID)
	id, _ := v.(uuid.UUID)
	return id
}

func GetRoleID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(ctxRoleID)
	id, _ := v.(uuid.UUID)
	return id
}
