package output

import "context"

type UnitOfWork interface {
    // Begin 开始事务
    Begin(ctx context.Context) (context.Context, error)
    // Commit 提交事务
    Commit(ctx context.Context) error
    // Rollback 回滚事务
    Rollback(ctx context.Context) error
    // WithTransaction 在事务中执行函数
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionManager interface {
    // GetTransaction 获取当前事务
    GetTransaction(ctx context.Context) (interface{}, error)
    // SetTransaction 设置当前事务
    SetTransaction(ctx context.Context, tx interface{}) context.Context
} 