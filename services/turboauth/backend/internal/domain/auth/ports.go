package auth

import (
	"context"
	"time"
)

// QubicPort defines the interface for interacting with Qubic blockchain
type QubicPort interface {
	// GetAuthStatus retrieves the authentication status from the smart contract
	GetAuthStatus(ctx context.Context, walletAddress string) (*WalletAuth, error)

	// SetAuthStatus updates the authentication status on the smart contract
	SetAuthStatus(ctx context.Context, req *SetStatusRequest) (txHash string, err error)

	// BatchGetAuthStatus retrieves multiple statuses in a single call (performance optimization)
	BatchGetAuthStatus(ctx context.Context, walletAddresses []string) ([]*WalletAuth, error)

	// GetContractAddress returns the current smart contract address
	GetContractAddress() string

	// HealthCheck verifies connection to the Qubic node
	HealthCheck(ctx context.Context) error
}

// WalletVerifierPort defines the interface for wallet signature verification
type WalletVerifierPort interface {
	// VerifySignature verifies that the signature was created by the wallet owner
	VerifySignature(ctx context.Context, walletAddress, message, signature string) (bool, error)

	// GenerateChallenge generates a challenge message for wallet verification
	GenerateChallenge(walletAddress string) string

	// ValidateAddress checks if a wallet address is valid
	ValidateAddress(walletAddress string) bool
}

// TrustStorePort defines the interface for caching authentication data
// Implements L1 (memory) and L2 (Redis) caching strategy
type TrustStorePort interface {
	// Get retrieves cached authentication status
	Get(ctx context.Context, walletAddress string) (*WalletAuth, error)

	// Set stores authentication status in cache
	Set(ctx context.Context, walletAddress string, data *WalletAuth, ttl time.Duration) error

	// Delete removes cached data
	Delete(ctx context.Context, walletAddress string) error

	// BatchGet retrieves multiple cached statuses (performance optimization)
	BatchGet(ctx context.Context, walletAddresses []string) (map[string]*WalletAuth, error)

	// BatchSet stores multiple statuses
	BatchSet(ctx context.Context, data map[string]*WalletAuth, ttl time.Duration) error

	// HealthCheck verifies cache connectivity
	HealthCheck(ctx context.Context) error
}
