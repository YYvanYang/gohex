package redis

import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/your-org/your-project/internal/application/port/output"
)

type redisCache struct {
	client  *redis.Client
	logger  Logger
	metrics MetricsReporter
}

func NewRedisCache(client *redis.Client, logger Logger, metrics MetricsReporter) output.Cache {
	return &redisCache{
		client:  client,
		logger:  logger,
		metrics: metrics,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	timer := c.metrics.StartTimer("redis_get_duration")
	defer timer.Stop()

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.metrics.IncrementCounter("cache_miss")
		return nil, nil
	}
	if err != nil {
		c.logger.Error("redis get failed", "error", err)
		c.metrics.IncrementCounter("cache_error")
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.logger.Error("failed to unmarshal cache value", "error", err)
		return nil, err
	}

	c.metrics.IncrementCounter("cache_hit")
	return result, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	timer := c.metrics.StartTimer("redis_set_duration")
	defer timer.Stop()

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("failed to marshal cache value", "error", err)
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Error("redis set failed", "error", err)
		c.metrics.IncrementCounter("cache_error")
		return err
	}

	c.metrics.IncrementCounter("cache_set")
	return nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	timer := c.metrics.StartTimer("cache_delete_duration")
	defer timer.Stop()

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("failed to delete from cache", "error", err)
		c.metrics.IncrementCounter("cache_error")
		return err
	}

	c.metrics.IncrementCounter("cache_delete")
	return nil
}

func (c *redisCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	timer := c.metrics.StartTimer("cache_increment_duration")
	defer timer.Stop()

	result, err := c.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		c.logger.Error("failed to increment cache", "error", err)
		c.metrics.IncrementCounter("cache_error")
		return 0, err
	}

	c.metrics.IncrementCounter("cache_increment")
	return result, nil
}

func (c *redisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
} 