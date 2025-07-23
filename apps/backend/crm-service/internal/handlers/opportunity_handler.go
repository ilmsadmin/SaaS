package handlers

import (
	"strconv"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type OpportunityHandler struct {
	service *services.OpportunityService
}

func NewOpportunityHandler(service *services.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{service: service}
}

// CreateOpportunity creates a new opportunity
func (h *OpportunityHandler) CreateOpportunity(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	var req models.CreateOpportunityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	opportunity, err := h.service.CreateOpportunity(tenantID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"error":       false,
		"message":     "Opportunity created successfully",
		"opportunity": opportunity,
	})
}

// GetOpportunity gets an opportunity by ID
func (h *OpportunityHandler) GetOpportunity(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid opportunity ID",
		})
	}

	opportunity, err := h.service.GetOpportunity(tenantID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":       false,
		"opportunity": opportunity,
	})
}

// GetOpportunities gets all opportunities with pagination
func (h *OpportunityHandler) GetOpportunities(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	opportunities, total, err := h.service.GetOpportunities(tenantID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":         false,
		"opportunities": opportunities,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateOpportunity updates an opportunity
func (h *OpportunityHandler) UpdateOpportunity(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid opportunity ID",
		})
	}

	var req models.UpdateOpportunityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	opportunity, err := h.service.UpdateOpportunity(tenantID, id, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":       false,
		"message":     "Opportunity updated successfully",
		"opportunity": opportunity,
	})
}

// DeleteOpportunity deletes an opportunity
func (h *OpportunityHandler) DeleteOpportunity(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid opportunity ID",
		})
	}

	err = h.service.DeleteOpportunity(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Opportunity deleted successfully",
	})
}

// CloseOpportunityWon marks an opportunity as won
func (h *OpportunityHandler) CloseOpportunityWon(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid opportunity ID",
		})
	}

	err = h.service.CloseOpportunityWon(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Opportunity closed as won",
	})
}

// CloseOpportunityLost marks an opportunity as lost
func (h *OpportunityHandler) CloseOpportunityLost(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid opportunity ID",
		})
	}

	err = h.service.CloseOpportunityLost(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Opportunity closed as lost",
	})
}

// GetOpportunitiesByStage gets opportunities by stage
func (h *OpportunityHandler) GetOpportunitiesByStage(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	stage := c.Params("stage")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	opportunities, err := h.service.GetOpportunitiesByStage(tenantID, stage, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":         false,
		"opportunities": opportunities,
		"stage":         stage,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	})
}

// GetOpportunitiesByCustomer gets opportunities for a customer
func (h *OpportunityHandler) GetOpportunitiesByCustomer(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	customerIDStr := c.Params("customerId")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid customer ID",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	opportunities, err := h.service.GetOpportunitiesByCustomer(tenantID, customerID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":         false,
		"opportunities": opportunities,
		"customer_id":   customerID,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	})
}

// GetOpportunityStats gets opportunity statistics
func (h *OpportunityHandler) GetOpportunityStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	stats, err := h.service.GetOpportunityStats(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"stats": stats,
	})
}

// GetSalesPipeline gets sales pipeline data
func (h *OpportunityHandler) GetSalesPipeline(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	pipeline, err := h.service.GetSalesPipeline(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"pipeline": pipeline,
	})
}
