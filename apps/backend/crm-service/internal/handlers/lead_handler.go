package handlers

import (
	"strconv"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LeadHandler struct {
	service *services.LeadService
}

func NewLeadHandler(service *services.LeadService) *LeadHandler {
	return &LeadHandler{service: service}
}

// CreateLead creates a new lead
func (h *LeadHandler) CreateLead(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	var req models.CreateLeadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	lead, err := h.service.CreateLead(tenantID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Lead created successfully",
		"lead":    lead,
	})
}

// GetLead gets a lead by ID
func (h *LeadHandler) GetLead(c *fiber.Ctx) error {
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
			"message": "Invalid lead ID",
		})
	}

	lead, err := h.service.GetLead(tenantID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"lead":  lead,
	})
}

// GetLeads gets all leads with pagination
func (h *LeadHandler) GetLeads(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	leads, total, err := h.service.GetLeads(tenantID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"leads": leads,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateLead updates a lead
func (h *LeadHandler) UpdateLead(c *fiber.Ctx) error {
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
			"message": "Invalid lead ID",
		})
	}

	var req models.UpdateLeadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	lead, err := h.service.UpdateLead(tenantID, id, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Lead updated successfully",
		"lead":    lead,
	})
}

// DeleteLead deletes a lead
func (h *LeadHandler) DeleteLead(c *fiber.Ctx) error {
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
			"message": "Invalid lead ID",
		})
	}

	err = h.service.DeleteLead(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Lead deleted successfully",
	})
}

// ConvertLead converts a lead to customer
func (h *LeadHandler) ConvertLead(c *fiber.Ctx) error {
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
			"message": "Invalid lead ID",
		})
	}

	err = h.service.ConvertLead(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Lead converted successfully",
	})
}

// GetLeadsByStatus gets leads by status
func (h *LeadHandler) GetLeadsByStatus(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	status := c.Params("status")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	leads, err := h.service.GetLeadsByStatus(tenantID, status, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":  false,
		"leads":  leads,
		"status": status,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	})
}

// ScoreLead updates lead score
func (h *LeadHandler) ScoreLead(c *fiber.Ctx) error {
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
			"message": "Invalid lead ID",
		})
	}

	var req struct {
		Score int `json:"score"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	err = h.service.ScoreLead(tenantID, id, req.Score)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Lead scored successfully",
	})
}

// GetLeadStats gets lead statistics
func (h *LeadHandler) GetLeadStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	stats, err := h.service.GetLeadStats(tenantID)
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
