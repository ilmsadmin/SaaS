package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"zplus-saas/apps/backend/auth-service/internal/handlers"
	"zplus-saas/apps/backend/auth-service/internal/middleware"
	"zplus-saas/apps/backend/auth-service/internal/repositories"
	"zplus-saas/apps/backend/auth-service/internal/services"
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

	// Connect to PostgreSQL
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tenantRepo := repositories.NewTenantRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// Initialize services
	jwtService := services.NewJWTService(cfg)
	authService := services.NewAuthService(userRepo, tenantRepo, refreshTokenRepo, jwtService, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName + " Auth Service",
		ServerHeader: "Auth Service",
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
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://localhost:8080",
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "auth-service",
			"version": cfg.AppVersion,
		})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	authProtected := auth.Use(authMiddleware.RequireAuth)
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Get("/profile", authHandler.GetProfile)
	authProtected.Put("/profile", authHandler.UpdateProfile)

	// Admin routes
	admin := api.Group("/admin")
	adminAuth := admin.Group("/auth")
	adminAuth.Post("/login", adminHandler.AdminLogin)
	adminAuth.Post("/create", adminHandler.CreateAdminUser)

	// Protected admin routes
	adminProtected := adminAuth.Use(authMiddleware.RequireAuth)
	adminProtected.Get("/validate", adminHandler.ValidateAdmin)

	// Admin dashboard routes
	adminDashboard := admin.Use(authMiddleware.RequireAuth)
	adminDashboard.Get("/stats", adminHandler.AdminStats)
	adminDashboard.Get("/activities", adminHandler.AdminActivities)
	adminDashboard.Get("/health", adminHandler.AdminHealth)

	// Legacy admin routes (for future use)
	adminProtected2 := api.Group("/admin")
	adminProtected2.Use(authMiddleware.RequireAuth)
	adminProtected2.Use(authMiddleware.RequireRole("admin", "super_admin"))

	// Start server
	go func() {
		port := os.Getenv("AUTH_SERVICE_PORT")
		if port == "" {
			port = cfg.Port
		}
		if port == "" {
			port = "8081"
		}

		log.Printf("🚀 Auth Service starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🔄 Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}
