package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db      *sql.DB
	logger  Logger
	metrics MetricsReporter
}

func NewDatabase(config DatabaseConfig, logger Logger, metrics MetricsReporter) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 验证连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}, nil
}

func (d *Database) Begin() (*sql.Tx, error) {
	timer := d.metrics.StartTimer("database_begin_transaction")
	defer timer.Stop()

	tx, err := d.db.Begin()
	if err != nil {
		d.logger.Error("failed to begin transaction", "error", err)
		d.metrics.IncrementCounter("database_transaction_failure")
		return nil, err
	}

	d.metrics.IncrementCounter("database_transaction_success")
	return tx, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// 监控指标收集
func (d *Database) collectMetrics() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			stats := d.db.Stats()
			d.metrics.Gauge("database_open_connections", float64(stats.OpenConnections))
			d.metrics.Gauge("database_in_use_connections", float64(stats.InUse))
			d.metrics.Gauge("database_idle_connections", float64(stats.Idle))
			d.metrics.Gauge("database_wait_count", float64(stats.WaitCount))
		}
	}()
} 