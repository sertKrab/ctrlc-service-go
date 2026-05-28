package usecase

import "github.com/gin-gonic/gin"

type RequestMeta struct {
	IPAddress string
	UserAgent string
}

func MetaFromGin(c *gin.Context) RequestMeta {
	return RequestMeta{
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
}
