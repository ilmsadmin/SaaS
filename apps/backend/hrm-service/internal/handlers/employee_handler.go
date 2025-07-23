package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/services"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

// Create employee
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req models.EmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	employee, err := h.service.CreateEmployee(tenantID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Employee created successfully",
		"data":    employee,
	})
}

// Get employee by ID
func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	employee, err := h.service.GetEmployee(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  employee,
	})
}

// Get employee by email
func (h *EmployeeHandler) GetEmployeeByEmail(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email parameter is required",
		})
	}

	employee, err := h.service.GetEmployeeByEmail(tenantID, email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  employee,
	})
}

// Get all employees
func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")

	var departmentID *int
	if depID := c.Query("department_id"); depID != "" {
		if id, err := strconv.Atoi(depID); err == nil {
			departmentID = &id
		}
	}

	employees, total, err := h.service.GetAllEmployees(tenantID, departmentID, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  employees,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Update employee
func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	var req models.EmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	employee, err := h.service.UpdateEmployee(tenantID, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Employee updated successfully",
		"data":    employee,
	})
}

// Delete employee
func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	err = h.service.DeleteEmployee(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Employee deleted successfully",
	})
}

// Search employees
func (h *EmployeeHandler) SearchEmployees(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	searchTerm := c.Query("q")
	if searchTerm == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Search term is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	employees, total, err := h.service.SearchEmployees(tenantID, searchTerm, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  employees,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Get HRM statistics
func (h *EmployeeHandler) GetHRMStatistics(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.service.GetHRMStatistics(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  stats,
	})
}
