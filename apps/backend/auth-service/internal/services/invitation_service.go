package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"zplus-saas/apps/backend/auth-service/internal/models"
)

type InvitationService struct {
	db    *sqlx.DB
	email EmailService
}

type EmailService interface {
	SendInvitationEmail(email, token, inviterName, tenantName string) error
	SendPasswordResetEmail(email, token, userName string) error
	SendVerificationEmail(email, token, userName string) error
}

func NewInvitationService(db *sqlx.DB, emailService EmailService) *InvitationService {
	return &InvitationService{
		db:    db,
		email: emailService,
	}
}

// InviteUser creates and sends a user invitation
func (s *InvitationService) InviteUser(tenantID, inviterID uuid.UUID, email, role string) (*models.UserInvitation, error) {
	// Check if user already exists
	var existingUser models.User
	err := s.db.Get(&existingUser,
		"SELECT id FROM users WHERE email = $1 AND tenant_id = $2",
		email, tenantID)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists in this tenant", email)
	}

	// Check if there's already a pending invitation
	var existingInvitation models.UserInvitation
	err = s.db.Get(&existingInvitation,
		"SELECT id FROM user_invitations WHERE email = $1 AND tenant_id = $2 AND status = $3",
		email, tenantID, models.InvitationStatusPending)
	if err == nil {
		return nil, fmt.Errorf("invitation already sent to %s", email)
	}

	// Generate invitation token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create invitation
	invitation := &models.UserInvitation{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		Token:     token,
		InvitedBy: inviterID,
		Status:    models.InvitationStatusPending,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	query := `
		INSERT INTO user_invitations (id, tenant_id, email, role, token, invited_by, status, expires_at, created_at, updated_at)
		VALUES (:id, :tenant_id, :email, :role, :token, :invited_by, :status, :expires_at, :created_at, :updated_at)`

	_, err = s.db.NamedExec(query, invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to save invitation: %w", err)
	}

	// Get inviter and tenant info for email
	var inviter models.User
	var tenant models.Tenant

	s.db.Get(&inviter, "SELECT first_name, last_name FROM users WHERE id = $1", inviterID)
	s.db.Get(&tenant, "SELECT name FROM tenants WHERE id = $1", tenantID)

	inviterName := fmt.Sprintf("%s %s", inviter.FirstName, inviter.LastName)

	// Send invitation email
	err = s.email.SendInvitationEmail(email, token, inviterName, tenant.Name)
	if err != nil {
		log.Printf("Failed to send invitation email to %s: %v", email, err)
		// Don't fail the invitation if email fails
	}

	return invitation, nil
}

// AcceptInvitation processes invitation acceptance
func (s *InvitationService) AcceptInvitation(token string, req *models.AcceptInvitationRequest) (*models.User, error) {
	// Find invitation
	var invitation models.UserInvitation
	err := s.db.Get(&invitation,
		"SELECT * FROM user_invitations WHERE token = $1 AND status = $2",
		token, models.InvitationStatusPending)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired invitation")
	}

	// Check if expired
	if time.Now().After(invitation.ExpiresAt) {
		// Update status to expired
		s.db.Exec("UPDATE user_invitations SET status = $1 WHERE id = $2",
			models.InvitationStatusExpired, invitation.ID)
		return nil, fmt.Errorf("invitation has expired")
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:         uuid.New(),
		TenantID:   invitation.TenantID,
		Email:      invitation.Email,
		Password:   hashedPassword,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Role:       invitation.Role,
		IsActive:   true,
		IsVerified: true, // Auto-verify invited users
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user
	userQuery := `
		INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, is_active, is_verified, created_at, updated_at)
		VALUES (:id, :tenant_id, :email, :password, :first_name, :last_name, :role, :is_active, :is_verified, :created_at, :updated_at)`

	_, err = tx.NamedExec(userQuery, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Update invitation status
	_, err = tx.Exec("UPDATE user_invitations SET status = $1, accepted_at = $2, updated_at = $3 WHERE id = $4",
		models.InvitationStatusAccepted, time.Now(), time.Now(), invitation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	// Create default user profile
	profile := &models.UserProfile{
		ID:        uuid.New(),
		UserID:    user.ID,
		Language:  "en",
		Timezone:  "UTC",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	profileQuery := `
		INSERT INTO user_profiles (id, user_id, language, timezone, created_at, updated_at)
		VALUES (:id, :user_id, :language, :timezone, :created_at, :updated_at)`

	_, err = tx.NamedExec(profileQuery, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user, nil
}

// GetInvitations returns invitations for a tenant
func (s *InvitationService) GetInvitations(tenantID uuid.UUID) ([]models.UserInvitation, error) {
	var invitations []models.UserInvitation

	query := `
		SELECT ui.*, u.first_name || ' ' || u.last_name as inviter_name
		FROM user_invitations ui
		LEFT JOIN users u ON ui.invited_by = u.id
		WHERE ui.tenant_id = $1
		ORDER BY ui.created_at DESC`

	err := s.db.Select(&invitations, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	return invitations, nil
}

// RevokeInvitation revokes a pending invitation
func (s *InvitationService) RevokeInvitation(invitationID uuid.UUID) error {
	result, err := s.db.Exec(
		"UPDATE user_invitations SET status = 'expired', updated_at = $1 WHERE id = $2 AND status = 'pending'",
		time.Now(), invitationID)
	if err != nil {
		return fmt.Errorf("failed to revoke invitation: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found or already processed")
	}

	return nil
}

// ResendInvitation resends an invitation email
func (s *InvitationService) ResendInvitation(invitationID uuid.UUID) error {
	var invitation models.UserInvitation
	err := s.db.Get(&invitation,
		"SELECT * FROM user_invitations WHERE id = $1 AND status = $2",
		invitationID, models.InvitationStatusPending)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	// Check if expired and update expiration
	if time.Now().After(invitation.ExpiresAt) {
		invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
		s.db.Exec("UPDATE user_invitations SET expires_at = $1, updated_at = $2 WHERE id = $3",
			invitation.ExpiresAt, time.Now(), invitation.ID)
	}

	// Get inviter and tenant info
	var inviter models.User
	var tenant models.Tenant

	s.db.Get(&inviter, "SELECT first_name, last_name FROM users WHERE id = $1", invitation.InvitedBy)
	s.db.Get(&tenant, "SELECT name FROM tenants WHERE id = $1", invitation.TenantID)

	inviterName := fmt.Sprintf("%s %s", inviter.FirstName, inviter.LastName)

	// Resend email
	return s.email.SendInvitationEmail(invitation.Email, invitation.Token, inviterName, tenant.Name)
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
