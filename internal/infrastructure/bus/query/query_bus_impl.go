package query

import (
	"context"
	"fmt"
	"reflect"
)

type queryBus struct {
	handlers   map[reflect.Type]Handler
	middleware []Middleware
	logger     Logger
	metrics    MetricsReporter
}

func NewQueryBus(logger Logger, metrics MetricsReporter, middleware ...Middleware) Bus {
	return &queryBus{
		handlers:   make(map[reflect.Type]Handler),
		middleware: middleware,
		logger:     logger,
		metrics:    metrics,
	}
}

func (b *queryBus) Execute(ctx context.Context, query interface{}) (interface{}, error) {
	queryType := reflect.TypeOf(query)
	handler, exists := b.handlers[queryType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %v", queryType)
	}

	// 构建中间件链
	var next Handler = handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		m := b.middleware[i]
		current := next
		next = middlewareHandler{
			middleware: m,
			next:      current,
			query:     query,
		}
	}

	return next.Handle(ctx, query)
}

func (b *queryBus) Register(queryType interface{}, handler Handler) {
	t := reflect.TypeOf(queryType)
	if _, exists := b.handlers[t]; exists {
		panic(fmt.Sprintf("handler already registered for query type: %v", t))
	}
	b.handlers[t] = handler
}

type middlewareHandler struct {
	middleware Middleware
	next       Handler
	query      interface{}
}

func (h middlewareHandler) Handle(ctx context.Context, q interface{}) (interface{}, error) {
	return h.middleware.Execute(ctx, q, h.next)
}

// 中间件实现
type CachingMiddleware struct {
	cache  Cache
	logger Logger
}

func (m *CachingMiddleware) Execute(ctx context.Context, q interface{}, next Handler) (interface{}, error) {
	// 检查是否可缓存
	if cacheable, ok := q.(Cacheable); ok {
		key := cacheable.CacheKey()
		
		// 尝试从缓存获取
		if cached, err := m.cache.Get(ctx, key); err == nil && cached != nil {
			return cached, nil
		}

		// 执行查询
		result, err := next.Handle(ctx, q)
		if err != nil {
			return nil, err
		}

		// 缓存结果
		if err := m.cache.Set(ctx, key, result, cacheable.TTL()); err != nil {
			m.logger.Error("failed to cache query result", "error", err)
		}

		return result, nil
	}

	return next.Handle(ctx, q)
}

type LoggingMiddleware struct {
	logger Logger
}

func (m *LoggingMiddleware) Execute(ctx context.Context, q interface{}, next Handler) (interface{}, error) {
	start := time.Now()
	queryType := reflect.TypeOf(q).String()

	m.logger.Debug("executing query", "type", queryType)

	result, err := next.Handle(ctx, q)

	duration := time.Since(start)
	if err != nil {
		m.logger.Error("query failed",
			"type", queryType,
			"duration", duration,
			"error", err,
		)
	} else {
		m.logger.Debug("query completed",
			"type", queryType,
			"duration", duration,
		)
	}

	return result, err
}

type MetricsMiddleware struct {
	metrics MetricsReporter
}

func (m *MetricsMiddleware) Execute(ctx context.Context, q interface{}, next Handler) (interface{}, error) {
	queryType := reflect.TypeOf(q).String()
	timer := m.metrics.StartTimer("query_duration", "type", queryType)
	defer timer.Stop()

	result, err := next.Handle(ctx, q)
	if err != nil {
		m.metrics.IncrementCounter("query_failure", "type", queryType)
	} else {
		m.metrics.IncrementCounter("query_success", "type", queryType)
	}

	return result, err
} 