package truststore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"turboauth/internal/domain/auth"
	"turboauth/pkg/metrics"
)

// MemoryStore implements in-memory caching (L1 cache)
type MemoryStore struct {
	data  sync.Map
	mutex sync.RWMutex
}

type cacheEntry struct {
	Data      *auth.WalletAuth
	ExpiresAt time.Time
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{}

	// Start cleanup goroutine
	go store.cleanupExpired()

	return store
}

// Get retrieves cached data
func (m *MemoryStore) Get(ctx context.Context, walletAddress string) (*auth.WalletAuth, error) {
	value, ok := m.data.Load(walletAddress)
	if !ok {
		metrics.CacheMisses.WithLabelValues("L1").Inc()
		return nil, fmt.Errorf("not found")
	}

	entry := value.(*cacheEntry)
	if time.Now().After(entry.ExpiresAt) {
		m.data.Delete(walletAddress)
		metrics.CacheMisses.WithLabelValues("L1").Inc()
		return nil, fmt.Errorf("expired")
	}

	metrics.CacheHits.WithLabelValues("L1").Inc()
	return entry.Data, nil
}

// Set stores data in cache
func (m *MemoryStore) Set(ctx context.Context, walletAddress string, data *auth.WalletAuth, ttl time.Duration) error {
	entry := &cacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	m.data.Store(walletAddress, entry)
	return nil
}

// Delete removes data from cache
func (m *MemoryStore) Delete(ctx context.Context, walletAddress string) error {
	m.data.Delete(walletAddress)
	return nil
}

// BatchGet retrieves multiple entries
func (m *MemoryStore) BatchGet(ctx context.Context, walletAddresses []string) (map[string]*auth.WalletAuth, error) {
	result := make(map[string]*auth.WalletAuth)

	for _, addr := range walletAddresses {
		if data, err := m.Get(ctx, addr); err == nil {
			result[addr] = data
		}
	}

	return result, nil
}

// BatchSet stores multiple entries
func (m *MemoryStore) BatchSet(ctx context.Context, data map[string]*auth.WalletAuth, ttl time.Duration) error {
	for addr, walletAuth := range data {
		_ = m.Set(ctx, addr, walletAuth, ttl)
	}
	return nil
}

// HealthCheck always returns nil for memory store
func (m *MemoryStore) HealthCheck(ctx context.Context) error {
	return nil
}

// cleanupExpired removes expired entries periodically
func (m *MemoryStore) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		m.data.Range(func(key, value interface{}) bool {
			entry := value.(*cacheEntry)
			if now.After(entry.ExpiresAt) {
				m.data.Delete(key)
			}
			return true
		})
	}
}
