package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// RedisQuoteStore implements QuoteStore using Redis
type RedisQuoteStore struct {
	client *redis.Client
}

// NewRedisQuoteStore creates a new Redis-based quote store
func NewRedisQuoteStore(client *redis.Client) *RedisQuoteStore {
	return &RedisQuoteStore{
		client: client,
	}
}

// Save stores a quote with TTL in Redis
func (s *RedisQuoteStore) Save(ctx context.Context, quoteID string, quote *backend.QuoteResponse, ttl time.Duration) error {
	// Serialize quote to JSON
	data, err := json.Marshal(quote)
	if err != nil {
		return fmt.Errorf("failed to marshal quote: %w", err)
	}

	// Store in Redis with TTL
	key := s.quoteKey(quoteID)
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save quote to Redis: %w", err)
	}

	return nil
}

// Get retrieves a quote by ID from Redis
func (s *RedisQuoteStore) Get(ctx context.Context, quoteID string) (*backend.QuoteResponse, error) {
	key := s.quoteKey(quoteID)

	// Get from Redis
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("quote not found: %s", quoteID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get quote from Redis: %w", err)
	}

	// Deserialize
	var quote backend.QuoteResponse
	if err := json.Unmarshal(data, &quote); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quote: %w", err)
	}

	return &quote, nil
}

// Delete removes a quote from Redis
func (s *RedisQuoteStore) Delete(ctx context.Context, quoteID string) error {
	key := s.quoteKey(quoteID)

	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete quote from Redis: %w", err)
	}

	return nil
}

// quoteKey generates the Redis key for a quote
func (s *RedisQuoteStore) quoteKey(quoteID string) string {
	return fmt.Sprintf("aggregator:quote:%s", quoteID)
}
