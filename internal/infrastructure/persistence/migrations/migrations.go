package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFiles embed.FS

type Migrator struct {
	db     *sql.DB
	logger Logger
}

func NewMigrator(db *sql.DB, logger Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

func (m *Migrator) Up() error {
	driver, err := mysql.WithInstance(m.db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	source, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	migrate, err := migrate.NewWithInstance(
		"iofs",
		source,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("database migrations completed")
	return nil
}

// 数据库表结构
const (
	createUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id         VARCHAR(36)  NOT NULL PRIMARY KEY,
			email      VARCHAR(255) NOT NULL UNIQUE,
			password   VARCHAR(255) NOT NULL,
			name       VARCHAR(100) NOT NULL,
			bio        TEXT,
			avatar     VARCHAR(255),
			status     VARCHAR(20)  NOT NULL DEFAULT 'active',
			created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	createEventsTable = `
		CREATE TABLE IF NOT EXISTS events (
			id           VARCHAR(36)  NOT NULL PRIMARY KEY,
			aggregate_id VARCHAR(36)  NOT NULL,
			type         VARCHAR(100) NOT NULL,
			version      INT         NOT NULL,
			data         JSON        NOT NULL,
			occurred_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			published_at TIMESTAMP   NULL,
			INDEX idx_aggregate_id (aggregate_id),
			INDEX idx_type (type),
			UNIQUE KEY uk_aggregate_version (aggregate_id, version)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	createUserRolesTable = `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id    VARCHAR(36) NOT NULL,
			role       VARCHAR(20) NOT NULL,
			created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, role),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
) 