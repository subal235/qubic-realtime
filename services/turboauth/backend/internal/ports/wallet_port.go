package ports

import (
	"context"
)

// WalletVerifierPort defines the interface for wallet signature verification
type WalletVerifierPort interface {
	// VerifySignature verifies that the signature was created by the wallet owner
	VerifySignature(ctx context.Context, walletAddress, message, signature string) (bool, error)

	// GenerateChallenge generates a challenge message for wallet verification
	GenerateChallenge(walletAddress string) string

	// ValidateAddress checks if a wallet address is valid
	ValidateAddress(walletAddress string) bool
}
