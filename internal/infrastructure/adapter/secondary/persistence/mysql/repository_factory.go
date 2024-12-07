package mysql

import (
	"database/sql"
	"github.com/gohex/gohex/internal/application/port"
)

// RepositoryFactory 创建仓储实例的工厂
type RepositoryFactory struct {
	db      *sql.DB
	logger  Logger
	metrics MetricsReporter
}

// NewRepositoryFactory 创建仓储工厂实例
func NewRepositoryFactory(db *sql.DB, logger Logger, metrics MetricsReporter) *RepositoryFactory {
	return &RepositoryFactory{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}
}

// CreateUserRepository 创建用户仓储实例
func (f *RepositoryFactory) CreateUserRepository() port.UserRepository {
	return &userRepositoryImpl{
		db:      f.db,
		logger:  f.logger,
		metrics: f.metrics,
	}
}

// CreateEventStore 创建事件存储实例
func (f *RepositoryFactory) CreateEventStore() port.EventStore {
	return &mysqlEventStore{
		db:      f.db,
		logger:  f.logger,
		metrics: f.metrics,
	}
}