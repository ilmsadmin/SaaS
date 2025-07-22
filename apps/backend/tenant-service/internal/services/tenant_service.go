package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zplus-saas/apps/backend/tenant-service/internal/models"
	"zplus-saas/apps/backend/tenant-service/internal/repositories"

	"github.com/google/uuid"
)

type TenantService interface {
	// Tenant management
	CreateTenant(ctx context.Context, req *models.CreateTenantRequest) (*models.Tenant, error)
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error)
	GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	ListTenants(ctx context.Context, limit, offset int) ([]*models.Tenant, int, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, req *models.UpdateTenantRequest) (*models.Tenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	ActivateTenant(ctx context.Context, id uuid.UUID) error
	SuspendTenant(ctx context.Context, id uuid.UUID) error

	// Plan management
	CreatePlan(ctx context.Context, req *models.CreatePlanRequest) (*models.Plan, error)
	GetPlan(ctx context.Context, id uuid.UUID) (*models.Plan, error)
	ListPlans(ctx context.Context, limit, offset int) ([]*models.Plan, int, error)
	ListActivePlans(ctx context.Context) ([]*models.Plan, error)
	UpdatePlan(ctx context.Context, id uuid.UUID, req *models.UpdatePlanRequest) (*models.Plan, error)
	DeletePlan(ctx context.Context, id uuid.UUID) error

	// Subscription management
	CreateSubscription(ctx context.Context, tenantID uuid.UUID, req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubscription(ctx context.Context, tenantID uuid.UUID) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, tenantID uuid.UUID, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, tenantID uuid.UUID) error
	GetExpiringSubscriptions(ctx context.Context, days int) ([]*models.Subscription, error)
}

type tenantService struct {
	tenantRepo       repositories.TenantRepository
	subscriptionRepo repositories.SubscriptionRepository
	planRepo         repositories.PlanRepository
}

func NewTenantService(
	tenantRepo repositories.TenantRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	planRepo repositories.PlanRepository,
) TenantService {
	return &tenantService{
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
	}
}

// Tenant management
func (s *tenantService) CreateTenant(ctx context.Context, req *models.CreateTenantRequest) (*models.Tenant, error) {
	// Validate subdomain
	if err := s.validateSubdomain(req.Subdomain); err != nil {
		return nil, err
	}

	// Check if subdomain already exists
	existing, _ := s.tenantRepo.GetBySubdomain(ctx, req.Subdomain)
	if existing != nil {
		return nil, fmt.Errorf("subdomain already exists")
	}

	// Check domain if provided
	if req.Domain != nil && *req.Domain != "" {
		existing, _ := s.tenantRepo.GetByDomain(ctx, *req.Domain)
		if existing != nil {
			return nil, fmt.Errorf("domain already exists")
		}
	}

	tenant := &models.Tenant{
		Name:      req.Name,
		Subdomain: strings.ToLower(req.Subdomain),
		Domain:    req.Domain,
		Status:    models.TenantStatusTrial, // Start with trial
		Settings:  "{}",
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Create default subscription if plan is provided
	if req.PlanID != nil && *req.PlanID != "" {
		planUUID, err := uuid.Parse(*req.PlanID)
		if err != nil {
			return nil, fmt.Errorf("invalid plan ID: %w", err)
		}

		// Verify plan exists
		_, err = s.planRepo.GetByID(ctx, planUUID)
		if err != nil {
			return nil, fmt.Errorf("plan not found: %w", err)
		}

		// Create subscription
		now := time.Now()
		subscription := &models.Subscription{
			TenantID:           tenant.ID,
			PlanID:             planUUID,
			Status:             models.SubscriptionStatusActive,
			TrialEndAt:         &time.Time{}, // 14 days trial
			CurrentPeriodStart: now,
			CurrentPeriodEnd:   now.AddDate(0, 1, 0), // 1 month
			CancelAtPeriodEnd:  false,
		}

		// Set trial end date (14 days)
		trialEnd := now.AddDate(0, 0, 14)
		subscription.TrialEndAt = &trialEnd

		if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
			// If subscription creation fails, we should still return the tenant
			// but log the error
			// TODO: Add proper logging
		}
	}

	return tenant, nil
}

func (s *tenantService) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return s.tenantRepo.GetByID(ctx, id)
}

func (s *tenantService) GetTenantBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error) {
	return s.tenantRepo.GetBySubdomain(ctx, subdomain)
}

func (s *tenantService) GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	return s.tenantRepo.GetByDomain(ctx, domain)
}

func (s *tenantService) ListTenants(ctx context.Context, limit, offset int) ([]*models.Tenant, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.tenantRepo.List(ctx, limit, offset)
}

func (s *tenantService) UpdateTenant(ctx context.Context, id uuid.UUID, req *models.UpdateTenantRequest) (*models.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Domain != nil {
		// Check if domain already exists for another tenant
		if *req.Domain != "" {
			existing, _ := s.tenantRepo.GetByDomain(ctx, *req.Domain)
			if existing != nil && existing.ID != tenant.ID {
				return nil, fmt.Errorf("domain already exists")
			}
		}
		tenant.Domain = req.Domain
	}
	if req.Logo != nil {
		tenant.Logo = req.Logo
	}
	if req.Settings != nil {
		tenant.Settings = *req.Settings
	}

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

func (s *tenantService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement soft delete and cleanup logic
	// - Cancel subscriptions
	// - Archive data
	// - Send notifications
	return s.tenantRepo.Delete(ctx, id)
}

func (s *tenantService) ActivateTenant(ctx context.Context, id uuid.UUID) error {
	return s.tenantRepo.UpdateStatus(ctx, id, models.TenantStatusActive)
}

func (s *tenantService) SuspendTenant(ctx context.Context, id uuid.UUID) error {
	return s.tenantRepo.UpdateStatus(ctx, id, models.TenantStatusSuspended)
}

// Plan management
func (s *tenantService) CreatePlan(ctx context.Context, req *models.CreatePlanRequest) (*models.Plan, error) {
	plan := &models.Plan{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Currency:     req.Currency,
		BillingCycle: req.BillingCycle,
		MaxUsers:     req.MaxUsers,
		MaxStorage:   req.MaxStorage,
		IsActive:     true,
	}

	if req.Features != nil {
		plan.Features = *req.Features
	} else {
		plan.Features = "[]"
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	return plan, nil
}

func (s *tenantService) GetPlan(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	return s.planRepo.GetByID(ctx, id)
}

func (s *tenantService) ListPlans(ctx context.Context, limit, offset int) ([]*models.Plan, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.planRepo.List(ctx, limit, offset)
}

func (s *tenantService) ListActivePlans(ctx context.Context) ([]*models.Plan, error) {
	return s.planRepo.ListActive(ctx)
}

func (s *tenantService) UpdatePlan(ctx context.Context, id uuid.UUID, req *models.UpdatePlanRequest) (*models.Plan, error) {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Description != nil {
		plan.Description = req.Description
	}
	if req.Price != nil {
		plan.Price = *req.Price
	}
	if req.Currency != nil {
		plan.Currency = *req.Currency
	}
	if req.BillingCycle != nil {
		plan.BillingCycle = *req.BillingCycle
	}
	if req.MaxUsers != nil {
		plan.MaxUsers = req.MaxUsers
	}
	if req.MaxStorage != nil {
		plan.MaxStorage = req.MaxStorage
	}
	if req.Features != nil {
		plan.Features = *req.Features
	}
	if req.IsActive != nil {
		plan.IsActive = *req.IsActive
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	return plan, nil
}

func (s *tenantService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	// TODO: Check if plan is being used by any subscriptions
	return s.planRepo.Delete(ctx, id)
}

// Subscription management
func (s *tenantService) CreateSubscription(ctx context.Context, tenantID uuid.UUID, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	// Verify tenant exists
	_, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Parse and verify plan
	planUUID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}

	_, err = s.planRepo.GetByID(ctx, planUUID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	// Cancel existing active subscription
	existing, _ := s.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if existing != nil && existing.Status == models.SubscriptionStatusActive {
		existing.Status = models.SubscriptionStatusCancelled
		s.subscriptionRepo.Update(ctx, existing)
	}

	subscription := &models.Subscription{
		TenantID:           tenantID,
		PlanID:             planUUID,
		Status:             models.SubscriptionStatusActive,
		TrialEndAt:         req.TrialEndAt,
		CurrentPeriodStart: req.CurrentPeriodStart,
		CurrentPeriodEnd:   req.CurrentPeriodEnd,
		CancelAtPeriodEnd:  false,
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

func (s *tenantService) GetSubscription(ctx context.Context, tenantID uuid.UUID) (*models.Subscription, error) {
	return s.subscriptionRepo.GetByTenantID(ctx, tenantID)
}

func (s *tenantService) UpdateSubscription(ctx context.Context, tenantID uuid.UUID, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.PlanID != nil {
		planUUID, err := uuid.Parse(*req.PlanID)
		if err != nil {
			return nil, fmt.Errorf("invalid plan ID: %w", err)
		}

		_, err = s.planRepo.GetByID(ctx, planUUID)
		if err != nil {
			return nil, fmt.Errorf("plan not found: %w", err)
		}

		subscription.PlanID = planUUID
	}

	if req.Status != nil {
		subscription.Status = *req.Status
	}
	if req.TrialEndAt != nil {
		subscription.TrialEndAt = req.TrialEndAt
	}
	if req.CurrentPeriodStart != nil {
		subscription.CurrentPeriodStart = *req.CurrentPeriodStart
	}
	if req.CurrentPeriodEnd != nil {
		subscription.CurrentPeriodEnd = *req.CurrentPeriodEnd
	}
	if req.CancelAtPeriodEnd != nil {
		subscription.CancelAtPeriodEnd = *req.CancelAtPeriodEnd
	}

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

func (s *tenantService) CancelSubscription(ctx context.Context, tenantID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	subscription.Status = models.SubscriptionStatusCancelled
	subscription.CancelAtPeriodEnd = true

	return s.subscriptionRepo.Update(ctx, subscription)
}

func (s *tenantService) GetExpiringSubscriptions(ctx context.Context, days int) ([]*models.Subscription, error) {
	return s.subscriptionRepo.GetExpiring(ctx, days)
}

// Helper functions
func (s *tenantService) validateSubdomain(subdomain string) error {
	// Basic validation
	if len(subdomain) < 3 || len(subdomain) > 63 {
		return fmt.Errorf("subdomain must be between 3 and 63 characters")
	}

	// Check reserved subdomains
	reserved := []string{"www", "api", "app", "admin", "dashboard", "mail", "ftp", "blog", "support", "help", "docs", "status"}
	for _, r := range reserved {
		if strings.ToLower(subdomain) == r {
			return fmt.Errorf("subdomain '%s' is reserved", subdomain)
		}
	}

	// Check alphanumeric and hyphens only
	for _, char := range subdomain {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("subdomain can only contain letters, numbers, and hyphens")
		}
	}

	// Cannot start or end with hyphen
	if strings.HasPrefix(subdomain, "-") || strings.HasSuffix(subdomain, "-") {
		return fmt.Errorf("subdomain cannot start or end with a hyphen")
	}

	return nil
}
