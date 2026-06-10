package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CtxRequestID = "request_id"

// RequestID reads the X-Request-ID header; generates a UUID v4 if absent.
// Echoes the final ID back in the response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(CtxRequestID, id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// GetRequestID returns the request ID stored in the Gin context by RequestID middleware.
func GetRequestID(c *gin.Context) string {
	v, _ := c.Get(CtxRequestID)
	s, _ := v.(string)
	return s
}
