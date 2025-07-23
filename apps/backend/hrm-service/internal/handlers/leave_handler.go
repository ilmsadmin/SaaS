package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"zplus-saas/apps/backend/hrm-service/internal/models"
	"zplus-saas/apps/backend/hrm-service/internal/services"
)

type LeaveHandler struct {
	service *services.LeaveService
}

func NewLeaveHandler(service *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// Create leave request
func (h *LeaveHandler) CreateLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req models.LeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	leave, err := h.service.CreateLeave(tenantID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Leave request created successfully",
		"data":    leave,
	})
}

// Get leave by ID
func (h *LeaveHandler) GetLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid leave ID",
		})
	}

	leave, err := h.service.GetLeave(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  leave,
	})
}

// Get leaves by employee ID
func (h *LeaveHandler) GetLeavesByEmployee(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	employeeID, err := strconv.Atoi(c.Params("employee_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	leaves, total, err := h.service.GetLeavesByEmployee(tenantID, employeeID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  leaves,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Get all leaves
func (h *LeaveHandler) GetAllLeaves(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	leaveType := c.Query("leave_type")

	var employeeID *int
	if empID := c.Query("employee_id"); empID != "" {
		if id, err := strconv.Atoi(empID); err == nil {
			employeeID = &id
		}
	}

	leaves, total, err := h.service.GetAllLeaves(tenantID, employeeID, status, leaveType, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  leaves,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Update leave
func (h *LeaveHandler) UpdateLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid leave ID",
		})
	}

	var req models.LeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	leave, err := h.service.UpdateLeave(tenantID, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Leave updated successfully",
		"data":    leave,
	})
}

// Approve leave
func (h *LeaveHandler) ApproveLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid leave ID",
		})
	}

	var req struct {
		ApproverID int    `json:"approver_id" validate:"required"`
		Comments   string `json:"comments"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	err = h.service.ApproveLeave(tenantID, id, req.ApproverID, req.Comments)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Leave approved successfully",
	})
}

// Reject leave
func (h *LeaveHandler) RejectLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid leave ID",
		})
	}

	var req struct {
		ApproverID int    `json:"approver_id" validate:"required"`
		Comments   string `json:"comments"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	err = h.service.RejectLeave(tenantID, id, req.ApproverID, req.Comments)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Leave rejected successfully",
	})
}

// Delete leave
func (h *LeaveHandler) DeleteLeave(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid leave ID",
		})
	}

	err = h.service.DeleteLeave(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Leave deleted successfully",
	})
}

// Get pending leaves count
func (h *LeaveHandler) GetPendingLeavesCount(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	count, err := h.service.GetPendingLeavesCount(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  fiber.Map{"pending_count": count},
	})
}

// Get leave balance
func (h *LeaveHandler) GetLeaveBalance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	employeeID, err := strconv.Atoi(c.Query("employee_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Employee ID is required",
		})
	}

	leaveType := c.Query("leave_type")
	if leaveType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Leave type is required",
		})
	}

	balance, err := h.service.GetLeaveBalance(tenantID, employeeID, leaveType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data": fiber.Map{
			"employee_id": employeeID,
			"leave_type":  leaveType,
			"balance":     balance,
		},
	})
}
