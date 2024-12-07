package command

import (
	"context"
	"fmt"
	"reflect"
)

type commandBus struct {
	handlers   map[reflect.Type]Handler
	middleware []Middleware
	logger     Logger
	metrics    MetricsReporter
}

func NewCommandBus(logger Logger, metrics MetricsReporter, middleware ...Middleware) Bus {
	return &commandBus{
		handlers:   make(map[reflect.Type]Handler),
		middleware: middleware,
		logger:     logger,
		metrics:    metrics,
	}
}

func (b *commandBus) Dispatch(ctx context.Context, cmd interface{}) (interface{}, error) {
	cmdType := reflect.TypeOf(cmd)
	handler, exists := b.handlers[cmdType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for command type: %v", cmdType)
	}

	// 构建中间件链
	var next Handler = handler
	for i := len(b.middleware) - 1; i >= 0; i-- {
		m := b.middleware[i]
		current := next
		next = middlewareHandler{
			middleware: m,
			next:      current,
			command:   cmd,
		}
	}

	return next.Handle(ctx, cmd)
}

func (b *commandBus) Register(cmdType interface{}, handler Handler) {
	t := reflect.TypeOf(cmdType)
	if _, exists := b.handlers[t]; exists {
		panic(fmt.Sprintf("handler already registered for command type: %v", t))
	}
	b.handlers[t] = handler
} 