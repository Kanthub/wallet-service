package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{
		client: client,
		ttl:    5 * time.Minute, // 信息级行情：默认 5 秒钟/ 分钟
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte) error {
	if c.client == nil {
		return errors.New("redis client is nil")
	}
	return c.client.Set(ctx, key, value, c.ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	if c.client == nil {
		return nil, false, errors.New("redis client is nil")
	}

	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}

	return val, true, err
}
