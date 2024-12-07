package mysql

import (
	"context"
	"database/sql"
	"github.com/gohex/gohex/internal/domain/aggregate"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/pkg/tracer"
)

type userRepositoryImpl struct {
	db      *sql.DB
	logger  Logger
	metrics MetricsReporter
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *aggregate.User) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Update")
	defer span.End()

	timer := r.metrics.StartTimer("repository_update_user")
	defer timer.Stop()

	query := `
		UPDATE users 
		SET email = ?, password = ?, name = ?, bio = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email().String(),
		user.Password().Hash(),
		user.Profile().Name(),
		user.Profile().Bio(),
		user.Status().String(),
		user.UpdatedAt(),
		user.ID(),
	)

	if err != nil {
		r.logger.Error("failed to update user", "error", err)
		r.metrics.IncrementCounter("repository_update_user_error")
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.ErrUserNotFound
	}

	r.metrics.IncrementCounter("repository_update_user_success")
	return nil
}

// Save 保存用户
func (r *userRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Save")
	defer span.End()

	timer := r.metrics.StartTimer("repository_save_user")
	defer timer.Stop()

	query :=
		 `INSERT INTO users (id, email, password, name, bio, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email().String(),
		user.Password().Hash(),
		user.Profile().Name(),
		user.Profile().Bio(),
		user.Status().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		r.logger.Error("failed to save user", "error", err)
		r.metrics.IncrementCounter("repository_save_user_error")
		return err
	}

	r.metrics.IncrementCounter("repository_save_user_success")
	return nil
}

// FindByEmail 通过邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email vo.Email) (*aggregate.User, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.FindByEmail")
	defer span.End()

	var model userModel
	query := `SELECT * FROM users WHERE email = ?`
	
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.Name,
		&model.Bio,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("failed to find user by email", "error", err)
		return nil, err
	}

	return r.toAggregate(&model)
}

// 其他方法实现... 

</rewritten_file>```
</```
rewritten_file>