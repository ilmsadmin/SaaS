package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware extracts and validates tenant information
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract tenant ID from header
		tenantID := c.Get("X-Tenant-ID")

		// For now, just add the tenant ID to the context
		// In a production system, you would validate the tenant exists
		// and the user has access to it
		if tenantID != "" {
			c.Locals("tenantID", tenantID)
		}

		return c.Next()
	}
}
