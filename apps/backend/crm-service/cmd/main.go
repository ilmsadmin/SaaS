package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"zplus-saas/apps/backend/crm-service/internal/handlers"
	"zplus-saas/apps/backend/crm-service/internal/repositories"
	"zplus-saas/apps/backend/crm-service/internal/routes"
	"zplus-saas/apps/backend/crm-service/internal/services"
	"zplus-saas/apps/backend/shared/config"
	"zplus-saas/apps/backend/shared/database"
	"zplus-saas/apps/backend/shared/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connections
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	// Initialize repositories
	customerRepo := repositories.NewCustomerRepository(db)
	leadRepo := repositories.NewLeadRepository(db)
	opportunityRepo := repositories.NewOpportunityRepository(db)

	// Initialize services
	customerService := services.NewCustomerService(customerRepo)
	leadService := services.NewLeadService(leadRepo)
	opportunityService := services.NewOpportunityService(opportunityRepo)

	// Initialize handlers
	customerHandler := handlers.NewCustomerHandler(customerService)
	leadHandler := handlers.NewLeadHandler(leadService)
	opportunityHandler := handlers.NewOpportunityHandler(opportunityService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Tenant-ID",
	}))

	// Add tenant middleware
	app.Use(middleware.TenantMiddleware())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "crm-service",
			"timestamp": time.Now().Unix(),
		})
	})

	// Setup routes
	api := app.Group("/api/v1")
	routes.SetupCustomerRoutes(api, customerHandler)
	routes.SetupLeadRoutes(api, leadHandler)
	routes.SetupOpportunityRoutes(api, opportunityHandler)

	// Start server in goroutine
	go func() {
		port := cfg.CRMServicePort
		if port == "" {
			port = "8082"
		}
		log.Printf("🚀 CRM Service starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down CRM Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}

	log.Println("✅ CRM Service stopped")
}
