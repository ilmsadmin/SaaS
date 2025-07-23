package services

import (
	"fmt"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/repositories"
)

type LeadService struct {
	repo *repositories.LeadRepository
}

func NewLeadService(repo *repositories.LeadRepository) *LeadService {
	return &LeadService{repo: repo}
}

// CreateLead creates a new lead
func (s *LeadService) CreateLead(tenantID string, req *models.CreateLeadRequest) (*models.Lead, error) {
	lead := &models.Lead{
		TenantID:   tenantID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Company:    req.Company,
		Title:      req.Title,
		Source:     req.Source,
		Status:     "new", // Default status
		Score:      0,     // Default score
		AssignedTo: req.AssignedTo,
		Value:      req.Value,
		Notes:      req.Notes,
	}

	err := s.repo.Create(lead)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}

	return lead, nil
}

// GetLead gets a lead by ID
func (s *LeadService) GetLead(tenantID string, id int) (*models.Lead, error) {
	lead, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}

	return lead, nil
}

// GetLeads gets all leads with pagination
func (s *LeadService) GetLeads(tenantID string, page, limit int) ([]*models.Lead, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	leads, err := s.repo.GetAll(tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get leads: %w", err)
	}

	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get lead count: %w", err)
	}

	return leads, total, nil
}

// UpdateLead updates a lead
func (s *LeadService) UpdateLead(tenantID string, id int, req *models.UpdateLeadRequest) (*models.Lead, error) {
	// Get existing lead
	lead, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		lead.Name = *req.Name
	}
	if req.Email != nil {
		lead.Email = *req.Email
	}
	if req.Phone != nil {
		lead.Phone = *req.Phone
	}
	if req.Company != nil {
		lead.Company = *req.Company
	}
	if req.Title != nil {
		lead.Title = *req.Title
	}
	if req.Source != nil {
		lead.Source = *req.Source
	}
	if req.Status != nil {
		lead.Status = *req.Status
	}
	if req.Score != nil {
		lead.Score = *req.Score
	}
	if req.AssignedTo != nil {
		lead.AssignedTo = *req.AssignedTo
	}
	if req.Value != nil {
		lead.Value = *req.Value
	}
	if req.Notes != nil {
		lead.Notes = *req.Notes
	}

	err = s.repo.Update(lead)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead: %w", err)
	}

	return lead, nil
}

// DeleteLead deletes a lead
func (s *LeadService) DeleteLead(tenantID string, id int) error {
	err := s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}

	return nil
}

// ConvertLead converts a lead to customer
func (s *LeadService) ConvertLead(tenantID string, id int) error {
	err := s.repo.ConvertToCustomer(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to convert lead: %w", err)
	}

	return nil
}

// GetLeadsByStatus gets leads by status
func (s *LeadService) GetLeadsByStatus(tenantID, status string, page, limit int) ([]*models.Lead, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	leads, err := s.repo.GetByStatus(tenantID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leads by status: %w", err)
	}

	return leads, nil
}

// GetLeadsByAssignedUser gets leads assigned to a user
func (s *LeadService) GetLeadsByAssignedUser(tenantID string, userID, page, limit int) ([]*models.Lead, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	leads, err := s.repo.GetByAssignedUser(tenantID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leads by assigned user: %w", err)
	}

	return leads, nil
}

// GetLeadStats gets lead statistics
func (s *LeadService) GetLeadStats(tenantID string) (map[string]interface{}, error) {
	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead count: %w", err)
	}

	// Get leads by status (simplified)
	newLeads, _ := s.repo.GetByStatus(tenantID, "new", 1000, 0)
	qualifiedLeads, _ := s.repo.GetByStatus(tenantID, "qualified", 1000, 0)
	convertedLeads, _ := s.repo.GetByStatus(tenantID, "converted", 1000, 0)

	stats := map[string]interface{}{
		"total_leads":     total,
		"new_leads":       len(newLeads),
		"qualified_leads": len(qualifiedLeads),
		"converted_leads": len(convertedLeads),
		"conversion_rate": s.calculateConversionRate(total, len(convertedLeads)),
	}

	return stats, nil
}

// ScoreLead updates lead score based on various factors
func (s *LeadService) ScoreLead(tenantID string, id int, score int) error {
	if score < 0 || score > 100 {
		return fmt.Errorf("score must be between 0 and 100")
	}

	lead, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to get lead: %w", err)
	}

	lead.Score = score
	err = s.repo.Update(lead)
	if err != nil {
		return fmt.Errorf("failed to update lead score: %w", err)
	}

	return nil
}

// calculateConversionRate calculates the conversion rate
func (s *LeadService) calculateConversionRate(total, converted int) float64 {
	if total == 0 {
		return 0
	}
	return float64(converted) / float64(total) * 100
}
