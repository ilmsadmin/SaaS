package repositories

import (
	"fmt"
	"time"

	"zplus-saas/apps/backend/pos-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(tenantID string, product *models.Product) error {
	query := `
		INSERT INTO products (tenant_id, category_id, name, description, sku, barcode, price, cost_price, 
		                     stock_quantity, min_stock_level, max_stock_level, unit, weight, dimensions, 
		                     image_url, is_active)
		VALUES (:tenant_id, :category_id, :name, :description, :sku, :barcode, :price, :cost_price,
		        :stock_quantity, :min_stock_level, :max_stock_level, :unit, :weight, :dimensions,
		        :image_url, :is_active)
		RETURNING id, created_at, updated_at
	`

	product.TenantID = tenantID

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(product, product)
}

// GetProductByID retrieves a product by ID
func (r *ProductRepository) GetProductByID(tenantID string, id int) (*models.Product, error) {
	query := `
		SELECT p.*, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.tenant_id = $1 AND p.id = $2
	`

	var product models.Product
	err := r.db.Get(&product, query, tenantID, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductBySKU retrieves a product by SKU
func (r *ProductRepository) GetProductBySKU(tenantID, sku string) (*models.Product, error) {
	query := `
		SELECT p.*, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.tenant_id = $1 AND p.sku = $2
	`

	var product models.Product
	err := r.db.Get(&product, query, tenantID, sku)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProducts retrieves all products with filtering and pagination
func (r *ProductRepository) GetProducts(tenantID string, filters map[string]interface{}, limit, offset int) ([]models.Product, int, error) {
	whereClause := "WHERE p.tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	// Build dynamic WHERE clause
	if categoryID, ok := filters["category_id"]; ok {
		whereClause += fmt.Sprintf(" AND p.category_id = $%d", argIndex)
		args = append(args, categoryID)
		argIndex++
	}

	if search, ok := filters["search"]; ok {
		whereClause += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d OR p.description ILIKE $%d)", argIndex, argIndex, argIndex)
		searchPattern := fmt.Sprintf("%%%s%%", search)
		args = append(args, searchPattern)
		argIndex++
	}

	if isActive, ok := filters["is_active"]; ok {
		whereClause += fmt.Sprintf(" AND p.is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if lowStock, ok := filters["low_stock"]; ok && lowStock.(bool) {
		whereClause += " AND p.stock_quantity <= p.min_stock_level"
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT p.*, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var products []models.Product
	err = r.db.Select(&products, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateProduct updates a product
func (r *ProductRepository) UpdateProduct(tenantID string, id int, product *models.Product) error {
	query := `
		UPDATE products 
		SET category_id = :category_id, name = :name, description = :description, sku = :sku,
		    barcode = :barcode, price = :price, cost_price = :cost_price, stock_quantity = :stock_quantity,
		    min_stock_level = :min_stock_level, max_stock_level = :max_stock_level, unit = :unit,
		    weight = :weight, dimensions = :dimensions, image_url = :image_url, is_active = :is_active,
		    updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = :tenant_id AND id = :id
	`

	product.TenantID = tenantID
	product.ID = id

	_, err := r.db.NamedExec(query, product)
	return err
}

// DeleteProduct soft deletes a product
func (r *ProductRepository) DeleteProduct(tenantID string, id int) error {
	query := `UPDATE products SET is_active = false WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// UpdateStock updates product stock quantity
func (r *ProductRepository) UpdateStock(tenantID string, productID int, quantity int, operation string) error {
	var query string
	if operation == "increase" {
		query = `UPDATE products SET stock_quantity = stock_quantity + $3 WHERE tenant_id = $1 AND id = $2`
	} else {
		query = `UPDATE products SET stock_quantity = stock_quantity - $3 WHERE tenant_id = $1 AND id = $2`
	}

	_, err := r.db.Exec(query, tenantID, productID, quantity)
	return err
}

// GetLowStockProducts gets products with low stock
func (r *ProductRepository) GetLowStockProducts(tenantID string) ([]models.Product, error) {
	query := `
		SELECT p.*, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.tenant_id = $1 AND p.stock_quantity <= p.min_stock_level AND p.is_active = true
		ORDER BY p.stock_quantity ASC
	`

	var products []models.Product
	err := r.db.Select(&products, query, tenantID)
	return products, err
}

// CategoryRepository handles category operations
type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(tenantID string, category *models.Category) error {
	query := `
		INSERT INTO categories (tenant_id, name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(query, tenantID, category.Name, category.Description, category.IsActive).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

// GetCategories retrieves all categories
func (r *CategoryRepository) GetCategories(tenantID string) ([]models.Category, error) {
	query := `
		SELECT * FROM categories 
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	var categories []models.Category
	err := r.db.Select(&categories, query, tenantID)
	return categories, err
}

// GetCategoryByID retrieves a category by ID
func (r *CategoryRepository) GetCategoryByID(tenantID string, id int) (*models.Category, error) {
	query := `SELECT * FROM categories WHERE tenant_id = $1 AND id = $2`

	var category models.Category
	err := r.db.Get(&category, query, tenantID, id)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// UpdateCategory updates a category
func (r *CategoryRepository) UpdateCategory(tenantID string, id int, category *models.Category) error {
	query := `
		UPDATE categories 
		SET name = $3, description = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $1 AND id = $2
	`

	_, err := r.db.Exec(query, tenantID, id, category.Name, category.Description, category.IsActive)
	return err
}

// DeleteCategory soft deletes a category
func (r *CategoryRepository) DeleteCategory(tenantID string, id int) error {
	query := `UPDATE categories SET is_active = false WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.Exec(query, tenantID, id)
	return err
}

// OrderRepository handles order operations
type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order with items
func (r *OrderRepository) CreateOrder(tenantID string, order *models.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), order.ID)

	// Insert order
	orderQuery := `
		INSERT INTO orders (tenant_id, order_number, customer_name, customer_phone, customer_email,
		                   subtotal, tax_amount, discount_amount, total_amount, payment_method,
		                   payment_status, order_status, served_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(orderQuery, tenantID, orderNumber, order.CustomerName, order.CustomerPhone,
		order.CustomerEmail, order.Subtotal, order.TaxAmount, order.DiscountAmount,
		order.TotalAmount, order.PaymentMethod, "paid", "completed", order.ServedBy, order.Notes).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	order.OrderNumber = orderNumber

	// Insert order items and update stock
	for i := range order.Items {
		item := &order.Items[i]
		item.TenantID = tenantID
		item.OrderID = order.ID

		itemQuery := `
			INSERT INTO order_items (tenant_id, order_id, product_id, product_name, product_sku,
			                        quantity, unit_price, discount_amount, total_price)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at
		`

		err = tx.QueryRow(itemQuery, item.TenantID, item.OrderID, item.ProductID, item.ProductName,
			item.ProductSKU, item.Quantity, item.UnitPrice, item.DiscountAmount, item.TotalPrice).
			Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}

		// Update product stock
		stockQuery := `UPDATE products SET stock_quantity = stock_quantity - $3 WHERE tenant_id = $1 AND id = $2`
		_, err = tx.Exec(stockQuery, tenantID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}

		// Create inventory transaction
		invQuery := `
			INSERT INTO inventory_transactions (tenant_id, product_id, transaction_type, quantity,
			                                   reference_type, reference_id, created_by)
			VALUES ($1, $2, 'out', $3, 'sale', $4, $5)
		`
		_, err = tx.Exec(invQuery, tenantID, item.ProductID, item.Quantity, order.ID, order.ServedBy)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrderByID retrieves an order by ID with items
func (r *OrderRepository) GetOrderByID(tenantID string, id int) (*models.Order, error) {
	query := `SELECT * FROM orders WHERE tenant_id = $1 AND id = $2`

	var order models.Order
	err := r.db.Get(&order, query, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `SELECT * FROM order_items WHERE tenant_id = $1 AND order_id = $2 ORDER BY id`
	err = r.db.Select(&order.Items, itemsQuery, tenantID, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrders retrieves orders with filtering and pagination
func (r *OrderRepository) GetOrders(tenantID string, filters map[string]interface{}, limit, offset int) ([]models.Order, int, error) {
	whereClause := "WHERE tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	// Build filters
	if status, ok := filters["status"]; ok {
		whereClause += fmt.Sprintf(" AND order_status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if dateFrom, ok := filters["date_from"]; ok {
		whereClause += fmt.Sprintf(" AND DATE(created_at) >= $%d", argIndex)
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo, ok := filters["date_to"]; ok {
		whereClause += fmt.Sprintf(" AND DATE(created_at) <= $%d", argIndex)
		args = append(args, dateTo)
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT * FROM orders %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var orders []models.Order
	err = r.db.Select(&orders, query, args...)
	return orders, total, err
}

// UpdateOrderStatus updates order status
func (r *OrderRepository) UpdateOrderStatus(tenantID string, id int, status string) error {
	query := `
		UPDATE orders 
		SET order_status = $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`

	_, err := r.db.Exec(query, status, tenantID, id)
	return err
}

// UpdatePaymentStatus updates payment status
func (r *OrderRepository) UpdatePaymentStatus(tenantID string, id int, status string) error {
	query := `
		UPDATE orders 
		SET payment_status = $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`

	_, err := r.db.Exec(query, status, tenantID, id)
	return err
}

// GetSalesAnalytics retrieves sales analytics
func (r *OrderRepository) GetSalesAnalytics(tenantID string, dateFrom, dateTo time.Time) (*models.SalesAnalytics, error) {
	analytics := &models.SalesAnalytics{}

	// Get total sales and orders
	summaryQuery := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			COUNT(*) as total_orders,
			COALESCE(AVG(total_amount), 0) as average_order_value
		FROM orders
		WHERE tenant_id = $1 AND order_status = 'completed'
		AND created_at BETWEEN $2 AND $3
	`
	err := r.db.Get(analytics, summaryQuery, tenantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Sales by date
	salesByDateQuery := `
		SELECT 
			DATE(created_at) as date,
			SUM(total_amount) as total_sales,
			COUNT(*) as total_orders
		FROM orders
		WHERE tenant_id = $1 AND order_status = 'completed'
		AND DATE(created_at) BETWEEN $2 AND $3
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`
	err = r.db.Select(&analytics.SalesByDate, salesByDateQuery, tenantID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

// GetTopSellingProducts retrieves top selling products
func (r *OrderRepository) GetTopSellingProducts(tenantID string, limit int) ([]models.TopProduct, error) {
	query := `
		SELECT 
			oi.product_id,
			p.name as product_name,
			p.sku as product_sku,
			SUM(oi.quantity) as total_sold,
			SUM(oi.total_price) as total_revenue
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.tenant_id = $1 AND o.order_status = 'completed'
		GROUP BY oi.product_id, p.name, p.sku
		ORDER BY total_sold DESC
		LIMIT $2
	`

	var products []models.TopProduct
	err := r.db.Select(&products, query, tenantID, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// InventoryRepository handles inventory-related database operations
type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// CreateInventoryTransaction creates a new inventory transaction
func (r *InventoryRepository) CreateInventoryTransaction(tenantID string, transaction *models.InventoryTransaction) error {
	query := `
		INSERT INTO inventory_transactions (tenant_id, product_id, transaction_type, quantity, 
		                                   reference_id, reference_type, notes, user_id)
		VALUES (:tenant_id, :product_id, :transaction_type, :quantity, 
		        :reference_id, :reference_type, :notes, :user_id)
		RETURNING id, created_at
	`

	transaction.TenantID = tenantID

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(transaction, transaction)
}

// GetInventoryTransactions retrieves inventory transactions with filtering
func (r *InventoryRepository) GetInventoryTransactions(tenantID string, filters map[string]interface{}, limit, offset int) ([]models.InventoryTransaction, int, error) {
	whereClause := "WHERE it.tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	// Build dynamic WHERE clause
	if productID, ok := filters["product_id"]; ok {
		whereClause += fmt.Sprintf(" AND it.product_id = $%d", argIndex)
		args = append(args, productID)
		argIndex++
	}

	if transactionType, ok := filters["transaction_type"]; ok {
		whereClause += fmt.Sprintf(" AND it.transaction_type = $%d", argIndex)
		args = append(args, transactionType)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM inventory_transactions it
		LEFT JOIN products p ON it.product_id = p.id
		%s
	`, whereClause)

	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get transactions with pagination
	query := fmt.Sprintf(`
		SELECT 
			it.*,
			p.name as product_name,
			p.sku as product_sku
		FROM inventory_transactions it
		LEFT JOIN products p ON it.product_id = p.id
		%s
		ORDER BY it.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var transactions []models.InventoryTransaction
	err = r.db.Select(&transactions, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// UpdateProductStock updates product stock quantity
func (r *InventoryRepository) UpdateProductStock(tenantID string, productID int, newQuantity int) error {
	query := `
		UPDATE products 
		SET stock_quantity = $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`

	_, err := r.db.Exec(query, newQuantity, tenantID, productID)
	return err
}

// GetLowStockProducts retrieves products with low stock levels
func (r *InventoryRepository) GetLowStockProducts(tenantID string) ([]models.Product, error) {
	query := `
		SELECT p.*, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.tenant_id = $1 
		AND p.stock_quantity <= p.min_stock_level
		AND p.is_active = true
		ORDER BY p.stock_quantity ASC
	`

	var products []models.Product
	err := r.db.Select(&products, query, tenantID)
	if err != nil {
		return nil, err
	}

	return products, nil
}
