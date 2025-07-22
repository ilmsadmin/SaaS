package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HTTPClient wraps http.Client with utility methods
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with timeout
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// ProxyRequest forwards a Fiber request to another service
func (h *HTTPClient) ProxyRequest(c *fiber.Ctx, targetURL string) error {
	// Prepare the request URL
	url := targetURL + c.Path()
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
	resp, err := h.client.Do(req)
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

// Forward is a simpler version that just forwards and returns response
func (h *HTTPClient) Forward(c *fiber.Ctx, targetURL string) error {
	return h.ProxyRequest(c, targetURL)
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
