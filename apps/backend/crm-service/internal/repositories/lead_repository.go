package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/crm-service/internal/models"
)

type LeadRepository struct {
	db *sql.DB
}

func NewLeadRepository(db *sql.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

// Create creates a new lead
func (r *LeadRepository) Create(lead *models.Lead) error {
	query := `
		INSERT INTO leads (tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	now := time.Now()

	err := r.db.QueryRow(
		query,
		lead.TenantID, lead.Name, lead.Email, lead.Phone,
		lead.Company, lead.Title, lead.Source, lead.Status,
		lead.Score, lead.AssignedTo, lead.Value, lead.Notes,
		now, now,
	).Scan(&lead.ID)

	if err != nil {
		return fmt.Errorf("failed to create lead: %w", err)
	}

	lead.CreatedAt = now
	lead.UpdatedAt = now

	return nil
}

// GetByID gets a lead by ID
func (r *LeadRepository) GetByID(tenantID string, id int) (*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes, created_at, updated_at, converted_at
		FROM leads
		WHERE tenant_id = $1 AND id = $2
	`

	lead := &models.Lead{}

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&lead.ID, &lead.TenantID, &lead.Name, &lead.Email,
		&lead.Phone, &lead.Company, &lead.Title, &lead.Source,
		&lead.Status, &lead.Score, &lead.AssignedTo, &lead.Value,
		&lead.Notes, &lead.CreatedAt, &lead.UpdatedAt, &lead.ConvertedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead not found")
		}
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}

	return lead, nil
}

// GetAll gets all leads for a tenant
func (r *LeadRepository) GetAll(tenantID string, limit, offset int) ([]*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes, created_at, updated_at, converted_at
		FROM leads
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leads: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead

	for rows.Next() {
		lead := &models.Lead{}

		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.Name, &lead.Email,
			&lead.Phone, &lead.Company, &lead.Title, &lead.Source,
			&lead.Status, &lead.Score, &lead.AssignedTo, &lead.Value,
			&lead.Notes, &lead.CreatedAt, &lead.UpdatedAt, &lead.ConvertedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}

		leads = append(leads, lead)
	}

	return leads, nil
}

// Update updates a lead
func (r *LeadRepository) Update(lead *models.Lead) error {
	query := `
		UPDATE leads
		SET name = $3, email = $4, phone = $5, company = $6, title = $7, source = $8, status = $9, score = $10, assigned_to = $11, value = $12, notes = $13, updated_at = $14
		WHERE tenant_id = $1 AND id = $2
	`

	lead.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		lead.TenantID, lead.ID, lead.Name, lead.Email,
		lead.Phone, lead.Company, lead.Title, lead.Source,
		lead.Status, lead.Score, lead.AssignedTo, lead.Value,
		lead.Notes, lead.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update lead: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

// Delete deletes a lead
func (r *LeadRepository) Delete(tenantID string, id int) error {
	query := `DELETE FROM leads WHERE tenant_id = $1 AND id = $2`

	result, err := r.db.Exec(query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

// ConvertToCustomer marks a lead as converted
func (r *LeadRepository) ConvertToCustomer(tenantID string, id int) error {
	query := `
		UPDATE leads
		SET status = 'converted', converted_at = $3, updated_at = $3
		WHERE tenant_id = $1 AND id = $2
	`

	now := time.Now()

	result, err := r.db.Exec(query, tenantID, id, now)
	if err != nil {
		return fmt.Errorf("failed to convert lead: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

// GetByStatus gets leads by status
func (r *LeadRepository) GetByStatus(tenantID, status string, limit, offset int) ([]*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes, created_at, updated_at, converted_at
		FROM leads
		WHERE tenant_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, tenantID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leads by status: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead

	for rows.Next() {
		lead := &models.Lead{}

		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.Name, &lead.Email,
			&lead.Phone, &lead.Company, &lead.Title, &lead.Source,
			&lead.Status, &lead.Score, &lead.AssignedTo, &lead.Value,
			&lead.Notes, &lead.CreatedAt, &lead.UpdatedAt, &lead.ConvertedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}

		leads = append(leads, lead)
	}

	return leads, nil
}

// GetByAssignedUser gets leads assigned to a specific user
func (r *LeadRepository) GetByAssignedUser(tenantID string, userID, limit, offset int) ([]*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, title, source, status, score, assigned_to, value, notes, created_at, updated_at, converted_at
		FROM leads
		WHERE tenant_id = $1 AND assigned_to = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, tenantID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leads by assigned user: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead

	for rows.Next() {
		lead := &models.Lead{}

		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.Name, &lead.Email,
			&lead.Phone, &lead.Company, &lead.Title, &lead.Source,
			&lead.Status, &lead.Score, &lead.AssignedTo, &lead.Value,
			&lead.Notes, &lead.CreatedAt, &lead.UpdatedAt, &lead.ConvertedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}

		leads = append(leads, lead)
	}

	return leads, nil
}

// Count gets total count of leads for a tenant
func (r *LeadRepository) Count(tenantID string) (int, error) {
	query := `SELECT COUNT(*) FROM leads WHERE tenant_id = $1`

	var count int
	err := r.db.QueryRow(query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count leads: %w", err)
	}

	return count, nil
}
