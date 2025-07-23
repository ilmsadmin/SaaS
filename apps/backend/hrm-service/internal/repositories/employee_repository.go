package repositories

import (
	"database/sql"
	"fmt"

	"zplus-saas/apps/backend/hrm-service/internal/models"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create employee
func (r *EmployeeRepository) Create(employee *models.Employee) error {
	query := `
		INSERT INTO employees (tenant_id, employee_code, first_name, last_name, email, phone, 
			department_id, position, hire_date, salary, status, manager_id, address, 
			date_of_birth, gender, emergency_name, emergency_phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, employee.TenantID, employee.EmployeeCode, employee.FirstName,
		employee.LastName, employee.Email, employee.Phone, employee.DepartmentID,
		employee.Position, employee.HireDate, employee.Salary, employee.Status,
		employee.ManagerID, employee.Address, employee.DateOfBirth, employee.Gender,
		employee.EmergencyName, employee.EmergencyPhone, employee.IsActive).Scan(
		&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)

	return err
}

// Get employee by ID
func (r *EmployeeRepository) GetByID(tenantID string, id int) (*models.Employee, error) {
	employee := &models.Employee{}
	query := `
		SELECT id, tenant_id, employee_code, first_name, last_name, email, phone,
			department_id, position, hire_date, salary, status, manager_id, address,
			date_of_birth, gender, emergency_name, emergency_phone, is_active, 
			created_at, updated_at
		FROM employees 
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&employee.ID, &employee.TenantID, &employee.EmployeeCode, &employee.FirstName,
		&employee.LastName, &employee.Email, &employee.Phone, &employee.DepartmentID,
		&employee.Position, &employee.HireDate, &employee.Salary, &employee.Status,
		&employee.ManagerID, &employee.Address, &employee.DateOfBirth, &employee.Gender,
		&employee.EmergencyName, &employee.EmergencyPhone, &employee.IsActive,
		&employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

// Get employee by email
func (r *EmployeeRepository) GetByEmail(tenantID, email string) (*models.Employee, error) {
	employee := &models.Employee{}
	query := `
		SELECT id, tenant_id, employee_code, first_name, last_name, email, phone,
			department_id, position, hire_date, salary, status, manager_id, address,
			date_of_birth, gender, emergency_name, emergency_phone, is_active, 
			created_at, updated_at
		FROM employees 
		WHERE tenant_id = $1 AND email = $2 AND is_active = true`

	err := r.db.QueryRow(query, tenantID, email).Scan(
		&employee.ID, &employee.TenantID, &employee.EmployeeCode, &employee.FirstName,
		&employee.LastName, &employee.Email, &employee.Phone, &employee.DepartmentID,
		&employee.Position, &employee.HireDate, &employee.Salary, &employee.Status,
		&employee.ManagerID, &employee.Address, &employee.DateOfBirth, &employee.Gender,
		&employee.EmergencyName, &employee.EmergencyPhone, &employee.IsActive,
		&employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

// Get all employees with filters
func (r *EmployeeRepository) GetAll(tenantID string, departmentID *int, status string, limit, offset int) ([]models.Employee, int, error) {
	var employees []models.Employee
	var totalCount int

	// Build query conditions
	whereClause := "WHERE tenant_id = $1 AND is_active = true"
	params := []interface{}{tenantID}
	paramCount := 1

	if departmentID != nil {
		paramCount++
		whereClause += fmt.Sprintf(" AND department_id = $%d", paramCount)
		params = append(params, *departmentID)
	}

	if status != "" {
		paramCount++
		whereClause += fmt.Sprintf(" AND status = $%d", paramCount)
		params = append(params, status)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM employees %s", whereClause)
	err := r.db.QueryRow(countQuery, params...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, tenant_id, employee_code, first_name, last_name, email, phone,
			department_id, position, hire_date, salary, status, manager_id, address,
			date_of_birth, gender, emergency_name, emergency_phone, is_active, 
			created_at, updated_at
		FROM employees 
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, paramCount+1, paramCount+2)

	params = append(params, limit, offset)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.ID, &employee.TenantID, &employee.EmployeeCode, &employee.FirstName,
			&employee.LastName, &employee.Email, &employee.Phone, &employee.DepartmentID,
			&employee.Position, &employee.HireDate, &employee.Salary, &employee.Status,
			&employee.ManagerID, &employee.Address, &employee.DateOfBirth, &employee.Gender,
			&employee.EmergencyName, &employee.EmergencyPhone, &employee.IsActive,
			&employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, employee)
	}

	return employees, totalCount, nil
}

// Update employee
func (r *EmployeeRepository) Update(employee *models.Employee) error {
	query := `
		UPDATE employees SET 
			employee_code = $2, first_name = $3, last_name = $4, email = $5, phone = $6,
			department_id = $7, position = $8, hire_date = $9, salary = $10, status = $11,
			manager_id = $12, address = $13, date_of_birth = $14, gender = $15,
			emergency_name = $16, emergency_phone = $17, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $18 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRow(query, employee.TenantID, employee.EmployeeCode, employee.FirstName,
		employee.LastName, employee.Email, employee.Phone, employee.DepartmentID,
		employee.Position, employee.HireDate, employee.Salary, employee.Status,
		employee.ManagerID, employee.Address, employee.DateOfBirth, employee.Gender,
		employee.EmergencyName, employee.EmergencyPhone, employee.ID).Scan(&employee.UpdatedAt)

	return err
}

// Delete employee (soft delete)
func (r *EmployeeRepository) Delete(tenantID string, id int) error {
	query := `UPDATE employees SET is_active = false, updated_at = NOW() 
			  WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// Search employees
func (r *EmployeeRepository) Search(tenantID, searchTerm string, limit, offset int) ([]models.Employee, int, error) {
	var employees []models.Employee
	var totalCount int

	searchPattern := "%" + searchTerm + "%"

	// Count query
	countQuery := `
		SELECT COUNT(*) FROM employees 
		WHERE tenant_id = $1 AND is_active = true 
		AND (first_name ILIKE $2 OR last_name ILIKE $2 OR email ILIKE $2 OR employee_code ILIKE $2)`

	err := r.db.QueryRow(countQuery, tenantID, searchPattern).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT id, tenant_id, employee_code, first_name, last_name, email, phone,
			department_id, position, hire_date, salary, status, manager_id, address,
			date_of_birth, gender, emergency_name, emergency_phone, is_active, 
			created_at, updated_at
		FROM employees 
		WHERE tenant_id = $1 AND is_active = true 
		AND (first_name ILIKE $2 OR last_name ILIKE $2 OR email ILIKE $2 OR employee_code ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, tenantID, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.ID, &employee.TenantID, &employee.EmployeeCode, &employee.FirstName,
			&employee.LastName, &employee.Email, &employee.Phone, &employee.DepartmentID,
			&employee.Position, &employee.HireDate, &employee.Salary, &employee.Status,
			&employee.ManagerID, &employee.Address, &employee.DateOfBirth, &employee.Gender,
			&employee.EmergencyName, &employee.EmergencyPhone, &employee.IsActive,
			&employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, employee)
	}

	return employees, totalCount, nil
}

// Get HRM statistics
func (r *EmployeeRepository) GetHRMStats(tenantID string) (*models.HRMStats, error) {
	stats := &models.HRMStats{}

	// Get employee statistics
	query := `
		SELECT 
			COUNT(*) as total_employees,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_employees,
			COUNT(CASE WHEN hire_date >= DATE_TRUNC('month', CURRENT_DATE) THEN 1 END) as new_hires_this_month
		FROM employees 
		WHERE tenant_id = $1 AND is_active = true`

	err := r.db.QueryRow(query, tenantID).Scan(&stats.TotalEmployees, &stats.ActiveEmployees, &stats.NewHiresThisMonth)
	if err != nil {
		return nil, err
	}

	// Calculate turnover rate (terminated employees in the last 12 months / average employees)
	query2 := `
		SELECT COALESCE(
			(COUNT(CASE WHEN status = 'terminated' AND updated_at >= CURRENT_DATE - INTERVAL '12 months' THEN 1 END)::float / 
			 NULLIF(COUNT(*)::float, 0)) * 100, 0)
		FROM employees 
		WHERE tenant_id = $1`

	err = r.db.QueryRow(query2, tenantID).Scan(&stats.TurnoverRate)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
