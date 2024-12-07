package mysql

import (
	"context"
	"database/sql"
	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
	"github.com/your-org/your-project/internal/application/port/output"
)

type userRepository struct {
	db      *sql.DB
	logger  Logger
	metrics MetricsReporter
}

func NewUserRepository(db *sql.DB, logger Logger, metrics MetricsReporter) output.UserRepository {
	return &userRepository{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}
}

func (r *userRepository) Save(ctx context.Context, user *aggregate.User) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Save")
	defer span.End()

	query := `
		INSERT INTO users (
			id, email, password, name, bio, avatar, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email().String(),
		user.Password().Hash(),
		user.Profile().Name(),
		user.Profile().Bio(),
		user.Profile().Avatar(),
		user.Status().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		r.logger.Error("failed to save user", "error", err)
		return err
	}

	// 保存用户角色
	for _, role := range user.Roles() {
		if err := r.saveUserRole(ctx, user.ID(), role); err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email vo.Email) (*aggregate.User, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.FindByEmail")
	defer span.End()

	query := `
		SELECT 
			u.id, u.email, u.password, u.name, u.bio, u.avatar, 
			u.status, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = ?
	`

	var model userModel
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.Name,
		&model.Bio,
		&model.Avatar,
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

	// 加载用户角色
	roles, err := r.getUserRoles(ctx, model.ID)
	if err != nil {
		return nil, err
	}
	model.Roles = roles

	return r.toAggregate(&model)
}

func (r *userRepository) FindAll(ctx context.Context, params output.FindAllParams) ([]*aggregate.User, int64, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.FindAll")
	defer span.End()

	// 构建查询条件
	where := "WHERE 1=1"
	args := []interface{}{}

	if params.Status != "" {
		where += " AND status = ?"
		args = append(args, params.Status)
	}

	if params.Role != "" {
		where += " AND id IN (SELECT user_id FROM user_roles WHERE role = ?)"
		args = append(args, params.Role)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM users " + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.email, u.password, u.name, u.bio, u.avatar, 
			u.status, u.created_at, u.updated_at
		FROM users u
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, where, params.SortBy, params.SortDir)

	args = append(args, params.Limit, params.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*aggregate.User
	for rows.Next() {
		var model userModel
		err := rows.Scan(
			&model.ID,
			&model.Email,
			&model.Password,
			&model.Name,
			&model.Bio,
			&model.Avatar,
			&model.Status,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// 加载用户角色
		roles, err := r.getUserRoles(ctx, model.ID)
		if err != nil {
			return nil, 0, err
		}
		model.Roles = roles

		user, err := r.toAggregate(&model)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// userModel 数据模型
type userModel struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	Bio       string    `db:"bio"`
	Avatar    string    `db:"avatar"`
	Status    string    `db:"status"`
	Roles     []string  `db:"roles"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// 实现剩余的仓储方法
func (r *userRepository) FindByID(ctx context.Context, id string) (*aggregate.User, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.FindByID")
	defer span.End()

	query := `
		SELECT 
			u.id, u.email, u.password, u.name, u.bio, u.avatar, 
			u.status, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = ?
	`

	var model userModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.Name,
		&model.Bio,
		&model.Avatar,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("failed to find user by id", "error", err)
		return nil, err
	}

	// 加载用户角色
	roles, err := r.getUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	model.Roles = roles

	return r.toAggregate(&model)
}

func (r *userRepository) Update(ctx context.Context, user *aggregate.User) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Update")
	defer span.End()

	query := `
		UPDATE users 
		SET 
			email = ?, 
			password = ?, 
			name = ?, 
			bio = ?, 
			avatar = ?, 
			status = ?, 
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email().String(),
		user.Password().Hash(),
		user.Profile().Name(),
		user.Profile().Bio(),
		user.Profile().Avatar(),
		user.Status().String(),
		user.UpdatedAt(),
		user.ID(),
	)

	if err != nil {
		r.logger.Error("failed to update user", "error", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrUserNotFound
	}

	// 更新用户角色
	if err := r.updateUserRoles(ctx, user.ID(), user.Roles()); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Delete")
	defer span.End()

	// 开启事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除用户角色
	if _, err := tx.ExecContext(ctx, "DELETE FROM user_roles WHERE user_id = ?", id); err != nil {
		return err
	}

	// 删除用户
	result, err := tx.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrUserNotFound
	}

	return tx.Commit()
}

// 辅助方法
func (r *userRepository) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT role FROM user_roles WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *userRepository) updateUserRoles(ctx context.Context, userID string, roles []vo.UserRole) error {
	// 开启事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除现有角色
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM user_roles WHERE user_id = ?",
		userID,
	); err != nil {
		return err
	}

	// 插入新角色
	for _, role := range roles {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO user_roles (user_id, role) VALUES (?, ?)",
			userID, role.String(),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *userRepository) toAggregate(model *userModel) (*aggregate.User, error) {
	// 1. 创建值对象
	email, err := vo.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	password := vo.NewPasswordFromHash(model.Password)

	profile, err := vo.NewUserProfile(model.Name, model.Bio)
	if err != nil {
		return nil, err
	}
	if model.Avatar != "" {
		profile, err = profile.WithAvatar(model.Avatar)
		if err != nil {
			return nil, err
		}
	}

	status := vo.UserStatus(model.Status)
	if !status.IsValid() {
		return nil, errors.ErrInvalidUserStatus
	}

	// 2. 创建角色列表
	roles := make([]vo.UserRole, len(model.Roles))
	for i, r := range model.Roles {
		role := vo.UserRole(r)
		if !role.IsValid() {
			return nil, errors.ErrInvalidRole
		}
		roles[i] = role
	}

	// 3. 重建聚合根
	user := &aggregate.User{
		BaseAggregate: aggregate.NewBaseAggregate(model.ID),
		email:         email,
		password:      password,
		profile:       profile,
		status:        status,
		roles:         roles,
	}

	return user, nil
}

// 其他方法实现... 

</```rewritten_file>