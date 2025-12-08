package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microauth_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "microauth_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// gRPC Metrics
	GRPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microauth_grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "microauth_grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Cache Metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microauth_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"layer"}, // L1 (memory) or L2 (redis)
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microauth_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"layer"},
	)

	// Blockchain Metrics
	BlockchainRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microauth_blockchain_requests_total",
			Help: "Total number of blockchain requests",
		},
		[]string{"operation", "status"},
	)

	BlockchainRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "microauth_blockchain_request_duration_seconds",
			Help:    "Blockchain request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)
)
