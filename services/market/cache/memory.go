package cache

import (
	"context"
	"sync"
	"time"
)

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	value []byte
	ts    time.Time
}

func NewMemoryCache() Cache {
	return &MemoryCache{
		store: make(map[string]item),
	}
}

func (c *MemoryCache) Set(ctx context.Context, key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 防御性拷贝，避免外部 slice 被修改
	buf := make([]byte, len(val))
	copy(buf, val)

	c.store[key] = item{
		value: buf,
		ts:    time.Now(),
	}
	return nil
}

func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, ok := c.store[key]
	if !ok {
		return nil, false, nil
	}

	// 防御性拷贝，避免调用方改脏 cache
	buf := make([]byte, len(it.value))
	copy(buf, it.value)

	return buf, true, nil
}
