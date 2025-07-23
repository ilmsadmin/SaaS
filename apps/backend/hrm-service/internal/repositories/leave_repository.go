package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"zplus-saas/apps/backend/hrm-service/internal/models"
)

type LeaveRepository struct {
	db *sql.DB
}

func NewLeaveRepository(db *sql.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// Create leave
func (r *LeaveRepository) Create(leave *models.Leave) error {
	query := `
		INSERT INTO leaves (tenant_id, employee_id, leave_type, start_date, end_date, days, reason, status, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, leave.TenantID, leave.EmployeeID, leave.LeaveType,
		leave.StartDate, leave.EndDate, leave.Days, leave.Reason, leave.Status,
		leave.IsActive).Scan(&leave.ID, &leave.CreatedAt, &leave.UpdatedAt)

	return err
}

// Get leave by ID
func (r *LeaveRepository) GetByID(tenantID string, id int) (*models.Leave, error) {
	leave := &models.Leave{}
	query := `
		SELECT id, tenant_id, employee_id, leave_type, start_date, end_date, days, reason, status,
			   approved_by, approved_at, comments, is_active, created_at, updated_at
		FROM leaves 
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&leave.ID, &leave.TenantID, &leave.EmployeeID, &leave.LeaveType,
		&leave.StartDate, &leave.EndDate, &leave.Days, &leave.Reason, &leave.Status,
		&leave.ApprovedBy, &leave.ApprovedAt, &leave.Comments, &leave.IsActive,
		&leave.CreatedAt, &leave.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return leave, nil
}

// Get leaves by employee ID
func (r *LeaveRepository) GetByEmployeeID(tenantID string, employeeID int, limit, offset int) ([]models.Leave, int, error) {
	var leaves []models.Leave
	var totalCount int

	// Count query
	countQuery := "SELECT COUNT(*) FROM leaves WHERE tenant_id = $1 AND employee_id = $2 AND is_active = true"
	err := r.db.QueryRow(countQuery, tenantID, employeeID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT id, tenant_id, employee_id, leave_type, start_date, end_date, days, reason, status,
			   approved_by, approved_at, comments, is_active, created_at, updated_at
		FROM leaves 
		WHERE tenant_id = $1 AND employee_id = $2 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, tenantID, employeeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var leave models.Leave
		err := rows.Scan(
			&leave.ID, &leave.TenantID, &leave.EmployeeID, &leave.LeaveType,
			&leave.StartDate, &leave.EndDate, &leave.Days, &leave.Reason, &leave.Status,
			&leave.ApprovedBy, &leave.ApprovedAt, &leave.Comments, &leave.IsActive,
			&leave.CreatedAt, &leave.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, leave)
	}

	return leaves, totalCount, nil
}

// Get all leaves with filters
func (r *LeaveRepository) GetAll(tenantID string, employeeID *int, status string, leaveType string, limit, offset int) ([]models.Leave, int, error) {
	var leaves []models.Leave
	var totalCount int

	// Build query conditions
	whereClause := "WHERE tenant_id = $1 AND is_active = true"
	params := []interface{}{tenantID}
	paramCount := 1

	if employeeID != nil {
		paramCount++
		whereClause += fmt.Sprintf(" AND employee_id = $%d", paramCount)
		params = append(params, *employeeID)
	}

	if status != "" {
		paramCount++
		whereClause += fmt.Sprintf(" AND status = $%d", paramCount)
		params = append(params, status)
	}

	if leaveType != "" {
		paramCount++
		whereClause += fmt.Sprintf(" AND leave_type = $%d", paramCount)
		params = append(params, leaveType)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM leaves %s", whereClause)
	err := r.db.QueryRow(countQuery, params...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, tenant_id, employee_id, leave_type, start_date, end_date, days, reason, status,
			   approved_by, approved_at, comments, is_active, created_at, updated_at
		FROM leaves 
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
		var leave models.Leave
		err := rows.Scan(
			&leave.ID, &leave.TenantID, &leave.EmployeeID, &leave.LeaveType,
			&leave.StartDate, &leave.EndDate, &leave.Days, &leave.Reason, &leave.Status,
			&leave.ApprovedBy, &leave.ApprovedAt, &leave.Comments, &leave.IsActive,
			&leave.CreatedAt, &leave.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, leave)
	}

	return leaves, totalCount, nil
}

// Update leave
func (r *LeaveRepository) Update(leave *models.Leave) error {
	query := `
		UPDATE leaves SET 
			leave_type = $2, start_date = $3, end_date = $4, days = $5, reason = $6, 
			status = $7, approved_by = $8, approved_at = $9, comments = $10, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $11 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRow(query, leave.TenantID, leave.LeaveType, leave.StartDate,
		leave.EndDate, leave.Days, leave.Reason, leave.Status, leave.ApprovedBy,
		leave.ApprovedAt, leave.Comments, leave.ID).Scan(&leave.UpdatedAt)

	return err
}

// Approve leave
func (r *LeaveRepository) Approve(tenantID string, leaveID, approverID int, comments string) error {
	now := time.Now()
	query := `
		UPDATE leaves SET 
			status = 'approved', approved_by = $3, approved_at = $4, comments = $5, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	_, err := r.db.Exec(query, tenantID, leaveID, approverID, now, comments)
	return err
}

// Reject leave
func (r *LeaveRepository) Reject(tenantID string, leaveID, approverID int, comments string) error {
	now := time.Now()
	query := `
		UPDATE leaves SET 
			status = 'rejected', approved_by = $3, approved_at = $4, comments = $5, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	_, err := r.db.Exec(query, tenantID, leaveID, approverID, now, comments)
	return err
}

// Delete leave (soft delete)
func (r *LeaveRepository) Delete(tenantID string, id int) error {
	query := `UPDATE leaves SET is_active = false, updated_at = NOW() 
			  WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// Get pending leaves count
func (r *LeaveRepository) GetPendingCount(tenantID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM leaves WHERE tenant_id = $1 AND status = 'pending' AND is_active = true`
	err := r.db.QueryRow(query, tenantID).Scan(&count)
	return count, err
}

// Get leave balance for employee
func (r *LeaveRepository) GetLeaveBalance(tenantID string, employeeID int, leaveType string, year int) (int, error) {
	var usedDays int

	// Calculate used days for the year
	query := `
		SELECT COALESCE(SUM(days), 0) 
		FROM leaves 
		WHERE tenant_id = $1 AND employee_id = $2 AND leave_type = $3 
		AND EXTRACT(YEAR FROM start_date) = $4 
		AND status = 'approved' AND is_active = true`

	err := r.db.QueryRow(query, tenantID, employeeID, leaveType, year).Scan(&usedDays)
	if err != nil {
		return 0, err
	}

	// Annual leave entitlement (should be configurable)
	entitlement := 21 // Default annual leave days
	if leaveType == "sick" {
		entitlement = 10 // Default sick leave days
	}

	balance := entitlement - usedDays
	if balance < 0 {
		balance = 0
	}

	return balance, nil
}
