package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// RedisSwapStore implements SwapStore using Redis
type RedisSwapStore struct {
	client *redis.Client
}

// NewRedisSwapStore creates a new Redis-based swap store
func NewRedisSwapStore(client *redis.Client) *RedisSwapStore {
	return &RedisSwapStore{
		client: client,
	}
}

// CreateSwap creates a new swap in Redis
func (s *RedisSwapStore) CreateSwap(ctx context.Context, swap *backend.Swap) error {
	key := s.swapKey(swap.SwapID)

	// Check if swap already exists
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check swap existence: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("swap already exists: %s", swap.SwapID)
	}

	// Set timestamps
	swap.CreatedAt = time.Now()
	swap.UpdatedAt = time.Now()

	// Serialize and save
	data, err := json.Marshal(swap)
	if err != nil {
		return fmt.Errorf("failed to marshal swap: %w", err)
	}

	// Store with 24 hour TTL (swaps should complete within a day)
	if err := s.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to save swap to Redis: %w", err)
	}

	return nil
}

// GetSwap retrieves a swap by ID from Redis
func (s *RedisSwapStore) GetSwap(ctx context.Context, swapID string) (*backend.Swap, error) {
	key := s.swapKey(swapID)

	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("swap not found: %s", swapID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get swap from Redis: %w", err)
	}

	var swap backend.Swap
	if err := json.Unmarshal(data, &swap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal swap: %w", err)
	}

	return &swap, nil
}

// UpdateSwap updates an existing swap in Redis
func (s *RedisSwapStore) UpdateSwap(ctx context.Context, swap *backend.Swap) error {
	key := s.swapKey(swap.SwapID)

	// Check if swap exists
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check swap existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("swap not found: %s", swap.SwapID)
	}

	// Update timestamp
	swap.UpdatedAt = time.Now()

	// Serialize and save
	data, err := json.Marshal(swap)
	if err != nil {
		return fmt.Errorf("failed to marshal swap: %w", err)
	}

	// Keep existing TTL
	if err := s.client.Set(ctx, key, data, redis.KeepTTL).Err(); err != nil {
		return fmt.Errorf("failed to update swap in Redis: %w", err)
	}

	return nil
}

// AddStep adds a new step to a swap
func (s *RedisSwapStore) AddStep(ctx context.Context, swapID string, step *backend.Step) error {
	swap, err := s.GetSwap(ctx, swapID)
	if err != nil {
		return err
	}

	swap.Steps = append(swap.Steps, step)
	return s.UpdateSwap(ctx, swap)
}

// UpdateStep updates an existing step
func (s *RedisSwapStore) UpdateStep(ctx context.Context, swapID string, stepIndex int, step *backend.Step) error {
	swap, err := s.GetSwap(ctx, swapID)
	if err != nil {
		return err
	}

	if stepIndex < 0 || stepIndex >= len(swap.Steps) {
		return fmt.Errorf("invalid step index: %d", stepIndex)
	}

	swap.Steps[stepIndex] = step
	return s.UpdateSwap(ctx, swap)
}

// GetStep retrieves a specific step
func (s *RedisSwapStore) GetStep(ctx context.Context, swapID string, stepIndex int) (*backend.Step, error) {
	swap, err := s.GetSwap(ctx, swapID)
	if err != nil {
		return nil, err
	}

	if stepIndex < 0 || stepIndex >= len(swap.Steps) {
		return nil, fmt.Errorf("invalid step index: %d", stepIndex)
	}

	return swap.Steps[stepIndex], nil
}

// CheckIdempotency checks if a request is duplicate and returns existing txHash if found
func (s *RedisSwapStore) CheckIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string) (string, bool) {
	key := s.idempotencyKey(swapID, stepIndex, idempotencyKey)

	txHash, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		return "", false
	}

	return txHash, true
}

// RecordIdempotency records an idempotency key with its txHash
func (s *RedisSwapStore) RecordIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string, txHash string) error {
	key := s.idempotencyKey(swapID, stepIndex, idempotencyKey)

	// Store with 24 hour TTL
	if err := s.client.Set(ctx, key, txHash, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to record idempotency: %w", err)
	}

	return nil
}

// swapKey generates the Redis key for a swap
func (s *RedisSwapStore) swapKey(swapID string) string {
	return fmt.Sprintf("aggregator:swap:%s", swapID)
}

// idempotencyKey generates the Redis key for idempotency tracking
func (s *RedisSwapStore) idempotencyKey(swapID string, stepIndex int, idempotencyKey string) string {
	return fmt.Sprintf("aggregator:idempotency:%s:%d:%s", swapID, stepIndex, idempotencyKey)
}
