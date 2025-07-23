package handlers

import (
	"strconv"

	"zplus-saas/apps/backend/crm-service/internal/models"
	"zplus-saas/apps/backend/crm-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	var req models.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	customer, err := h.service.CreateCustomer(tenantID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"error":    false,
		"message":  "Customer created successfully",
		"customer": customer,
	})
}

// GetCustomer gets a customer by ID
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
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
			"message": "Invalid customer ID",
		})
	}

	customer, err := h.service.GetCustomer(tenantID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"customer": customer,
	})
}

// GetCustomers gets all customers with pagination
func (h *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	customers, total, err := h.service.GetCustomers(tenantID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":     false,
		"customers": customers,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateCustomer updates a customer
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
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
			"message": "Invalid customer ID",
		})
	}

	var req models.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	customer, err := h.service.UpdateCustomer(tenantID, id, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"message":  "Customer updated successfully",
		"customer": customer,
	})
}

// DeleteCustomer deletes a customer
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
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
			"message": "Invalid customer ID",
		})
	}

	err = h.service.DeleteCustomer(tenantID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Customer deleted successfully",
	})
}

// SearchCustomers searches customers
func (h *CustomerHandler) SearchCustomers(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	query := c.Query("q", "")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Search query is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	customers, err := h.service.SearchCustomers(tenantID, query, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":     false,
		"customers": customers,
		"query":     query,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	})
}

// GetCustomerStats gets customer statistics
func (h *CustomerHandler) GetCustomerStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenantID").(string)
	if tenantID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant ID is required",
		})
	}

	stats, err := h.service.GetCustomerStats(tenantID)
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
