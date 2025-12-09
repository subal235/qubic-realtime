package qubic

import (
	"context"
	"fmt"
)

// MockPaymentClient is a mock implementation for development
type MockPaymentClient struct {
	nodeURL string
}

// NewMockPaymentClient creates a new mock payment client
func NewMockPaymentClient(nodeURL string) *MockPaymentClient {
	return &MockPaymentClient{
		nodeURL: nodeURL,
	}
}

// ExecuteTransfer simulates a blockchain transfer
func (c *MockPaymentClient) ExecuteTransfer(ctx context.Context, from, to string, amount int64) (string, error) {
	// Simulate successful transfer
	txHash := fmt.Sprintf("0x%s%s%d", from[:8], to[:8], amount)
	return txHash, nil
}

// GetTransactionStatus returns the status of a transaction
func (c *MockPaymentClient) GetTransactionStatus(ctx context.Context, txHash string) (string, error) {
	// Mock: always return confirmed
	return "confirmed", nil
}

// EstimateFee estimates the fee for a transfer
func (c *MockPaymentClient) EstimateFee(ctx context.Context, from, to string, amount int64) (int64, error) {
	// Mock: simple fee calculation
	baseFee := int64(2)
	if amount > 10000 {
		return baseFee + (amount / 10000), nil
	}
	return baseFee, nil
}

// GetBalance returns the balance of a wallet
func (c *MockPaymentClient) GetBalance(ctx context.Context, wallet string) (int64, error) {
	// Mock: always return sufficient balance
	return 1000000, nil
}
