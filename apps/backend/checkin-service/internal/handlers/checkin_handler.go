package handlers

import (
	"strconv"
	"time"

	"../models"
	"../services"
	"github.com/gofiber/fiber/v2"
)

type CheckinHandler struct {
	service *services.CheckinService
}

func NewCheckinHandler(service *services.CheckinService) *CheckinHandler {
	return &CheckinHandler{service: service}
}

// CreateCheckin creates a new checkin record
func (h *CheckinHandler) CreateCheckin(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req models.CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Get client info
	ipAddress := c.IP()
	deviceInfo := c.Get("User-Agent")

	record, err := h.service.CreateCheckin(tenantID, &req, ipAddress, deviceInfo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Checkin recorded successfully",
		"data":    record,
	})
}

// GetCheckinRecords gets checkin records with filters
func (h *CheckinHandler) GetCheckinRecords(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	var employeeID *int
	if employeeIDStr := c.Query("employee_id"); employeeIDStr != "" {
		if id, err := strconv.Atoi(employeeIDStr); err == nil {
			employeeID = &id
		}
	}

	checkinType := c.Query("checkin_type")
	status := c.Query("status")

	var dateFrom, dateTo time.Time
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if date, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = date
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if date, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = date
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	records, total, err := h.service.GetCheckinRecords(tenantID, employeeID, checkinType, status, dateFrom, dateTo, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  records,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetCheckinRecordByID gets a checkin record by ID
func (h *CheckinHandler) GetCheckinRecordByID(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin record ID",
		})
	}

	record, err := h.service.GetCheckinRecordByID(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  record,
	})
}

// UpdateCheckinRecord updates a checkin record
func (h *CheckinHandler) UpdateCheckinRecord(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin record ID",
		})
	}

	var req struct {
		Location string `json:"location"`
		Photo    string `json:"photo"`
		Notes    string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	record, err := h.service.UpdateCheckinRecord(tenantID, id, req.Location, req.Photo, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Checkin record updated successfully",
		"data":    record,
	})
}

// ApproveCheckin approves a checkin record
func (h *CheckinHandler) ApproveCheckin(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin record ID",
		})
	}

	// TODO: Get approver ID from JWT token
	approverID := 1 // Placeholder

	err = h.service.ApproveCheckin(tenantID, id, approverID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Checkin record approved successfully",
	})
}

// RejectCheckin rejects a checkin record
func (h *CheckinHandler) RejectCheckin(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin record ID",
		})
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// TODO: Get approver ID from JWT token
	approverID := 1 // Placeholder

	err = h.service.RejectCheckin(tenantID, id, approverID, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Checkin record rejected successfully",
	})
}

// DeleteCheckinRecord deletes a checkin record
func (h *CheckinHandler) DeleteCheckinRecord(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid checkin record ID",
		})
	}

	err = h.service.DeleteCheckinRecord(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Checkin record deleted successfully",
	})
}

// GetTodayCheckinRecords gets today's checkin records for an employee
func (h *CheckinHandler) GetTodayCheckinRecords(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	employeeID, err := strconv.Atoi(c.Params("employee_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid employee ID",
		})
	}

	records, err := h.service.GetTodayCheckinRecords(tenantID, employeeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  records,
	})
}

// GetAttendanceStats gets attendance statistics
func (h *CheckinHandler) GetAttendanceStats(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var date time.Time
	if dateStr := c.Query("date"); dateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsedDate
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid date format. Use YYYY-MM-DD",
			})
		}
	} else {
		date = time.Now()
	}

	stats, err := h.service.GetAttendanceStats(tenantID, date)
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
