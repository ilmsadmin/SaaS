package repositories

import (
	"database/sql"
	"fmt"

	"zplus-saas/apps/backend/hrm-service/internal/models"
)

type PerformanceRepository struct {
	db *sql.DB
}

func NewPerformanceRepository(db *sql.DB) *PerformanceRepository {
	return &PerformanceRepository{db: db}
}

// Create performance review
func (r *PerformanceRepository) Create(performance *models.Performance) error {
	query := `
		INSERT INTO performance_reviews (tenant_id, employee_id, reviewer_id, period, review_type, 
			overall_rating, goals, achievements, strengths, areas_for_improvement, comments, status, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, performance.TenantID, performance.EmployeeID, performance.ReviewerID,
		performance.Period, performance.ReviewType, performance.OverallRating, performance.Goals,
		performance.Achievements, performance.Strengths, performance.Areas, performance.Comments,
		performance.Status, performance.IsActive).Scan(&performance.ID, &performance.CreatedAt, &performance.UpdatedAt)

	return err
}

// Get performance review by ID
func (r *PerformanceRepository) GetByID(tenantID string, id int) (*models.Performance, error) {
	performance := &models.Performance{}
	query := `
		SELECT id, tenant_id, employee_id, reviewer_id, period, review_type, overall_rating,
			   goals, achievements, strengths, areas_for_improvement, comments, status, is_active, created_at, updated_at
		FROM performance_reviews 
		WHERE tenant_id = $1 AND id = $2 AND is_active = true`

	err := r.db.QueryRow(query, tenantID, id).Scan(
		&performance.ID, &performance.TenantID, &performance.EmployeeID, &performance.ReviewerID,
		&performance.Period, &performance.ReviewType, &performance.OverallRating, &performance.Goals,
		&performance.Achievements, &performance.Strengths, &performance.Areas, &performance.Comments,
		&performance.Status, &performance.IsActive, &performance.CreatedAt, &performance.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return performance, nil
}

// Get performance reviews by employee ID
func (r *PerformanceRepository) GetByEmployeeID(tenantID string, employeeID int, limit, offset int) ([]models.Performance, int, error) {
	var reviews []models.Performance
	var totalCount int

	// Count query
	countQuery := "SELECT COUNT(*) FROM performance_reviews WHERE tenant_id = $1 AND employee_id = $2 AND is_active = true"
	err := r.db.QueryRow(countQuery, tenantID, employeeID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	query := `
		SELECT id, tenant_id, employee_id, reviewer_id, period, review_type, overall_rating,
			   goals, achievements, strengths, areas_for_improvement, comments, status, is_active, created_at, updated_at
		FROM performance_reviews 
		WHERE tenant_id = $1 AND employee_id = $2 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, tenantID, employeeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var review models.Performance
		err := rows.Scan(
			&review.ID, &review.TenantID, &review.EmployeeID, &review.ReviewerID,
			&review.Period, &review.ReviewType, &review.OverallRating, &review.Goals,
			&review.Achievements, &review.Strengths, &review.Areas, &review.Comments,
			&review.Status, &review.IsActive, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, review)
	}

	return reviews, totalCount, nil
}

// Get all performance reviews with filters
func (r *PerformanceRepository) GetAll(tenantID string, employeeID *int, reviewerID *int, reviewType string, status string, limit, offset int) ([]models.Performance, int, error) {
	var reviews []models.Performance
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

	if reviewerID != nil {
		paramCount++
		whereClause += fmt.Sprintf(" AND reviewer_id = $%d", paramCount)
		params = append(params, *reviewerID)
	}

	if reviewType != "" {
		paramCount++
		whereClause += fmt.Sprintf(" AND review_type = $%d", paramCount)
		params = append(params, reviewType)
	}

	if status != "" {
		paramCount++
		whereClause += fmt.Sprintf(" AND status = $%d", paramCount)
		params = append(params, status)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM performance_reviews %s", whereClause)
	err := r.db.QueryRow(countQuery, params...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, tenant_id, employee_id, reviewer_id, period, review_type, overall_rating,
			   goals, achievements, strengths, areas_for_improvement, comments, status, is_active, created_at, updated_at
		FROM performance_reviews 
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
		var review models.Performance
		err := rows.Scan(
			&review.ID, &review.TenantID, &review.EmployeeID, &review.ReviewerID,
			&review.Period, &review.ReviewType, &review.OverallRating, &review.Goals,
			&review.Achievements, &review.Strengths, &review.Areas, &review.Comments,
			&review.Status, &review.IsActive, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, review)
	}

	return reviews, totalCount, nil
}

// Update performance review
func (r *PerformanceRepository) Update(performance *models.Performance) error {
	query := `
		UPDATE performance_reviews SET 
			period = $2, review_type = $3, overall_rating = $4, goals = $5, achievements = $6,
			strengths = $7, areas_for_improvement = $8, comments = $9, status = $10, updated_at = NOW()
		WHERE tenant_id = $1 AND id = $11 AND is_active = true
		RETURNING updated_at`

	err := r.db.QueryRow(query, performance.TenantID, performance.Period, performance.ReviewType,
		performance.OverallRating, performance.Goals, performance.Achievements, performance.Strengths,
		performance.Areas, performance.Comments, performance.Status, performance.ID).Scan(&performance.UpdatedAt)

	return err
}

// Delete performance review (soft delete)
func (r *PerformanceRepository) Delete(tenantID string, id int) error {
	query := `UPDATE performance_reviews SET is_active = false, updated_at = NOW() 
			  WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// Get average performance rating
func (r *PerformanceRepository) GetAverageRating(tenantID string) (float64, error) {
	var avgRating sql.NullFloat64
	query := `
		SELECT AVG(overall_rating) 
		FROM performance_reviews 
		WHERE tenant_id = $1 AND status = 'completed' AND is_active = true`

	err := r.db.QueryRow(query, tenantID).Scan(&avgRating)
	if err != nil {
		return 0, err
	}

	if !avgRating.Valid {
		return 0, nil
	}

	return avgRating.Float64, nil
}

// Get performance statistics by department
func (r *PerformanceRepository) GetStatsByDepartment(tenantID string) ([]map[string]interface{}, error) {
	query := `
		SELECT d.name as department_name, 
			   AVG(pr.overall_rating) as avg_rating,
			   COUNT(pr.id) as total_reviews
		FROM performance_reviews pr
		JOIN employees e ON pr.employee_id = e.id
		JOIN departments d ON e.department_id = d.id
		WHERE pr.tenant_id = $1 AND pr.status = 'completed' AND pr.is_active = true
		GROUP BY d.id, d.name
		ORDER BY avg_rating DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var departmentName string
		var avgRating sql.NullFloat64
		var totalReviews int

		err := rows.Scan(&departmentName, &avgRating, &totalReviews)
		if err != nil {
			return nil, err
		}

		rating := 0.0
		if avgRating.Valid {
			rating = avgRating.Float64
		}

		stats = append(stats, map[string]interface{}{
			"department_name": departmentName,
			"avg_rating":      rating,
			"total_reviews":   totalReviews,
		})
	}

	return stats, nil
}
