package http

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(app *fiber.App, handler *Handler) {
	// Health check
	app.Get("/health", handler.Health)

	// API v1
	v1 := app.Group("/api/v1")

	// Payment routes
	v1.Post("/pay", handler.Pay)
	v1.Post("/route", handler.FindRoute)
	v1.Get("/health/:routeID", handler.GetRouteHealth)
}
