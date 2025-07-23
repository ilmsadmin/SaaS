package routes

import (
	"zplus-saas/apps/backend/crm-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupCustomerRoutes sets up customer routes
func SetupCustomerRoutes(api fiber.Router, handler *handlers.CustomerHandler) {
	customers := api.Group("/customers")

	customers.Post("/", handler.CreateCustomer)
	customers.Get("/", handler.GetCustomers)
	customers.Get("/search", handler.SearchCustomers)
	customers.Get("/stats", handler.GetCustomerStats)
	customers.Get("/:id", handler.GetCustomer)
	customers.Put("/:id", handler.UpdateCustomer)
	customers.Delete("/:id", handler.DeleteCustomer)
}

// SetupLeadRoutes sets up lead routes
func SetupLeadRoutes(api fiber.Router, handler *handlers.LeadHandler) {
	leads := api.Group("/leads")

	leads.Post("/", handler.CreateLead)
	leads.Get("/", handler.GetLeads)
	leads.Get("/stats", handler.GetLeadStats)
	leads.Get("/status/:status", handler.GetLeadsByStatus)
	leads.Get("/:id", handler.GetLead)
	leads.Put("/:id", handler.UpdateLead)
	leads.Delete("/:id", handler.DeleteLead)
	leads.Post("/:id/convert", handler.ConvertLead)
	leads.Post("/:id/score", handler.ScoreLead)
}

// SetupOpportunityRoutes sets up opportunity routes
func SetupOpportunityRoutes(api fiber.Router, handler *handlers.OpportunityHandler) {
	opportunities := api.Group("/opportunities")

	opportunities.Post("/", handler.CreateOpportunity)
	opportunities.Get("/", handler.GetOpportunities)
	opportunities.Get("/stats", handler.GetOpportunityStats)
	opportunities.Get("/pipeline", handler.GetSalesPipeline)
	opportunities.Get("/stage/:stage", handler.GetOpportunitiesByStage)
	opportunities.Get("/customer/:customerId", handler.GetOpportunitiesByCustomer)
	opportunities.Get("/:id", handler.GetOpportunity)
	opportunities.Put("/:id", handler.UpdateOpportunity)
	opportunities.Delete("/:id", handler.DeleteOpportunity)
	opportunities.Post("/:id/close-won", handler.CloseOpportunityWon)
	opportunities.Post("/:id/close-lost", handler.CloseOpportunityLost)
}
