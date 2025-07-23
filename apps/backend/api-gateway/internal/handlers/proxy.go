package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"
	"zplus-saas/apps/backend/shared/config"

	"github.com/gofiber/fiber/v2"
)

type ProxyHandler struct {
	cfg    *config.Config
	client *http.Client
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	return &ProxyHandler{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *ProxyHandler) CRM(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.CRMServiceURL, "/crm")
}

// proxyRequest handles the actual proxying of requests to microservices
func (h *ProxyHandler) proxyRequest(c *fiber.Ctx, serviceURL, pathPrefix string) error {
	// Remove the path prefix and create the target URL
	path := strings.TrimPrefix(c.Path(), "/api/v1"+pathPrefix)
	targetURL := serviceURL + "/api/v1" + path

	// Create new request
	req, err := http.NewRequest(c.Method(), targetURL, strings.NewReader(string(c.Body())))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create proxy request",
		})
	}

	// Copy headers from original request
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add tenant ID header if available
	if tenantID := c.Locals("tenantID"); tenantID != nil {
		req.Header.Set("X-Tenant-ID", tenantID.(string))
	}

	// Make the request
	resp, err := h.client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service unavailable",
		})
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read response",
		})
	}

	return c.Status(resp.StatusCode).Send(body)
}

func (h *ProxyHandler) HRM(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.HRMServiceURL, "/hrm")
}

func (h *ProxyHandler) POS(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.POSServiceURL, "/pos")
}

func (h *ProxyHandler) LMS(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.LMSServiceURL, "/lms")
}

func (h *ProxyHandler) Checkin(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.CheckinServiceURL, "/checkin")
}

func (h *ProxyHandler) Payment(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.PaymentServiceURL, "/payment")
}

func (h *ProxyHandler) Files(c *fiber.Ctx) error {
	return h.proxyRequest(c, h.cfg.FileServiceURL, "/files")
}
