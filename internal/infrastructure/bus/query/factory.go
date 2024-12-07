package query

import (
	"github.com/your-org/your-project/internal/application/port/input/query"
)

type QueryBusFactory interface {
	CreateQueryBus() query.Bus
}

type queryBusFactory struct {
	config   QueryBusConfig
	logger   Logger
	metrics  MetricsReporter
	cache    Cache
	tracer   Tracer
}

func NewQueryBusFactory(
	config QueryBusConfig,
	logger Logger,
	metrics MetricsReporter,
	cache Cache,
	tracer Tracer,
) QueryBusFactory {
	return &queryBusFactory{
		config:  config,
		logger:  logger,
		metrics: metrics,
		cache:   cache,
		tracer:  tracer,
	}
}

func (f *queryBusFactory) CreateQueryBus() query.Bus {
	// 创建中间件
	middleware := f.createMiddleware()
	
	// 创建查询总线
	return NewQueryBus(f.logger, f.metrics, middleware...)
}

func (f *queryBusFactory) createMiddleware() []query.Middleware {
	var middleware []query.Middleware
	
	// 按配置添加中间件
	if f.config.Validation.Enabled {
		middleware = append(middleware, NewValidationMiddleware(f.config.Validation))
	}
	
	if f.config.Cache.Enabled {
		middleware = append(middleware, NewCacheMiddleware(f.cache, f.logger))
	}
	
	// ... 添加其他中间件
	
	return middleware
} 