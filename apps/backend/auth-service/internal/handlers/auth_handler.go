package handlers

import (
	"log"

	"zplus-saas/apps/backend/auth-service/internal/models"
	"zplus-saas/apps/backend/auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required fields",
			Message: "Email, password, first_name, and last_name are required",
		})
	}

	tokens, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		log.Printf("Registration error: %v", err)
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tokens)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	tokens, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		log.Printf("Login error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Authentication failed",
			Message: "Invalid credentials",
		})
	}

	return c.JSON(tokens)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing refresh token",
			Message: "refresh_token is required",
		})
	}

	tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		log.Printf("Refresh token error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "Token refresh failed",
			Message: err.Error(),
		})
	}

	return c.JSON(tokens)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate all refresh tokens for the user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
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

	err = h.authService.Logout(c.Context(), uid)
	if err != nil {
		log.Printf("Logout error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Logout failed",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Successfully logged out",
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
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
		log.Printf("Get profile error: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	return c.JSON(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Update profile request"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
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

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Get current user
	user, err := h.authService.GetUserByID(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	err = h.authService.UpdateUser(c.Context(), user)
	if err != nil {
		log.Printf("Update profile error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Profile update failed",
			Message: err.Error(),
		})
	}

	return c.JSON(user)
}

// SetupUserManagementRoutes sets up user management routes
func SetupUserManagementRoutes(app *fiber.App, handler *UserManagementHandler) {
	api := app.Group("/api")

	// User invitation routes
	invitations := api.Group("/invitations")
	invitations.Post("/", handler.InviteUser)                 // POST /api/invitations
	invitations.Get("/", handler.GetInvitations)              // GET /api/invitations
	invitations.Post("/accept", handler.AcceptInvitation)     // POST /api/invitations/accept?token=xxx
	invitations.Delete("/:id", handler.RevokeInvitation)      // DELETE /api/invitations/:id
	invitations.Post("/:id/resend", handler.ResendInvitation) // POST /api/invitations/:id/resend

	// Password management routes
	password := api.Group("/password")
	password.Post("/reset/request", handler.RequestPasswordReset) // POST /api/password/reset/request
	password.Post("/reset/confirm", handler.ResetPassword)        // POST /api/password/reset/confirm
	password.Post("/change", handler.ChangePassword)              // POST /api/password/change (authenticated)

	// Email verification routes
	email := api.Group("/email")
	email.Post("/verify/send", handler.SendVerificationEmail) // POST /api/email/verify/send (authenticated)
	email.Post("/verify", handler.VerifyEmail)                // POST /api/email/verify?token=xxx

	// User profile routes
	profile := api.Group("/profile")
	profile.Get("/", handler.GetProfile)    // GET /api/profile (authenticated)
	profile.Put("/", handler.UpdateProfile) // PUT /api/profile (authenticated)
}

// Request/Response models
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
