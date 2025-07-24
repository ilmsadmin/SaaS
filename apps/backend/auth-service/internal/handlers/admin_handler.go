package handlers

import (
	"log"

	"zplus-saas/apps/backend/auth-service/internal/models"
	"zplus-saas/apps/backend/auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	authService services.AuthService
}

// AdminCreateRequest represents admin creation request
type AdminCreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Role      string `json:"role" validate:"required,oneof=admin super_admin"`
}

func NewAdminHandler(authService services.AuthService) *AdminHandler {
	return &AdminHandler{
		authService: authService,
	}
}

// AdminLogin godoc
// @Summary Admin login
// @Description Authenticate admin user and return tokens
// @Tags admin
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Admin login request"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/auth/login [post]
func (h *AdminHandler) AdminLogin(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required fields",
			Message: "Email and password are required",
		})
	}

	// Login through auth service
	tokens, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		log.Printf("Admin login error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Authentication failed",
			Message: "Invalid credentials",
		})
	}

	// Additional check for admin role
	user, err := h.authService.GetUserByID(c.Context(), tokens.User.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to verify user",
			Message: err.Error(),
		})
	}

	if user.Role != "admin" && user.Role != "super_admin" {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Error:   "Access denied",
			Message: "Admin access required",
		})
	}

	return c.JSON(tokens)
}

// CreateAdminUser godoc
// @Summary Create admin user
// @Description Create a new admin user (super_admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param request body AdminCreateRequest true "Admin creation request"
// @Success 201 {object} models.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/admin/auth/create [post]
func (h *AdminHandler) CreateAdminUser(c *fiber.Ctx) error {
	var req AdminCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Convert to RegisterRequest
	registerReq := &models.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Register through auth service
	tokens, err := h.authService.Register(c.Context(), registerReq)
	if err != nil {
		log.Printf("Admin creation error: %v", err)
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:   "Admin creation failed",
			Message: err.Error(),
		})
	}

	// Update user role to super_admin after creation
	user := tokens.User
	user.Role = req.Role
	err = h.authService.UpdateUser(c.Context(), user)
	if err != nil {
		log.Printf("Failed to update admin role: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to set admin role",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tokens)
}

// ValidateAdmin godoc
// @Summary Validate admin token
// @Description Validate admin token and return user info
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/auth/validate [get]
func (h *AdminHandler) ValidateAdmin(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
	}

	user, err := h.authService.GetUserByID(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	// Check admin role
	if user.Role != "admin" && user.Role != "super_admin" {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Error:   "Access denied",
			Message: "Admin access required",
		})
	}

	return c.JSON(user)
}

// AdminStats godoc
// @Summary Get admin dashboard stats
// @Description Get dashboard statistics for admin panel
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/stats [get]
func (h *AdminHandler) AdminStats(c *fiber.Ctx) error {
	// For now, return mock data
	stats := map[string]interface{}{
		"totalTenants":        25,
		"totalUsers":          1500,
		"totalRevenue":        125000,
		"activeSubscriptions": 23,
		"growth": map[string]interface{}{
			"tenants": 15.5,
			"users":   23.2,
			"revenue": 18.7,
		},
	}

	return c.JSON(stats)
}

// AdminActivities godoc
// @Summary Get recent admin activities
// @Description Get recent activities for admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/activities [get]
func (h *AdminHandler) AdminActivities(c *fiber.Ctx) error {
	// For now, return mock data
	activities := []map[string]interface{}{
		{
			"id":          "1",
			"type":        "tenant_created",
			"title":       "New Tenant Created",
			"description": "Acme Corp has been created",
			"timestamp":   "2025-07-24T10:30:00Z",
		},
		{
			"id":          "2",
			"type":        "user_registered",
			"title":       "User Registration",
			"description": "john.doe@example.com registered",
			"timestamp":   "2025-07-24T10:15:00Z",
		},
		{
			"id":          "3",
			"type":        "payment_received",
			"title":       "Payment Received",
			"description": "$99.00 payment from Tech Solutions",
			"timestamp":   "2025-07-24T09:45:00Z",
		},
	}

	return c.JSON(activities)
}

// AdminHealth godoc
// @Summary Get system health status
// @Description Get health status of all system services
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/health [get]
func (h *AdminHandler) AdminHealth(c *fiber.Ctx) error {
	// For now, return mock data
	health := map[string]interface{}{
		"status": "healthy",
		"services": []map[string]interface{}{
			{
				"name":          "API Gateway",
				"status":        "healthy",
				"response_time": 45,
				"last_check":    "2025-07-24T11:30:00Z",
			},
			{
				"name":          "Auth Service",
				"status":        "healthy",
				"response_time": 32,
				"last_check":    "2025-07-24T11:30:00Z",
			},
			{
				"name":          "Tenant Service",
				"status":        "healthy",
				"response_time": 28,
				"last_check":    "2025-07-24T11:30:00Z",
			},
			{
				"name":          "Database",
				"status":        "healthy",
				"response_time": 15,
				"last_check":    "2025-07-24T11:30:00Z",
			},
			{
				"name":          "Redis",
				"status":        "healthy",
				"response_time": 8,
				"last_check":    "2025-07-24T11:30:00Z",
			},
		},
	}

	return c.JSON(health)
}
