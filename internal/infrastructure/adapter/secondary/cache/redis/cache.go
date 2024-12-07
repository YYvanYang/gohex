package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/your-org/your-project/internal/application/port"
)

type redisCache struct {
	client  *redis.Client
	logger  Logger
	metrics MetricsReporter
}

func NewRedisCache(client *redis.Client, logger Logger, metrics MetricsReporter) port.Cache {
	return &redisCache{
		client:  client,
		logger:  logger,
		metrics: metrics,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	span, ctx := tracer.StartSpan(ctx, "redisCache.Get")
	defer span.End()

	timer := c.metrics.StartTimer("cache_get_duration")
	defer timer.Stop()

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.metrics.IncrementCounter("cache_miss")
		return nil, nil
	}
	if err != nil {
		c.logger.Error("failed to get from cache", "key", key, "error", err)
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.logger.Error("failed to unmarshal cached value", "error", err)
		return nil, err
	}

	c.metrics.IncrementCounter("cache_hit")
	return result, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	span, ctx := tracer.StartSpan(ctx, "redisCache.Set")
	defer span.End()

	timer := c.metrics.StartTimer("cache_set_duration")
	defer timer.Stop()

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("failed to marshal value", "error", err)
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Error("failed to set cache", "key", key, "error", err)
		return err
	}

	c.metrics.IncrementCounter("cache_set_success")
	return nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	span, ctx := tracer.StartSpan(ctx, "redisCache.Delete")
	defer span.End()

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("failed to delete from cache", "key", key, "error", err)
		return err
	}

	return nil
} 