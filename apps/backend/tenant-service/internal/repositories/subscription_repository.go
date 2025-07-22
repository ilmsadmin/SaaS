package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/tenant-service/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetExpiring(ctx context.Context, days int) ([]*models.Subscription, error)
}

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, tenant_id, plan_id, status, trial_end_at, current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	subscription.ID = uuid.New()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	if subscription.Status == "" {
		subscription.Status = models.SubscriptionStatusActive
	}

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.TenantID,
		subscription.PlanID,
		subscription.Status,
		subscription.TrialEndAt,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	return err
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT s.id, s.tenant_id, s.plan_id, s.status, s.trial_end_at, 
			   s.current_period_start, s.current_period_end, s.cancel_at_period_end, 
			   s.created_at, s.updated_at,
			   t.id, t.name, t.subdomain, t.domain, t.logo, t.status, t.settings, t.created_at, t.updated_at,
			   p.id, p.name, p.description, p.price, p.currency, p.billing_cycle, 
			   p.max_users, p.max_storage, p.features, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		LEFT JOIN tenants t ON s.tenant_id = t.id
		LEFT JOIN plans p ON s.plan_id = p.id
		WHERE s.id = $1
	`

	subscription := &models.Subscription{
		Tenant: &models.Tenant{},
		Plan:   &models.Plan{},
	}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.TrialEndAt,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Tenant.ID,
		&subscription.Tenant.Name,
		&subscription.Tenant.Subdomain,
		&subscription.Tenant.Domain,
		&subscription.Tenant.Logo,
		&subscription.Tenant.Status,
		&subscription.Tenant.Settings,
		&subscription.Tenant.CreatedAt,
		&subscription.Tenant.UpdatedAt,
		&subscription.Plan.ID,
		&subscription.Plan.Name,
		&subscription.Plan.Description,
		&subscription.Plan.Price,
		&subscription.Plan.Currency,
		&subscription.Plan.BillingCycle,
		&subscription.Plan.MaxUsers,
		&subscription.Plan.MaxStorage,
		&subscription.Plan.Features,
		&subscription.Plan.IsActive,
		&subscription.Plan.CreatedAt,
		&subscription.Plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}

	return subscription, nil
}

func (r *subscriptionRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT s.id, s.tenant_id, s.plan_id, s.status, s.trial_end_at, 
			   s.current_period_start, s.current_period_end, s.cancel_at_period_end, 
			   s.created_at, s.updated_at,
			   p.id, p.name, p.description, p.price, p.currency, p.billing_cycle, 
			   p.max_users, p.max_storage, p.features, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		LEFT JOIN plans p ON s.plan_id = p.id
		WHERE s.tenant_id = $1 AND s.status = 'active'
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	subscription := &models.Subscription{
		Plan: &models.Plan{},
	}

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.TrialEndAt,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Plan.ID,
		&subscription.Plan.Name,
		&subscription.Plan.Description,
		&subscription.Plan.Price,
		&subscription.Plan.Currency,
		&subscription.Plan.BillingCycle,
		&subscription.Plan.MaxUsers,
		&subscription.Plan.MaxStorage,
		&subscription.Plan.Features,
		&subscription.Plan.IsActive,
		&subscription.Plan.CreatedAt,
		&subscription.Plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}

	return subscription, nil
}

func (r *subscriptionRepository) List(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM subscriptions`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT s.id, s.tenant_id, s.plan_id, s.status, s.trial_end_at, 
			   s.current_period_start, s.current_period_end, s.cancel_at_period_end, 
			   s.created_at, s.updated_at,
			   t.id, t.name, t.subdomain, t.domain, t.logo, t.status, t.settings, t.created_at, t.updated_at,
			   p.id, p.name, p.description, p.price, p.currency, p.billing_cycle, 
			   p.max_users, p.max_storage, p.features, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		LEFT JOIN tenants t ON s.tenant_id = t.id
		LEFT JOIN plans p ON s.plan_id = p.id
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{
			Tenant: &models.Tenant{},
			Plan:   &models.Plan{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.TenantID,
			&subscription.PlanID,
			&subscription.Status,
			&subscription.TrialEndAt,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.CancelAtPeriodEnd,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Tenant.ID,
			&subscription.Tenant.Name,
			&subscription.Tenant.Subdomain,
			&subscription.Tenant.Domain,
			&subscription.Tenant.Logo,
			&subscription.Tenant.Status,
			&subscription.Tenant.Settings,
			&subscription.Tenant.CreatedAt,
			&subscription.Tenant.UpdatedAt,
			&subscription.Plan.ID,
			&subscription.Plan.Name,
			&subscription.Plan.Description,
			&subscription.Plan.Price,
			&subscription.Plan.Currency,
			&subscription.Plan.BillingCycle,
			&subscription.Plan.MaxUsers,
			&subscription.Plan.MaxStorage,
			&subscription.Plan.Features,
			&subscription.Plan.IsActive,
			&subscription.Plan.CreatedAt,
			&subscription.Plan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, total, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET plan_id = $2, status = $3, trial_end_at = $4, 
			current_period_start = $5, current_period_end = $6, 
			cancel_at_period_end = $7, updated_at = $8
		WHERE id = $1
	`

	subscription.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.PlanID,
		subscription.Status,
		subscription.TrialEndAt,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		subscription.UpdatedAt,
	)

	return err
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *subscriptionRepository) GetExpiring(ctx context.Context, days int) ([]*models.Subscription, error) {
	query := `
		SELECT s.id, s.tenant_id, s.plan_id, s.status, s.trial_end_at, 
			   s.current_period_start, s.current_period_end, s.cancel_at_period_end, 
			   s.created_at, s.updated_at,
			   t.id, t.name, t.subdomain, t.domain, t.logo, t.status, t.settings, t.created_at, t.updated_at,
			   p.id, p.name, p.description, p.price, p.currency, p.billing_cycle, 
			   p.max_users, p.max_storage, p.features, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		LEFT JOIN tenants t ON s.tenant_id = t.id
		LEFT JOIN plans p ON s.plan_id = p.id
		WHERE s.status = 'active' 
		AND s.current_period_end <= NOW() + INTERVAL '%d days'
		ORDER BY s.current_period_end ASC
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(query, days))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{
			Tenant: &models.Tenant{},
			Plan:   &models.Plan{},
		}

		err := rows.Scan(
			&subscription.ID,
			&subscription.TenantID,
			&subscription.PlanID,
			&subscription.Status,
			&subscription.TrialEndAt,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.CancelAtPeriodEnd,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Tenant.ID,
			&subscription.Tenant.Name,
			&subscription.Tenant.Subdomain,
			&subscription.Tenant.Domain,
			&subscription.Tenant.Logo,
			&subscription.Tenant.Status,
			&subscription.Tenant.Settings,
			&subscription.Tenant.CreatedAt,
			&subscription.Tenant.UpdatedAt,
			&subscription.Plan.ID,
			&subscription.Plan.Name,
			&subscription.Plan.Description,
			&subscription.Plan.Price,
			&subscription.Plan.Currency,
			&subscription.Plan.BillingCycle,
			&subscription.Plan.MaxUsers,
			&subscription.Plan.MaxStorage,
			&subscription.Plan.Features,
			&subscription.Plan.IsActive,
			&subscription.Plan.CreatedAt,
			&subscription.Plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}
