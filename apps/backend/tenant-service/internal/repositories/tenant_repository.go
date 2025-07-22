package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/tenant-service/internal/models"

	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant *models.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	List(ctx context.Context, limit, offset int) ([]*models.Tenant, int, error)
	Update(ctx context.Context, tenant *models.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type tenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, subdomain, domain, logo, status, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	tenant.ID = uuid.New()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	if tenant.Status == "" {
		tenant.Status = models.TenantStatusTrial
	}

	if tenant.Settings == "" {
		tenant.Settings = "{}"
	}

	_, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Subdomain,
		tenant.Domain,
		tenant.Logo,
		tenant.Status,
		tenant.Settings,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)

	return err
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	query := `
		SELECT id, name, subdomain, domain, logo, status, settings, created_at, updated_at
		FROM tenants 
		WHERE id = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Logo,
		&tenant.Status,
		&tenant.Settings,
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
		SELECT id, name, subdomain, domain, logo, status, settings, created_at, updated_at
		FROM tenants 
		WHERE subdomain = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, subdomain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Logo,
		&tenant.Status,
		&tenant.Settings,
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
		SELECT id, name, subdomain, domain, logo, status, settings, created_at, updated_at
		FROM tenants 
		WHERE domain = $1
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, domain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Logo,
		&tenant.Status,
		&tenant.Settings,
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

func (r *tenantRepository) List(ctx context.Context, limit, offset int) ([]*models.Tenant, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM tenants`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, name, subdomain, domain, logo, status, settings, created_at, updated_at
		FROM tenants 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tenants []*models.Tenant
	for rows.Next() {
		tenant := &models.Tenant{}
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Subdomain,
			&tenant.Domain,
			&tenant.Logo,
			&tenant.Status,
			&tenant.Settings,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, total, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
		UPDATE tenants 
		SET name = $2, subdomain = $3, domain = $4, logo = $5, 
			status = $6, settings = $7, updated_at = $8
		WHERE id = $1
	`

	tenant.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Subdomain,
		tenant.Domain,
		tenant.Logo,
		tenant.Status,
		tenant.Settings,
		tenant.UpdatedAt,
	)

	return err
}

func (r *tenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *tenantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE tenants SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	return err
}
