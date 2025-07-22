package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TenantID    uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password_hash"` // Don't expose password in JSON
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	Role        string     `json:"role" db:"role"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	IsVerified  bool       `json:"is_verified" db:"is_verified"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Tenant represents a tenant/organization in the system
type Tenant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Subdomain string    `json:"subdomain" db:"subdomain"`
	Domain    *string   `json:"domain" db:"domain"`
	Settings  string    `json:"settings" db:"settings"` // JSON string
	PlanType  string    `json:"plan_type" db:"plan_type"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RefreshToken represents a refresh token for JWT
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	TenantID  string `json:"tenant_id,omitempty"`
}

// TokenResponse represents token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

// UserInvitation represents a user invitation
type UserInvitation struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	TenantID   uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Email      string     `json:"email" db:"email"`
	Role       string     `json:"role" db:"role"`
	Token      string     `json:"token" db:"token"`
	InvitedBy  uuid.UUID  `json:"invited_by" db:"invited_by"`
	Status     string     `json:"status" db:"status"` // pending, accepted, expired
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	AcceptedAt *time.Time `json:"accepted_at" db:"accepted_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// EmailVerification represents an email verification token
type EmailVerification struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Email      string     `json:"email" db:"email"`
	Token      string     `json:"token" db:"token"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at" db:"verified_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Avatar      *string    `json:"avatar" db:"avatar"`
	Phone       *string    `json:"phone" db:"phone"`
	Address     *string    `json:"address" db:"address"`
	City        *string    `json:"city" db:"city"`
	Country     *string    `json:"country" db:"country"`
	PostalCode  *string    `json:"postal_code" db:"postal_code"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
	Bio         *string    `json:"bio" db:"bio"`
	Language    string     `json:"language" db:"language" default:"en"`
	Timezone    string     `json:"timezone" db:"timezone" default:"UTC"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response DTOs for User Management
type InviteUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin manager user"`
}

type AcceptInvitationRequest struct {
	Token     string `json:"token" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

type UpdateProfileRequest struct {
	FirstName   *string    `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName    *string    `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Avatar      *string    `json:"avatar,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Address     *string    `json:"address,omitempty"`
	City        *string    `json:"city,omitempty"`
	Country     *string    `json:"country,omitempty"`
	PostalCode  *string    `json:"postal_code,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	Language    *string    `json:"language,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Constants for invitations and verification
const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusExpired  = "expired"
)

// UserRole constants
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleUser       = "user"
)

// PlanType constants
const (
	PlanFree         = "free"
	PlanBasic        = "basic"
	PlanProfessional = "professional"
	PlanEnterprise   = "enterprise"
)
