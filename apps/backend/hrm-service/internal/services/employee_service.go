package services

import (
	"errors"
	"fmt"
	"time"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/repositories"
)

type EmployeeService struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeService(repo *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// Create employee
func (s *EmployeeService) CreateEmployee(tenantID string, req *models.EmployeeRequest) (*models.Employee, error) {
	// Validate employee code uniqueness (this would need another repo method)
	// For now, we'll assume it's handled by database constraints

	employee := &models.Employee{
		TenantID:       tenantID,
		EmployeeCode:   req.EmployeeCode,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DepartmentID:   req.DepartmentID,
		Position:       req.Position,
		HireDate:       req.HireDate,
		Salary:         req.Salary,
		Status:         req.Status,
		ManagerID:      req.ManagerID,
		Address:        req.Address,
		DateOfBirth:    req.DateOfBirth,
		Gender:         req.Gender,
		EmergencyName:  req.EmergencyName,
		EmergencyPhone: req.EmergencyPhone,
		IsActive:       true,
	}

	err := s.repo.Create(employee)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee, nil
}

// Get employee by ID
func (s *EmployeeService) GetEmployee(tenantID string, id int) (*models.Employee, error) {
	employee, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return employee, nil
}

// Get employee by email
func (s *EmployeeService) GetEmployeeByEmail(tenantID, email string) (*models.Employee, error) {
	employee, err := s.repo.GetByEmail(tenantID, email)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return employee, nil
}

// Get all employees with pagination and filters
func (s *EmployeeService) GetAllEmployees(tenantID string, departmentID *int, status string, page, limit int) ([]models.Employee, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	employees, total, err := s.repo.GetAll(tenantID, departmentID, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employees: %w", err)
	}

	return employees, total, nil
}

// Update employee
func (s *EmployeeService) UpdateEmployee(tenantID string, id int, req *models.EmployeeRequest) (*models.Employee, error) {
	// Get existing employee
	employee, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Update fields
	employee.EmployeeCode = req.EmployeeCode
	employee.FirstName = req.FirstName
	employee.LastName = req.LastName
	employee.Email = req.Email
	employee.Phone = req.Phone
	employee.DepartmentID = req.DepartmentID
	employee.Position = req.Position
	employee.HireDate = req.HireDate
	employee.Salary = req.Salary
	employee.Status = req.Status
	employee.ManagerID = req.ManagerID
	employee.Address = req.Address
	employee.DateOfBirth = req.DateOfBirth
	employee.Gender = req.Gender
	employee.EmergencyName = req.EmergencyName
	employee.EmergencyPhone = req.EmergencyPhone

	err = s.repo.Update(employee)
	if err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	return employee, nil
}

// Delete employee
func (s *EmployeeService) DeleteEmployee(tenantID string, id int) error {
	// Check if employee exists
	_, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}

	err = s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}

// Search employees
func (s *EmployeeService) SearchEmployees(tenantID, searchTerm string, page, limit int) ([]models.Employee, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	employees, total, err := s.repo.Search(tenantID, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search employees: %w", err)
	}

	return employees, total, nil
}

// Get HRM statistics
func (s *EmployeeService) GetHRMStatistics(tenantID string) (*models.HRMStats, error) {
	stats, err := s.repo.GetHRMStats(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get HRM statistics: %w", err)
	}

	return stats, nil
}

// Validate employee data
func (s *EmployeeService) validateEmployeeRequest(req *models.EmployeeRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.EmployeeCode == "" {
		return errors.New("employee code is required")
	}
	if req.Position == "" {
		return errors.New("position is required")
	}
	if req.DepartmentID == 0 {
		return errors.New("department ID is required")
	}
	if req.HireDate.IsZero() {
		return errors.New("hire date is required")
	}
	if req.HireDate.After(time.Now()) {
		return errors.New("hire date cannot be in the future")
	}
	if req.Salary < 0 {
		return errors.New("salary cannot be negative")
	}
	if req.Status != "" && req.Status != models.EmployeeStatusActive &&
		req.Status != models.EmployeeStatusInactive && req.Status != models.EmployeeStatusTerminated {
		return errors.New("invalid employee status")
	}

	return nil
}
