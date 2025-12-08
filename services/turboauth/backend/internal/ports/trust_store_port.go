package ports

import (
	"context"
	"turboauth/internal/domain/auth"
	"time"
)

// TrustStorePort defines the interface for caching authentication data
// Implements L1 (memory) and L2 (Redis) caching strategy
type TrustStorePort interface {
	// Get retrieves cached authentication status
	Get(ctx context.Context, walletAddress string) (*auth.WalletAuth, error)

	// Set stores authentication status in cache
	Set(ctx context.Context, walletAddress string, data *auth.WalletAuth, ttl time.Duration) error

	// Delete removes cached data
	Delete(ctx context.Context, walletAddress string) error

	// BatchGet retrieves multiple cached statuses (performance optimization)
	BatchGet(ctx context.Context, walletAddresses []string) (map[string]*auth.WalletAuth, error)

	// BatchSet stores multiple statuses
	BatchSet(ctx context.Context, data map[string]*auth.WalletAuth, ttl time.Duration) error

	// HealthCheck verifies cache connectivity
	HealthCheck(ctx context.Context) error
}
