package cache

import "context"

// Cache 是行情系统使用的缓存抽象
// 注意：这是“信息级行情缓存”，不负责复杂 TTL / 原子操作
type Cache interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, bool, error)
}
