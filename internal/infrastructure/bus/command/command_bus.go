package command

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/gohex/gohex/internal/application/command"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/pkg/tracer"
)

type commandBus struct {
	handlers sync.Map
	logger   Logger
	metrics  MetricsReporter
}

func NewCommandBus(logger Logger, metrics MetricsReporter) command.CommandBus {
	return &commandBus{
		logger:  logger,
		metrics: metrics,
	}
}

func (b *commandBus) Dispatch(ctx context.Context, cmd interface{}) (interface{}, error) {
	span, ctx := tracer.StartSpan(ctx, "commandBus.Dispatch")
	defer span.End()

	timer := b.metrics.StartTimer("command_dispatch_duration")
	defer timer.Stop()

	// 获取命令类型
	cmdType := reflect.TypeOf(cmd)

	// 查找对应的处理器
	handler, ok := b.handlers.Load(cmdType)
	if !ok {
		err := fmt.Errorf("no handler registered for command type: %v", cmdType)
		b.logger.Error("command dispatch failed", "error", err)
		return nil, err
	}

	// 执行命令
	result, err := handler.(command.CommandHandler).Handle(ctx, cmd)
	if err != nil {
		b.metrics.IncrementCounter("command_dispatch_failure")
		b.logger.Error("command execution failed",
			"command_type", cmdType.String(),
			"error", err,
		)
		return nil, err
	}

	b.metrics.IncrementCounter("command_dispatch_success")
	return result, nil
}

func (b *commandBus) Register(cmdType interface{}, handler command.CommandHandler) {
	b.handlers.Store(reflect.TypeOf(cmdType), handler)
} 