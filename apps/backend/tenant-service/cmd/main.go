package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"zplus-saas/apps/backend/shared/config"
	"zplus-saas/apps/backend/shared/database"
	"zplus-saas/apps/backend/tenant-service/internal/handlers"
	"zplus-saas/apps/backend/tenant-service/internal/repositories"
	"zplus-saas/apps/backend/tenant-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	tenantRepo := repositories.NewTenantRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	planRepo := repositories.NewPlanRepository(db)

	// Initialize services
	tenantService := services.NewTenantService(tenantRepo, subscriptionRepo, planRepo)

	// Initialize handlers
	tenantHandler := handlers.NewTenantHandler(tenantService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName + " Tenant Service",
		ServerHeader: "Tenant Service",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "tenant-service",
			"version": cfg.AppVersion,
		})
	})

	// API routes
	api := app.Group("/api")

	// Tenant routes
	tenants := api.Group("/tenants")
	tenants.Post("/", tenantHandler.Create)
	tenants.Get("/", tenantHandler.List)
	tenants.Get("/:id", tenantHandler.GetByID)
	tenants.Put("/:id", tenantHandler.Update)
	tenants.Delete("/:id", tenantHandler.Delete)
	tenants.Post("/:id/activate", tenantHandler.Activate)
	tenants.Post("/:id/suspend", tenantHandler.Suspend)

	// Subscription routes
	subscriptions := api.Group("/subscriptions")
	subscriptions.Get("/tenant/:tenant_id", tenantHandler.GetSubscription)
	subscriptions.Post("/tenant/:tenant_id", tenantHandler.CreateSubscription)
	subscriptions.Put("/tenant/:tenant_id", tenantHandler.UpdateSubscription)

	// Plans routes
	plans := api.Group("/plans")
	plans.Get("/", tenantHandler.ListPlans)
	plans.Get("/:id", tenantHandler.GetPlan)
	plans.Post("/", tenantHandler.CreatePlan)
	plans.Put("/:id", tenantHandler.UpdatePlan)

	// Start server
	go func() {
		port := cfg.TenantServicePort
		if port == "" {
			port = "8089"
		}

		log.Printf("🏢 Tenant Service starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down tenant service...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Tenant service stopped")
}
