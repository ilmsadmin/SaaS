package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/auth-service/internal/models"
)

type EmailVerificationService struct {
	db    *sqlx.DB
	email EmailService
}

func NewEmailVerificationService(db *sqlx.DB, emailService EmailService) *EmailVerificationService {
	return &EmailVerificationService{
		db:    db,
		email: emailService,
	}
}

// SendVerificationEmail sends verification email to user
func (s *EmailVerificationService) SendVerificationEmail(userID uuid.UUID) error {
	// Get user info
	var user models.User
	err := s.db.Get(&user, "SELECT id, email, first_name, last_name, is_verified FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.IsVerified {
		return fmt.Errorf("user email is already verified")
	}

	// Generate verification token
	token, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Invalidate any existing verification tokens
	_, err = s.db.Exec("UPDATE email_verifications SET verified_at = $1 WHERE user_id = $2 AND verified_at IS NULL",
		time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to invalidate existing tokens: %w", err)
	}

	// Create new verification record
	verification := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO email_verifications (id, user_id, email, token, expires_at, created_at)
		VALUES (:id, :user_id, :email, :token, :expires_at, :created_at)`

	_, err = s.db.NamedExec(query, verification)
	if err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}

	// Send verification email
	userName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	err = s.email.SendVerificationEmail(user.Email, token, userName)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// VerifyEmail verifies user email with token
func (s *EmailVerificationService) VerifyEmail(token string) error {
	// Find valid verification token
	var verification models.EmailVerification
	err := s.db.Get(&verification,
		"SELECT * FROM email_verifications WHERE token = $1 AND verified_at IS NULL AND expires_at > $2",
		token, time.Now())
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user verification status
	_, err = tx.Exec("UPDATE users SET is_verified = true, updated_at = $1 WHERE id = $2",
		time.Now(), verification.UserID)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	// Mark verification as completed
	_, err = tx.Exec("UPDATE email_verifications SET verified_at = $1 WHERE id = $2",
		time.Now(), verification.ID)
	if err != nil {
		return fmt.Errorf("failed to mark verification as completed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ResendVerificationEmail resends verification email
func (s *EmailVerificationService) ResendVerificationEmail(email string) error {
	// Find user by email
	var user models.User
	err := s.db.Get(&user, "SELECT id, email, first_name, last_name, is_verified FROM users WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.IsVerified {
		return fmt.Errorf("user email is already verified")
	}

	return s.SendVerificationEmail(user.ID)
}

// CheckVerificationStatus checks if user email is verified
func (s *EmailVerificationService) CheckVerificationStatus(userID uuid.UUID) (bool, error) {
	var isVerified bool
	err := s.db.Get(&isVerified, "SELECT is_verified FROM users WHERE id = $1", userID)
	if err != nil {
		return false, fmt.Errorf("user not found")
	}
	return isVerified, nil
}

// CleanupExpiredTokens removes expired verification tokens
func (s *EmailVerificationService) CleanupExpiredTokens() error {
	_, err := s.db.Exec("DELETE FROM email_verifications WHERE expires_at < $1 OR verified_at IS NOT NULL",
		time.Now().Add(-24*time.Hour))
	return err
}
