package query

import (
	"context"
	"reflect"
	"time"

	"github.com/your-project/port"
)

type Middleware interface {
	Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error)
}

type ValidationMiddleware struct {
	validator Validator
	logger    Logger
}

func (m *ValidationMiddleware) Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error) {
	if v, ok := query.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, errors.NewValidationError(err.Error())
		}
	}
	return next.Handle(ctx, query)
}

type CacheMiddleware struct {
	cache    port.Cache
	logger   Logger
	metrics  MetricsReporter
}

func (m *CacheMiddleware) Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error) {
	if cacheable, ok := query.(Cacheable); ok {
		// 尝试从缓存获取
		if result, err := m.cache.Get(ctx, cacheable.CacheKey()); err == nil && result != nil {
			m.metrics.IncrementCounter("cache_hit", "type", reflect.TypeOf(query).String())
			return result, nil
		}
		m.metrics.IncrementCounter("cache_miss", "type", reflect.TypeOf(query).String())

		// 执行查询
		result, err := next.Handle(ctx, query)
		if err != nil {
			return nil, err
		}

		// 更新缓存
		if err := m.cache.Set(ctx, cacheable.CacheKey(), result, cacheable.TTL()); err != nil {
			m.logger.Error("failed to cache query result", "error", err)
		}

		return result, nil
	}
	return next.Handle(ctx, query)
}

// 添加工厂方法
func NewValidationMiddleware(validator Validator, logger Logger) *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator,
		logger:    logger,
	}
}

func NewCacheMiddleware(cache port.Cache, logger Logger, metrics MetricsReporter) *CacheMiddleware {
	return &CacheMiddleware{
		cache:   cache,
		logger:  logger,
		metrics: metrics,
	}
}

// 添加日志中间件
type LoggingMiddleware struct {
	logger  Logger
	metrics MetricsReporter
}

func NewLoggingMiddleware(logger Logger, metrics MetricsReporter) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger:  logger,
		metrics: metrics,
	}
}

func (m *LoggingMiddleware) Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error) {
	start := time.Now()
	queryType := reflect.TypeOf(query).String()

	m.logger.Debug("executing query", "type", queryType)

	result, err := next.Handle(ctx, query)

	duration := time.Since(start)
	if err != nil {
		m.logger.Error("query failed",
			"type", queryType,
			"duration", duration,
			"error", err,
		)
		m.metrics.IncrementCounter("query_failure", "type", queryType)
	} else {
		m.logger.Debug("query completed",
			"type", queryType,
			"duration", duration,
		)
		m.metrics.IncrementCounter("query_success", "type", queryType)
	}

	return result, err
}

// 添加新的中间件
type RetryMiddleware struct {
	maxRetries int
	backoff    time.Duration
	logger     Logger
}

func (m *RetryMiddleware) Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error) {
	var lastErr error
	for i := 0; i < m.maxRetries; i++ {
		result, err := next.Handle(ctx, query)
		if err == nil {
			return result, nil
		}
		lastErr = err
		time.Sleep(m.backoff * time.Duration(i+1))
	}
	return nil, lastErr
}

type TimeoutMiddleware struct {
	timeout time.Duration
	logger  Logger
}

func (m *TimeoutMiddleware) Execute(ctx context.Context, query interface{}, next Handler) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	done := make(chan struct {
		result interface{}
		err    error
	})

	go func() {
		result, err := next.Handle(ctx, query)
		done <- struct {
			result interface{}
			err    error
		}{result, err}
	}()

	select {
	case <-ctx.Done():
		return nil, errors.ErrQueryTimeout
	case res := <-done:
		return res.result, res.err
	}
} 