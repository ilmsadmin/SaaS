package middleware

import (
	"strings"

	"zplus-saas/apps/backend/auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtService services.JWTService
}

func NewAuthMiddleware(jwtService services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authorization header required",
			"message": "Missing Authorization header",
		})
	}

	// Check if it starts with Bearer
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid authorization header",
			"message": "Authorization header must start with 'Bearer '",
		})
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Token required",
			"message": "Missing authentication token",
		})
	}

	// Validate token
	jwtToken, err := m.jwtService.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token",
			"message": err.Error(),
		})
	}

	// Parse claims
	claims, err := m.jwtService.ParseUserClaims(jwtToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token claims",
			"message": err.Error(),
		})
	}

	// Check if user is active
	if !claims.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Account deactivated",
			"message": "Your account has been deactivated",
		})
	}

	// Set user context
	c.Locals("user_id", claims.UserID)
	c.Locals("tenant_id", claims.TenantID)
	c.Locals("user_email", claims.Email)
	c.Locals("user_role", claims.Role)
	c.Locals("is_verified", claims.IsVerified)

	return c.Next()
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "User role not found",
				"message": "User authentication required",
			})
		}

		role := userRole.(string)
		for _, requiredRole := range roles {
			if role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Insufficient permissions",
			"message": "You don't have permission to access this resource",
		})
	}
}

// RequireVerification middleware checks if user email is verified
func (m *AuthMiddleware) RequireVerification(c *fiber.Ctx) error {
	isVerified := c.Locals("is_verified")
	if isVerified == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "User verification status not found",
			"message": "User authentication required",
		})
	}

	if !isVerified.(bool) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Email verification required",
			"message": "Please verify your email address to access this resource",
		})
	}

	return c.Next()
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	// Check if it starts with Bearer
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Next()
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return c.Next()
	}

	// Validate token
	jwtToken, err := m.jwtService.ValidateToken(token)
	if err != nil {
		return c.Next()
	}

	// Parse claims
	claims, err := m.jwtService.ParseUserClaims(jwtToken)
	if err != nil {
		return c.Next()
	}

	// Set user context if token is valid
	c.Locals("user_id", claims.UserID)
	c.Locals("tenant_id", claims.TenantID)
	c.Locals("user_email", claims.Email)
	c.Locals("user_role", claims.Role)
	c.Locals("is_verified", claims.IsVerified)

	return c.Next()
}
