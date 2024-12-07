package mysql

import (
	"database/sql"
	"github.com/your-org/your-project/internal/application/port"
)

type repositoryFactory struct {
	logger  Logger
	metrics MetricsReporter
}

func NewRepositoryFactory(logger Logger, metrics MetricsReporter) *repositoryFactory {
	return &repositoryFactory{
		logger:  logger,
		metrics: metrics,
	}
}

func (f *repositoryFactory) CreateUserRepository(tx *sql.Tx) port.UserRepository {
	return NewUserRepository(tx, f.logger, f.metrics)
}

func (f *repositoryFactory) CreateEventStore(tx *sql.Tx) port.EventStore {
	return NewEventStore(tx, f.logger, f.metrics)
} 