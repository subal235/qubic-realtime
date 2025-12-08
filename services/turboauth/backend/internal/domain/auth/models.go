package auth

import (
	"errors"
	"time"
)

// AuthStatus represents the authentication status of a wallet
type AuthStatus string

const (
	StatusActive  AuthStatus = "ACTIVE"
	StatusBlocked AuthStatus = "BLOCKED"
	StatusReview  AuthStatus = "REVIEW"
	StatusUnknown AuthStatus = "UNKNOWN"
)

// WalletAuth represents the complete authentication state
type WalletAuth struct {
	WalletAddress   string     `json:"wallet_address"`
	Status          AuthStatus `json:"status"`
	TrustScore      int        `json:"trust_score"` // 0-100
	ContractAddress string     `json:"contract_address"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// VerifyRequest represents a wallet verification request
type VerifyRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
	Message       string `json:"message" validate:"required"`
}

// SetStatusRequest represents a request to update wallet status
type SetStatusRequest struct {
	WalletAddress  string     `json:"wallet_address" validate:"required"`
	Status         AuthStatus `json:"status" validate:"required,oneof=ACTIVE BLOCKED REVIEW"`
	TrustScore     int        `json:"trust_score" validate:"min=0,max=100"`
	AdminSignature string     `json:"admin_signature" validate:"required"`
}

// Common errors
var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInvalidSignature  = errors.New("invalid signature")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrBlockchainFailure = errors.New("blockchain operation failed")
	ErrCacheFailure      = errors.New("cache operation failed")
	ErrInvalidTrustScore = errors.New("trust score must be between 0 and 100")
)

// IsValid checks if the trust score is valid
func (w *WalletAuth) IsValid() error {
	if w.TrustScore < 0 || w.TrustScore > 100 {
		return ErrInvalidTrustScore
	}
	return nil
}

// IsActive returns true if the wallet is in active status
func (w *WalletAuth) IsActive() bool {
	return w.Status == StatusActive
}
