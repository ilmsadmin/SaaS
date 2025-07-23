# POS Service

Point of Sale (POS) Service for the Zplus SaaS Platform. This service handles product catalog management, order processing, inventory tracking, and sales analytics.

## Features

### 🛍️ Product Management
- Product catalog with categories
- SKU and barcode management
- Stock level tracking with low stock alerts
- Price and cost management
- Product search and filtering

### 📦 Order Management
- Sales order creation and processing
- Order item management
- Payment method tracking
- Order status management
- Customer information storage

### 📊 Inventory Tracking
- Real-time stock updates
- Inventory transaction logging
- Stock level alerts
- Purchase order management
- Supplier management

### 📈 Analytics & Reporting
- Sales analytics and dashboards
- Top selling products
- Sales by date/period
- Order statistics
- Product performance metrics

## API Endpoints

### Health Check
- `GET /health` - Service health check
- `GET /api/v1/pos/health` - POS service health check

### Categories
- `GET /api/v1/pos/categories` - List all categories
- `POST /api/v1/pos/categories` - Create new category
- `GET /api/v1/pos/categories/:id` - Get category by ID
- `PUT /api/v1/pos/categories/:id` - Update category
- `DELETE /api/v1/pos/categories/:id` - Delete category

### Products
- `GET /api/v1/pos/products` - List products with filtering
- `POST /api/v1/pos/products` - Create new product
- `GET /api/v1/pos/products/:id` - Get product by ID
- `PUT /api/v1/pos/products/:id` - Update product
- `DELETE /api/v1/pos/products/:id` - Delete product
- `GET /api/v1/pos/products/low-stock` - Get low stock products

### Orders
- `GET /api/v1/pos/orders` - List orders with filtering
- `POST /api/v1/pos/orders` - Create new order
- `GET /api/v1/pos/orders/:id` - Get order by ID

### Analytics
- `GET /api/v1/pos/analytics` - Get POS analytics
- `GET /api/v1/pos/analytics/dashboard` - Get dashboard analytics

### Inventory
- `GET /api/v1/pos/inventory/transactions` - Get inventory transactions
- `POST /api/v1/pos/inventory/adjust` - Adjust inventory levels

## Database Schema

The service uses PostgreSQL with the following main tables:

- `categories` - Product categories
- `products` - Product catalog
- `orders` - Sales orders
- `order_items` - Order line items
- `inventory_transactions` - Stock movements
- `suppliers` - Supplier information
- `purchase_orders` - Purchase orders from suppliers
- `purchase_order_items` - Purchase order line items

## Configuration

### Environment Variables

- `PORT` - Service port (default: 8084)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password

## Multi-tenant Support

The service supports multi-tenancy through:
- Tenant ID in all database operations
- Tenant-specific data isolation
- Header-based tenant identification (`X-Tenant-ID`)

## Running the Service

### Development
```bash
go run apps/backend/pos-service/cmd/main.go
```

### Docker
```bash
docker build -t pos-service -f apps/backend/pos-service/Dockerfile .
docker run -p 8084:8084 pos-service
```

### With Docker Compose
The service is included in the main docker-compose.yml file:
```bash
docker-compose up pos-service
```

## Request/Response Format

### Standard Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Paginated Response Format
```json
{
  "success": true,
  "message": "Results retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "code": 400
}
```

## Example Usage

### Create Product
```bash
curl -X POST http://localhost:8084/api/v1/pos/products \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant123" \
  -d '{
    "name": "Coffee Mug",
    "sku": "MUG001",
    "price": 15.99,
    "cost_price": 8.00,
    "stock_quantity": 100,
    "min_stock_level": 10,
    "category_id": 1
  }'
```

### Create Order
```bash
curl -X POST http://localhost:8084/api/v1/pos/orders \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant123" \
  -d '{
    "customer_name": "John Doe",
    "payment_method": "cash",
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

### Get Analytics
```bash
curl http://localhost:8084/api/v1/pos/analytics?date_from=2025-07-01&date_to=2025-07-23 \
  -H "X-Tenant-ID: tenant123"
```

## Development Status

- ✅ Database schema design
- ✅ Basic service structure
- ✅ API endpoint stubs
- ✅ Health check endpoints
- 🚧 Repository layer implementation
- 🚧 Service layer implementation
- 🚧 Handler implementation
- 🚧 Database integration
- 🚧 Authentication middleware
- ⏳ Testing suite
- ⏳ API documentation
- ⏳ Frontend integration

## Integration

The POS service integrates with:
- **API Gateway** - Routing and authentication
- **Auth Service** - User authentication and authorization
- **Tenant Service** - Multi-tenant support
- **File Service** - Product image management (future)
- **Payment Service** - Payment processing integration (future)

## License

Part of the Zplus SaaS Platform - Internal Development Project
