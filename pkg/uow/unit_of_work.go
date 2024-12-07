package uow

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gohex/gohex/internal/application/port"
	"github.com/gohex/gohex/pkg/errors"
)

type UnitOfWork struct {
	db      *sql.DB
	tx      *sql.Tx
	logger  Logger
	metrics MetricsReporter
}

func NewUnitOfWork(db *sql.DB, logger Logger, metrics MetricsReporter) *UnitOfWork {
	return &UnitOfWork{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}
}

func (u *UnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	timer := u.metrics.StartTimer("uow_begin_duration")
	defer timer.Stop()

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		u.metrics.IncrementCounter("uow_begin_error")
		return ctx, fmt.Errorf("failed to begin transaction: %w", err)
	}

	u.tx = tx
	ctx = context.WithValue(ctx, txKey{}, tx)
	u.metrics.IncrementCounter("uow_begin_success")
	return ctx, nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	timer := u.metrics.StartTimer("uow_commit_duration")
	defer timer.Stop()

	if u.tx == nil {
		return fmt.Errorf("no active transaction")
	}

	if err := u.tx.Commit(); err != nil {
		u.metrics.IncrementCounter("uow_commit_error")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.tx = nil
	u.metrics.IncrementCounter("uow_commit_success")
	return nil
}

func (u *UnitOfWork) Rollback(ctx context.Context) error {
	timer := u.metrics.StartTimer("uow_rollback_duration")
	defer timer.Stop()

	if u.tx == nil {
		return nil
	}

	if err := u.tx.Rollback(); err != nil {
		u.metrics.IncrementCounter("uow_rollback_error")
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	u.tx = nil
	u.metrics.IncrementCounter("uow_rollback_success")
	return nil
}

// 添加事务上下文
type txKey struct{}

func FromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

func WithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// 添加事务包装方法
func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	// 如果已经在事务中，直接执行
	if tx, ok := FromContext(ctx); ok {
		return fn(ctx)
	}

	// 开始新事务
	ctx, err := u.Begin(ctx)
	if err != nil {
		return err
	}

	// 确保事务结束
	defer func() {
		if r := recover(); r != nil {
			u.Rollback(ctx)
			panic(r)
		}
	}()

	// 执行业务逻辑
	if err := fn(ctx); err != nil {
		u.Rollback(ctx)
		return err
	}

	// 提交事务
	return u.Commit(ctx)
} 