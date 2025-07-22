package handlers

import (
	"strconv"

	"zplus-saas/apps/backend/tenant-service/internal/models"
	"zplus-saas/apps/backend/tenant-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TenantHandler struct {
	tenantService services.TenantService
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func NewTenantHandler(tenantService services.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// Tenant endpoints

// Create godoc
// @Summary Create a new tenant
// @Description Create a new tenant with subdomain and optional custom domain
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body models.CreateTenantRequest true "Create tenant request"
// @Success 201 {object} models.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/tenants [post]
func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var req models.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	tenant, err := h.tenantService.CreateTenant(c.Context(), &req)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "subdomain already exists" || err.Error() == "domain already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(ErrorResponse{
			Error:   "Failed to create tenant",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

// List godoc
// @Summary List all tenants
// @Description Get paginated list of tenants
// @Tags tenants
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} ListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/tenants [get]
func (h *TenantHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	tenants, total, err := h.tenantService.ListTenants(c.Context(), perPage, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list tenants",
			Message: err.Error(),
		})
	}

	totalPages := (total + perPage - 1) / perPage

	return c.JSON(ListResponse{
		Data:       tenants,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get tenant by ID
// @Description Get tenant information by ID
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} models.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/tenants/{id} [get]
func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	tenant, err := h.tenantService.GetTenant(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Tenant not found",
			Message: err.Error(),
		})
	}

	return c.JSON(tenant)
}

// Update godoc
// @Summary Update tenant
// @Description Update tenant information
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param request body models.UpdateTenantRequest true "Update tenant request"
// @Success 200 {object} models.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/tenants/{id} [put]
func (h *TenantHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	var req models.UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	tenant, err := h.tenantService.UpdateTenant(c.Context(), id, &req)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "tenant not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "domain already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(ErrorResponse{
			Error:   "Failed to update tenant",
			Message: err.Error(),
		})
	}

	return c.JSON(tenant)
}

// Delete godoc
// @Summary Delete tenant
// @Description Delete tenant by ID
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/tenants/{id} [delete]
func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	err = h.tenantService.DeleteTenant(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete tenant",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Tenant deleted successfully",
	})
}

// Activate godoc
// @Summary Activate tenant
// @Description Activate a suspended tenant
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/tenants/{id}/activate [post]
func (h *TenantHandler) Activate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	err = h.tenantService.ActivateTenant(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to activate tenant",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Tenant activated successfully",
	})
}

// Suspend godoc
// @Summary Suspend tenant
// @Description Suspend an active tenant
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/tenants/{id}/suspend [post]
func (h *TenantHandler) Suspend(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	err = h.tenantService.SuspendTenant(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to suspend tenant",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Tenant suspended successfully",
	})
}

// Subscription endpoints

// GetSubscription godoc
// @Summary Get tenant subscription
// @Description Get current subscription for a tenant
// @Tags subscriptions
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/subscriptions/tenant/{tenant_id} [get]
func (h *TenantHandler) GetSubscription(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	subscription, err := h.tenantService.GetSubscription(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Subscription not found",
			Message: err.Error(),
		})
	}

	return c.JSON(subscription)
}

// CreateSubscription godoc
// @Summary Create tenant subscription
// @Description Create a new subscription for a tenant
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body models.CreateSubscriptionRequest true "Create subscription request"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} ErrorResponse
// @Router /api/subscriptions/tenant/{tenant_id} [post]
func (h *TenantHandler) CreateSubscription(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	var req models.CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	subscription, err := h.tenantService.CreateSubscription(c.Context(), tenantID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to create subscription",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(subscription)
}

// UpdateSubscription godoc
// @Summary Update tenant subscription
// @Description Update subscription for a tenant
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param request body models.UpdateSubscriptionRequest true "Update subscription request"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/subscriptions/tenant/{tenant_id} [put]
func (h *TenantHandler) UpdateSubscription(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid tenant ID",
			Message: "ID must be a valid UUID",
		})
	}

	var req models.UpdateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	subscription, err := h.tenantService.UpdateSubscription(c.Context(), tenantID, &req)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "subscription not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(ErrorResponse{
			Error:   "Failed to update subscription",
			Message: err.Error(),
		})
	}

	return c.JSON(subscription)
}

// Plan endpoints

// ListPlans godoc
// @Summary List all plans
// @Description Get paginated list of service plans
// @Tags plans
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param active_only query bool false "Show only active plans"
// @Success 200 {object} ListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/plans [get]
func (h *TenantHandler) ListPlans(c *fiber.Ctx) error {
	activeOnly := c.Query("active_only") == "true"

	if activeOnly {
		plans, err := h.tenantService.ListActivePlans(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to list plans",
				Message: err.Error(),
			})
		}

		return c.JSON(ListResponse{
			Data:       plans,
			Total:      len(plans),
			Page:       1,
			PerPage:    len(plans),
			TotalPages: 1,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	plans, total, err := h.tenantService.ListPlans(c.Context(), perPage, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list plans",
			Message: err.Error(),
		})
	}

	totalPages := (total + perPage - 1) / perPage

	return c.JSON(ListResponse{
		Data:       plans,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetPlan godoc
// @Summary Get plan by ID
// @Description Get plan information by ID
// @Tags plans
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} models.Plan
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/plans/{id} [get]
func (h *TenantHandler) GetPlan(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid plan ID",
			Message: "ID must be a valid UUID",
		})
	}

	plan, err := h.tenantService.GetPlan(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Plan not found",
			Message: err.Error(),
		})
	}

	return c.JSON(plan)
}

// CreatePlan godoc
// @Summary Create a new plan
// @Description Create a new service plan
// @Tags plans
// @Accept json
// @Produce json
// @Param request body models.CreatePlanRequest true "Create plan request"
// @Success 201 {object} models.Plan
// @Failure 400 {object} ErrorResponse
// @Router /api/plans [post]
func (h *TenantHandler) CreatePlan(c *fiber.Ctx) error {
	var req models.CreatePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	plan, err := h.tenantService.CreatePlan(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to create plan",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(plan)
}

// UpdatePlan godoc
// @Summary Update plan
// @Description Update plan information
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Param request body models.UpdatePlanRequest true "Update plan request"
// @Success 200 {object} models.Plan
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/plans/{id} [put]
func (h *TenantHandler) UpdatePlan(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid plan ID",
			Message: "ID must be a valid UUID",
		})
	}

	var req models.UpdatePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	plan, err := h.tenantService.UpdatePlan(c.Context(), id, &req)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "plan not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(ErrorResponse{
			Error:   "Failed to update plan",
			Message: err.Error(),
		})
	}

	return c.JSON(plan)
}
