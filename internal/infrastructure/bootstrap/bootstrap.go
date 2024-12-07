package bootstrap

import (
	"context"
	"net/http"
	"github.com/gohex/gohex/internal/infrastructure/config"
	"github.com/gohex/gohex/internal/infrastructure/container"
	"github.com/gohex/gohex/internal/application/command"
	"github.com/gohex/gohex/internal/application/query"
	"github.com/gohex/gohex/internal/domain/event"
	"github.com/gohex/gohex/internal/domain/service"
)

type Application struct {
	config      *config.Config
	logger      Logger
	metrics     MetricsReporter
	tracer      Tracer
	commandBus  command.Bus
	queryBus    query.Bus
	eventBus    event.Bus
	httpServer  *http.Server
}

func NewApplication(configPath string) (*Application, error) {
	// 1. 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 2. 初始化基础设施
	logger := initLogger(cfg.Log)
	metrics := initMetrics(cfg.Metrics)
	tracer := initTracer(cfg.Tracing)

	// 3. 初始化数据库连接
	db := initDatabase(cfg.Database)
	cache := initCache(cfg.Redis)

	// 4. 创建仓储
	userRepo := mysql.NewUserRepository(db, logger, metrics)
	eventStore := mysql.NewEventStore(db, logger, metrics)

	// 5. 创建服务
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewAuthService(cfg.JWT, logger)
	emailService := service.NewEmailService(cfg.SMTP, logger)

	// 6. 创建命令和查询总线
	commandBus := initCommandBus(cfg, logger, metrics, db)
	queryBus := initQueryBus(cfg, logger, metrics, cache)
	eventBus := initEventBus(cfg, logger, metrics)

	// 7. 创建 HTTP 服务器
	httpServer := initHTTPServer(cfg.HTTP, logger, metrics)

	return &Application{
		config:     cfg,
		logger:     logger,
		metrics:    metrics,
		tracer:     tracer,
		commandBus: commandBus,
		queryBus:   queryBus,
		eventBus:   eventBus,
		httpServer: httpServer,
	}, nil
}

func (app *Application) Start(ctx context.Context) error {
	// 1. 启动追踪器
	if err := app.tracer.Start(ctx); err != nil {
		return err
	}

	// 2. 启动指标收集器
	if err := app.metrics.Start(ctx); err != nil {
		return err
	}

	// 3. 启动事件总线
	if err := app.eventBus.Start(ctx); err != nil {
		return err
	}

	// 4. 启动 HTTP 服务器
	return app.httpServer.Start()
}

func (app *Application) Stop(ctx context.Context) error {
	// 1. 停止 HTTP 服务器
	if err := app.httpServer.Stop(ctx); err != nil {
		app.logger.Error("failed to stop http server", "error", err)
	}

	// 2. 停止事件总线
	if err := app.eventBus.Stop(ctx); err != nil {
		app.logger.Error("failed to stop event bus", "error", err)
	}

	// 3. 停止指标收集器
	if err := app.metrics.Stop(ctx); err != nil {
		app.logger.Error("failed to stop metrics reporter", "error", err)
	}

	// 4. 停止追踪器
	if err := app.tracer.Stop(ctx); err != nil {
		app.logger.Error("failed to stop tracer", "error", err)
	}

	return nil
} 