package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"../models"
)

type CheckinRepository struct {
	db *sqlx.DB
}

func NewCheckinRepository(db *sqlx.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// CreateCheckinRecord creates a new checkin record
func (r *CheckinRepository) CreateCheckinRecord(record *models.CheckinRecord) error {
	query := `
		INSERT INTO checkin_records (
			tenant_id, employee_id, employee_name, checkin_type, timestamp,
			location, latitude, longitude, ip_address, device_info, photo, notes, status
		) VALUES (
			:tenant_id, :employee_id, :employee_name, :checkin_type, :timestamp,
			:location, :latitude, :longitude, :ip_address, :device_info, :photo, :notes, :status
		) RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, record)
	if err != nil {
		return fmt.Errorf("failed to create checkin record: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)
	}

	return fmt.Errorf("failed to retrieve created record data")
}

// GetCheckinRecords gets checkin records with filters
func (r *CheckinRepository) GetCheckinRecords(tenantID string, employeeID *int, checkinType, status string, dateFrom, dateTo time.Time, page, limit int) ([]models.CheckinRecord, int, error) {
	whereClause := "WHERE tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	if employeeID != nil {
		whereClause += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}

	if checkinType != "" {
		whereClause += fmt.Sprintf(" AND checkin_type = $%d", argIndex)
		args = append(args, checkinType)
		argIndex++
	}

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if !dateFrom.IsZero() {
		whereClause += fmt.Sprintf(" AND timestamp >= $%d", argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if !dateTo.IsZero() {
		whereClause += fmt.Sprintf(" AND timestamp <= $%d", argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM checkin_records " + whereClause
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count checkin records: %w", err)
	}

	if limit > 0 {
		offset := (page - 1) * limit
		whereClause += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, limit, offset)
	} else {
		whereClause += " ORDER BY timestamp DESC"
	}

	// Data query
	query := `
		SELECT id, tenant_id, employee_id, employee_name, checkin_type, timestamp,
			   location, latitude, longitude, ip_address, device_info, photo, notes,
			   status, approved_by, approved_at, created_at, updated_at
		FROM checkin_records ` + whereClause

	var records []models.CheckinRecord
	err = r.db.Select(&records, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get checkin records: %w", err)
	}

	return records, total, nil
}

// GetCheckinRecordByID gets a checkin record by ID
func (r *CheckinRepository) GetCheckinRecordByID(tenantID string, id int) (*models.CheckinRecord, error) {
	query := `
		SELECT id, tenant_id, employee_id, employee_name, checkin_type, timestamp,
			   location, latitude, longitude, ip_address, device_info, photo, notes,
			   status, approved_by, approved_at, created_at, updated_at
		FROM checkin_records
		WHERE tenant_id = $1 AND id = $2
	`

	var record models.CheckinRecord
	err := r.db.Get(&record, query, tenantID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("checkin record not found")
		}
		return nil, fmt.Errorf("failed to get checkin record: %w", err)
	}

	return &record, nil
}

// UpdateCheckinRecord updates a checkin record
func (r *CheckinRepository) UpdateCheckinRecord(record *models.CheckinRecord) error {
	query := `
		UPDATE checkin_records 
		SET employee_name = :employee_name, location = :location, latitude = :latitude,
			longitude = :longitude, photo = :photo, notes = :notes, status = :status,
			approved_by = :approved_by, approved_at = :approved_at, updated_at = NOW()
		WHERE tenant_id = :tenant_id AND id = :id
	`

	result, err := r.db.NamedExec(query, record)
	if err != nil {
		return fmt.Errorf("failed to update checkin record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("checkin record not found")
	}

	return nil
}

// DeleteCheckinRecord deletes a checkin record
func (r *CheckinRepository) DeleteCheckinRecord(tenantID string, id int) error {
	query := "DELETE FROM checkin_records WHERE tenant_id = $1 AND id = $2"

	result, err := r.db.Exec(query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete checkin record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("checkin record not found")
	}

	return nil
}

// GetTodayCheckinRecords gets today's checkin records for an employee
func (r *CheckinRepository) GetTodayCheckinRecords(tenantID string, employeeID int) ([]models.CheckinRecord, error) {
	today := time.Now().Format("2006-01-02")
	query := `
		SELECT id, tenant_id, employee_id, employee_name, checkin_type, timestamp,
			   location, latitude, longitude, ip_address, device_info, photo, notes,
			   status, approved_by, approved_at, created_at, updated_at
		FROM checkin_records
		WHERE tenant_id = $1 AND employee_id = $2 AND DATE(timestamp) = $3
		ORDER BY timestamp ASC
	`

	var records []models.CheckinRecord
	err := r.db.Select(&records, query, tenantID, employeeID, today)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's checkin records: %w", err)
	}

	return records, nil
}

// CreateAttendanceSummary creates attendance summary
func (r *CheckinRepository) CreateAttendanceSummary(summary *models.AttendanceSummary) error {
	query := `
		INSERT INTO attendance_summary (
			tenant_id, employee_id, employee_name, date, checkin_time, checkout_time,
			work_hours, break_hours, overtime_hours, status, notes
		) VALUES (
			:tenant_id, :employee_id, :employee_name, :date, :checkin_time, :checkout_time,
			:work_hours, :break_hours, :overtime_hours, :status, :notes
		) ON CONFLICT (tenant_id, employee_id, date) 
		DO UPDATE SET
			checkin_time = EXCLUDED.checkin_time,
			checkout_time = EXCLUDED.checkout_time,
			work_hours = EXCLUDED.work_hours,
			break_hours = EXCLUDED.break_hours,
			overtime_hours = EXCLUDED.overtime_hours,
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, summary)
	if err != nil {
		return fmt.Errorf("failed to create/update attendance summary: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&summary.ID, &summary.CreatedAt, &summary.UpdatedAt)
	}

	return fmt.Errorf("failed to retrieve summary data")
}

// GetAttendanceStats gets attendance statistics
func (r *CheckinRepository) GetAttendanceStats(tenantID string, date time.Time) (*models.AttendanceStats, error) {
	dateStr := date.Format("2006-01-02")
	
	query := `
		SELECT 
			COUNT(DISTINCT employee_id) as total_employees,
			COUNT(CASE WHEN status IN ('present', 'late') THEN 1 END) as present_today,
			COUNT(CASE WHEN status = 'absent' THEN 1 END) as absent_today,
			COUNT(CASE WHEN status = 'late' THEN 1 END) as late_today,
			COALESCE(AVG(work_hours), 0) as avg_work_hours
		FROM attendance_summary
		WHERE tenant_id = $1 AND date = $2
	`

	var stats models.AttendanceStats
	err := r.db.Get(&stats, query, tenantID, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance stats: %w", err)
	}

	// Calculate attendance rate
	if stats.TotalEmployees > 0 {
		stats.AttendanceRate = float64(stats.PresentToday) / float64(stats.TotalEmployees) * 100
	}

	return &stats, nil
}
