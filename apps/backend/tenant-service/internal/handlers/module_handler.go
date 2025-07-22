package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"zplus-saas/apps/backend/tenant-service/internal/models"
	"zplus-saas/apps/backend/tenant-service/internal/services"
)

type ModuleHandler struct {
	moduleService       *services.ModuleService
	tenantConfigService *services.TenantConfigService
	validator           *validator.Validate
}

func NewModuleHandler(
	moduleService *services.ModuleService,
	tenantConfigService *services.TenantConfigService,
) *ModuleHandler {
	return &ModuleHandler{
		moduleService:       moduleService,
		tenantConfigService: tenantConfigService,
		validator:           validator.New(),
	}
}

// GetAvailableModules returns all available modules
func (h *ModuleHandler) GetAvailableModules(c *fiber.Ctx) error {
	modules, err := h.moduleService.GetAvailableModules()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"modules": modules})
}

// GetTenantModules returns modules for a specific tenant
func (h *ModuleHandler) GetTenantModules(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	modules, err := h.moduleService.GetTenantModules(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"modules": modules})
}

// InstallModule installs a module for a tenant
func (h *ModuleHandler) InstallModule(c *fiber.Ctx) error {
	var req models.InstallModuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	moduleID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid module ID"})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	installation, err := h.moduleService.InstallModule(tenantID, moduleID, userID, req.Version, req.Config)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":      "Module installed successfully",
		"installation": installation,
	})
}

// UninstallModule uninstalls a module for a tenant
func (h *ModuleHandler) UninstallModule(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	moduleID, err := uuid.Parse(c.Params("moduleId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid module ID"})
	}

	err = h.moduleService.UninstallModule(tenantID, moduleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Module uninstalled successfully"})
}

// EnableModule enables a module for a tenant
func (h *ModuleHandler) EnableModule(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	moduleID, err := uuid.Parse(c.Params("moduleId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid module ID"})
	}

	err = h.moduleService.EnableModule(tenantID, moduleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Module enabled successfully"})
}

// DisableModule disables a module for a tenant
func (h *ModuleHandler) DisableModule(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	moduleID, err := uuid.Parse(c.Params("moduleId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid module ID"})
	}

	err = h.moduleService.DisableModule(tenantID, moduleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Module disabled successfully"})
}

// UpdateModuleConfig updates module configuration
func (h *ModuleHandler) UpdateModuleConfig(c *fiber.Ctx) error {
	var req models.UpdateModuleConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	moduleID, err := uuid.Parse(c.Params("moduleId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid module ID"})
	}

	err = h.moduleService.UpdateModuleConfig(tenantID, moduleID, req.Config)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Module configuration updated successfully"})
}

// GetTenantConfiguration gets tenant configuration
func (h *ModuleHandler) GetTenantConfiguration(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	config, err := h.tenantConfigService.GetTenantConfiguration(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"configuration": config})
}

// UpdateTenantConfiguration updates tenant configuration
func (h *ModuleHandler) UpdateTenantConfiguration(c *fiber.Ctx) error {
	var req models.UpdateTenantConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	err = h.tenantConfigService.UpdateTenantConfiguration(tenantID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Configuration updated successfully"})
}

// SetupCustomDomain sets up custom domain for tenant
func (h *ModuleHandler) SetupCustomDomain(c *fiber.Ctx) error {
	var req struct {
		Domain     string `json:"domain" validate:"required"`
		SSLEnabled bool   `json:"ssl_enabled"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	err = h.tenantConfigService.SetupCustomDomain(tenantID, req.Domain, req.SSLEnabled)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Custom domain configured successfully"})
}

// RemoveCustomDomain removes custom domain from tenant
func (h *ModuleHandler) RemoveCustomDomain(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	err = h.tenantConfigService.RemoveCustomDomain(tenantID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Custom domain removed successfully"})
}

// GetFeatureFlags gets feature flags for tenant
func (h *ModuleHandler) GetFeatureFlags(c *fiber.Ctx) error {
	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	flags, err := h.tenantConfigService.GetFeatureFlags(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"feature_flags": flags})
}

// UpdateFeatureFlag updates a specific feature flag
func (h *ModuleHandler) UpdateFeatureFlag(c *fiber.Ctx) error {
	var req struct {
		Value interface{} `json:"value" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tenantID, err := uuid.Parse(c.Locals("tenant_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	flagName := c.Params("flagName")
	if flagName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Flag name is required"})
	}

	err = h.tenantConfigService.UpdateFeatureFlag(tenantID, flagName, req.Value)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Feature flag updated successfully"})
}
