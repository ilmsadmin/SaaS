package services

import (
	"context"
	"fmt"
	"time"

	"zplus-saas/apps/backend/auth-service/internal/models"
	"zplus-saas/apps/backend/auth-service/internal/repositories"
	"zplus-saas/apps/backend/shared/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.TokenResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type authService struct {
	userRepo         repositories.UserRepository
	tenantRepo       repositories.TenantRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	jwtService       JWTService
	config           *config.Config
}

func NewAuthService(
	userRepo repositories.UserRepository,
	tenantRepo repositories.TenantRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtService JWTService,
	config *config.Config,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		config:           config,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.TokenResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Parse tenant ID if provided
	var tenantID uuid.UUID
	if req.TenantID != "" {
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("invalid tenant ID: %w", err)
		}

		// Verify tenant exists
		_, err = s.tenantRepo.GetByID(ctx, tenantID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
	} else {
		// Create a new tenant for this user (first user becomes admin)
		tenant := &models.Tenant{
			Name:      fmt.Sprintf("%s %s's Organization", req.FirstName, req.LastName),
			Subdomain: fmt.Sprintf("org-%s", uuid.New().String()[:8]),
			PlanType:  models.PlanFree,
		}

		err = s.tenantRepo.Create(ctx, tenant)
		if err != nil {
			return nil, fmt.Errorf("failed to create tenant: %w", err)
		}
		tenantID = tenant.ID
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		TenantID:  tenantID,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
	}

	// If this is the first user in tenant, make them admin
	if req.TenantID == "" {
		user.Role = models.RoleAdmin
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.jwtService.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshTokenRecord := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(168 * time.Hour), // 7 days
	}

	err = s.refreshTokenRepo.Create(ctx, refreshTokenRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokens, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	err = s.verifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login for user %s: %v\n", user.ID, err)
	}

	// Generate tokens
	tokens, err := s.jwtService.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshTokenRecord := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(168 * time.Hour), // 7 days
	}

	err = s.refreshTokenRepo.Create(ctx, refreshTokenRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokens, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	token, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, err := s.jwtService.ParseUserClaims(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	// Check if refresh token exists in database
	storedToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}

	// Check if token is expired or used
	if storedToken.IsUsed || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token is expired or already used")
	}

	// Mark old token as used
	err = s.refreshTokenRepo.MarkAsUsed(ctx, storedToken.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new tokens
	tokens, err := s.jwtService.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Store new refresh token
	newRefreshTokenRecord := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(168 * time.Hour), // 7 days
	}

	err = s.refreshTokenRepo.Create(ctx, newRefreshTokenRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return tokens, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	// Invalidate all refresh tokens for the user
	err := s.refreshTokenRepo.DeleteByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate refresh tokens: %w", err)
	}

	return nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *authService) UpdateUser(ctx context.Context, user *models.User) error {
	err := s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Helper methods
func (s *authService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *authService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
