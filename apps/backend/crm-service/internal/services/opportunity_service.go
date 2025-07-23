package services

import (
	"fmt"
	"time"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/repositories"
)

type OpportunityService struct {
	repo *repositories.OpportunityRepository
}

func NewOpportunityService(repo *repositories.OpportunityRepository) *OpportunityService {
	return &OpportunityService{repo: repo}
}

// CreateOpportunity creates a new opportunity
func (s *OpportunityService) CreateOpportunity(tenantID string, req *models.CreateOpportunityRequest) (*models.Opportunity, error) {
	// Parse expected date
	expectedDate, err := time.Parse("2006-01-02", req.ExpectedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expected date format: %w", err)
	}

	opportunity := &models.Opportunity{
		TenantID:     tenantID,
		CustomerID:   req.CustomerID,
		Name:         req.Name,
		Description:  req.Description,
		Value:        req.Value,
		Currency:     req.Currency,
		Stage:        "prospecting", // Default stage
		Probability:  10,            // Default probability
		Source:       req.Source,
		AssignedTo:   req.AssignedTo,
		ExpectedDate: expectedDate,
	}

	if opportunity.Currency == "" {
		opportunity.Currency = "USD"
	}

	err = s.repo.Create(opportunity)
	if err != nil {
		return nil, fmt.Errorf("failed to create opportunity: %w", err)
	}

	return opportunity, nil
}

// GetOpportunity gets an opportunity by ID
func (s *OpportunityService) GetOpportunity(tenantID string, id int) (*models.Opportunity, error) {
	opportunity, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunity: %w", err)
	}

	return opportunity, nil
}

// GetOpportunities gets all opportunities with pagination
func (s *OpportunityService) GetOpportunities(tenantID string, page, limit int) ([]*models.Opportunity, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	opportunities, err := s.repo.GetAll(tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get opportunities: %w", err)
	}

	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get opportunity count: %w", err)
	}

	return opportunities, total, nil
}

// UpdateOpportunity updates an opportunity
func (s *OpportunityService) UpdateOpportunity(tenantID string, id int, req *models.UpdateOpportunityRequest) (*models.Opportunity, error) {
	// Get existing opportunity
	opportunity, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunity: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		opportunity.Name = *req.Name
	}
	if req.Description != nil {
		opportunity.Description = *req.Description
	}
	if req.Value != nil {
		opportunity.Value = *req.Value
	}
	if req.Currency != nil {
		opportunity.Currency = *req.Currency
	}
	if req.Stage != nil {
		opportunity.Stage = *req.Stage
	}
	if req.Probability != nil {
		if *req.Probability < 0 || *req.Probability > 100 {
			return nil, fmt.Errorf("probability must be between 0 and 100")
		}
		opportunity.Probability = *req.Probability
	}
	if req.Source != nil {
		opportunity.Source = *req.Source
	}
	if req.AssignedTo != nil {
		opportunity.AssignedTo = *req.AssignedTo
	}
	if req.ExpectedDate != nil {
		expectedDate, err := time.Parse("2006-01-02", *req.ExpectedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expected date format: %w", err)
		}
		opportunity.ExpectedDate = expectedDate
	}

	err = s.repo.Update(opportunity)
	if err != nil {
		return nil, fmt.Errorf("failed to update opportunity: %w", err)
	}

	return opportunity, nil
}

// DeleteOpportunity deletes an opportunity
func (s *OpportunityService) DeleteOpportunity(tenantID string, id int) error {
	err := s.repo.Delete(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete opportunity: %w", err)
	}

	return nil
}

// CloseOpportunityWon marks an opportunity as won
func (s *OpportunityService) CloseOpportunityWon(tenantID string, id int) error {
	err := s.repo.CloseWon(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to close opportunity as won: %w", err)
	}

	return nil
}

// CloseOpportunityLost marks an opportunity as lost
func (s *OpportunityService) CloseOpportunityLost(tenantID string, id int) error {
	err := s.repo.CloseLost(tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to close opportunity as lost: %w", err)
	}

	return nil
}

// GetOpportunitiesByStage gets opportunities by stage
func (s *OpportunityService) GetOpportunitiesByStage(tenantID, stage string, page, limit int) ([]*models.Opportunity, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	opportunities, err := s.repo.GetByStage(tenantID, stage, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunities by stage: %w", err)
	}

	return opportunities, nil
}

// GetOpportunitiesByCustomer gets opportunities for a customer
func (s *OpportunityService) GetOpportunitiesByCustomer(tenantID string, customerID, page, limit int) ([]*models.Opportunity, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	opportunities, err := s.repo.GetByCustomer(tenantID, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunities by customer: %w", err)
	}

	return opportunities, nil
}

// GetOpportunityStats gets opportunity statistics
func (s *OpportunityService) GetOpportunityStats(tenantID string) (map[string]interface{}, error) {
	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get opportunity count: %w", err)
	}

	totalValue, err := s.repo.GetTotalValue(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total value: %w", err)
	}

	// Get opportunities by stage (simplified)
	prospecting, _ := s.repo.GetByStage(tenantID, "prospecting", 1000, 0)
	qualification, _ := s.repo.GetByStage(tenantID, "qualification", 1000, 0)
	proposal, _ := s.repo.GetByStage(tenantID, "proposal", 1000, 0)
	negotiation, _ := s.repo.GetByStage(tenantID, "negotiation", 1000, 0)
	closedWon, _ := s.repo.GetByStage(tenantID, "closed-won", 1000, 0)
	closedLost, _ := s.repo.GetByStage(tenantID, "closed-lost", 1000, 0)

	stats := map[string]interface{}{
		"total_opportunities": total,
		"total_value":         totalValue,
		"pipeline": map[string]interface{}{
			"prospecting":   len(prospecting),
			"qualification": len(qualification),
			"proposal":      len(proposal),
			"negotiation":   len(negotiation),
			"closed_won":    len(closedWon),
			"closed_lost":   len(closedLost),
		},
		"win_rate": s.calculateWinRate(len(closedWon), len(closedLost)),
	}

	return stats, nil
}

// GetSalesPipeline gets sales pipeline data
func (s *OpportunityService) GetSalesPipeline(tenantID string) (map[string]interface{}, error) {
	stages := []string{"prospecting", "qualification", "proposal", "negotiation", "closed-won", "closed-lost"}
	pipeline := make(map[string]interface{})

	for _, stage := range stages {
		opportunities, err := s.repo.GetByStage(tenantID, stage, 1000, 0)
		if err != nil {
			continue
		}

		var stageValue float64
		for _, opp := range opportunities {
			stageValue += opp.Value
		}

		pipeline[stage] = map[string]interface{}{
			"count": len(opportunities),
			"value": stageValue,
		}
	}

	return pipeline, nil
}

// calculateWinRate calculates the win rate
func (s *OpportunityService) calculateWinRate(won, lost int) float64 {
	total := won + lost
	if total == 0 {
		return 0
	}
	return float64(won) / float64(total) * 100
}
