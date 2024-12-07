package command

type CommandBusFactory interface {
    CreateCommandBus() Bus
}

type commandBusFactory struct {
    config   CommandBusConfig
    logger   Logger
    metrics  MetricsReporter
    uow      UnitOfWork
    tracer   Tracer
}

func NewCommandBusFactory(
    config CommandBusConfig,
    logger Logger,
    metrics MetricsReporter,
    uow UnitOfWork,
    tracer Tracer,
) CommandBusFactory {
    return &commandBusFactory{
        config:  config,
        logger:  logger,
        metrics: metrics,
        uow:     uow,
        tracer:  tracer,
    }
}

func (f *commandBusFactory) CreateCommandBus() Bus {
    // 创建中间件
    middleware := f.createMiddleware()
    
    // 创建命令总线
    return NewCommandBus(f.logger, f.metrics, middleware...)
}

func (f *commandBusFactory) createMiddleware() []Middleware {
    var middleware []Middleware
    
    // 按配置添加中间件
    if f.config.Validation.Enabled {
        middleware = append(middleware, NewValidationMiddleware(f.config.Validation))
    }
    
    if f.config.Transaction.Enabled {
        middleware = append(middleware, NewTransactionMiddleware(f.uow, f.logger))
    }
    
    // ... 添加其他中间件
    
    return middleware
} 