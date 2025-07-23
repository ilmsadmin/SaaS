package services

import (
	"fmt"
	"strconv"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/repositories"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(tenantID string, req *models.CreateCustomerRequest) (*models.Customer, error) {
	customer := &models.Customer{
		TenantID: tenantID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Company:  req.Company,
		Address:  req.Address,
		City:     req.City,
		State:    req.State,
		Country:  req.Country,
		ZipCode:  req.ZipCode,
		Status:   "active", // Default status
		Source:   req.Source,
		Tags:     req.Tags,
		Notes:    req.Notes,
	}

	err := s.repo.Create(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetCustomer gets a customer by ID
func (s *CustomerService) GetCustomer(tenantID string, id int) (*models.Customer, error) {
	customer, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// GetCustomers gets all customers with pagination
func (s *CustomerService) GetCustomers(tenantID string, page, limit int) ([]*models.Customer, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	customers, err := s.repo.GetAll(tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get customers: %w", err)
	}

	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get customer count: %w", err)
	}

	return customers, total, nil
}

// UpdateCustomer updates a customer
func (s *CustomerService) UpdateCustomer(tenantID string, id int, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	// Get existing customer
	customer, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Email != nil {
		customer.Email = *req.Email
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Company != nil {
		customer.Company = *req.Company
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.City != nil {
		customer.City = *req.City
	}
	if req.State != nil {
		customer.State = *req.State
	}
	if req.Country != nil {
		customer.Country = *req.Country
	}
	if req.ZipCode != nil {
		customer.ZipCode = *req.ZipCode
	}
	if req.Status != nil {
		customer.Status = *req.Status
	}
	if req.Source != nil {
		customer.Source = *req.Source
	}
	if req.Tags != nil {
		customer.Tags = req.Tags
	}
	if req.Notes != nil {
		customer.Notes = *req.Notes
	}

	err = s.repo.Update(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

// DeleteCustomer deletes a customer
func (s *CustomerService) DeleteCustomer(tenantID string, id int) error {
	err := s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// SearchCustomers searches customers by query
func (s *CustomerService) SearchCustomers(tenantID, query string, page, limit int) ([]*models.Customer, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	customers, err := s.repo.Search(tenantID, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, nil
}

// GetCustomerStats gets customer statistics
func (s *CustomerService) GetCustomerStats(tenantID string) (map[string]interface{}, error) {
	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer count: %w", err)
	}

	stats := map[string]interface{}{
		"total_customers":  total,
		"active_customers": total, // For now, assume all are active
		"new_this_month":   0,     // Would need additional queries
	}

	return stats, nil
}

// ValidateCustomerID validates if a customer ID exists for a tenant
func (s *CustomerService) ValidateCustomerID(tenantID string, customerID int) error {
	_, err := s.repo.GetByID(tenantID, customerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	return nil
}

// ParseCustomerID parses customer ID from string
func (s *CustomerService) ParseCustomerID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid customer ID: %w", err)
	}
	return id, nil
}
