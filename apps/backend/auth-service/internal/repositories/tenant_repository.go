package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/auth-service/internal/models"

	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant *models.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	Update(ctx context.Context, tenant *models.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type tenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, subdomain, domain, settings, plan_type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	tenant.ID = uuid.New()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now
	tenant.IsActive = true

	if tenant.PlanType == "" {
		tenant.PlanType = models.PlanFree
	}

	if tenant.Settings == "" {
		tenant.Settings = "{}"
	}

	_, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Subdomain,
		tenant.Domain,
		tenant.Settings,
		tenant.PlanType,
		tenant.IsActive,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)

	return err
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, domain, settings, plan_type, is_active, created_at, updated_at
		FROM tenants 
		WHERE id = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Settings,
		&tenant.PlanType,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, err
	}

	return tenant, nil
}

func (r *tenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, domain, settings, plan_type, is_active, created_at, updated_at
		FROM tenants 
		WHERE subdomain = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, subdomain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Settings,
		&tenant.PlanType,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, err
	}

	return tenant, nil
}

func (r *tenantRepository) GetByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, domain, settings, plan_type, is_active, created_at, updated_at
		FROM tenants 
		WHERE domain = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, domain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Settings,
		&tenant.PlanType,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, err
	}

	return tenant, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
		UPDATE tenants 
		SET name = $2, subdomain = $3, domain = $4, settings = $5, 
			plan_type = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	tenant.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Subdomain,
		tenant.Domain,
		tenant.Settings,
		tenant.PlanType,
		tenant.IsActive,
		tenant.UpdatedAt,
	)

	return err
}

func (r *tenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
