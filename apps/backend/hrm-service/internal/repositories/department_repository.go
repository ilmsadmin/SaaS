package repositories

import (
	"database/sql"

	"zplus-saas/apps/backend/hrm-service/internal/models"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// Create department
func (r *DepartmentRepository) Create(department *models.Department) error {
	query := `
		INSERT INTO departments (tenant_id, name, description, manager_id, budget, location, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, department.TenantID, department.Name, department.Description,
		department.ManagerID, department.Budget, department.Location, department.IsActive).Scan(
		&department.ID, &department.CreatedAt, &department.UpdatedAt)

	return err
}

// Get department by ID
func (r *DepartmentRepository) GetByID(tenantID string, id int) (*models.Department, error) {
	department := &models.Department{}
	query := `
		SELECT id, tenant_id, name, description, manager_id, budget, location, is_active, created_at, updated_at
		FROM departments 
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&department.ID, &department.TenantID, &department.Name, &department.Description,
		&department.ManagerID, &department.Budget, &department.Location, &department.IsActive,
		&department.CreatedAt, &department.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return department, nil
}

// Get all departments
func (r *DepartmentRepository) GetAll(tenantID string, limit, offset int) ([]models.Department, int, error) {
	var departments []models.Department
	var totalCount int

	// Count query
	countQuery := "SELECT COUNT(*) FROM departments WHERE tenant_id = $1 AND is_active = true"
	err := r.db.QueryRow(countQuery, tenantID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT id, tenant_id, name, description, manager_id, budget, location, is_active, created_at, updated_at
		FROM departments 
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var department models.Department
		err := rows.Scan(
			&department.ID, &department.TenantID, &department.Name, &department.Description,
			&department.ManagerID, &department.Budget, &department.Location, &department.IsActive,
			&department.CreatedAt, &department.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		departments = append(departments, department)
	}

	return departments, totalCount, nil
}

// Update department
func (r *DepartmentRepository) Update(department *models.Department) error {
	query := `
		UPDATE departments SET 
			name = $2, description = $3, manager_id = $4, budget = $5, location = $6, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $7 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRow(query, department.TenantID, department.Name, department.Description,
		department.ManagerID, department.Budget, department.Location, department.ID).Scan(&department.UpdatedAt)

	return err
}

// Delete department (soft delete)
func (r *DepartmentRepository) Delete(tenantID string, id int) error {
	query := `UPDATE departments SET is_active = false, updated_at = NOW() 
			  WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// Get departments with employee count
func (r *DepartmentRepository) GetWithEmployeeCount(tenantID string) ([]map[string]interface{}, error) {
	query := `
		SELECT d.id, d.name, d.description, d.manager_id, d.budget, d.location,
			   COUNT(e.id) as employee_count
		FROM departments d
		LEFT JOIN employees e ON d.id = e.department_id AND e.is_active = true AND e.status = 'active'
		WHERE d.tenant_id = $1 AND d.is_active = true
		GROUP BY d.id, d.name, d.description, d.manager_id, d.budget, d.location
		ORDER BY d.name ASC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []map[string]interface{}
	for rows.Next() {
		var id, managerID sql.NullInt64
		var name, description, location string
		var budget float64
		var employeeCount int

		err := rows.Scan(&id, &name, &description, &managerID, &budget, &location, &employeeCount)
		if err != nil {
			return nil, err
		}

		department := map[string]interface{}{
			"id":             id.Int64,
			"name":           name,
			"description":    description,
			"manager_id":     nil,
			"budget":         budget,
			"location":       location,
			"employee_count": employeeCount,
		}

		if managerID.Valid {
			department["manager_id"] = managerID.Int64
		}

		departments = append(departments, department)
	}

	return departments, nil
}
