package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zplus-saas/apps/backend/shared/config"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	db         *sql.DB
	redis      *redis.Client
	cfg        *config.Config
	httpClient *http.Client
}

func NewAuthHandler(db *sql.DB, redis *redis.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:    db,
		redis: redis,
		cfg:   cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Proxy to auth service with /api prefix
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	path := "/api" + strings.TrimPrefix(c.Path(), "/api/v1")
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	// Convert /api/v1/admin/auth/login to /api/admin/auth/login
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) CreateAdmin(c *fiber.Ctx) error {
	// Convert /api/v1/admin/auth/create to /api/admin/auth/create
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) ValidateAdmin(c *fiber.Ctx) error {
	// Convert /api/v1/admin/auth/validate to /api/admin/auth/validate
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) AdminStats(c *fiber.Ctx) error {
	// Convert /api/v1/admin/stats to /api/admin/stats
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) AdminActivities(c *fiber.Ctx) error {
	// Convert /api/v1/admin/activities to /api/admin/activities
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

func (h *AuthHandler) AdminHealth(c *fiber.Ctx) error {
	// Convert /api/v1/admin/health to /api/admin/health
	path := strings.Replace(c.Path(), "/api/v1/admin", "/api/admin", 1)
	return h.proxyToAuthService(c, path)
}

// proxyToAuthService proxies request to auth service with custom path
func (h *AuthHandler) proxyToAuthService(c *fiber.Ctx, customPath string) error {
	// Prepare the request URL with custom path
	url := h.cfg.AuthServiceURL + customPath
	if c.Request().URI().QueryString() != nil {
		url += "?" + string(c.Request().URI().QueryString())
	}

	// Create request body reader
	var bodyReader io.Reader
	if c.Body() != nil {
		bodyReader = bytes.NewReader(c.Body())
	}

	// Create HTTP request
	req, err := http.NewRequest(c.Method(), url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers from original request
	c.Request().Header.VisitAll(func(key, value []byte) {
		// Skip hop-by-hop headers
		keyStr := string(key)
		if !isHopByHopHeader(keyStr) {
			req.Header.Set(keyStr, string(value))
		}
	})

	// Make the request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to proxy request: %w", err)
	}
	defer resp.Body.Close()

	// Copy response status
	c.Status(resp.StatusCode)

	// Copy response headers
	for key, values := range resp.Header {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				c.Set(key, value)
			}
		}
	}

	// Copy response body
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

// isHopByHopHeader checks if header should be stripped when proxying
func isHopByHopHeader(header string) bool {
	hopByHop := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	header = strings.ToLower(header)
	for _, h := range hopByHop {
		if strings.ToLower(h) == header {
			return true
		}
	}
	return false
}
