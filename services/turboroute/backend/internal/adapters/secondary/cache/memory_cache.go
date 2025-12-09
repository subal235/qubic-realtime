package cache

import (
	"context"
	"sync"
	"time"
	"turboroute/internal/domain/route"
)

// MemoryCache is an in-memory cache for routes
type MemoryCache struct {
	routes map[string]*cachedRoute
	health map[string]*route.RouteHealth
	mu     sync.RWMutex
}

type cachedRoute struct {
	route     *route.RouteOption
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		routes: make(map[string]*cachedRoute),
		health: make(map[string]*route.RouteHealth),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// GetHealth returns health metrics for a route
func (c *MemoryCache) GetHealth(ctx context.Context, routeID string) (*route.RouteHealth, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if health, ok := c.health[routeID]; ok {
		return health, nil
	}

	// Return default health if not found
	return &route.RouteHealth{
		RouteID:         routeID,
		IsActive:        true,
		SuccessRate:     0.99,
		AverageFee:      2,
		AverageTime:     150,
		TotalExecutions: 0,
		LastUpdated:     time.Now(),
	}, nil
}

// UpdateMetrics updates route health metrics
func (c *MemoryCache) UpdateMetrics(ctx context.Context, metrics route.RouteMetrics) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	health, ok := c.health[metrics.RouteID]
	if !ok {
		health = &route.RouteHealth{
			RouteID:         metrics.RouteID,
			IsActive:        true,
			SuccessRate:     0.99,
			AverageFee:      0,
			AverageTime:     0,
			TotalExecutions: 0,
			LastExecuted:    time.Now(),
			LastUpdated:     time.Now(),
		}
		c.health[metrics.RouteID] = health
	}

	// Update metrics (simple moving average)
	if metrics.Success {
		health.TotalExecutions++
		health.SuccessRate = (health.SuccessRate*float64(health.TotalExecutions-1) + 1.0) / float64(health.TotalExecutions)

		if metrics.ActualFee > 0 {
			health.AverageFee = (health.AverageFee*int64(health.TotalExecutions-1) + metrics.ActualFee) / int64(health.TotalExecutions)
		}

		if metrics.ActualTime > 0 {
			health.AverageTime = (health.AverageTime*int64(health.TotalExecutions-1) + metrics.ActualTime.Milliseconds()) / int64(health.TotalExecutions)
		}
	} else {
		health.TotalExecutions++
		health.SuccessRate = (health.SuccessRate * float64(health.TotalExecutions-1)) / float64(health.TotalExecutions)
	}

	health.LastExecuted = metrics.ExecutedAt
	health.LastUpdated = time.Now()

	return nil
}

// GetCachedRoute retrieves a cached route
func (c *MemoryCache) GetCachedRoute(ctx context.Context, from, to string) (*route.RouteOption, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := from + "->" + to
	if cached, ok := c.routes[key]; ok {
		if time.Now().Before(cached.expiresAt) {
			return cached.route, nil
		}
	}

	return nil, route.ErrNoRoutesFound
}

// CacheRoute stores a route in cache
func (c *MemoryCache) CacheRoute(ctx context.Context, from, to string, r *route.RouteOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := from + "->" + to
	c.routes[key] = &cachedRoute{
		route:     r,
		expiresAt: time.Now().Add(5 * time.Minute),
	}

	return nil
}

// cleanup removes expired routes
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, cached := range c.routes {
			if now.After(cached.expiresAt) {
				delete(c.routes, key)
			}
		}
		c.mu.Unlock()
	}
}
