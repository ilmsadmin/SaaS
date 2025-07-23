package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"zplus-saas/apps/backend/crm-service/internal/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create creates a new customer
func (r *CustomerRepository) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (tenant_id, name, email, phone, company, address, city, state, country, zip_code, status, source, tags, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`

	now := time.Now()
	tagsStr := strings.Join(customer.Tags, ",")

	err := r.db.QueryRow(
		query,
		customer.TenantID, customer.Name, customer.Email, customer.Phone,
		customer.Company, customer.Address, customer.City, customer.State,
		customer.Country, customer.ZipCode, customer.Status, customer.Source,
		tagsStr, customer.Notes, now, now,
	).Scan(&customer.ID)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	customer.CreatedAt = now
	customer.UpdatedAt = now

	return nil
}

// GetByID gets a customer by ID
func (r *CustomerRepository) GetByID(tenantID string, id int) (*models.Customer, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, address, city, state, country, zip_code, status, source, tags, notes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1 AND id = $2
	`

	customer := &models.Customer{}
	var tagsStr string

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&customer.ID, &customer.TenantID, &customer.Name, &customer.Email,
		&customer.Phone, &customer.Company, &customer.Address, &customer.City,
		&customer.State, &customer.Country, &customer.ZipCode, &customer.Status,
		&customer.Source, &tagsStr, &customer.Notes, &customer.CreatedAt, &customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if tagsStr != "" {
		customer.Tags = strings.Split(tagsStr, ",")
	}

	return customer, nil
}

// GetAll gets all customers for a tenant
func (r *CustomerRepository) GetAll(tenantID string, limit, offset int) ([]*models.Customer, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, company, address, city, state, country, zip_code, status, source, tags, notes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	defer rows.Close()

	var customers []*models.Customer

	for rows.Next() {
		customer := &models.Customer{}
		var tagsStr string

		err := rows.Scan(
			&customer.ID, &customer.TenantID, &customer.Name, &customer.Email,
			&customer.Phone, &customer.Company, &customer.Address, &customer.City,
			&customer.State, &customer.Country, &customer.ZipCode, &customer.Status,
			&customer.Source, &tagsStr, &customer.Notes, &customer.CreatedAt, &customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		if tagsStr != "" {
			customer.Tags = strings.Split(tagsStr, ",")
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

// Update updates a customer
func (r *CustomerRepository) Update(customer *models.Customer) error {
	query := `
		UPDATE customers
		SET name = $3, email = $4, phone = $5, company = $6, address = $7, city = $8, state = $9, country = $10, zip_code = $11, status = $12, source = $13, tags = $14, notes = $15, updated_at = $16
		WHERE tenant_id = $1 AND id = $2
	`

	tagsStr := strings.Join(customer.Tags, ",")
	customer.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		customer.TenantID, customer.ID, customer.Name, customer.Email,
		customer.Phone, customer.Company, customer.Address, customer.City,
		customer.State, customer.Country, customer.ZipCode, customer.Status,
		customer.Source, tagsStr, customer.Notes, customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Delete deletes a customer
func (r *CustomerRepository) Delete(tenantID string, id int) error {
	query := `DELETE FROM customers WHERE tenant_id = $1 AND id = $2`

	result, err := r.db.Exec(query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Search searches customers by name, email, or company
func (r *CustomerRepository) Search(tenantID, query string, limit, offset int) ([]*models.Customer, error) {
	searchQuery := `
		SELECT id, tenant_id, name, email, phone, company, address, city, state, country, zip_code, status, source, tags, notes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1 AND (
			name ILIKE '%' || $2 || '%' OR
			email ILIKE '%' || $2 || '%' OR
			company ILIKE '%' || $2 || '%'
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(searchQuery, tenantID, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	var customers []*models.Customer

	for rows.Next() {
		customer := &models.Customer{}
		var tagsStr string

		err := rows.Scan(
			&customer.ID, &customer.TenantID, &customer.Name, &customer.Email,
			&customer.Phone, &customer.Company, &customer.Address, &customer.City,
			&customer.State, &customer.Country, &customer.ZipCode, &customer.Status,
			&customer.Source, &tagsStr, &customer.Notes, &customer.CreatedAt, &customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		if tagsStr != "" {
			customer.Tags = strings.Split(tagsStr, ",")
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

// Count gets total count of customers for a tenant
func (r *CustomerRepository) Count(tenantID string) (int, error) {
	query := `SELECT COUNT(*) FROM customers WHERE tenant_id = $1`

	var count int
	err := r.db.QueryRow(query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}

	return count, nil
}
