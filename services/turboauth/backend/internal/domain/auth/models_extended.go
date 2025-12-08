package auth

import (
	"errors"
	"time"
)

// Session represents an authenticated session after wallet verification
type Session struct {
	SessionID     string    `json:"session_id"`
	WalletAddress string    `json:"wallet_address"`
	Token         string    `json:"token"` // JWT token
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	LastActivity  time.Time `json:"last_activity"`
	IPAddress     string    `json:"ip_address,omitempty"`
	UserAgent     string    `json:"user_agent,omitempty"`
}

// SessionRequest represents a request to create a session
type SessionRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
	Message       string `json:"message" validate:"required"`
	TTL           int    `json:"ttl,omitempty"` // Session TTL in seconds (default: 3600)
}

// SessionResponse represents the response after session creation
type SessionResponse struct {
	Session *Session    `json:"session"`
	Status  *WalletAuth `json:"status"`
}

// RefreshSessionRequest represents a request to refresh a session
type RefreshSessionRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Token     string `json:"token" validate:"required"`
}

// RateLimitInfo represents rate limiting information for a wallet
type RateLimitInfo struct {
	WalletAddress string    `json:"wallet_address"`
	RequestCount  int       `json:"request_count"`
	Limit         int       `json:"limit"`
	ResetAt       time.Time `json:"reset_at"`
	Remaining     int       `json:"remaining"`
}

// WebhookEvent represents an event to be sent via webhook
type WebhookEvent struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"` // status_changed, session_created, etc.
	WalletAddress string                 `json:"wallet_address"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
}

// WebhookSubscription represents a webhook subscription
type WebhookSubscription struct {
	SubscriptionID string    `json:"subscription_id"`
	URL            string    `json:"url"`
	Events         []string  `json:"events"` // Which events to subscribe to
	Secret         string    `json:"secret"` // For HMAC signature verification
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
}

// BatchVerifyRequest represents a batch verification request
type BatchVerifyRequest struct {
	Verifications []VerifyRequest `json:"verifications" validate:"required,min=1,max=100"`
}

// BatchVerifyResponse represents a batch verification response
type BatchVerifyResponse struct {
	Results []VerifyResult `json:"results"`
}

// VerifyResult represents a single verification result in a batch
type VerifyResult struct {
	WalletAddress string      `json:"wallet_address"`
	Verified      bool        `json:"verified"`
	Status        *WalletAuth `json:"status,omitempty"`
	Error         string      `json:"error,omitempty"`
}

// Additional errors
var (
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionNotFound   = errors.New("session not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrWebhookFailed     = errors.New("webhook delivery failed")
)
