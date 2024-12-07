package command

import (
	"context"
	"reflect"
	"time"
	"github.com/your-org/your-project/pkg/errors"
)

// Bus 定义命令总线接口
type Bus interface {
	// Dispatch 分发命令并返回结果
	Dispatch(ctx context.Context, command interface{}) (interface{}, error)
	// Register 注册命令处理器
	Register(commandType interface{}, handler Handler)
}

// Handler 定义命令处理器接口
type Handler interface {
	// Handle 处理命令并返回结果
	Handle(ctx context.Context, command interface{}) (interface{}, error)
}

// Middleware 定义命令中间件接口
type Middleware interface {
	// Execute 执行中间件逻辑
	Execute(ctx context.Context, command interface{}, next Handler) (interface{}, error)
}

// Validator 定义命令验证器接口
type Validator interface {
	// Validate 验证命令
	Validate(command interface{}) error
}

// TransactionHandler 定义事务处理器接口
type TransactionHandler interface {
	Handler
	// WithTransaction 指示该处理器需要在事务中执行
	WithTransaction() bool
}

// EventHandler 定义事件处理器接口
type EventHandler interface {
	Handler
	// Events 返回命令执行后产生的事件
	Events() []event.Event
}

// ValidatorMiddleware 验证中间件
type ValidatorMiddleware struct {
	validator Validator
	logger    Logger
}

func (m *ValidatorMiddleware) Execute(ctx context.Context, cmd interface{}, next Handler) (interface{}, error) {
	if err := m.validator.Struct(cmd); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}
	return next.Handle(ctx, cmd)
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger Logger
}

func (m *LoggingMiddleware) Execute(ctx context.Context, cmd interface{}, next Handler) (interface{}, error) {
	cmdType := reflect.TypeOf(cmd).String()
	m.logger.Info("executing command", "type", cmdType)

	start := time.Now()
	result, err := next.Handle(ctx, cmd)
	duration := time.Since(start)

	if err != nil {
		m.logger.Error("command failed",
			"type", cmdType,
			"duration", duration,
			"error", err,
		)
	} else {
		m.logger.Info("command completed",
			"type", cmdType,
			"duration", duration,
		)
	}

	return result, err
}

// MetricsMiddleware 指标中间件
type MetricsMiddleware struct {
	metrics MetricsReporter
}

func (m *MetricsMiddleware) Execute(ctx context.Context, cmd interface{}, next Handler) (interface{}, error) {
	cmdType := reflect.TypeOf(cmd).String()
	timer := m.metrics.StartTimer("command_duration", "type", cmdType)
	defer timer.Stop()

	result, err := next.Handle(ctx, cmd)
	if err != nil {
		m.metrics.IncrementCounter("command_failure", "type", cmdType)
	} else {
		m.metrics.IncrementCounter("command_success", "type", cmdType)
	}

	return result, err
} 