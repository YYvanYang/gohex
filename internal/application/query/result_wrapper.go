package query

import "time"

type QueryResult struct {
    Data      interface{}
    Error     error
    Duration  time.Duration
    CacheHit  bool
    TraceID   string
    Timestamp time.Time
}

func NewQueryResult(data interface{}, err error, duration time.Duration) *QueryResult {
    return &QueryResult{
        Data:      data,
        Error:     err,
        Duration:  duration,
        Timestamp: time.Now(),
    }
}

type PagedQueryResult struct {
    *QueryResult
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalItems int64       `json:"total_items"`
    TotalPages int         `json:"total_pages"`
    HasMore    bool        `json:"has_more"`
    Items      interface{} `json:"items"`
}

func NewPagedQueryResult(items interface{}, total int64, page, pageSize int, duration time.Duration) *PagedQueryResult {
    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
    return &PagedQueryResult{
        QueryResult: NewQueryResult(items, nil, duration),
        Page:       page,
        PageSize:   pageSize,
        TotalItems: total,
        TotalPages: totalPages,
        HasMore:    page < totalPages,
        Items:      items,
    }
} 