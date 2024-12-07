package output

import (
	"context"
	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
)

type UserRepository interface {
	// 基本操作
	Save(ctx context.Context, user *aggregate.User) error
	Update(ctx context.Context, user *aggregate.User) error
	Delete(ctx context.Context, id string) error

	// 查询方法
	FindByID(ctx context.Context, id string) (*aggregate.User, error)
	FindByEmail(ctx context.Context, email vo.Email) (*aggregate.User, error)
	FindAll(ctx context.Context, params FindAllParams) ([]*aggregate.User, int64, error)
	ExistsByEmail(ctx context.Context, email vo.Email) (bool, error)

	// 批量操作
	SaveBatch(ctx context.Context, users []*aggregate.User) error
	FindByIDs(ctx context.Context, ids []string) ([]*aggregate.User, error)

	// 统计方法
	Count(ctx context.Context, status string) (int64, error)
	CountByRole(ctx context.Context, role vo.UserRole) (int64, error)
}

type FindAllParams struct {
	Status   string
	Role     string
	Offset   int
	Limit    int
	SortBy   string
	SortDir  string
} 