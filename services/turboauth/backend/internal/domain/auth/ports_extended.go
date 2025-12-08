package auth

import (
	"context"
	"time"
)

// SessionPort defines the interface for session management
type SessionPort interface {
	// CreateSession creates a new authenticated session
	CreateSession(ctx context.Context, session *Session) error

	// GetSession retrieves a session by ID
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// RefreshSession extends the session expiration
	RefreshSession(ctx context.Context, sessionID string) (*Session, error)

	// DeleteSession invalidates a session
	DeleteSession(ctx context.Context, sessionID string) error

	// GetActiveSessions returns all active sessions for a wallet
	GetActiveSessions(ctx context.Context, walletAddress string) ([]*Session, error)

	// CleanupExpiredSessions removes expired sessions
	CleanupExpiredSessions(ctx context.Context) (int, error)
}

// RateLimitPort defines the interface for rate limiting
type RateLimitPort interface {
	// CheckRateLimit checks if a wallet has exceeded rate limits
	CheckRateLimit(ctx context.Context, walletAddress string) (*RateLimitInfo, error)

	// IncrementCounter increments the request counter for a wallet
	IncrementCounter(ctx context.Context, walletAddress string) error

	// ResetCounter resets the counter for a wallet
	ResetCounter(ctx context.Context, walletAddress string) error

	// GetLimitInfo returns current rate limit info
	GetLimitInfo(ctx context.Context, walletAddress string) (*RateLimitInfo, error)
}

// WebhookPort defines the interface for webhook notifications
type WebhookPort interface {
	// SendWebhook sends a webhook event
	SendWebhook(ctx context.Context, event *WebhookEvent) error

	// RegisterWebhook registers a new webhook subscription
	RegisterWebhook(ctx context.Context, subscription *WebhookSubscription) error

	// UnregisterWebhook removes a webhook subscription
	UnregisterWebhook(ctx context.Context, subscriptionID string) error

	// GetWebhooks returns all active webhooks
	GetWebhooks(ctx context.Context) ([]*WebhookSubscription, error)

	// RetryFailedWebhooks retries failed webhook deliveries
	RetryFailedWebhooks(ctx context.Context) error
}

// TokenPort defines the interface for JWT token operations
type TokenPort interface {
	// GenerateToken creates a new JWT token
	GenerateToken(walletAddress string, expiresAt time.Time) (string, error)

	// ValidateToken validates a JWT token
	ValidateToken(token string) (walletAddress string, err error)

	// RefreshToken creates a new token from an existing one
	RefreshToken(token string) (string, error)
}
