package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"zplus-saas/apps/backend/auth-service/internal/models"
)

type PasswordResetService struct {
	db    *sqlx.DB
	email EmailService
}

func NewPasswordResetService(db *sqlx.DB, emailService EmailService) *PasswordResetService {
	return &PasswordResetService{
		db:    db,
		email: emailService,
	}
}

// RequestPasswordReset initiates password reset process
func (s *PasswordResetService) RequestPasswordReset(email string) error {
	// Find user by email
	var user models.User
	err := s.db.Get(&user, "SELECT id, first_name, last_name, email FROM users WHERE email = $1 AND is_active = true", email)
	if err != nil {
		// Don't reveal if user exists for security
		return nil
	}

	// Generate reset token
	token, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Invalidate any existing reset tokens
	_, err = s.db.Exec("UPDATE password_resets SET used_at = $1 WHERE user_id = $2 AND used_at IS NULL",
		time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to invalidate existing tokens: %w", err)
	}

	// Create new reset record
	reset := &models.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO password_resets (id, user_id, token, expires_at, created_at)
		VALUES (:id, :user_id, :token, :expires_at, :created_at)`

	_, err = s.db.NamedExec(query, reset)
	if err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// Send reset email
	userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	err = s.email.SendPasswordResetEmail(user.Email, token, userName)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword completes the password reset process
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	// Find valid reset token
	var reset models.PasswordReset
	err := s.db.Get(&reset,
		"SELECT * FROM password_resets WHERE token = $1 AND used_at IS NULL AND expires_at > $2",
		token, time.Now())
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user password
	_, err = tx.Exec("UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3",
		string(hashedPassword), time.Now(), reset.UserID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset token as used
	_, err = tx.Exec("UPDATE password_resets SET used_at = $1 WHERE id = $2",
		time.Now(), reset.ID)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Invalidate all refresh tokens for this user
	_, err = tx.Exec("UPDATE refresh_tokens SET is_used = true WHERE user_id = $1",
		reset.UserID)
	if err != nil {
		return fmt.Errorf("failed to invalidate refresh tokens: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ChangePassword changes password for authenticated user
func (s *PasswordResetService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	// Get current password hash
	var user models.User
	err := s.db.Get(&user, "SELECT password_hash FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update password
	_, err = tx.Exec("UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3",
		string(hashedPassword), time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all refresh tokens except current session
	_, err = tx.Exec("UPDATE refresh_tokens SET is_used = true WHERE user_id = $1",
		userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate refresh tokens: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired reset tokens
func (s *PasswordResetService) CleanupExpiredTokens() error {
	_, err := s.db.Exec("DELETE FROM password_resets WHERE expires_at < $1 OR used_at IS NOT NULL",
		time.Now().Add(-24*time.Hour))
	return err
}
