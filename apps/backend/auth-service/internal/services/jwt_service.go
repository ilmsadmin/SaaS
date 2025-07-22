package services

import (
	"fmt"
	"time"

	"zplus-saas/apps/backend/auth-service/internal/models"
	"zplus-saas/apps/backend/shared/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateTokens(user *models.User) (*models.TokenResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ParseUserClaims(token *jwt.Token) (*UserClaims, error)
	GenerateRefreshToken() (string, error)
}

type UserClaims struct {
	UserID     string `json:"user_id"`
	TenantID   string `json:"tenant_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	jwt.RegisteredClaims
}

type jwtService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		config: cfg,
	}
}

func (s *jwtService) GenerateTokens(user *models.User) (*models.TokenResponse, error) {
	// Parse expiration durations
	accessDuration, err := time.ParseDuration(s.config.JWTExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration: %w", err)
	}

	refreshDuration, err := time.ParseDuration(s.config.JWTRefreshExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration: %w", err)
	}

	now := time.Now()
	accessExpiry := now.Add(accessDuration)
	refreshExpiry := now.Add(refreshDuration)

	// Create access token claims
	accessClaims := &UserClaims{
		UserID:     user.ID.String(),
		TenantID:   user.TenantID.String(),
		Email:      user.Email,
		Role:       user.Role,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.JWTIssuer,
			Audience:  []string{s.config.JWTAudience},
			Subject:   user.ID.String(),
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token claims
	refreshClaims := &UserClaims{
		UserID:   user.ID.String(),
		TenantID: user.TenantID.String(),
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.JWTIssuer,
			Audience:  []string{s.config.JWTAudience},
			Subject:   user.ID.String(),
		},
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessDuration.Seconds()),
		User:         user,
	}, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return token, nil
}

func (s *jwtService) ParseUserClaims(token *jwt.Token) (*UserClaims, error) {
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *jwtService) GenerateRefreshToken() (string, error) {
	// For simplicity, we use JWT for refresh tokens too
	// In production, you might want to use random tokens stored in database
	refreshID := uuid.New().String()
	now := time.Now()

	refreshDuration, err := time.ParseDuration(s.config.JWTRefreshExpiresIn)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token duration: %w", err)
	}

	claims := &jwt.RegisteredClaims{
		ID:        refreshID,
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshDuration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    s.config.JWTIssuer,
		Audience:  []string{s.config.JWTAudience},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}
