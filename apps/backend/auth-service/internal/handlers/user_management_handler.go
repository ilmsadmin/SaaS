package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"zplus-saas/apps/backend/auth-service/internal/models"
	"zplus-saas/apps/backend/auth-service/internal/services"
)

type UserManagementHandler struct {
	invitationService        *services.InvitationService
	passwordResetService     *services.PasswordResetService
	emailVerificationService *services.EmailVerificationService
	userProfileService       *services.UserProfileService
	validator                *validator.Validate
}

func NewUserManagementHandler(
	invitationService *services.InvitationService,
	passwordResetService *services.PasswordResetService,
	emailVerificationService *services.EmailVerificationService,
	userProfileService *services.UserProfileService,
) *UserManagementHandler {
	return &UserManagementHandler{
		invitationService:        invitationService,
		passwordResetService:     passwordResetService,
		emailVerificationService: emailVerificationService,
		userProfileService:       userProfileService,
		validator:                validator.New(),
	}
}

// InviteUser invites a new user to the tenant
func (h *UserManagementHandler) InviteUser(c *fiber.Ctx) error {
	var req models.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Get tenant ID from context (set by middleware)
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	// Get user ID from context (authenticated user)
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	invitation, err := h.invitationService.InviteUser(tenantID, userID, req.Email, req.Role)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":    "Invitation sent successfully",
		"invitation": invitation,
	})
}

// GetInvitations returns all invitations for the tenant
func (h *UserManagementHandler) GetInvitations(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	invitations, err := h.invitationService.GetInvitations(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"invitations": invitations})
}

// AcceptInvitation accepts a user invitation
func (h *UserManagementHandler) AcceptInvitation(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Token is required"})
	}

	var req models.AcceptInvitationRequest
	req.Token = token
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.invitationService.AcceptInvitation(token, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Invitation accepted successfully",
		"user":    user,
	})
}

// RevokeInvitation revokes a pending invitation
func (h *UserManagementHandler) RevokeInvitation(c *fiber.Ctx) error {
	invitationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid invitation ID"})
	}

	err = h.invitationService.RevokeInvitation(invitationID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Invitation revoked successfully"})
}

// ResendInvitation resends an invitation email
func (h *UserManagementHandler) ResendInvitation(c *fiber.Ctx) error {
	invitationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid invitation ID"})
	}

	err = h.invitationService.ResendInvitation(invitationID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Invitation resent successfully"})
}

// RequestPasswordReset initiates password reset
func (h *UserManagementHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.passwordResetService.RequestPasswordReset(req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password reset email sent"})
}

// ResetPassword completes password reset
func (h *UserManagementHandler) ResetPassword(c *fiber.Ctx) error {
	var req models.ConfirmResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.passwordResetService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

// ChangePassword changes password for authenticated user
func (h *UserManagementHandler) ChangePassword(c *fiber.Ctx) error {
	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.passwordResetService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

// SendVerificationEmail sends email verification
func (h *UserManagementHandler) SendVerificationEmail(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.emailVerificationService.SendVerificationEmail(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Verification email sent"})
}

// VerifyEmail verifies user email
func (h *UserManagementHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Token is required"})
	}

	err := h.emailVerificationService.VerifyEmail(token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Email verified successfully"})
}

// GetProfile gets user profile
func (h *UserManagementHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, profile, err := h.userProfileService.GetUserWithProfile(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"user":    user,
		"profile": profile,
	})
}

// UpdateProfile updates user profile
func (h *UserManagementHandler) UpdateProfile(c *fiber.Ctx) error {
	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.userProfileService.UpdateProfile(userID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}
