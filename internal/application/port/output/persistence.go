package output

import (
	"context"
	"github.com/your-org/your-project/internal/domain/aggregate"
	"github.com/your-org/your-project/internal/domain/vo"
)

// UserRepository 定义用户仓储接口
type UserRepository interface {
	Save(ctx context.Context, user *aggregate.User) error
	FindByID(ctx context.Context, id string) (*aggregate.User, error)
	FindByEmail(ctx context.Context, email vo.Email) (*aggregate.User, error)
	ExistsByEmail(ctx context.Context, email vo.Email) (bool, error)
	Update(ctx context.Context, user *aggregate.User) error
	Delete(ctx context.Context, id string) error
}

// EventStore 定义事件存储接口
type EventStore interface {
	SaveEvents(ctx context.Context, aggregateID string, events []event.Event, expectedVersion int) error
	GetEvents(ctx context.Context, aggregateID string) ([]event.Event, error)
	GetEventsByType(ctx context.Context, eventType string) ([]event.Event, error)
}

// Cache 定义缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type Repository interface {
	Transaction
	Save(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id string) error
}

type Transaction interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Queryable interface {
	FindByID(ctx context.Context, id string) (interface{}, error)
	FindAll(ctx context.Context, params QueryParams) ([]interface{}, int64, error)
}

type QueryParams struct {
	Filters map[string]interface{}
	Sort    map[string]string
	Offset  int
	Limit   int
} 