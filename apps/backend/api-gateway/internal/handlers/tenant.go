package handlers

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"time"

	"zplus-saas/apps/backend/shared/config"

	"github.com/gofiber/fiber/v2"
)

type TenantHandler struct {
	db         *sql.DB
	cfg        *config.Config
	httpClient *http.Client
}

func NewTenantHandler(db *sql.DB, cfg *config.Config) *TenantHandler {
	return &TenantHandler{
		db:  db,
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *TenantHandler) List(c *fiber.Ctx) error {
	return h.proxyToTenantService(c, "/api/tenants")
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
	return h.proxyToTenantService(c, "/api/tenants")
}

func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	return h.proxyToTenantService(c, "/api/tenants/"+tenantID)
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	return h.proxyToTenantService(c, "/api/tenants/"+tenantID)
}

func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	return h.proxyToTenantService(c, "/api/tenants/"+tenantID)
}

// proxyToTenantService proxies request to tenant service
func (h *TenantHandler) proxyToTenantService(c *fiber.Ctx, path string) error {
	// Build URL with query parameters
	url := h.cfg.TenantServiceURL + path
	if len(c.Request().URI().QueryString()) > 0 {
		url += "?" + string(c.Request().URI().QueryString())
	}

	// Create request body reader
	var bodyReader io.Reader
	if c.Body() != nil {
		bodyReader = bytes.NewReader(c.Body())
	}

	// Create HTTP request
	req, err := http.NewRequest(c.Method(), url, bodyReader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create request",
			"message": err.Error(),
		})
	}

	// Copy headers from original request
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if !isHopByHopHeader(keyStr) {
			req.Header.Set(keyStr, string(value))
		}
	})

	// Make the request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "Tenant service unavailable",
			"message": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				c.Response().Header.Add(key, value)
			}
		}
	}

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read response",
			"message": err.Error(),
		})
	}

	c.Response().SetStatusCode(resp.StatusCode)
	return c.Send(body)
}
