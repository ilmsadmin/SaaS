package services

import (
	"errors"
	"fmt"
	"time"

	"zplus-saas/apps/backend/pos-service/internal/models"
	"zplus-saas/apps/backend/pos-service/internal/repositories"
)

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	CustomerName   string            `json:"customer_name"`
	CustomerPhone  string            `json:"customer_phone"`
	CustomerEmail  string            `json:"customer_email"`
	PaymentMethod  string            `json:"payment_method"`
	DiscountAmount float64           `json:"discount_amount"`
	TaxAmount      float64           `json:"tax_amount"`
	Notes          string            `json:"notes"`
	Items          []CreateOrderItem `json:"items"`
}

// CreateOrderItem represents an item in the create order request
type CreateOrderItem struct {
	ProductID      int     `json:"product_id"`
	Quantity       int     `json:"quantity"`
	DiscountAmount float64 `json:"discount_amount"`
}

// POSService handles POS business logic
type POSService struct {
	productRepo   *repositories.ProductRepository
	categoryRepo  *repositories.CategoryRepository
	orderRepo     *repositories.OrderRepository
	inventoryRepo *repositories.InventoryRepository
}

// NewPOSService creates a new POS service
func NewPOSService(productRepo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository, orderRepo *repositories.OrderRepository, inventoryRepo *repositories.InventoryRepository) *POSService {
	return &POSService{
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreateProduct creates a new product
func (s *POSService) CreateProduct(tenantID string, product *models.Product) error {
	// Validate required fields
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.SKU == "" {
		return errors.New("product SKU is required")
	}
	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	// Set defaults
	if product.Unit == "" {
		product.Unit = "pcs"
	}
	product.IsActive = true

	return s.productRepo.CreateProduct(tenantID, product)
}

// GetProducts retrieves products with filtering
func (s *POSService) GetProducts(tenantID string, filters map[string]interface{}, page, limit int) ([]models.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.productRepo.GetProducts(tenantID, filters, limit, offset)
}

// GetProductByID retrieves a product by ID
func (s *POSService) GetProductByID(tenantID string, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	return s.productRepo.GetProductByID(tenantID, id)
}

// UpdateProduct updates an existing product
func (s *POSService) UpdateProduct(tenantID string, id int, product *models.Product) error {
	// Validate required fields
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.SKU == "" {
		return errors.New("product SKU is required")
	}
	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	return s.productRepo.UpdateProduct(tenantID, id, product)
}

// DeleteProduct marks a product as inactive (soft delete)
func (s *POSService) DeleteProduct(tenantID string, id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}

	return s.productRepo.DeleteProduct(tenantID, id)
}

// Category methods

// CreateCategory creates a new category
func (s *POSService) CreateCategory(tenantID string, category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	category.IsActive = true
	return s.categoryRepo.CreateCategory(tenantID, category)
}

// GetCategories retrieves all categories
func (s *POSService) GetCategories(tenantID string) ([]models.Category, error) {
	return s.categoryRepo.GetCategories(tenantID)
}

// Order methods

// CreateOrder creates a new order
func (s *POSService) CreateOrder(tenantID, userID string, req *CreateOrderRequest) (*models.Order, error) {
	// Validate request
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), time.Now().Nanosecond()/1000)

	// Calculate totals
	var subtotal float64
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than 0")
		}

		// Get product details
		product, err := s.productRepo.GetProductByID(tenantID, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID %d: %w", item.ProductID, err)
		}

		// Check stock
		if product.StockQuantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		itemTotal := (product.Price * float64(item.Quantity)) - item.DiscountAmount
		subtotal += itemTotal

		orderItem := models.OrderItem{
			TenantID:       tenantID,
			ProductID:      item.ProductID,
			ProductName:    product.Name,
			ProductSKU:     product.SKU,
			Quantity:       item.Quantity,
			UnitPrice:      product.Price,
			DiscountAmount: item.DiscountAmount,
			TotalPrice:     itemTotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	totalAmount := subtotal + req.TaxAmount - req.DiscountAmount

	// Create order
	order := &models.Order{
		TenantID:       tenantID,
		OrderNumber:    orderNumber,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		CustomerEmail:  req.CustomerEmail,
		Subtotal:       subtotal,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    totalAmount,
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  "pending",
		OrderStatus:    "pending",
		ServedBy:       userID,
		Notes:          req.Notes,
		Items:          orderItems,
	}

	// Create order with items
	err := s.orderRepo.CreateOrder(tenantID, order)
	if err != nil {
		return nil, err
	}

	// Update product stock
	for _, item := range req.Items {
		product, _ := s.productRepo.GetProductByID(tenantID, item.ProductID)
		newStock := product.StockQuantity - item.Quantity
		s.inventoryRepo.UpdateProductStock(tenantID, item.ProductID, newStock)

		// Create inventory transaction
		transaction := &models.InventoryTransaction{
			TenantID:        tenantID,
			ProductID:       item.ProductID,
			TransactionType: "sale",
			Quantity:        -item.Quantity,
			ReferenceID:     &order.ID,
			ReferenceType:   "order",
			Notes:           fmt.Sprintf("Sale - Order #%s", order.OrderNumber),
			UserID:          userID,
		}
		s.inventoryRepo.CreateInventoryTransaction(tenantID, transaction)
	}

	return order, nil
}

// GetOrders retrieves orders with filtering
func (s *POSService) GetOrders(tenantID string, filters map[string]interface{}, page, limit int) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.orderRepo.GetOrders(tenantID, filters, limit, offset)
}

// GetOrderByID retrieves an order by ID
func (s *POSService) GetOrderByID(tenantID string, id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}

	return s.orderRepo.GetOrderByID(tenantID, id)
}

// UpdateOrderStatus updates order status
func (s *POSService) UpdateOrderStatus(tenantID string, id int, status string) error {
	if id <= 0 {
		return errors.New("invalid order ID")
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"preparing": true,
		"ready":     true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid order status")
	}

	return s.orderRepo.UpdateOrderStatus(tenantID, id, status)
}

// UpdatePaymentStatus updates payment status
func (s *POSService) UpdatePaymentStatus(tenantID string, id int, status string) error {
	if id <= 0 {
		return errors.New("invalid order ID")
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"partial":   true,
		"refunded":  true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid payment status")
	}

	return s.orderRepo.UpdatePaymentStatus(tenantID, id, status)
}

// Analytics methods

// GetDashboardAnalytics retrieves dashboard analytics
func (s *POSService) GetDashboardAnalytics(tenantID string, dateFrom, dateTo time.Time) (*models.SalesAnalytics, error) {
	return s.orderRepo.GetSalesAnalytics(tenantID, dateFrom, dateTo)
}

// GetTopSellingProducts retrieves top selling products
func (s *POSService) GetTopSellingProducts(tenantID string, limit int) ([]models.TopProduct, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.orderRepo.GetTopSellingProducts(tenantID, limit)
}

// Inventory methods

// GetLowStockProducts retrieves products with low stock
func (s *POSService) GetLowStockProducts(tenantID string) ([]models.Product, error) {
	return s.inventoryRepo.GetLowStockProducts(tenantID)
}

// GetInventoryTransactions retrieves inventory transactions
func (s *POSService) GetInventoryTransactions(tenantID string, filters map[string]interface{}, page, limit int) ([]models.InventoryTransaction, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.inventoryRepo.GetInventoryTransactions(tenantID, filters, limit, offset)
}

// CreateInventoryAdjustment creates inventory adjustment
func (s *POSService) CreateInventoryAdjustment(tenantID, userID string, productID int, quantity int, reason string) error {
	if productID <= 0 {
		return errors.New("invalid product ID")
	}

	if quantity == 0 {
		return errors.New("quantity cannot be zero")
	}

	// Get current product stock
	product, err := s.productRepo.GetProductByID(tenantID, productID)
	if err != nil {
		return err
	}

	// Update stock
	newStock := product.StockQuantity + quantity
	if newStock < 0 {
		return errors.New("adjustment would result in negative stock")
	}

	err = s.inventoryRepo.UpdateProductStock(tenantID, productID, newStock)
	if err != nil {
		return err
	}

	// Create inventory transaction
	transactionType := "adjustment_in"
	if quantity < 0 {
		transactionType = "adjustment_out"
	}

	transaction := &models.InventoryTransaction{
		TenantID:        tenantID,
		ProductID:       productID,
		TransactionType: transactionType,
		Quantity:        quantity,
		Notes:           reason,
		UserID:          userID,
	}

	return s.inventoryRepo.CreateInventoryTransaction(tenantID, transaction)
}
