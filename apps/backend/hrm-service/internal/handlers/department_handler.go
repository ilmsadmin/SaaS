package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"zplus-saas/apps/backend/hrm-service/internal/services"
)

type DepartmentHandler struct {
	service *services.DepartmentService
}

func NewDepartmentHandler(service *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// Create department
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req struct {
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description"`
		Location    string  `json:"location"`
		ManagerID   *int    `json:"manager_id"`
		Budget      float64 `json:"budget"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	department, err := h.service.CreateDepartment(tenantID, req.Name, req.Description, req.Location, req.ManagerID, req.Budget)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Department created successfully",
		"data":    department,
	})
}

// Get department by ID
func (h *DepartmentHandler) GetDepartment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid department ID",
		})
	}

	department, err := h.service.GetDepartment(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  department,
	})
}

// Get all departments
func (h *DepartmentHandler) GetAllDepartments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	departments, total, err := h.service.GetAllDepartments(tenantID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  departments,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Update department
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid department ID",
		})
	}

	var req struct {
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description"`
		Location    string  `json:"location"`
		ManagerID   *int    `json:"manager_id"`
		Budget      float64 `json:"budget"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	department, err := h.service.UpdateDepartment(tenantID, id, req.Name, req.Description, req.Location, req.ManagerID, req.Budget)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Department updated successfully",
		"data":    department,
	})
}

// Delete department
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid department ID",
		})
	}

	err = h.service.DeleteDepartment(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Department deleted successfully",
	})
}

// Get departments with employee count
func (h *DepartmentHandler) GetDepartmentsWithEmployeeCount(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	departments, err := h.service.GetDepartmentsWithEmployeeCount(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  departments,
	})
}
