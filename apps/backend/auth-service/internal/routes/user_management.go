package routes

import (
	"zplus-saas/apps/backend/auth-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupUserManagementRoutes sets up user management routes
func SetupUserManagementRoutes(app *fiber.App, handler *handlers.UserManagementHandler) {
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
