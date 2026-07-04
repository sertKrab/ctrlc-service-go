package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/logging"
)

const CtxRequestID = "request_id"

// RequestID reads the X-Request-ID header; generates a UUID v4 if absent.
// Echoes the final ID back in the response header, and attaches it to the
// request's context.Context so logging.FromContext(ctx) includes it in
// every log call made from handlers, usecases, and repositories.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(CtxRequestID, id)
		c.Header("X-Request-ID", id)
		c.Request = c.Request.WithContext(logging.WithRequestID(c.Request.Context(), id))
		c.Next()
	}
}

// GetRequestID returns the request ID stored in the Gin context by RequestID middleware.
func GetRequestID(c *gin.Context) string {
	v, _ := c.Get(CtxRequestID)
	s, _ := v.(string)
	return s
}
