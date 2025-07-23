package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"zplus-saas/apps/backend/hrm-service/internal/services"
)

type PerformanceHandler struct {
	service *services.PerformanceService
}

func NewPerformanceHandler(service *services.PerformanceService) *PerformanceHandler {
	return &PerformanceHandler{service: service}
}

// Create performance review
func (h *PerformanceHandler) CreatePerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req struct {
		EmployeeID    int     `json:"employee_id" validate:"required"`
		ReviewerID    int     `json:"reviewer_id" validate:"required"`
		Period        string  `json:"period" validate:"required"`
		ReviewType    string  `json:"review_type" validate:"required"`
		OverallRating float64 `json:"overall_rating" validate:"required"`
		Goals         string  `json:"goals"`
		Achievements  string  `json:"achievements"`
		Strengths     string  `json:"strengths"`
		Areas         string  `json:"areas_for_improvement"`
		Comments      string  `json:"comments"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	performance, err := h.service.CreatePerformance(tenantID, req.EmployeeID, req.ReviewerID,
		req.Period, req.ReviewType, req.OverallRating, req.Goals, req.Achievements,
		req.Strengths, req.Areas, req.Comments)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Performance review created successfully",
		"data":    performance,
	})
}

// Get performance review by ID
func (h *PerformanceHandler) GetPerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid performance review ID",
		})
	}

	performance, err := h.service.GetPerformance(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  performance,
	})
}

// Get performance reviews by employee ID
func (h *PerformanceHandler) GetPerformanceByEmployee(c *fiber.Ctx) error {
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

	reviews, total, err := h.service.GetPerformanceByEmployee(tenantID, employeeID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  reviews,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Get all performance reviews
func (h *PerformanceHandler) GetAllPerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	reviewType := c.Query("review_type")
	status := c.Query("status")

	var employeeID, reviewerID *int
	if empID := c.Query("employee_id"); empID != "" {
		if id, err := strconv.Atoi(empID); err == nil {
			employeeID = &id
		}
	}
	if revID := c.Query("reviewer_id"); revID != "" {
		if id, err := strconv.Atoi(revID); err == nil {
			reviewerID = &id
		}
	}

	reviews, total, err := h.service.GetAllPerformance(tenantID, employeeID, reviewerID, reviewType, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"error": false,
		"data":  reviews,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

// Update performance review
func (h *PerformanceHandler) UpdatePerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid performance review ID",
		})
	}

	var req struct {
		Period        string  `json:"period" validate:"required"`
		ReviewType    string  `json:"review_type" validate:"required"`
		OverallRating float64 `json:"overall_rating" validate:"required"`
		Goals         string  `json:"goals"`
		Achievements  string  `json:"achievements"`
		Strengths     string  `json:"strengths"`
		Areas         string  `json:"areas_for_improvement"`
		Comments      string  `json:"comments"`
		Status        string  `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	performance, err := h.service.UpdatePerformance(tenantID, id, req.Period, req.ReviewType,
		req.OverallRating, req.Goals, req.Achievements, req.Strengths, req.Areas,
		req.Comments, req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Performance review updated successfully",
		"data":    performance,
	})
}

// Delete performance review
func (h *PerformanceHandler) DeletePerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid performance review ID",
		})
	}

	err = h.service.DeletePerformance(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Performance review deleted successfully",
	})
}

// Submit performance review
func (h *PerformanceHandler) SubmitPerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid performance review ID",
		})
	}

	err = h.service.SubmitPerformance(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Performance review submitted successfully",
	})
}

// Complete performance review
func (h *PerformanceHandler) CompletePerformance(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid performance review ID",
		})
	}

	err = h.service.CompletePerformance(tenantID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Performance review completed successfully",
	})
}

// Get average performance rating
func (h *PerformanceHandler) GetAveragePerformanceRating(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	avgRating, err := h.service.GetAveragePerformanceRating(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  fiber.Map{"average_rating": avgRating},
	})
}

// Get performance statistics by department
func (h *PerformanceHandler) GetPerformanceStatsByDepartment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	stats, err := h.service.GetPerformanceStatsByDepartment(tenantID)
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
