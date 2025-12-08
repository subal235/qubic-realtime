package route

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

// Service implements the core routing business logic
type Service struct {
	paymentPort QubicPaymentPort
	catalogPort RouteCatalogPort
	healthPort  RouteHealthPort
}

// NewService creates a new routing service
func NewService(
	paymentPort QubicPaymentPort,
	catalogPort RouteCatalogPort,
	healthPort RouteHealthPort,
) *Service {
	return &Service{
		paymentPort: paymentPort,
		catalogPort: catalogPort,
		healthPort:  healthPort,
	}
}

// FindRoute discovers and selects the best route for a payment
func (s *Service) FindRoute(ctx context.Context, intent PaymentIntent) (*RouteDecision, error) {
	start := time.Now()

	// Check cached route first
	if cached, err := s.healthPort.GetCachedRoute(ctx, intent.From, intent.To); err == nil && cached != nil {
		log.Debug().Str("from", intent.From).Str("to", intent.To).Msg("Using cached route")
		return &RouteDecision{
			SelectedRoute: *cached,
			Reason:        "cached_route",
			Timestamp:     time.Now(),
		}, nil
	}

	// Discover all available routes
	routes, err := s.catalogPort.ListRoutes(ctx, intent.From, intent.To)
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 {
		return nil, ErrNoRoutesFound
	}

	// Score each route based on preferences
	for i := range routes {
		routes[i].Score = s.catalogPort.ScoreRoute(routes[i], intent.Preferences)
	}

	// Sort by score (highest first)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Score > routes[j].Score
	})

	// Select best route
	bestRoute := routes[0]

	// Validate against constraints
	if intent.Preferences.MaxFee > 0 && bestRoute.EstimatedFee > intent.Preferences.MaxFee {
		return nil, ErrFeeExceedsMax
	}

	// Cache the selected route
	_ = s.healthPort.CacheRoute(ctx, intent.From, intent.To, &bestRoute)

	log.Info().
		Str("route_id", bestRoute.RouteID).
		Int64("fee", bestRoute.EstimatedFee).
		Dur("discovery_time", time.Since(start)).
		Msg("Route selected")

	return &RouteDecision{
		SelectedRoute: bestRoute,
		Reason:        fmt.Sprintf("best_score_%.2f", bestRoute.Score),
		Alternatives:  routes[1:min(len(routes), 4)], // Top 3 alternatives
		Timestamp:     time.Now(),
	}, nil
}

// ExecutePayment executes a payment along the selected route
func (s *Service) ExecutePayment(ctx context.Context, intent PaymentIntent, route RouteOption) (*PaymentExecution, error) {
	start := time.Now()
	executionID := generateExecutionID()

	log.Info().
		Str("execution_id", executionID).
		Str("route_id", route.RouteID).
		Int64("amount", intent.Amount).
		Msg("Executing payment")

	// Check balance
	balance, err := s.paymentPort.GetBalance(ctx, intent.From)
	if err != nil {
		return nil, err
	}

	if balance < intent.Amount+route.EstimatedFee {
		return nil, ErrInsufficientFunds
	}

	// Execute the transfer
	txHash, err := s.paymentPort.ExecuteTransfer(ctx, intent.From, intent.To, intent.Amount)
	if err != nil {
		// Update metrics with failure
		_ = s.healthPort.UpdateMetrics(ctx, RouteMetrics{
			RouteID:    route.RouteID,
			Success:    false,
			ExecutedAt: time.Now(),
		})

		return &PaymentExecution{
			ExecutionID: executionID,
			Intent:      intent,
			Route:       route,
			Status:      "failed",
			Error:       err.Error(),
			ExecutedAt:  time.Now(),
		}, ErrPaymentFailed
	}

	actualTime := time.Since(start)

	// Get actual fee (for now, use estimated)
	actualFee := route.EstimatedFee

	// Update metrics with success
	_ = s.healthPort.UpdateMetrics(ctx, RouteMetrics{
		RouteID:    route.RouteID,
		Success:    true,
		ActualFee:  actualFee,
		ActualTime: actualTime,
		ExecutedAt: time.Now(),
	})

	completedAt := time.Now()

	log.Info().
		Str("execution_id", executionID).
		Str("tx_hash", txHash).
		Dur("execution_time", actualTime).
		Msg("Payment executed successfully")

	return &PaymentExecution{
		ExecutionID: executionID,
		Intent:      intent,
		Route:       route,
		TxHash:      txHash,
		Status:      "executed",
		ActualFee:   actualFee,
		ActualTime:  actualTime,
		ExecutedAt:  start,
		CompletedAt: &completedAt,
	}, nil
}

// Pay is a convenience method that finds and executes in one call
func (s *Service) Pay(ctx context.Context, intent PaymentIntent) (*PaymentExecution, error) {
	// Find best route
	decision, err := s.FindRoute(ctx, intent)
	if err != nil {
		return nil, err
	}

	// Execute payment
	return s.ExecutePayment(ctx, intent, decision.SelectedRoute)
}

// GetRouteHealth returns health metrics for a route
func (s *Service) GetRouteHealth(ctx context.Context, routeID string) (*RouteHealth, error) {
	return s.healthPort.GetHealth(ctx, routeID)
}

// Helper functions
func generateExecutionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
