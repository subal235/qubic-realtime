package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/rs/zerolog/log"
)

// CreateSession creates a new authenticated session after wallet verification
func (s *Service) CreateSession(ctx context.Context, req *SessionRequest) (*SessionResponse, error) {
	// Verify wallet first
	verifyReq := &VerifyRequest{
		WalletAddress: req.WalletAddress,
		Signature:     req.Signature,
		Message:       req.Message,
	}

	status, verified, err := s.VerifyWallet(ctx, verifyReq)
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, ErrInvalidSignature
	}

	// Check rate limit
	if s.rateLimitPort != nil {
		limitInfo, err := s.rateLimitPort.CheckRateLimit(ctx, req.WalletAddress)
		if err == nil && limitInfo.Remaining <= 0 {
			return nil, ErrRateLimitExceeded
		}
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	// Set TTL (default 1 hour)
	ttl := time.Hour
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	expiresAt := time.Now().Add(ttl)

	// Generate JWT token
	var token string
	if s.tokenPort != nil {
		token, err = s.tokenPort.GenerateToken(req.WalletAddress, expiresAt)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to generate token, continuing without it")
		}
	}

	// Create session
	session := &Session{
		SessionID:     sessionID,
		WalletAddress: req.WalletAddress,
		Token:         token,
		ExpiresAt:     expiresAt,
		CreatedAt:     time.Now(),
		LastActivity:  time.Now(),
	}

	// Store session
	if s.sessionPort != nil {
		if err := s.sessionPort.CreateSession(ctx, session); err != nil {
			return nil, err
		}
	}

	// Send webhook notification
	if s.webhookPort != nil {
		event := &WebhookEvent{
			EventID:       generateEventID(),
			EventType:     "session_created",
			WalletAddress: req.WalletAddress,
			Timestamp:     time.Now(),
			Data: map[string]interface{}{
				"session_id": sessionID,
				"expires_at": expiresAt,
			},
		}
		go s.webhookPort.SendWebhook(context.Background(), event)
	}

	log.Info().
		Str("wallet", req.WalletAddress).
		Str("session_id", sessionID).
		Msg("Session created")

	return &SessionResponse{
		Session: session,
		Status:  status,
	}, nil
}

// RefreshSession extends an existing session
func (s *Service) RefreshSession(ctx context.Context, req *RefreshSessionRequest) (*Session, error) {
	if s.sessionPort == nil {
		return nil, ErrSessionNotFound
	}

	// Validate token if provided
	if req.Token != "" && s.tokenPort != nil {
		_, err := s.tokenPort.ValidateToken(req.Token)
		if err != nil {
			return nil, ErrInvalidToken
		}
	}

	// Refresh session
	session, err := s.sessionPort.RefreshSession(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// BatchVerify verifies multiple wallets in a single request
func (s *Service) BatchVerify(ctx context.Context, req *BatchVerifyRequest) (*BatchVerifyResponse, error) {
	results := make([]VerifyResult, len(req.Verifications))

	for i, verifyReq := range req.Verifications {
		status, verified, err := s.VerifyWallet(ctx, &verifyReq)

		result := VerifyResult{
			WalletAddress: verifyReq.WalletAddress,
			Verified:      verified,
		}

		if err != nil {
			result.Error = err.Error()
		} else {
			result.Status = status
		}

		results[i] = result
	}

	return &BatchVerifyResponse{
		Results: results,
	}, nil
}

// Helper functions
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateEventID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
