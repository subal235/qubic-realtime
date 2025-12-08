package route

import (
	"context"
)

// QubicPaymentPort defines the interface for blockchain payments
type QubicPaymentPort interface {
	ExecuteTransfer(ctx context.Context, from, to string, amount int64) (txHash string, err error)
	GetTransactionStatus(ctx context.Context, txHash string) (string, error)
	EstimateFee(ctx context.Context, from, to string, amount int64) (int64, error)
	GetBalance(ctx context.Context, wallet string) (int64, error)
}

// RouteCatalogPort defines the interface for route discovery
type RouteCatalogPort interface {
	ListRoutes(ctx context.Context, from, to string) ([]RouteOption, error)
	ScoreRoute(route RouteOption, preferences RoutePreferences) float64
	RegisterRoute(ctx context.Context, route RouteOption) error
	GetRoute(ctx context.Context, routeID string) (*RouteOption, error)
}

// RouteHealthPort defines the interface for route health tracking
type RouteHealthPort interface {
	GetHealth(ctx context.Context, routeID string) (*RouteHealth, error)
	UpdateMetrics(ctx context.Context, metrics RouteMetrics) error
	GetCachedRoute(ctx context.Context, from, to string) (*RouteOption, error)
	CacheRoute(ctx context.Context, from, to string, route *RouteOption) error
}
