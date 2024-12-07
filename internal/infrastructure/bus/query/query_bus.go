package query

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/your-org/your-project/internal/application/port/input/query"
)

type queryBus struct {
	handlers sync.Map
	logger   Logger
	metrics  MetricsReporter
}

func NewQueryBus(logger Logger, metrics MetricsReporter) query.QueryBus {
	return &queryBus{
		logger:  logger,
		metrics: metrics,
	}
}

func (b *queryBus) Execute(ctx context.Context, q interface{}) (interface{}, error) {
	span, ctx := tracer.StartSpan(ctx, "queryBus.Execute")
	defer span.End()

	timer := b.metrics.StartTimer("query_execution_duration")
	defer timer.Stop()

	// 获取查询类型
	queryType := reflect.TypeOf(q)

	// 查找对应的处理器
	handler, ok := b.handlers.Load(queryType)
	if !ok {
		err := fmt.Errorf("no handler registered for query type: %v", queryType)
		b.logger.Error("query execution failed", "error", err)
		return nil, err
	}

	// 执行查询
	result, err := handler.(query.QueryHandler).Handle(ctx, q)
	if err != nil {
		b.metrics.IncrementCounter("query_execution_failure")
		b.logger.Error("query execution failed",
			"query_type", queryType.String(),
			"error", err,
		)
		return nil, err
	}

	b.metrics.IncrementCounter("query_execution_success")
	return result, nil
}

func (b *queryBus) Register(queryType interface{}, handler query.QueryHandler) {
	b.handlers.Store(reflect.TypeOf(queryType), handler)
} 