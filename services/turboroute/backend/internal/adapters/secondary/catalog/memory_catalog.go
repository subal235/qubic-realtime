package catalog

import (
	"context"
	"sync"
	"turboroute/internal/domain/route"
)

// MemoryCatalog is an in-memory route catalog
type MemoryCatalog struct {
	routes map[string][]route.RouteOption
	mu     sync.RWMutex
}

// NewMemoryCatalog creates a new in-memory catalog
func NewMemoryCatalog() *MemoryCatalog {
	catalog := &MemoryCatalog{
		routes: make(map[string][]route.RouteOption),
	}

	// Add some default routes for testing
	catalog.addDefaultRoutes()

	return catalog
}

// ListRoutes returns all routes between two wallets
func (c *MemoryCatalog) ListRoutes(ctx context.Context, from, to string) ([]route.RouteOption, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := from + "->" + to
	if routes, ok := c.routes[key]; ok {
		return routes, nil
	}

	// Return default direct route
	return []route.RouteOption{
		{
			RouteID:       "direct",
			Hops:          []string{from, to},
			EstimatedFee:  2,
			EstimatedTime: 150000000, // 150ms in nanoseconds
			SuccessRate:   0.99,
			RouteType:     "direct",
		},
	}, nil
}

// ScoreRoute calculates a score for a route based on preferences
func (c *MemoryCatalog) ScoreRoute(r route.RouteOption, prefs route.RoutePreferences) float64 {
	score := 100.0

	switch prefs.Priority {
	case "speed":
		// Lower time = higher score
		score -= float64(r.EstimatedTime.Milliseconds()) / 10.0
	case "cost":
		// Lower fee = higher score
		score -= float64(r.EstimatedFee) * 2.0
	case "privacy":
		// More hops = higher score (for privacy)
		score += float64(len(r.Hops)) * 5.0
	default: // "reliability"
		// Higher success rate = higher score
		score += r.SuccessRate * 50.0
	}

	// Success rate always matters
	score += r.SuccessRate * 20.0

	return score
}

// RegisterRoute adds a new route to the catalog
func (c *MemoryCatalog) RegisterRoute(ctx context.Context, r route.RouteOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(r.Hops) < 2 {
		return route.ErrNoRoutesFound
	}

	key := r.Hops[0] + "->" + r.Hops[len(r.Hops)-1]
	c.routes[key] = append(c.routes[key], r)

	return nil
}

// GetRoute returns a specific route by ID
func (c *MemoryCatalog) GetRoute(ctx context.Context, routeID string) (*route.RouteOption, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, routes := range c.routes {
		for _, r := range routes {
			if r.RouteID == routeID {
				return &r, nil
			}
		}
	}

	return nil, route.ErrRouteNotAvailable
}

// addDefaultRoutes adds some test routes
func (c *MemoryCatalog) addDefaultRoutes() {
	// These are just examples for testing
	defaultRoutes := []route.RouteOption{
		{
			RouteID:       "direct",
			Hops:          []string{"WALLET_A", "WALLET_B"},
			EstimatedFee:  2,
			EstimatedTime: 150000000,
			SuccessRate:   0.99,
			RouteType:     "direct",
		},
		{
			RouteID:       "fast_lane",
			Hops:          []string{"WALLET_A", "WALLET_B"},
			EstimatedFee:  5,
			EstimatedTime: 50000000,
			SuccessRate:   0.995,
			RouteType:     "direct",
		},
	}

	for _, r := range defaultRoutes {
		key := r.Hops[0] + "->" + r.Hops[len(r.Hops)-1]
		c.routes[key] = append(c.routes[key], r)
	}
}
