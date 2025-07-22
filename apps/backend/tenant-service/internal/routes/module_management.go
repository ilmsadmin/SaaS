package routes

import (
	"zplus-saas/apps/backend/tenant-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupModuleRoutes sets up module management routes
func SetupModuleRoutes(app *fiber.App, handler *handlers.ModuleHandler) {
	api := app.Group("/api")

	// Module routes
	modules := api.Group("/modules")
	modules.Get("/", handler.GetAvailableModules)                // GET /api/modules
	modules.Get("/tenant", handler.GetTenantModules)             // GET /api/modules/tenant (authenticated)
	modules.Post("/install", handler.InstallModule)              // POST /api/modules/install (authenticated)
	modules.Delete("/:moduleId", handler.UninstallModule)        // DELETE /api/modules/:moduleId (authenticated)
	modules.Post("/:moduleId/enable", handler.EnableModule)      // POST /api/modules/:moduleId/enable (authenticated)
	modules.Post("/:moduleId/disable", handler.DisableModule)    // POST /api/modules/:moduleId/disable (authenticated)
	modules.Put("/:moduleId/config", handler.UpdateModuleConfig) // PUT /api/modules/:moduleId/config (authenticated)

	// Tenant configuration routes
	config := api.Group("/tenant/config")
	config.Get("/", handler.GetTenantConfiguration)              // GET /api/tenant/config (authenticated)
	config.Put("/", handler.UpdateTenantConfiguration)           // PUT /api/tenant/config (authenticated)
	config.Post("/domain", handler.SetupCustomDomain)            // POST /api/tenant/config/domain (authenticated)
	config.Delete("/domain", handler.RemoveCustomDomain)         // DELETE /api/tenant/config/domain (authenticated)
	config.Get("/features", handler.GetFeatureFlags)             // GET /api/tenant/config/features (authenticated)
	config.Put("/features/:flagName", handler.UpdateFeatureFlag) // PUT /api/tenant/config/features/:flagName (authenticated)
}
