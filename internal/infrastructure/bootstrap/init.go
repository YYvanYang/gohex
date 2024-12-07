package bootstrap

import (
	"database/sql"
	"github.com/gohex/gohex/internal/infrastructure/config"
	"github.com/gohex/gohex/internal/infrastructure/logger"
	"github.com/gohex/gohex/internal/infrastructure/metrics"
	"github.com/gohex/gohex/internal/infrastructure/tracing"
	"github.com/gohex/gohex/internal/infrastructure/adapter/secondary/persistence/mysql"
	"github.com/gohex/gohex/internal/infrastructure/adapter/secondary/cache/redis"
	"github.com/gohex/gohex/internal/application/command"
	"github.com/gohex/gohex/internal/application/query"
	"github.com/gohex/gohex/internal/domain/event"
)

func initLogger(cfg config.LogConfig) Logger {
	logger, err := logger.NewZapLogger(cfg)
	if err != nil {
		panic(err)
	}
	return logger
}

func initMetrics(cfg config.MetricsConfig) MetricsReporter {
	metrics, err := metrics.NewPrometheusMetrics(cfg)
	if err != nil {
		panic(err)
	}
	return metrics
}

func initTracer(cfg config.TracingConfig) Tracer {
	tracer, err := tracing.NewJaegerTracer(cfg)
	if err != nil {
		panic(err)
	}
	return tracer
}

func initDatabase(cfg config.DatabaseConfig) *sql.DB {
	db, err := mysql.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func initCache(cfg config.RedisConfig) Cache {
	cache, err := redis.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return cache
}

func initCommandBus(
	cfg *config.Config,
	logger Logger,
	metrics MetricsReporter,
	db *sql.DB,
) command.Bus {
	uow := mysql.NewUnitOfWork(db, logger)
	factory := command.NewCommandBusFactory(
		cfg.CommandBus,
		logger,
		metrics,
		uow,
		tracer,
	)
	return factory.CreateCommandBus()
}

func initQueryBus(
	cfg *config.Config,
	logger Logger,
	metrics MetricsReporter,
	cache Cache,
) query.Bus {
	factory := query.NewQueryBusFactory(
		cfg.QueryBus,
		logger,
		metrics,
		cache,
		tracer,
	)
	return factory.CreateQueryBus()
}

func initEventBus(
	cfg *config.Config,
	logger Logger,
	metrics MetricsReporter,
) event.Bus {
	factory := event.NewEventBusFactory(
		cfg.EventBus,
		logger,
		metrics,
	)
	return factory.CreateEventBus()
} 