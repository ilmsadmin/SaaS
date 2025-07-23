package services

import (
	"fmt"
	"time"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/repositories"
)

type LeaveService struct {
	repo *repositories.LeaveRepository
}

func NewLeaveService(repo *repositories.LeaveRepository) *LeaveService {
	return &LeaveService{repo: repo}
}

// Create leave request
func (s *LeaveService) CreateLeave(tenantID string, req *models.LeaveRequest) (*models.Leave, error) {
	// Validate request
	if err := s.validateLeaveRequest(req); err != nil {
		return nil, err
	}

	// Calculate days
	days := s.calculateLeaveDays(req.StartDate, req.EndDate)

	leave := &models.Leave{
		TenantID:   tenantID,
		EmployeeID: req.EmployeeID,
		LeaveType:  req.LeaveType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Days:       days,
		Reason:     req.Reason,
		Status:     models.LeaveStatusPending,
		IsActive:   true,
	}

	err := s.repo.Create(leave)
	if err != nil {
		return nil, fmt.Errorf("failed to create leave request: %w", err)
	}

	return leave, nil
}

// Get leave by ID
func (s *LeaveService) GetLeave(tenantID string, id int) (*models.Leave, error) {
	leave, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("leave not found: %w", err)
	}

	return leave, nil
}

// Get leaves by employee ID
func (s *LeaveService) GetLeavesByEmployee(tenantID string, employeeID int, page, limit int) ([]models.Leave, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	leaves, total, err := s.repo.GetByEmployeeID(tenantID, employeeID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employee leaves: %w", err)
	}

	return leaves, total, nil
}

// Get all leaves with filters
func (s *LeaveService) GetAllLeaves(tenantID string, employeeID *int, status, leaveType string, page, limit int) ([]models.Leave, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	leaves, total, err := s.repo.GetAll(tenantID, employeeID, status, leaveType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get leaves: %w", err)
	}

	return leaves, total, nil
}

// Update leave
func (s *LeaveService) UpdateLeave(tenantID string, id int, req *models.LeaveRequest) (*models.Leave, error) {
	// Get existing leave
	leave, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("leave not found: %w", err)
	}

	// Only allow updates for pending leaves
	if leave.Status != models.LeaveStatusPending {
		return nil, fmt.Errorf("can only update pending leave requests")
	}

	// Validate request
	if err := s.validateLeaveRequest(req); err != nil {
		return nil, err
	}

	// Calculate days
	days := s.calculateLeaveDays(req.StartDate, req.EndDate)

	// Update fields
	leave.LeaveType = req.LeaveType
	leave.StartDate = req.StartDate
	leave.EndDate = req.EndDate
	leave.Days = days
	leave.Reason = req.Reason

	err = s.repo.Update(leave)
	if err != nil {
		return nil, fmt.Errorf("failed to update leave: %w", err)
	}

	return leave, nil
}

// Approve leave
func (s *LeaveService) ApproveLeave(tenantID string, leaveID, approverID int, comments string) error {
	// Check if leave exists and is pending
	leave, err := s.repo.GetByID(tenantID, leaveID)
	if err != nil {
		return fmt.Errorf("leave not found: %w", err)
	}

	if leave.Status != models.LeaveStatusPending {
		return fmt.Errorf("leave is not pending approval")
	}

	err = s.repo.Approve(tenantID, leaveID, approverID, comments)
	if err != nil {
		return fmt.Errorf("failed to approve leave: %w", err)
	}

	return nil
}

// Reject leave
func (s *LeaveService) RejectLeave(tenantID string, leaveID, approverID int, comments string) error {
	// Check if leave exists and is pending
	leave, err := s.repo.GetByID(tenantID, leaveID)
	if err != nil {
		return fmt.Errorf("leave not found: %w", err)
	}

	if leave.Status != models.LeaveStatusPending {
		return fmt.Errorf("leave is not pending approval")
	}

	err = s.repo.Reject(tenantID, leaveID, approverID, comments)
	if err != nil {
		return fmt.Errorf("failed to reject leave: %w", err)
	}

	return nil
}

// Delete leave
func (s *LeaveService) DeleteLeave(tenantID string, id int) error {
	// Check if leave exists
	leave, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("leave not found: %w", err)
	}

	// Only allow deletion for pending leaves
	if leave.Status != models.LeaveStatusPending {
		return fmt.Errorf("can only delete pending leave requests")
	}

	err = s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete leave: %w", err)
	}

	return nil
}

// Get pending leaves count
func (s *LeaveService) GetPendingLeavesCount(tenantID string) (int, error) {
	count, err := s.repo.GetPendingCount(tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending leaves count: %w", err)
	}

	return count, nil
}

// Get leave balance
func (s *LeaveService) GetLeaveBalance(tenantID string, employeeID int, leaveType string) (int, error) {
	currentYear := time.Now().Year()
	balance, err := s.repo.GetLeaveBalance(tenantID, employeeID, leaveType, currentYear)
	if err != nil {
		return 0, fmt.Errorf("failed to get leave balance: %w", err)
	}

	return balance, nil
}

// Validate leave request
func (s *LeaveService) validateLeaveRequest(req *models.LeaveRequest) error {
	if req.EmployeeID == 0 {
		return fmt.Errorf("employee ID is required")
	}

	if req.LeaveType == "" {
		return fmt.Errorf("leave type is required")
	}

	// Validate leave type
	validTypes := []string{
		models.LeaveTypeAnnual, models.LeaveTypeSick, models.LeaveTypeMaternity,
		models.LeaveTypePaternity, models.LeaveTypePersonal, models.LeaveTypeEmergency,
	}
	isValidType := false
	for _, validType := range validTypes {
		if req.LeaveType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid leave type")
	}

	if req.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}

	if req.EndDate.IsZero() {
		return fmt.Errorf("end date is required")
	}

	if req.EndDate.Before(req.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	if req.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("leave cannot be requested for past dates")
	}

	if req.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	return nil
}

// Calculate leave days (excluding weekends)
func (s *LeaveService) calculateLeaveDays(startDate, endDate time.Time) int {
	days := 0
	current := startDate

	for current.Before(endDate) || current.Equal(endDate) {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}

	return days
}
