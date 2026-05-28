package utils

import "github.com/gin-gonic/gin"

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type PageRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type PagedResult struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func ParsePage(c *gin.Context) PageRequest {
	var req PageRequest
	_ = c.ShouldBindQuery(&req)
	if req.Page < 1 {
		req.Page = DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = DefaultPageSize
	}
	if req.PageSize > MaxPageSize {
		req.PageSize = MaxPageSize
	}
	return req
}

func NewPagedResult(items interface{}, total int64, req PageRequest) PagedResult {
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}
	return PagedResult{
		Items: items, Total: total,
		Page: req.Page, PageSize: req.PageSize,
		TotalPages: totalPages,
	}
}

func Offset(req PageRequest) int {
	return (req.Page - 1) * req.PageSize
}
