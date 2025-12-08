package ports

import (
	"context"
	"turboauth/internal/domain/auth"
)

// QubicPort defines the interface for interacting with Qubic blockchain
type QubicPort interface {
	// GetAuthStatus retrieves the authentication status from the smart contract
	GetAuthStatus(ctx context.Context, walletAddress string) (*auth.WalletAuth, error)

	// SetAuthStatus updates the authentication status on the smart contract
	SetAuthStatus(ctx context.Context, req *auth.SetStatusRequest) (txHash string, err error)

	// BatchGetAuthStatus retrieves multiple statuses in a single call (performance optimization)
	BatchGetAuthStatus(ctx context.Context, walletAddresses []string) ([]*auth.WalletAuth, error)

	// GetContractAddress returns the current smart contract address
	GetContractAddress() string

	// HealthCheck verifies connection to the Qubic node
	HealthCheck(ctx context.Context) error
}
