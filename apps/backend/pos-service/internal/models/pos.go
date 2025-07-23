package models

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID          int       `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Product represents a product in the catalog
type Product struct {
	ID            int       `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	CategoryID    *int      `json:"category_id" db:"category_id"`
	Name          string    `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description   string    `json:"description" db:"description"`
	SKU           string    `json:"sku" db:"sku" validate:"required,min=1,max=100"`
	Barcode       string    `json:"barcode" db:"barcode"`
	Price         float64   `json:"price" db:"price" validate:"gte=0"`
	CostPrice     float64   `json:"cost_price" db:"cost_price" validate:"gte=0"`
	StockQuantity int       `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel int       `json:"min_stock_level" db:"min_stock_level"`
	MaxStockLevel int       `json:"max_stock_level" db:"max_stock_level"`
	Unit          string    `json:"unit" db:"unit"`
	Weight        *float64  `json:"weight" db:"weight"`
	Dimensions    string    `json:"dimensions" db:"dimensions"`
	ImageURL      string    `json:"image_url" db:"image_url"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	CategoryName string `json:"category_name,omitempty" db:"category_name"`
}

// Order represents a sales order
type Order struct {
	ID             int       `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	OrderNumber    string    `json:"order_number" db:"order_number"`
	CustomerName   string    `json:"customer_name" db:"customer_name"`
	CustomerPhone  string    `json:"customer_phone" db:"customer_phone"`
	CustomerEmail  string    `json:"customer_email" db:"customer_email"`
	Subtotal       float64   `json:"subtotal" db:"subtotal"`
	TaxAmount      float64   `json:"tax_amount" db:"tax_amount"`
	DiscountAmount float64   `json:"discount_amount" db:"discount_amount"`
	TotalAmount    float64   `json:"total_amount" db:"total_amount" validate:"gte=0"`
	PaymentMethod  string    `json:"payment_method" db:"payment_method"`
	PaymentStatus  string    `json:"payment_status" db:"payment_status"`
	OrderStatus    string    `json:"order_status" db:"order_status"`
	ServedBy       string    `json:"served_by" db:"served_by"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	// Order items
	Items []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID             int       `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	OrderID        int       `json:"order_id" db:"order_id"`
	ProductID      int       `json:"product_id" db:"product_id"`
	ProductName    string    `json:"product_name" db:"product_name"`
	ProductSKU     string    `json:"product_sku" db:"product_sku"`
	Quantity       int       `json:"quantity" db:"quantity" validate:"gt=0"`
	UnitPrice      float64   `json:"unit_price" db:"unit_price" validate:"gte=0"`
	DiscountAmount float64   `json:"discount_amount" db:"discount_amount" validate:"gte=0"`
	TotalPrice     float64   `json:"total_price" db:"total_price" validate:"gte=0"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// InventoryTransaction represents stock movement
type InventoryTransaction struct {
	ID              int       `json:"id" db:"id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	ProductID       int       `json:"product_id" db:"product_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type" validate:"required,oneof=in out adjustment"`
	Quantity        int       `json:"quantity" db:"quantity" validate:"required"`
	UnitCost        *float64  `json:"unit_cost" db:"unit_cost"`
	TotalCost       *float64  `json:"total_cost" db:"total_cost"`
	ReferenceType   string    `json:"reference_type" db:"reference_type"`
	ReferenceID     *int      `json:"reference_id" db:"reference_id"`
	Notes           string    `json:"notes" db:"notes"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`

	// Joined fields
	ProductName string `json:"product_name,omitempty" db:"product_name"`
	ProductSKU  string `json:"product_sku,omitempty" db:"product_sku"`
}

// Supplier represents a supplier for inventory
type Supplier struct {
	ID            int       `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	Name          string    `json:"name" db:"name" validate:"required,min=2,max=255"`
	ContactPerson string    `json:"contact_person" db:"contact_person"`
	Email         string    `json:"email" db:"email" validate:"omitempty,email"`
	Phone         string    `json:"phone" db:"phone"`
	Address       string    `json:"address" db:"address"`
	PaymentTerms  string    `json:"payment_terms" db:"payment_terms"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// PurchaseOrder represents a purchase order to suppliers
type PurchaseOrder struct {
	ID           int        `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	SupplierID   *int       `json:"supplier_id" db:"supplier_id"`
	PONumber     string     `json:"po_number" db:"po_number"`
	OrderDate    time.Time  `json:"order_date" db:"order_date"`
	ExpectedDate *time.Time `json:"expected_date" db:"expected_date"`
	Status       string     `json:"status" db:"status"`
	Subtotal     float64    `json:"subtotal" db:"subtotal"`
	TaxAmount    float64    `json:"tax_amount" db:"tax_amount"`
	TotalAmount  float64    `json:"total_amount" db:"total_amount"`
	Notes        string     `json:"notes" db:"notes"`
	CreatedBy    string     `json:"created_by" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Joined fields
	SupplierName string              `json:"supplier_name,omitempty" db:"supplier_name"`
	Items        []PurchaseOrderItem `json:"items,omitempty"`
}

// PurchaseOrderItem represents an item in a purchase order
type PurchaseOrderItem struct {
	ID               int       `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	PurchaseOrderID  int       `json:"purchase_order_id" db:"purchase_order_id"`
	ProductID        int       `json:"product_id" db:"product_id"`
	Quantity         int       `json:"quantity" db:"quantity" validate:"gt=0"`
	UnitCost         float64   `json:"unit_cost" db:"unit_cost" validate:"gte=0"`
	TotalCost        float64   `json:"total_cost" db:"total_cost" validate:"gte=0"`
	ReceivedQuantity int       `json:"received_quantity" db:"received_quantity"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`

	// Joined fields
	ProductName string `json:"product_name,omitempty" db:"product_name"`
	ProductSKU  string `json:"product_sku,omitempty" db:"product_sku"`
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	CustomerName   string            `json:"customer_name"`
	CustomerPhone  string            `json:"customer_phone"`
	CustomerEmail  string            `json:"customer_email"`
	PaymentMethod  string            `json:"payment_method" validate:"required,oneof=cash card bank_transfer"`
	DiscountAmount float64           `json:"discount_amount" validate:"gte=0"`
	TaxAmount      float64           `json:"tax_amount" validate:"gte=0"`
	Notes          string            `json:"notes"`
	Items          []CreateOrderItem `json:"items" validate:"required,min=1"`
}

// CreateOrderItem represents an item in the create order request
type CreateOrderItem struct {
	ProductID      int     `json:"product_id" validate:"required,gt=0"`
	Quantity       int     `json:"quantity" validate:"required,gt=0"`
	DiscountAmount float64 `json:"discount_amount" validate:"gte=0"`
}

// POSAnalytics represents POS analytics data
type POSAnalytics struct {
	TotalSales         float64       `json:"total_sales"`
	TotalOrders        int           `json:"total_orders"`
	AverageOrderValue  float64       `json:"average_order_value"`
	TotalProducts      int           `json:"total_products"`
	LowStockProducts   int           `json:"low_stock_products"`
	TopSellingProducts []TopProduct  `json:"top_selling_products"`
	SalesByDate        []SalesByDate `json:"sales_by_date"`
}

// SalesAnalytics represents sales analytics data
type SalesAnalytics struct {
	TotalSales        float64       `json:"total_sales"`
	TotalOrders       int           `json:"total_orders"`
	AverageOrderValue float64       `json:"average_order_value"`
	TopProducts       []TopProduct  `json:"top_products"`
	RecentOrders      []Order       `json:"recent_orders"`
	SalesByDate       []SalesByDate `json:"sales_by_date"`
}

// TopProduct represents top selling product data
type TopProduct struct {
	ProductID    int     `json:"product_id" db:"product_id"`
	ProductName  string  `json:"product_name" db:"product_name"`
	ProductSKU   string  `json:"product_sku" db:"product_sku"`
	TotalSold    int     `json:"total_sold" db:"total_sold"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
}

// SalesByDate represents sales data by date
type SalesByDate struct {
	Date        string  `json:"date" db:"date"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
	TotalOrders int     `json:"total_orders" db:"total_orders"`
}
