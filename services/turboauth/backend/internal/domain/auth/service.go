package auth

import (
	"context"
	"time"

	"turboauth/pkg/metrics"

	"github.com/rs/zerolog/log"
)

// Service implements the core authentication business logic (hexagonal core)
type Service struct {
	qubicPort      QubicPort
	walletPort     WalletVerifierPort
	trustStorePort TrustStorePort
	cacheTTL       time.Duration
}

// NewService creates a new authentication service
func NewService(
	qubicPort QubicPort,
	walletPort WalletVerifierPort,
	trustStorePort TrustStorePort,
	cacheTTL time.Duration,
) *Service {
	return &Service{
		qubicPort:      qubicPort,
		walletPort:     walletPort,
		trustStorePort: trustStorePort,
		cacheTTL:       cacheTTL,
	}
}

// GetStatus retrieves authentication status with L1/L2/L3 caching strategy
// L1: In-memory (future optimization)
// L2: Redis cache (~5-10ms)
// L3: Qubic blockchain (~100-500ms)
func (s *Service) GetStatus(ctx context.Context, walletAddress string) (*WalletAuth, error) {
	start := time.Now()

	// Validate wallet address
	if !s.walletPort.ValidateAddress(walletAddress) {
		return nil, ErrInvalidSignature
	}

	// L2: Try Redis cache first
	cached, err := s.trustStorePort.Get(ctx, walletAddress)
	if err == nil && cached != nil {
		metrics.CacheHits.WithLabelValues("L2").Inc()
		log.Debug().
			Str("wallet", walletAddress).
			Dur("duration_ms", time.Since(start)).
			Msg("Cache hit (L2)")
		return cached, nil
	}
	metrics.CacheMisses.WithLabelValues("L2").Inc()

	// L3: Query blockchain
	log.Debug().Str("wallet", walletAddress).Msg("Querying blockchain")
	status, err := s.qubicPort.GetAuthStatus(ctx, walletAddress)
	if err != nil {
		metrics.BlockchainRequestsTotal.WithLabelValues("get_status", "error").Inc()
		return nil, err
	}
	metrics.BlockchainRequestsTotal.WithLabelValues("get_status", "success").Inc()
	metrics.BlockchainRequestDuration.WithLabelValues("get_status").Observe(time.Since(start).Seconds())

	// Cache for next time
	if err := s.trustStorePort.Set(ctx, walletAddress, status, s.cacheTTL); err != nil {
		log.Warn().Err(err).Msg("Failed to cache status")
	}

	return status, nil
}

// BatchGetStatus retrieves multiple statuses efficiently
func (s *Service) BatchGetStatus(ctx context.Context, walletAddresses []string) ([]*WalletAuth, error) {
	// Try cache first
	cachedData, _ := s.trustStorePort.BatchGet(ctx, walletAddresses)

	// Determine which addresses need blockchain lookup
	var missingAddresses []string
	for _, addr := range walletAddresses {
		if _, found := cachedData[addr]; !found {
			missingAddresses = append(missingAddresses, addr)
		}
	}

	// Fetch missing from blockchain
	var blockchainData []*WalletAuth
	if len(missingAddresses) > 0 {
		var err error
		blockchainData, err = s.qubicPort.BatchGetAuthStatus(ctx, missingAddresses)
		if err != nil {
			return nil, err
		}

		// Cache the newly fetched data
		cacheMap := make(map[string]*WalletAuth)
		for _, data := range blockchainData {
			cacheMap[data.WalletAddress] = data
		}
		_ = s.trustStorePort.BatchSet(ctx, cacheMap, s.cacheTTL)
	}

	// Combine cached and fresh data
	result := make([]*WalletAuth, 0, len(walletAddresses))
	for _, addr := range walletAddresses {
		if cached, ok := cachedData[addr]; ok {
			result = append(result, cached)
		} else {
			// Find in blockchain data
			for _, data := range blockchainData {
				if data.WalletAddress == addr {
					result = append(result, data)
					break
				}
			}
		}
	}

	return result, nil
}

// SetStatus updates the authentication status (admin only)
func (s *Service) SetStatus(ctx context.Context, req *SetStatusRequest) (string, error) {
	start := time.Now()

	// Validate request
	if !s.walletPort.ValidateAddress(req.WalletAddress) {
		return "", ErrInvalidSignature
	}

	// TODO: Verify admin signature
	// For now, we'll assume the caller is authorized

	// Update on blockchain
	txHash, err := s.qubicPort.SetAuthStatus(ctx, req)
	if err != nil {
		metrics.BlockchainRequestsTotal.WithLabelValues("set_status", "error").Inc()
		return "", err
	}
	metrics.BlockchainRequestsTotal.WithLabelValues("set_status", "success").Inc()
	metrics.BlockchainRequestDuration.WithLabelValues("set_status").Observe(time.Since(start).Seconds())

	// Invalidate cache
	if err := s.trustStorePort.Delete(ctx, req.WalletAddress); err != nil {
		log.Warn().Err(err).Msg("Failed to invalidate cache")
	}

	log.Info().
		Str("wallet", req.WalletAddress).
		Str("status", string(req.Status)).
		Str("tx_hash", txHash).
		Msg("Status updated")

	return txHash, nil
}

// VerifyWallet verifies a wallet signature and returns its auth status
func (s *Service) VerifyWallet(ctx context.Context, req *VerifyRequest) (*WalletAuth, bool, error) {
	// Verify signature
	verified, err := s.walletPort.VerifySignature(ctx, req.WalletAddress, req.Message, req.Signature)
	if err != nil {
		return nil, false, err
	}

	if !verified {
		return nil, false, ErrInvalidSignature
	}

	// Get current status
	status, err := s.GetStatus(ctx, req.WalletAddress)
	if err != nil {
		return nil, true, err // Signature is valid, but status lookup failed
	}

	return status, true, nil
}

// HealthCheck verifies all dependencies are healthy
func (s *Service) HealthCheck(ctx context.Context) error {
	if err := s.qubicPort.HealthCheck(ctx); err != nil {
		return err
	}
	if err := s.trustStorePort.HealthCheck(ctx); err != nil {
		return err
	}
	return nil
}
