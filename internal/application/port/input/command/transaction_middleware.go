package command

import (
    "context"
)

type TransactionMiddleware struct {
    uow    UnitOfWork
    logger Logger
}

func NewTransactionMiddleware(uow UnitOfWork, logger Logger) *TransactionMiddleware {
    return &TransactionMiddleware{
        uow:    uow,
        logger: logger,
    }
}

func (m *TransactionMiddleware) Execute(ctx context.Context, cmd interface{}, next Handler) (interface{}, error) {
    // 检查是否需要事务
    if th, ok := next.(TransactionHandler); ok && th.WithTransaction() {
        var result interface{}
        err := m.uow.WithTransaction(ctx, func(ctx context.Context) error {
            var err error
            result, err = next.Handle(ctx, cmd)
            return err
        })
        return result, err
    }
    return next.Handle(ctx, cmd)
} 