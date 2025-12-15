package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// QuoteStore defines the interface for storing quotes
type QuoteStore interface {
	Save(ctx context.Context, quoteID string, quote *backend.QuoteStore, ttl time.Duration) error
	Update(ctx context.Context, quoteID string, quote *backend.QuoteStore, ttl time.Duration) error
	Get(ctx context.Context, quoteID string) (*backend.QuoteStore, error)
	Delete(ctx context.Context, quoteID string) error
}

// quoteEntry represents a quote with expiration
type quoteEntry struct {
	quote     *backend.QuoteStore
	expiresAt time.Time
}

// InMemoryQuoteStore implements QuoteStore using in-memory storage
type InMemoryQuoteStore struct {
	mu     sync.RWMutex
	quotes map[string]*quoteEntry
}

// NewInMemoryQuoteStore creates a new in-memory quote store
func NewInMemoryQuoteStore() *InMemoryQuoteStore {
	store := &InMemoryQuoteStore{
		quotes: make(map[string]*quoteEntry),
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// Save stores a quote with TTL
func (s *InMemoryQuoteStore) Save(ctx context.Context, quoteID string, quote *backend.QuoteStore, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.quotes[quoteID] = &quoteEntry{
		quote:     quote,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Get retrieves a quote by ID
func (s *InMemoryQuoteStore) Get(ctx context.Context, quoteID string) (*backend.QuoteStore, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.quotes[quoteID]
	if !exists {
		return nil, fmt.Errorf("quote not found: %s", quoteID)
	}

	if time.Now().After(entry.expiresAt) {
		return nil, fmt.Errorf("quote expired: %s", quoteID)
	}

	return entry.quote, nil
}

func (s *InMemoryQuoteStore) Update(ctx context.Context, quoteID string, quote *backend.QuoteStore, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.quotes[quoteID]; !exists {
		return fmt.Errorf("quote not found: %s", quoteID)
	}

	s.quotes[quoteID] = &quoteEntry{
		quote:     quote,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// Delete removes a quote
func (s *InMemoryQuoteStore) Delete(ctx context.Context, quoteID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.quotes, quoteID)
	return nil
}

// cleanup periodically removes expired quotes
func (s *InMemoryQuoteStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, entry := range s.quotes {
			if now.After(entry.expiresAt) {
				delete(s.quotes, id)
			}
		}
		s.mu.Unlock()
	}
}
