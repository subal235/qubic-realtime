package route

import (
	"errors"
	"time"
)

// PaymentIntent represents a payment request
type PaymentIntent struct {
	From        string            `json:"from" validate:"required"`
	To          string            `json:"to" validate:"required"`
	Amount      int64             `json:"amount" validate:"required,gt=0"`
	Preferences RoutePreferences  `json:"preferences"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RoutePreferences defines user preferences for routing
type RoutePreferences struct {
	Priority  string `json:"priority"`   // "speed", "cost", "privacy", "reliability"
	MaxFee    int64  `json:"max_fee"`    // Maximum acceptable fee
	TimeoutMs int    `json:"timeout_ms"` // Maximum time to wait
	MinHops   int    `json:"min_hops"`   // Minimum hops (for privacy)
	MaxHops   int    `json:"max_hops"`   // Maximum hops
}

// RouteOption represents a possible payment route
type RouteOption struct {
	RouteID       string        `json:"route_id"`
	Hops          []string      `json:"hops"`
	EstimatedFee  int64         `json:"estimated_fee"`
	EstimatedTime time.Duration `json:"estimated_time_ms"`
	SuccessRate   float64       `json:"success_rate"` // 0.0 to 1.0
	Score         float64       `json:"score"`        // Calculated score
	RouteType     string        `json:"route_type"`   // "direct", "multi-hop", "liquidity"
}

// RouteDecision represents the selected route
type RouteDecision struct {
	SelectedRoute RouteOption   `json:"selected_route"`
	Reason        string        `json:"reason"`
	Alternatives  []RouteOption `json:"alternatives,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
}

// PaymentExecution represents an executed payment
type PaymentExecution struct {
	ExecutionID string        `json:"execution_id"`
	Intent      PaymentIntent `json:"intent"`
	Route       RouteOption   `json:"route"`
	TxHash      string        `json:"tx_hash"`
	Status      string        `json:"status"` // "pending", "executed", "failed"
	ActualFee   int64         `json:"actual_fee"`
	ActualTime  time.Duration `json:"actual_time_ms"`
	ExecutedAt  time.Time     `json:"executed_at"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// RouteHealth represents the health metrics of a route
type RouteHealth struct {
	RouteID         string    `json:"route_id"`
	IsActive        bool      `json:"is_active"`
	SuccessRate     float64   `json:"success_rate"`
	AverageFee      int64     `json:"average_fee"`
	AverageTime     int64     `json:"average_time_ms"`
	TotalExecutions int64     `json:"total_executions"`
	LastExecuted    time.Time `json:"last_executed"`
	LastUpdated     time.Time `json:"last_updated"`
}

// RouteMetrics represents metrics for route scoring
type RouteMetrics struct {
	RouteID    string
	Success    bool
	ActualFee  int64
	ActualTime time.Duration
	ExecutedAt time.Time
}

// Common errors
var (
	ErrNoRoutesFound     = errors.New("no routes found")
	ErrRouteNotAvailable = errors.New("route not available")
	ErrPaymentFailed     = errors.New("payment execution failed")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrTimeoutExceeded   = errors.New("timeout exceeded")
	ErrFeeExceedsMax     = errors.New("fee exceeds maximum")
)
