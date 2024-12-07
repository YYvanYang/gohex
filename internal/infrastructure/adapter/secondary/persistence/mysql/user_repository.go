package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
)

type userRepository struct {
	db        *sql.DB
	logger    Logger
	metrics   MetricsReporter
}

type userModel struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	Bio       string    `db:"bio"`
	Avatar    string    `db:"avatar"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUserRepository(db *sql.DB, logger Logger, metrics MetricsReporter) *userRepository {
	return &userRepository{
		db:      db,
		logger:  logger,
		metrics: metrics,
	}
}

func (r *userRepository) Save(ctx context.Context, user *aggregate.User) error {
	span, ctx := tracer.StartSpan(ctx, "userRepository.Save")
	defer span.End()

	timer := r.metrics.StartTimer("repository_save_user")
	defer timer.Stop()

	query := `
		INSERT INTO users (id, email, password, name, bio, avatar, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		r.metrics.IncrementCounter("repository_save_user_error")
		return err
	}

	r.metrics.IncrementCounter("repository_save_user_success")
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*aggregate.User, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.FindByID")
	defer span.End()

	var model userModel
	query := `SELECT * FROM users WHERE id = ?`
	
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
		return nil, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("failed to find user", "error", err)
		return nil, err
	}

	return r.toAggregate(&model)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	span, ctx := tracer.StartSpan(ctx, "userRepository.ExistsByEmail")
	defer span.End()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		r.logger.Error("failed to check email existence", "error", err)
		return false, err
	}

	return exists, nil
}

// toAggregate 将数据模型转换为聚合根
func (r *userRepository) toAggregate(model *userModel) (*aggregate.User, error) {
	email, err := vo.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	password := vo.NewPasswordFromHash(model.Password)

	profile, err := vo.NewUserProfile(model.Name, model.Bio)
	if err != nil {
		return nil, err
	}
	profile = profile.UpdateAvatar(model.Avatar)

	status := vo.UserStatus(model.Status)
	if !status.IsValid() {
		return nil, ErrInvalidUserStatus
	}

	user := &aggregate.User{
		BaseAggregate: aggregate.NewBaseAggregate(model.ID),
		email:         email,
		password:      password,
		profile:       profile,
		status:        status,
	}

	return user, nil
} 