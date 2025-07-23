package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/crm-service/internal/models"
)

type OpportunityRepository struct {
	db *sql.DB
}

func NewOpportunityRepository(db *sql.DB) *OpportunityRepository {
	return &OpportunityRepository{db: db}
}

// Create creates a new opportunity
func (r *OpportunityRepository) Create(opportunity *models.Opportunity) error {
	query := `
		INSERT INTO opportunities (tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	now := time.Now()

	err := r.db.QueryRow(
		query,
		opportunity.TenantID, opportunity.CustomerID, opportunity.Name,
		opportunity.Description, opportunity.Value, opportunity.Currency,
		opportunity.Stage, opportunity.Probability, opportunity.Source,
		opportunity.AssignedTo, opportunity.ExpectedDate, now, now,
	).Scan(&opportunity.ID)

	if err != nil {
		return fmt.Errorf("failed to create opportunity: %w", err)
	}

	opportunity.CreatedAt = now
	opportunity.UpdatedAt = now

	return nil
}

// GetByID gets an opportunity by ID
func (r *OpportunityRepository) GetByID(tenantID string, id int) (*models.Opportunity, error) {
	query := `
		SELECT id, tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date, closed_date, created_at, updated_at
		FROM opportunities
		WHERE tenant_id = $1 AND id = $2
	`

	opportunity := &models.Opportunity{}

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&opportunity.ID, &opportunity.TenantID, &opportunity.CustomerID,
		&opportunity.Name, &opportunity.Description, &opportunity.Value,
		&opportunity.Currency, &opportunity.Stage, &opportunity.Probability,
		&opportunity.Source, &opportunity.AssignedTo, &opportunity.ExpectedDate,
		&opportunity.ClosedDate, &opportunity.CreatedAt, &opportunity.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("opportunity not found")
		}
		return nil, fmt.Errorf("failed to get opportunity: %w", err)
	}

	return opportunity, nil
}

// GetAll gets all opportunities for a tenant
func (r *OpportunityRepository) GetAll(tenantID string, limit, offset int) ([]*models.Opportunity, error) {
	query := `
		SELECT id, tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date, closed_date, created_at, updated_at
		FROM opportunities
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunities: %w", err)
	}
	defer rows.Close()

	var opportunities []*models.Opportunity

	for rows.Next() {
		opportunity := &models.Opportunity{}

		err := rows.Scan(
			&opportunity.ID, &opportunity.TenantID, &opportunity.CustomerID,
			&opportunity.Name, &opportunity.Description, &opportunity.Value,
			&opportunity.Currency, &opportunity.Stage, &opportunity.Probability,
			&opportunity.Source, &opportunity.AssignedTo, &opportunity.ExpectedDate,
			&opportunity.ClosedDate, &opportunity.CreatedAt, &opportunity.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan opportunity: %w", err)
		}

		opportunities = append(opportunities, opportunity)
	}

	return opportunities, nil
}

// Update updates an opportunity
func (r *OpportunityRepository) Update(opportunity *models.Opportunity) error {
	query := `
		UPDATE opportunities
		SET name = $3, description = $4, value = $5, currency = $6, stage = $7, probability = $8, source = $9, assigned_to = $10, expected_date = $11, updated_at = $12
		WHERE tenant_id = $1 AND id = $2
	`

	opportunity.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		opportunity.TenantID, opportunity.ID, opportunity.Name,
		opportunity.Description, opportunity.Value, opportunity.Currency,
		opportunity.Stage, opportunity.Probability, opportunity.Source,
		opportunity.AssignedTo, opportunity.ExpectedDate, opportunity.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("opportunity not found")
	}

	return nil
}

// Delete deletes an opportunity
func (r *OpportunityRepository) Delete(tenantID string, id int) error {
	query := `DELETE FROM opportunities WHERE tenant_id = $1 AND id = $2`

	result, err := r.db.Exec(query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete opportunity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("opportunity not found")
	}

	return nil
}

// CloseWon marks an opportunity as closed won
func (r *OpportunityRepository) CloseWon(tenantID string, id int) error {
	query := `
		UPDATE opportunities
		SET stage = 'closed-won', probability = 100, closed_date = $3, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`

	now := time.Now()

	result, err := r.db.Exec(query, tenantID, id, now)
	if err != nil {
		return fmt.Errorf("failed to close opportunity as won: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("opportunity not found")
	}

	return nil
}

// CloseLost marks an opportunity as closed lost
func (r *OpportunityRepository) CloseLost(tenantID string, id int) error {
	query := `
		UPDATE opportunities
		SET stage = 'closed-lost', probability = 0, closed_date = $3, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`

	now := time.Now()

	result, err := r.db.Exec(query, tenantID, id, now)
	if err != nil {
		return fmt.Errorf("failed to close opportunity as lost: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("opportunity not found")
	}

	return nil
}

// GetByStage gets opportunities by stage
func (r *OpportunityRepository) GetByStage(tenantID, stage string, limit, offset int) ([]*models.Opportunity, error) {
	query := `
		SELECT id, tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date, closed_date, created_at, updated_at
		FROM opportunities
		WHERE tenant_id = $1 AND stage = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, tenantID, stage, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunities by stage: %w", err)
	}
	defer rows.Close()

	var opportunities []*models.Opportunity

	for rows.Next() {
		opportunity := &models.Opportunity{}

		err := rows.Scan(
			&opportunity.ID, &opportunity.TenantID, &opportunity.CustomerID,
			&opportunity.Name, &opportunity.Description, &opportunity.Value,
			&opportunity.Currency, &opportunity.Stage, &opportunity.Probability,
			&opportunity.Source, &opportunity.AssignedTo, &opportunity.ExpectedDate,
			&opportunity.ClosedDate, &opportunity.CreatedAt, &opportunity.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan opportunity: %w", err)
		}

		opportunities = append(opportunities, opportunity)
	}

	return opportunities, nil
}

// GetByCustomer gets opportunities for a customer
func (r *OpportunityRepository) GetByCustomer(tenantID string, customerID, limit, offset int) ([]*models.Opportunity, error) {
	query := `
		SELECT id, tenant_id, customer_id, name, description, value, currency, stage, probability, source, assigned_to, expected_date, closed_date, created_at, updated_at
		FROM opportunities
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, tenantID, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunities by customer: %w", err)
	}
	defer rows.Close()

	var opportunities []*models.Opportunity

	for rows.Next() {
		opportunity := &models.Opportunity{}

		err := rows.Scan(
			&opportunity.ID, &opportunity.TenantID, &opportunity.CustomerID,
			&opportunity.Name, &opportunity.Description, &opportunity.Value,
			&opportunity.Currency, &opportunity.Stage, &opportunity.Probability,
			&opportunity.Source, &opportunity.AssignedTo, &opportunity.ExpectedDate,
			&opportunity.ClosedDate, &opportunity.CreatedAt, &opportunity.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan opportunity: %w", err)
		}

		opportunities = append(opportunities, opportunity)
	}

	return opportunities, nil
}

// Count gets total count of opportunities for a tenant
func (r *OpportunityRepository) Count(tenantID string) (int, error) {
	query := `SELECT COUNT(*) FROM opportunities WHERE tenant_id = $1`

	var count int
	err := r.db.QueryRow(query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count opportunities: %w", err)
	}

	return count, nil
}

// GetTotalValue gets the total value of opportunities for a tenant
func (r *OpportunityRepository) GetTotalValue(tenantID string) (float64, error) {
	query := `SELECT COALESCE(SUM(value), 0) FROM opportunities WHERE tenant_id = $1 AND stage NOT IN ('closed-lost')`

	var total float64
	err := r.db.QueryRow(query, tenantID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total opportunities value: %w", err)
	}

	return total, nil
}
