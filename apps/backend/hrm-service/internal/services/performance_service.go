package services

import (
	"fmt"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/repositories"
)

type PerformanceService struct {
	repo *repositories.PerformanceRepository
}

func NewPerformanceService(repo *repositories.PerformanceRepository) *PerformanceService {
	return &PerformanceService{repo: repo}
}

// Create performance review
func (s *PerformanceService) CreatePerformance(tenantID string, employeeID, reviewerID int, period, reviewType string, overallRating float64, goals, achievements, strengths, areas, comments string) (*models.Performance, error) {
	// Validate input
	if err := s.validatePerformanceData(employeeID, reviewerID, period, reviewType, overallRating); err != nil {
		return nil, err
	}

	performance := &models.Performance{
		TenantID:      tenantID,
		EmployeeID:    employeeID,
		ReviewerID:    reviewerID,
		Period:        period,
		ReviewType:    reviewType,
		OverallRating: overallRating,
		Goals:         goals,
		Achievements:  achievements,
		Strengths:     strengths,
		Areas:         areas,
		Comments:      comments,
		Status:        models.PerformanceStatusDraft,
		IsActive:      true,
	}

	err := s.repo.Create(performance)
	if err != nil {
		return nil, fmt.Errorf("failed to create performance review: %w", err)
	}

	return performance, nil
}

// Get performance review by ID
func (s *PerformanceService) GetPerformance(tenantID string, id int) (*models.Performance, error) {
	performance, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("performance review not found: %w", err)
	}

	return performance, nil
}

// Get performance reviews by employee ID
func (s *PerformanceService) GetPerformanceByEmployee(tenantID string, employeeID int, page, limit int) ([]models.Performance, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	reviews, total, err := s.repo.GetByEmployeeID(tenantID, employeeID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employee performance reviews: %w", err)
	}

	return reviews, total, nil
}

// Get all performance reviews with filters
func (s *PerformanceService) GetAllPerformance(tenantID string, employeeID, reviewerID *int, reviewType, status string, page, limit int) ([]models.Performance, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	reviews, total, err := s.repo.GetAll(tenantID, employeeID, reviewerID, reviewType, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get performance reviews: %w", err)
	}

	return reviews, total, nil
}

// Update performance review
func (s *PerformanceService) UpdatePerformance(tenantID string, id int, period, reviewType string, overallRating float64, goals, achievements, strengths, areas, comments, status string) (*models.Performance, error) {
	// Get existing performance review
	performance, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("performance review not found: %w", err)
	}

	// Validate input
	if err := s.validatePerformanceUpdate(period, reviewType, overallRating, status); err != nil {
		return nil, err
	}

	// Update fields
	performance.Period = period
	performance.ReviewType = reviewType
	performance.OverallRating = overallRating
	performance.Goals = goals
	performance.Achievements = achievements
	performance.Strengths = strengths
	performance.Areas = areas
	performance.Comments = comments
	performance.Status = status

	err = s.repo.Update(performance)
	if err != nil {
		return nil, fmt.Errorf("failed to update performance review: %w", err)
	}

	return performance, nil
}

// Delete performance review
func (s *PerformanceService) DeletePerformance(tenantID string, id int) error {
	// Check if performance review exists
	performance, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("performance review not found: %w", err)
	}

	// Only allow deletion for draft reviews
	if performance.Status != models.PerformanceStatusDraft {
		return fmt.Errorf("can only delete draft performance reviews")
	}

	err = s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete performance review: %w", err)
	}

	return nil
}

// Submit performance review
func (s *PerformanceService) SubmitPerformance(tenantID string, id int) error {
	// Get existing performance review
	performance, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("performance review not found: %w", err)
	}

	if performance.Status != models.PerformanceStatusDraft {
		return fmt.Errorf("can only submit draft performance reviews")
	}

	// Update status to submitted
	performance.Status = models.PerformanceStatusSubmitted

	err = s.repo.Update(performance)
	if err != nil {
		return fmt.Errorf("failed to submit performance review: %w", err)
	}

	return nil
}

// Complete performance review
func (s *PerformanceService) CompletePerformance(tenantID string, id int) error {
	// Get existing performance review
	performance, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("performance review not found: %w", err)
	}

	if performance.Status != models.PerformanceStatusSubmitted {
		return fmt.Errorf("can only complete submitted performance reviews")
	}

	// Update status to completed
	performance.Status = models.PerformanceStatusCompleted

	err = s.repo.Update(performance)
	if err != nil {
		return fmt.Errorf("failed to complete performance review: %w", err)
	}

	return nil
}

// Get average performance rating
func (s *PerformanceService) GetAveragePerformanceRating(tenantID string) (float64, error) {
	avgRating, err := s.repo.GetAverageRating(tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average performance rating: %w", err)
	}

	return avgRating, nil
}

// Get performance statistics by department
func (s *PerformanceService) GetPerformanceStatsByDepartment(tenantID string) ([]map[string]interface{}, error) {
	stats, err := s.repo.GetStatsByDepartment(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance statistics by department: %w", err)
	}

	return stats, nil
}

// Validate performance data
func (s *PerformanceService) validatePerformanceData(employeeID, reviewerID int, period, reviewType string, overallRating float64) error {
	if employeeID == 0 {
		return fmt.Errorf("employee ID is required")
	}

	if reviewerID == 0 {
		return fmt.Errorf("reviewer ID is required")
	}

	if period == "" {
		return fmt.Errorf("period is required")
	}

	if reviewType == "" {
		return fmt.Errorf("review type is required")
	}

	// Validate review type
	validTypes := []string{models.ReviewTypeQuarterly, models.ReviewTypeAnnual, models.ReviewTypeProbation}
	isValidType := false
	for _, validType := range validTypes {
		if reviewType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid review type")
	}

	if overallRating < 1 || overallRating > 5 {
		return fmt.Errorf("overall rating must be between 1 and 5")
	}

	return nil
}

// Validate performance update
func (s *PerformanceService) validatePerformanceUpdate(period, reviewType string, overallRating float64, status string) error {
	if period == "" {
		return fmt.Errorf("period is required")
	}

	if reviewType == "" {
		return fmt.Errorf("review type is required")
	}

	// Validate review type
	validTypes := []string{models.ReviewTypeQuarterly, models.ReviewTypeAnnual, models.ReviewTypeProbation}
	isValidType := false
	for _, validType := range validTypes {
		if reviewType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid review type")
	}

	if overallRating < 1 || overallRating > 5 {
		return fmt.Errorf("overall rating must be between 1 and 5")
	}

	// Validate status
	if status != "" {
		validStatuses := []string{models.PerformanceStatusDraft, models.PerformanceStatusSubmitted, models.PerformanceStatusCompleted}
		isValidStatus := false
		for _, validStatus := range validStatuses {
			if status == validStatus {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			return fmt.Errorf("invalid status")
		}
	}

	return nil
}
