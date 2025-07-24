package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"zplus-saas/apps/backend/api-gateway/internal/handlers"
	"zplus-saas/apps/backend/api-gateway/internal/middleware"
	"zplus-saas/apps/backend/shared/config"
	"zplus-saas/apps/backend/shared/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	redis, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redis.Close()

	// TODO: Add MongoDB when needed
	// mongo, err := database.NewMongoDB(cfg.MongoURL)
	// if err != nil {
	// 	log.Fatal("Failed to connect to MongoDB:", err)
	// }
	// defer mongo.Client().Disconnect(context.Background())

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Zplus SaaS API Gateway v" + cfg.AppVersion,
		ServerHeader: "Zplus-Gateway",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"code":    code,
				"success": false,
			})
		},
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://localhost:3002",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Custom middleware
	app.Use(middleware.TenantResolver())
	app.Use(middleware.SecurityHeaders())

	// Initialize handlers
	handlers := handlers.New(db, redis, nil, cfg)

	// Routes
	setupRoutes(app, handlers)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "api-gateway",
			"version":   cfg.AppVersion,
		})
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	log.Printf("🚀 Zplus SaaS API Gateway starting on port %s", port)
	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}

	fmt.Println("Server shutdown complete.")
}

func setupRoutes(app *fiber.App, h *handlers.Handlers) {
	api := app.Group("/api/v1")

	// Authentication routes
	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)
	auth.Post("/logout", h.Auth.Logout)
	auth.Get("/profile", middleware.AuthRequired(), h.Auth.Profile)
	auth.Put("/profile", middleware.AuthRequired(), h.Auth.UpdateProfile)

	// Admin authentication routes
	admin := api.Group("/admin")
	adminAuth := admin.Group("/auth")
	adminAuth.Post("/login", h.Auth.AdminLogin)
	adminAuth.Post("/create", h.Auth.CreateAdmin)
	adminAuth.Get("/validate", middleware.AuthRequired(), h.Auth.ValidateAdmin)

	// Admin dashboard routes
	admin.Get("/stats", middleware.AuthRequired(), h.Auth.AdminStats)
	admin.Get("/activities", middleware.AuthRequired(), h.Auth.AdminActivities)
	admin.Get("/health", middleware.AuthRequired(), h.Auth.AdminHealth)

	// Tenant management (System level)
	tenants := api.Group("/tenants", middleware.AuthRequired(), middleware.SystemAdminRequired())
	tenants.Get("/", h.Tenant.List)
	tenants.Post("/", h.Tenant.Create)
	tenants.Get("/:id", h.Tenant.GetByID)
	tenants.Put("/:id", h.Tenant.Update)
	tenants.Delete("/:id", h.Tenant.Delete)

	// Module management
	modules := api.Group("/modules", middleware.AuthRequired())
	modules.Get("/", h.Module.List)
	modules.Get("/:module", h.Module.GetStatus)
	modules.Post("/:module/enable", middleware.TenantAdminRequired(), h.Module.Enable)
	modules.Post("/:module/disable", middleware.TenantAdminRequired(), h.Module.Disable)

	// Proxy routes to microservices
	api.All("/crm/*", h.Proxy.CRM)
	api.All("/hrm/*", h.Proxy.HRM)
	api.All("/pos/*", h.Proxy.POS)
	api.All("/lms/*", h.Proxy.LMS)
	api.All("/checkin/*", h.Proxy.Checkin)
	api.All("/payment/*", h.Proxy.Payment)
	api.All("/files/*", h.Proxy.Files)
}
