package query

import (
	"context"
	"time"
)

// Bus 定义查询总线接口
type Bus interface {
	// Execute 执行查询并返回结果
	Execute(ctx context.Context, query interface{}) (interface{}, error)
	// Register 注册查询处理器
	Register(queryType interface{}, handler Handler)
}

// Handler 定义查询处理器接口
type Handler interface {
	// Handle 处理查询并返回结果
	Handle(ctx context.Context, query interface{}) (interface{}, error)
}

// Middleware 定义查询中间件接口
type Middleware interface {
	// Execute 执行中间件逻辑
	Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error)
}

// 添加查询上下文
type QueryContext struct {
	Context    context.Context
	StartTime  time.Time
	Timeout    time.Duration
	RetryCount int
	CacheKey   string
	TraceID    string
}

func NewQueryContext(ctx context.Context) *QueryContext {
	return &QueryContext{
		Context:   ctx,
		StartTime: time.Now(),
	}
} 