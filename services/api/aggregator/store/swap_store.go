package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// SwapStore defines the interface for storing swaps
type SwapStore interface {
	CreateSwap(ctx context.Context, swap *backend.Swap) error
	GetSwap(ctx context.Context, swapID string) (*backend.Swap, error)
	UpdateSwap(ctx context.Context, swap *backend.Swap) error
	AddStep(ctx context.Context, swapID string, step *backend.Step) error
	UpdateStep(ctx context.Context, swapID string, stepIndex int, step *backend.Step) error
	GetStep(ctx context.Context, swapID string, stepIndex int) (*backend.Step, error)
	CheckIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string) (string, bool)
}

// InMemorySwapStore implements SwapStore using in-memory storage
type InMemorySwapStore struct {
	mu    sync.RWMutex
	swaps map[string]*backend.Swap
	// idempotencyMap: swapID+stepIndex+idempotencyKey -> txHash
	idempotencyMap map[string]string
}

// NewInMemorySwapStore creates a new in-memory swap store
func NewInMemorySwapStore() *InMemorySwapStore {
	return &InMemorySwapStore{
		swaps:          make(map[string]*backend.Swap),
		idempotencyMap: make(map[string]string),
	}
}

// CreateSwap creates a new swap
func (s *InMemorySwapStore) CreateSwap(ctx context.Context, swap *backend.Swap) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.swaps[swap.SwapID]; exists {
		return fmt.Errorf("swap already exists: %s", swap.SwapID)
	}

	swap.CreatedAt = time.Now()
	swap.UpdatedAt = time.Now()
	s.swaps[swap.SwapID] = swap

	return nil
}

// GetSwap retrieves a swap by ID
func (s *InMemorySwapStore) GetSwap(ctx context.Context, swapID string) (*backend.Swap, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	swap, exists := s.swaps[swapID]
	if !exists {
		return nil, fmt.Errorf("swap not found: %s", swapID)
	}

	return swap, nil
}

// UpdateSwap updates an existing swap
func (s *InMemorySwapStore) UpdateSwap(ctx context.Context, swap *backend.Swap) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.swaps[swap.SwapID]; !exists {
		return fmt.Errorf("swap not found: %s", swap.SwapID)
	}

	swap.UpdatedAt = time.Now()
	s.swaps[swap.SwapID] = swap

	return nil
}

// AddStep adds a new step to a swap
func (s *InMemorySwapStore) AddStep(ctx context.Context, swapID string, step *backend.Step) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	swap, exists := s.swaps[swapID]
	if !exists {
		return fmt.Errorf("swap not found: %s", swapID)
	}

	swap.Steps = append(swap.Steps, step)
	swap.UpdatedAt = time.Now()

	return nil
}

// UpdateStep updates an existing step
func (s *InMemorySwapStore) UpdateStep(ctx context.Context, swapID string, stepIndex int, step *backend.Step) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	swap, exists := s.swaps[swapID]
	if !exists {
		return fmt.Errorf("swap not found: %s", swapID)
	}

	if stepIndex < 0 || stepIndex >= len(swap.Steps) {
		return fmt.Errorf("invalid step index: %d", stepIndex)
	}

	swap.Steps[stepIndex] = step
	swap.UpdatedAt = time.Now()

	return nil
}

// GetStep retrieves a specific step
func (s *InMemorySwapStore) GetStep(ctx context.Context, swapID string, stepIndex int) (*backend.Step, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	swap, exists := s.swaps[swapID]
	if !exists {
		return nil, fmt.Errorf("swap not found: %s", swapID)
	}

	if stepIndex < 0 || stepIndex >= len(swap.Steps) {
		return nil, fmt.Errorf("invalid step index: %d", stepIndex)
	}

	return swap.Steps[stepIndex], nil
}

// CheckIdempotency checks if a request is duplicate and returns existing txHash if found
func (s *InMemorySwapStore) CheckIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s:%d:%s", swapID, stepIndex, idempotencyKey)
	txHash, exists := s.idempotencyMap[key]

	return txHash, exists
}

// RecordIdempotency records an idempotency key with its txHash
func (s *InMemorySwapStore) RecordIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string, txHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%d:%s", swapID, stepIndex, idempotencyKey)
	s.idempotencyMap[key] = txHash

	return nil
}
