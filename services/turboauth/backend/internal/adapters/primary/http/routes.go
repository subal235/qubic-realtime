package http

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(app *fiber.App, handler *Handler) {
	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", handler.HealthCheck)

	// Metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// API v1 routes
	v1 := app.Group("/api/v1")
	{
		// Auth status endpoints
		v1.Get("/status/:wallet", handler.GetStatus)
		v1.Post("/status/batch", handler.BatchGetStatus)
		v1.Post("/status", handler.SetStatus)

		// Wallet verification
		v1.Post("/verify", handler.VerifyWallet)
	}
}
