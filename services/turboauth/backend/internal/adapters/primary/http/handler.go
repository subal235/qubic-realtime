package http

import (
	"time"

	"turboauth/internal/domain/auth"
	"turboauth/pkg/metrics"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	authService *auth.Service
}

// NewHandler creates a new HTTP handler
func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// GetStatus handles GET /api/v1/status/:wallet
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	start := time.Now()
	walletAddress := c.Params("wallet")

	status, err := h.authService.GetStatus(c.Context(), walletAddress)
	if err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("GET", "/status", "500").Inc()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	metrics.HTTPRequestsTotal.WithLabelValues("GET", "/status", "200").Inc()
	metrics.HTTPRequestDuration.WithLabelValues("GET", "/status").Observe(time.Since(start).Seconds())

	return c.JSON(status)
}

// BatchGetStatus handles POST /api/v1/status/batch
func (h *Handler) BatchGetStatus(c *fiber.Ctx) error {
	start := time.Now()

	var req struct {
		WalletAddresses []string `json:"wallet_addresses"`
	}

	if err := c.BodyParser(&req); err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status/batch", "400").Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	statuses, err := h.authService.BatchGetStatus(c.Context(), req.WalletAddresses)
	if err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status/batch", "500").Inc()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status/batch", "200").Inc()
	metrics.HTTPRequestDuration.WithLabelValues("POST", "/status/batch").Observe(time.Since(start).Seconds())

	return c.JSON(fiber.Map{
		"statuses": statuses,
	})
}

// SetStatus handles POST /api/v1/status
func (h *Handler) SetStatus(c *fiber.Ctx) error {
	start := time.Now()

	var req auth.SetStatusRequest
	if err := c.BodyParser(&req); err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status", "400").Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	txHash, err := h.authService.SetStatus(c.Context(), &req)
	if err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status", "500").Inc()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	metrics.HTTPRequestsTotal.WithLabelValues("POST", "/status", "200").Inc()
	metrics.HTTPRequestDuration.WithLabelValues("POST", "/status").Observe(time.Since(start).Seconds())

	return c.JSON(fiber.Map{
		"success": true,
		"tx_hash": txHash,
	})
}

// VerifyWallet handles POST /api/v1/verify
func (h *Handler) VerifyWallet(c *fiber.Ctx) error {
	start := time.Now()

	var req auth.VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/verify", "400").Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	status, verified, err := h.authService.VerifyWallet(c.Context(), &req)
	if err != nil {
		metrics.HTTPRequestsTotal.WithLabelValues("POST", "/verify", "500").Inc()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	metrics.HTTPRequestsTotal.WithLabelValues("POST", "/verify", "200").Inc()
	metrics.HTTPRequestDuration.WithLabelValues("POST", "/verify").Observe(time.Since(start).Seconds())

	return c.JSON(fiber.Map{
		"verified":    verified,
		"status":      status.Status,
		"trust_score": status.TrustScore,
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	if err := h.authService.HealthCheck(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}
