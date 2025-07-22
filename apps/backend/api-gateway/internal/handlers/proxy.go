package handlers

import (
	"zplus-saas/apps/backend/shared/config"

	"github.com/gofiber/fiber/v2"
)

type ProxyHandler struct {
	cfg *config.Config
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	return &ProxyHandler{
		cfg: cfg,
	}
}

func (h *ProxyHandler) CRM(c *fiber.Ctx) error {
	// TODO: Proxy to CRM service
	return c.JSON(fiber.Map{
		"message": "CRM service proxy - TODO: Implement",
		"service": "crm",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) HRM(c *fiber.Ctx) error {
	// TODO: Proxy to HRM service
	return c.JSON(fiber.Map{
		"message": "HRM service proxy - TODO: Implement",
		"service": "hrm",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) POS(c *fiber.Ctx) error {
	// TODO: Proxy to POS service
	return c.JSON(fiber.Map{
		"message": "POS service proxy - TODO: Implement",
		"service": "pos",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) LMS(c *fiber.Ctx) error {
	// TODO: Proxy to LMS service
	return c.JSON(fiber.Map{
		"message": "LMS service proxy - TODO: Implement",
		"service": "lms",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) Checkin(c *fiber.Ctx) error {
	// TODO: Proxy to Checkin service
	return c.JSON(fiber.Map{
		"message": "Checkin service proxy - TODO: Implement",
		"service": "checkin",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) Payment(c *fiber.Ctx) error {
	// TODO: Proxy to Payment service
	return c.JSON(fiber.Map{
		"message": "Payment service proxy - TODO: Implement",
		"service": "payment",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}

func (h *ProxyHandler) Files(c *fiber.Ctx) error {
	// TODO: Proxy to File service
	return c.JSON(fiber.Map{
		"message": "File service proxy - TODO: Implement",
		"service": "files",
		"path":    c.Path(),
		"method":  c.Method(),
		"status":  "coming_soon",
	})
}
