package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

// PaginatedResponse wraps paginated data
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// GetPaginationParams extracts pagination from query parameters
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Enforce limits
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Max 100 items per page
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// CreatePaginatedResponse creates a paginated response
func CreatePaginatedResponse(data interface{}, page, pageSize int, totalCount int64) PaginatedResponse {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
