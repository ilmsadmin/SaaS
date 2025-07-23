# Payment Service

Payment processing service for the Zplus SaaS Platform. This service handles payment transactions, refunds, payment methods, and financial operations.

## Features

### 💳 Payment Processing
- Multiple payment methods (credit/debit cards, bank transfer, digital wallets, cash)
- Real-time payment processing with status tracking
- Automatic payment gateway integration
- Payment retry mechanisms
- Multi-currency support

### 🔒 Security Features
- PCI DSS compliant data handling
- Secure payment tokenization
- Encrypted payment data storage
- Gateway transaction validation
- Fraud detection capabilities

### 📊 Financial Management
- Payment analytics and reporting
- Revenue tracking and insights
- Success rate monitoring
- Payment method analysis
- Transaction history management

### 🔄 Refund Management
- Full and partial refunds
- Refund status tracking
- Automatic gateway refund processing
- Refund analytics and reporting

### 🏪 Multi-Gateway Support
- Configurable payment gateways per tenant
- Gateway failover and load balancing
- Webhook handling for gateway notifications
- Gateway-specific features and configurations

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Payments
- `POST /api/v1/payment` - Create new payment
- `GET /api/v1/payment` - List payments with filtering
- `GET /api/v1/payment/:id` - Get payment by ID
- `GET /api/v1/payment/stats` - Get payment statistics

### Request/Response Format

#### Create Payment
```json
POST /api/v1/payment
{
  "order_id": "ORDER_123",
  "customer_id": 456,
  "customer_name": "John Doe",
  "amount": 99.99,
  "currency": "USD",
  "payment_method": "credit_card",
  "gateway": "stripe",
  "description": "Product purchase",
  "metadata": "{\"product_id\": 789}"
}
```

#### Response
```json
{
  "error": false,
  "message": "Payment created successfully",
  "data": {
    "id": 1,
    "tenant_id": "demo-tenant",
    "order_id": "ORDER_123",
    "customer_id": 456,
    "customer_name": "John Doe",
    "amount": 99.99,
    "currency": "USD",
    "payment_method": "credit_card",
    "status": "pending",
    "gateway": "stripe",
    "description": "Product purchase",
    "metadata": "{\"product_id\": 789}",
    "created_at": "2025-07-23T09:00:00Z",
    "updated_at": "2025-07-23T09:00:00Z"
  }
}
```

## Database Schema

### payments
- Core payment records with amount, status, and gateway information
- Multi-tenant support with tenant_id isolation
- Payment method and currency tracking
- Gateway transaction ID mapping

### payment_methods
- Stored payment methods for customers
- Tokenized card information (PCI compliant)
- Default payment method management
- Support for multiple payment types

### refunds
- Refund transactions linked to original payments
- Partial and full refund support
- Refund status tracking and processing
- Gateway refund ID mapping

### payment_logs
- Comprehensive audit trail for all payment events
- Status change tracking
- Gateway response logging
- Error message capture

### payment_gateways
- Configurable payment gateway settings per tenant
- API key and configuration management
- Supported payment methods per gateway
- Gateway-specific webhook URLs

## Configuration

### Environment Variables
- `PAYMENT_SERVICE_PORT` - Service port (default: 8087)
- `DATABASE_HOST` - PostgreSQL host
- `DATABASE_PORT` - PostgreSQL port
- `DATABASE_USER` - Database user
- `DATABASE_PASSWORD` - Database password
- `DATABASE_NAME` - Database name

### Payment Methods Supported
- `credit_card` - Credit card payments
- `debit_card` - Debit card payments
- `bank_transfer` - Bank transfer/ACH
- `digital_wallet` - PayPal, Apple Pay, Google Pay
- `cash` - Cash payments (for POS systems)

### Payment Statuses
- `pending` - Payment initiated, awaiting processing
- `processing` - Payment being processed by gateway
- `completed` - Payment successfully completed
- `failed` - Payment failed or declined
- `cancelled` - Payment cancelled by user
- `refunded` - Payment refunded

## Multi-tenant Support

- All data isolated by `tenant_id`
- Tenant-specific payment gateway configurations
- Per-tenant payment method preferences
- Isolated financial reporting and analytics

## Running the Service

### Prerequisites
- Go 1.21+
- PostgreSQL database
- Required environment variables

### Database Setup
```sql
-- Run the migration file
psql -d zplus_saas -f migrations/001_create_payment_tables.sql
```

### Start Service
```bash
# Set environment variables
export PAYMENT_SERVICE_PORT=8087
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres123
export DATABASE_NAME=zplus_saas

# Run the service
go run cmd/main.go
```

### Docker
```bash
docker build -t payment-service .
docker run -p 8087:8087 --env-file .env payment-service
```

## Request/Response Format

### Standard Response Format
```json
{
  "error": false,
  "message": "Operation successful",
  "data": {...}
}
```

### Error Response Format
```json
{
  "error": true,
  "message": "Error description"
}
```

## Example Usage

### Process Payment
```bash
curl -X POST "http://localhost:8087/api/v1/payment" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo-tenant" \
  -d '{
    "order_id": "ORDER_123",
    "customer_name": "John Doe",
    "amount": 99.99,
    "currency": "USD",
    "payment_method": "credit_card",
    "gateway": "stripe",
    "description": "Product purchase"
  }'
```

### Get Payment Statistics
```bash
curl -X GET "http://localhost:8087/api/v1/payment/stats" \
  -H "X-Tenant-ID: demo-tenant"
```

### Health Check
```bash
curl http://localhost:8087/health
```

## Development Status

- ✅ Database schema design
- ✅ Basic service structure
- ✅ Core payment processing
- ✅ Payment status tracking
- ✅ Multi-tenant support
- ✅ Payment analytics
- ✅ Health check endpoint
- 🚧 Gateway integrations (Stripe, PayPal)
- 🚧 Refund processing
- 🚧 Payment method management
- ⏳ Webhook handling
- ⏳ PCI compliance features
- ⏳ Fraud detection
- ⏳ Recurring payments
- ⏳ Payment plans and subscriptions

## Integration

The Payment service integrates with:
- **API Gateway** - Routing and authentication
- **Auth Service** - User authentication and authorization
- **POS Service** - Point of sale payment processing
- **CRM Service** - Customer payment history
- **Tenant Service** - Multi-tenant configuration
- **External Gateways** - Stripe, PayPal, Square (future)

## Business Logic

### Payment Processing Flow
1. Payment request validation
2. Customer and order verification
3. Payment method validation
4. Gateway processing (simulated)
5. Status updates and notifications
6. Transaction logging and audit trail

### Payment Simulation
- 90% success rate for testing
- 3-second processing delay
- Random gateway transaction ID generation
- Automatic status updates

## Security Features

- Multi-tenant data isolation
- Input validation and sanitization
- SQL injection prevention
- Payment data encryption (planned)
- PCI DSS compliance measures (planned)
- Secure API endpoints

## Performance Features

- Database indexing for payment queries
- Optimized transaction lookups
- Efficient status filtering
- Payment analytics caching (planned)
- Background payment processing

## License

Part of the Zplus SaaS Platform - Internal Development Project
