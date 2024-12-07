package query

import (
	"math"
)

type Result struct {
	Data  interface{}
	Error error
}

type PagedResult struct {
	Items      interface{} `json:"items"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

func NewPagedResult(items interface{}, total int64, page, pageSize int) *PagedResult {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &PagedResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
} 