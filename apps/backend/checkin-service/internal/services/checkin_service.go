package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../models"
	"../repositories"
)

type CheckinService struct {
	repo *repositories.CheckinRepository
}

func NewCheckinService(repo *repositories.CheckinRepository) *CheckinService {
	return &CheckinService{repo: repo}
}

// CreateCheckin creates a new checkin record
func (s *CheckinService) CreateCheckin(tenantID string, req *models.CheckinRequest, ipAddress, deviceInfo string) (*models.CheckinRecord, error) {
	// Validate request
	if err := s.validateCheckinRequest(req); err != nil {
		return nil, err
	}

	// Check if employee has already checked in today for the same type
	todayRecords, err := s.repo.GetTodayCheckinRecords(tenantID, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's records: %w", err)
	}

	// Validate checkin sequence
	if err := s.validateCheckinSequence(req.CheckinType, todayRecords); err != nil {
		return nil, err
	}

	// Create checkin record
	record := &models.CheckinRecord{
		TenantID:     tenantID,
		EmployeeID:   req.EmployeeID,
		EmployeeName: fmt.Sprintf("Employee_%d", req.EmployeeID), // TODO: Get actual name from HRM service
		CheckinType:  req.CheckinType,
		Timestamp:    time.Now(),
		Location:     req.Location,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IPAddress:    ipAddress,
		DeviceInfo:   deviceInfo,
		Photo:        req.Photo,
		Notes:        req.Notes,
		Status:       models.StatusApproved, // Auto-approve for now
	}

	if err := s.repo.CreateCheckinRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create checkin record: %w", err)
	}

	// Update attendance summary if this is a checkout
	if req.CheckinType == models.CheckinTypeCheckout {
		if err := s.updateAttendanceSummary(tenantID, req.EmployeeID, todayRecords); err != nil {
			// Log error but don't fail the checkin
			fmt.Printf("Failed to update attendance summary: %v\n", err)
		}
	}

	return record, nil
}

// GetCheckinRecords gets checkin records with filters
func (s *CheckinService) GetCheckinRecords(tenantID string, employeeID *int, checkinType, status string, dateFrom, dateTo time.Time, page, limit int) ([]models.CheckinRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetCheckinRecords(tenantID, employeeID, checkinType, status, dateFrom, dateTo, page, limit)
}

// GetCheckinRecordByID gets a checkin record by ID
func (s *CheckinService) GetCheckinRecordByID(tenantID string, id int) (*models.CheckinRecord, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid checkin record ID")
	}

	return s.repo.GetCheckinRecordByID(tenantID, id)
}

// UpdateCheckinRecord updates a checkin record
func (s *CheckinService) UpdateCheckinRecord(tenantID string, id int, location, photo, notes string) (*models.CheckinRecord, error) {
	record, err := s.repo.GetCheckinRecordByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if location != "" {
		record.Location = location
	}
	if photo != "" {
		record.Photo = photo
	}
	if notes != "" {
		record.Notes = notes
	}

	if err := s.repo.UpdateCheckinRecord(record); err != nil {
		return nil, err
	}

	return record, nil
}

// ApproveCheckin approves a checkin record
func (s *CheckinService) ApproveCheckin(tenantID string, id, approverID int) error {
	record, err := s.repo.GetCheckinRecordByID(tenantID, id)
	if err != nil {
		return err
	}

	if record.Status == models.StatusApproved {
		return fmt.Errorf("checkin record is already approved")
	}

	record.Status = models.StatusApproved
	record.ApprovedBy = &approverID
	now := time.Now()
	record.ApprovedAt = &now

	return s.repo.UpdateCheckinRecord(record)
}

// RejectCheckin rejects a checkin record
func (s *CheckinService) RejectCheckin(tenantID string, id, approverID int, reason string) error {
	record, err := s.repo.GetCheckinRecordByID(tenantID, id)
	if err != nil {
		return err
	}

	if record.Status == models.StatusRejected {
		return fmt.Errorf("checkin record is already rejected")
	}

	record.Status = models.StatusRejected
	record.ApprovedBy = &approverID
	now := time.Now()
	record.ApprovedAt = &now
	record.Notes = reason

	return s.repo.UpdateCheckinRecord(record)
}

// DeleteCheckinRecord deletes a checkin record
func (s *CheckinService) DeleteCheckinRecord(tenantID string, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid checkin record ID")
	}

	return s.repo.DeleteCheckinRecord(tenantID, id)
}

// GetTodayCheckinRecords gets today's checkin records for an employee
func (s *CheckinService) GetTodayCheckinRecords(tenantID string, employeeID int) ([]models.CheckinRecord, error) {
	if employeeID <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}

	return s.repo.GetTodayCheckinRecords(tenantID, employeeID)
}

// GetAttendanceStats gets attendance statistics
func (s *CheckinService) GetAttendanceStats(tenantID string, date time.Time) (*models.AttendanceStats, error) {
	return s.repo.GetAttendanceStats(tenantID, date)
}

// validateCheckinRequest validates checkin request
func (s *CheckinService) validateCheckinRequest(req *models.CheckinRequest) error {
	if req.EmployeeID <= 0 {
		return fmt.Errorf("employee ID is required")
	}

	validTypes := []string{
		models.CheckinTypeCheckin,
		models.CheckinTypeCheckout,
		models.CheckinTypeBreakStart,
		models.CheckinTypeBreakEnd,
	}

	isValidType := false
	for _, validType := range validTypes {
		if req.CheckinType == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid checkin type")
	}

	return nil
}

// validateCheckinSequence validates checkin sequence
func (s *CheckinService) validateCheckinSequence(checkinType string, todayRecords []models.CheckinRecord) error {
	if len(todayRecords) == 0 {
		// First record of the day must be checkin
		if checkinType != models.CheckinTypeCheckin {
			return fmt.Errorf("first action of the day must be checkin")
		}
		return nil
	}

	// Get last record
	lastRecord := todayRecords[len(todayRecords)-1]

	switch checkinType {
	case models.CheckinTypeCheckin:
		// Cannot checkin twice without checkout
		for _, record := range todayRecords {
			if record.CheckinType == models.CheckinTypeCheckin {
				return fmt.Errorf("already checked in today")
			}
		}

	case models.CheckinTypeCheckout:
		// Must have checkin first
		hasCheckin := false
		for _, record := range todayRecords {
			if record.CheckinType == models.CheckinTypeCheckin {
				hasCheckin = true
				break
			}
		}
		if !hasCheckin {
			return fmt.Errorf("must check in first")
		}

		// Cannot checkout twice
		for _, record := range todayRecords {
			if record.CheckinType == models.CheckinTypeCheckout {
				return fmt.Errorf("already checked out today")
			}
		}

	case models.CheckinTypeBreakStart:
		// Must be checked in and not on break
		if lastRecord.CheckinType == models.CheckinTypeBreakStart {
			return fmt.Errorf("already on break")
		}

	case models.CheckinTypeBreakEnd:
		// Must be on break
		if lastRecord.CheckinType != models.CheckinTypeBreakStart {
			return fmt.Errorf("not on break")
		}
	}

	return nil
}

// updateAttendanceSummary updates daily attendance summary
func (s *CheckinService) updateAttendanceSummary(tenantID string, employeeID int, todayRecords []models.CheckinRecord) error {
	var checkinTime, checkoutTime *time.Time
	var breakHours float64

	// Find checkin and checkout times
	for _, record := range todayRecords {
		switch record.CheckinType {
		case models.CheckinTypeCheckin:
			checkinTime = &record.Timestamp
		case models.CheckinTypeCheckout:
			checkoutTime = &record.Timestamp
		}
	}

	// Calculate break hours (simplified)
	breakHours = s.calculateBreakHours(todayRecords)

	// Calculate work hours
	var workHours float64
	if checkinTime != nil && checkoutTime != nil {
		totalHours := checkoutTime.Sub(*checkinTime).Hours()
		workHours = totalHours - breakHours
	}

	// Determine status
	status := models.AttendanceStatusPresent
	if checkinTime == nil {
		status = models.AttendanceStatusAbsent
	} else {
		// Check if late (simplified - assumes 9 AM start)
		expectedStart := time.Date(checkinTime.Year(), checkinTime.Month(), checkinTime.Day(), 9, 0, 0, 0, checkinTime.Location())
		if checkinTime.After(expectedStart.Add(15 * time.Minute)) {
			status = models.AttendanceStatusLate
		}
	}

	// Create summary
	summary := &models.AttendanceSummary{
		TenantID:      tenantID,
		EmployeeID:    employeeID,
		EmployeeName:  fmt.Sprintf("Employee_%d", employeeID),
		Date:          time.Now().Truncate(24 * time.Hour),
		CheckinTime:   checkinTime,
		CheckoutTime:  checkoutTime,
		WorkHours:     workHours,
		BreakHours:    breakHours,
		OvertimeHours: 0, // TODO: Calculate overtime
		Status:        status,
	}

	return s.repo.CreateAttendanceSummary(summary)
}

// calculateBreakHours calculates total break hours from records
func (s *CheckinService) calculateBreakHours(records []models.CheckinRecord) float64 {
	var totalBreakHours float64
	var breakStart *time.Time

	for _, record := range records {
		switch record.CheckinType {
		case models.CheckinTypeBreakStart:
			breakStart = &record.Timestamp
		case models.CheckinTypeBreakEnd:
			if breakStart != nil {
				breakDuration := record.Timestamp.Sub(*breakStart).Hours()
				totalBreakHours += breakDuration
				breakStart = nil
			}
		}
	}

	return totalBreakHours
}
