package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Error     *ErrorInfo  `json:"error"`
	Timestamp int64       `json:"timestamp"`
}

type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func now() int64 { return time.Now().UnixMilli() }

func reqID(c *gin.Context) string {
	v, _ := c.Get("request_id")
	s, _ := v.(string)
	return s
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data, Timestamp: now()})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data, Timestamp: now()})
}

func BadRequest(c *gin.Context, code, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message, RequestID: reqID(c)},
		Timestamp: now(),
	})
}

func Unauthorized(c *gin.Context, code, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message, RequestID: reqID(c)},
		Timestamp: now(),
	})
}

func Forbidden(c *gin.Context, code, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message, RequestID: reqID(c)},
		Timestamp: now(),
	})
}

func NotFound(c *gin.Context, code, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message, RequestID: reqID(c)},
		Timestamp: now(),
	})
}

func InternalError(c *gin.Context, code, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message, RequestID: reqID(c)},
		Timestamp: now(),
	})
}
