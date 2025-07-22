package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TenantResolver middleware để xác định tenant từ subdomain hoặc custom domain
func TenantResolver() fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := c.Get("Host")

		// Extract subdomain from host
		// Format: subdomain.domain.com or custom.domain.com
		parts := strings.Split(host, ".")

		var tenantID string
		if len(parts) >= 2 {
			// Assume first part is subdomain/tenant identifier
			subdomain := parts[0]

			// Skip common subdomains
			if subdomain != "www" && subdomain != "api" && subdomain != "admin" {
				tenantID = subdomain
			}
		}

		// Set tenant context
		c.Locals("tenant_id", tenantID)
		c.Locals("host", host)

		return c.Next()
	}
}

// SecurityHeaders middleware để thêm security headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self'")

		return c.Next()
	}
}

// AuthRequired middleware để kiểm tra authentication
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
				"code":  "UNAUTHORIZED",
			})
		}

		// Extract token from "Bearer <token>"
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
				"code":  "INVALID_AUTH_FORMAT",
			})
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token required",
				"code":  "TOKEN_REQUIRED",
			})
		}

		// TODO: Validate JWT token and extract user info
		// For now, we'll just set some dummy user data
		c.Locals("user_id", "user-123")
		c.Locals("user_role", "user")
		c.Locals("token", token)

		return c.Next()
	}
}

// SystemAdminRequired middleware để kiểm tra quyền system admin
func SystemAdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole != "system_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "System admin role required",
				"code":  "INSUFFICIENT_PRIVILEGES",
			})
		}

		return c.Next()
	}
}

// TenantAdminRequired middleware để kiểm tra quyền tenant admin
func TenantAdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole != "tenant_admin" && userRole != "system_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Tenant admin role required",
				"code":  "INSUFFICIENT_PRIVILEGES",
			})
		}

		return c.Next()
	}
}
