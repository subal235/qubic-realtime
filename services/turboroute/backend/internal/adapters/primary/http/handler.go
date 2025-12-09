package http

import (
	"turboroute/internal/domain/route"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Handler handles HTTP requests for routing
type Handler struct {
	service *route.Service
}

// NewHandler creates a new HTTP handler
func NewHandler(service *route.Service) *Handler {
	return &Handler{service: service}
}

// Pay handles payment execution with auto-routing
func (h *Handler) Pay(c *fiber.Ctx) error {
	var intent route.PaymentIntent
	if err := c.BodyParser(&intent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	execution, err := h.service.Pay(c.Context(), intent)
	if err != nil {
		log.Error().Err(err).Msg("Payment failed")
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(execution)
}

// FindRoute finds the best route without executing
func (h *Handler) FindRoute(c *fiber.Ctx) error {
	var intent route.PaymentIntent
	if err := c.BodyParser(&intent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	decision, err := h.service.FindRoute(c.Context(), intent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(decision)
}

// GetRouteHealth returns health metrics for a route
func (h *Handler) GetRouteHealth(c *fiber.Ctx) error {
	routeID := c.Params("routeID")

	health, err := h.service.GetRouteHealth(c.Context(), routeID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "route not found"})
	}

	return c.JSON(health)
}

// Health check endpoint
func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "turboroute"})
}
