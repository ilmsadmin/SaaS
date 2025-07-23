package services

import (
	"fmt"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/repositories"
)

type DepartmentService struct {
	repo *repositories.DepartmentRepository
}

func NewDepartmentService(repo *repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// Create department
func (s *DepartmentService) CreateDepartment(tenantID string, name, description, location string, managerID *int, budget float64) (*models.Department, error) {
	department := &models.Department{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		ManagerID:   managerID,
		Budget:      budget,
		Location:    location,
		IsActive:    true,
	}

	err := s.repo.Create(department)
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return department, nil
}

// Get department by ID
func (s *DepartmentService) GetDepartment(tenantID string, id int) (*models.Department, error) {
	department, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	return department, nil
}

// Get all departments with pagination
func (s *DepartmentService) GetAllDepartments(tenantID string, page, limit int) ([]models.Department, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	departments, total, err := s.repo.GetAll(tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get departments: %w", err)
	}

	return departments, total, nil
}

// Update department
func (s *DepartmentService) UpdateDepartment(tenantID string, id int, name, description, location string, managerID *int, budget float64) (*models.Department, error) {
	// Get existing department
	department, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	// Update fields
	department.Name = name
	department.Description = description
	department.ManagerID = managerID
	department.Budget = budget
	department.Location = location

	err = s.repo.Update(department)
	if err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return department, nil
}

// Delete department
func (s *DepartmentService) DeleteDepartment(tenantID string, id int) error {
	// Check if department exists
	_, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}

	// TODO: Check if there are employees in this department before deletion
	// For now, we'll allow deletion but should add this validation

	err = s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}

// Get departments with employee count
func (s *DepartmentService) GetDepartmentsWithEmployeeCount(tenantID string) ([]map[string]interface{}, error) {
	departments, err := s.repo.GetWithEmployeeCount(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments with employee count: %w", err)
	}

	return departments, nil
}
