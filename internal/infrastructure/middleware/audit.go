package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/audit"
)

type auditResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var writeMethods = map[string]string{
	http.MethodPost:   "CREATE",
	http.MethodPut:    "UPDATE",
	http.MethodPatch:  "UPDATE",
	http.MethodDelete: "DELETE",
}

func AuditWrites(auditRepo audit.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		action, isWrite := writeMethods[c.Request.Method]
		if !isWrite {
			c.Next()
			return
		}

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		aw := &auditResponseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = aw

		c.Next()

		// Only audit successful writes (2xx)
		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			return
		}

		userID := GetUserID(c)
		entityType, entityID := parseEntity(c.FullPath(), c.Param("id"))

		entry := audit.Entry{
			Action:     action,
			EntityType: entityType,
			EntityID:   entityID,
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
			Changes:    sanitize(bodyBytes),
		}
		if userID != uuid.Nil {
			entry.UserID = &userID
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = auditRepo.Log(ctx, entry)
		}()
	}
}

func parseEntity(path, id string) (string, string) {
	if path == "" {
		return "unknown", ""
	}
	return path, id
}

func sanitize(body []byte) datatypes.JSON {
	if len(body) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}
	for _, key := range []string{"password", "password_hash", "token", "secret", "refresh_token"} {
		delete(m, key)
	}
	b, _ := json.Marshal(m)
	return datatypes.JSON(b)
}
