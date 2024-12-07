package output

import (
	"context"
	"time"
)

// Cache 定义缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	DeleteMulti(ctx context.Context, keys []string) error
	Clear(ctx context.Context) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}

// Cacheable 定义可缓存接口
type Cacheable interface {
	CacheKey() string
	TTL() time.Duration
} 