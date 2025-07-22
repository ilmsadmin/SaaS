package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/tenant-service/internal/models"

	"github.com/google/uuid"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *models.Plan) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error)
	List(ctx context.Context, limit, offset int) ([]*models.Plan, int, error)
	ListActive(ctx context.Context) ([]*models.Plan, error)
	Update(ctx context.Context, plan *models.Plan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type planRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) Create(ctx context.Context, plan *models.Plan) error {
	query := `
		INSERT INTO plans (id, name, description, price, currency, billing_cycle, max_users, max_storage, features, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	plan.ID = uuid.New()
	plan.CreatedAt = now
	plan.UpdatedAt = now
	plan.IsActive = true

	if plan.Currency == "" {
		plan.Currency = "USD"
	}

	if plan.Features == "" {
		plan.Features = "[]"
	}

	_, err := r.db.ExecContext(ctx, query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.Currency,
		plan.BillingCycle,
		plan.MaxUsers,
		plan.MaxStorage,
		plan.Features,
		plan.IsActive,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	return err
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	query := `
		SELECT id, name, description, price, currency, billing_cycle, max_users, max_storage, features, is_active, created_at, updated_at
		FROM plans 
		WHERE id = $1
	`

	plan := &models.Plan{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.Price,
		&plan.Currency,
		&plan.BillingCycle,
		&plan.MaxUsers,
		&plan.MaxStorage,
		&plan.Features,
		&plan.IsActive,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, err
	}

	return plan, nil
}

func (r *planRepository) List(ctx context.Context, limit, offset int) ([]*models.Plan, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM plans`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, name, description, price, currency, billing_cycle, max_users, max_storage, features, is_active, created_at, updated_at
		FROM plans 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var plans []*models.Plan
	for rows.Next() {
		plan := &models.Plan{}
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.Price,
			&plan.Currency,
			&plan.BillingCycle,
			&plan.MaxUsers,
			&plan.MaxStorage,
			&plan.Features,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		plans = append(plans, plan)
	}

	return plans, total, nil
}

func (r *planRepository) ListActive(ctx context.Context) ([]*models.Plan, error) {
	query := `
		SELECT id, name, description, price, currency, billing_cycle, max_users, max_storage, features, is_active, created_at, updated_at
		FROM plans 
		WHERE is_active = true
		ORDER BY price ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*models.Plan
	for rows.Next() {
		plan := &models.Plan{}
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.Price,
			&plan.Currency,
			&plan.BillingCycle,
			&plan.MaxUsers,
			&plan.MaxStorage,
			&plan.Features,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *planRepository) Update(ctx context.Context, plan *models.Plan) error {
	query := `
		UPDATE plans 
		SET name = $2, description = $3, price = $4, currency = $5, 
			billing_cycle = $6, max_users = $7, max_storage = $8, 
			features = $9, is_active = $10, updated_at = $11
		WHERE id = $1
	`

	plan.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.Currency,
		plan.BillingCycle,
		plan.MaxUsers,
		plan.MaxStorage,
		plan.Features,
		plan.IsActive,
		plan.UpdatedAt,
	)

	return err
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM plans WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
