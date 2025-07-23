package routes

import (
	"zplus-saas/apps/backend/pos-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupPOSRoutes sets up all POS service routes
func SetupPOSRoutes(app *fiber.App, handler *handlers.POSHandler) {
	api := app.Group("/api/v1")
	pos := api.Group("/pos")

	// Health check
	pos.Get("/health", handler.HealthCheck)

	// Categories routes
	categories := pos.Group("/categories")
	categories.Post("/", handler.CreateCategory)
	categories.Get("/", handler.GetCategories)
	categories.Get("/:id", handler.GetCategoryByID)
	categories.Put("/:id", handler.UpdateCategory)
	categories.Delete("/:id", handler.DeleteCategory)

	// Products routes
	products := pos.Group("/products")
	products.Post("/", handler.CreateProduct)
	products.Get("/", handler.GetProducts)
	products.Get("/low-stock", handler.GetLowStockProducts)
	products.Get("/:id", handler.GetProductByID)
	products.Put("/:id", handler.UpdateProduct)
	products.Delete("/:id", handler.DeleteProduct)

	// Orders routes
	orders := pos.Group("/orders")
	orders.Post("/", handler.CreateOrder)
	orders.Get("/", handler.GetOrders)
	orders.Get("/:id", handler.GetOrderByID)
	orders.Put("/:id/status", handler.UpdateOrderStatus)
	orders.Put("/:id/payment", handler.UpdatePaymentStatus)

	// Inventory routes
	inventory := pos.Group("/inventory")
	inventory.Get("/transactions", handler.GetInventoryTransactions)
	inventory.Post("/adjustments", handler.CreateInventoryAdjustment)

	// Analytics routes
	analytics := pos.Group("/analytics")
	analytics.Get("/dashboard", handler.GetDashboardAnalytics)
	analytics.Get("/top-products", handler.GetTopProducts)
}
