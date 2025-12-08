package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "turboauth/api/proto"
	grpcAdapter "turboauth/internal/adapters/primary/grpc"
	httpAdapter "turboauth/internal/adapters/primary/http"
	"turboauth/internal/adapters/secondary/qubic"
	"turboauth/internal/adapters/secondary/truststore"
	"turboauth/internal/adapters/secondary/wallet"
	"turboauth/internal/domain/auth"
	"turboauth/pkg/config"
	"turboauth/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel, cfg.LogFormat)

	log.Info().
		Str("env", cfg.Env).
		Int("http_port", cfg.HTTPPort).
		Int("grpc_port", cfg.GRPCPort).
		Msg("Starting Qubic MicroAuth")

	// Initialize adapters (secondary/infrastructure)
	qubicClient := qubic.NewClient(cfg.QubicNodeURL, cfg.QubicContractAddr)
	walletVerifier := wallet.NewVerifier()

	// Initialize trust store (cache)
	var trustStore truststore.RedisStore
	if cfg.UseMemoryCache {
		log.Info().Msg("Using in-memory cache")
		memStore := truststore.NewMemoryStore()
		// Type assertion to satisfy the interface
		// In production, you might want a composite store (L1 + L2)
		_ = memStore
	}

	// Initialize Redis store
	redisStore, err := truststore.NewRedisStore(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Redis, using memory store only")
		// Fallback to memory store
	} else {
		trustStore = *redisStore
		log.Info().Msg("Connected to Redis")
	}

	// Initialize domain service (hexagonal core)
	authService := auth.NewService(
		qubicClient,
		walletVerifier,
		&trustStore,
		cfg.CacheTTL,
	)

	// Start HTTP server (Fiber)
	go startHTTPServer(cfg, authService)

	// Start gRPC server
	go startGRPCServer(cfg, authService)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down servers...")
	// TODO: Implement graceful shutdown for both servers
}

func startHTTPServer(cfg *config.Config, svc *auth.Service) {
	app := fiber.New(fiber.Config{
		Prefork:           false, // Set true for multi-process in production
		ServerHeader:      "MicroAuth",
		StrictRouting:     true,
		CaseSensitive:     true,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		DisableKeepalive:  false,
		ReduceMemoryUsage: false, // Set true if memory is constrained
	})

	// Compression middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Fast compression
	}))

	// Setup routes
	handler := httpAdapter.NewHandler(svc)
	httpAdapter.SetupRoutes(app, handler)

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Info().Msgf("🚀 HTTP server listening on %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}

func startGRPCServer(cfg *config.Config, svc *auth.Service) {
	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen for gRPC")
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(1000),
		grpc.ConnectionTimeout(10*time.Second),
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
	)

	// Register service
	grpcSvc := grpcAdapter.NewServer(svc)
	pb.RegisterAuthServiceServer(grpcServer, grpcSvc)

	log.Info().Msgf("🚀 gRPC server listening on %s", addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("gRPC server failed")
	}
}
