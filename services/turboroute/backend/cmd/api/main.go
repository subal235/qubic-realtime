package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"turboroute/internal/adapters/primary/http"
	"turboroute/internal/adapters/secondary/cache"
	"turboroute/internal/adapters/secondary/catalog"
	"turboroute/internal/adapters/secondary/qubic"
	"turboroute/internal/domain/route"
)

func main() {
	// Setup logger
	setupLogger()

	log.Info().Msg("Starting TurboRoute...")

	// Get configuration from environment
	httpPort := getEnv("HTTP_PORT", "8081")
	qubicNodeURL := getEnv("QUBIC_NODE_URL", "http://qubic-node:21841")

	// Initialize adapters (secondary/infrastructure)
	paymentClient := qubic.NewMockPaymentClient(qubicNodeURL)
	catalogAdapter := catalog.NewMemoryCatalog()
	cacheAdapter := cache.NewMemoryCache()

	// Initialize domain service (hexagonal core)
	routeService := route.NewService(
		paymentClient,
		catalogAdapter,
		cacheAdapter,
	)

	// Start HTTP server
	startHTTPServer(httpPort, routeService)
}

func startHTTPServer(port string, svc *route.Service) {
	app := fiber.New(fiber.Config{
		ServerHeader:     "TurboRoute",
		StrictRouting:    true,
		CaseSensitive:    true,
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      120 * time.Second,
		DisableKeepalive: false,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(cors.New())

	// Setup routes
	handler := http.NewHandler(svc)
	http.SetupRoutes(app, handler)

	addr := fmt.Sprintf(":%s", port)
	log.Info().Msgf("🚀 TurboRoute listening on %s", addr)

	// Graceful shutdown
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}
	log.Info().Msg("Server stopped")
}

func setupLogger() {
	logLevel := getEnv("LOG_LEVEL", "info")
	logFormat := getEnv("LOG_FORMAT", "json")

	// Set log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set format
	if logFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
