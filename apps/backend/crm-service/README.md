# CRM Service Documentation

## Overview
The CRM (Customer Relationship Management) Service is a core microservice in the Zplus SaaS platform that handles customer data, lead management, and sales pipeline operations.

## Features

### 🏢 Customer Management
- Complete CRUD operations for customers
- Customer search and filtering
- Customer status tracking (active, inactive, prospect)
- Customer analytics and statistics
- Multi-source tracking (website, referral, social, etc.)
- Custom tags and notes system

### 📈 Lead Management
- Lead creation and tracking
- Lead scoring system (0-100)
- Lead status management (new, qualified, contacted, converted, lost)
- Lead assignment to users
- Lead to customer conversion
- Lead analytics and conversion rates

### 💼 Opportunity Management
- Sales pipeline management
- Opportunity stages: prospecting → qualification → proposal → negotiation → closed
- Opportunity value and currency tracking
- Probability-based forecasting
- Win/loss ratio analytics
- Customer-specific opportunities

### 📊 Contact Activities (Future Enhancement)
- Call logging
- Email tracking
- Meeting scheduling
- Task management
- Activity history

## Technical Architecture

### Database Schema
```sql
-- Main Tables
- customers: Customer information and contact details
- leads: Sales leads and tracking
- opportunities: Sales opportunities and pipeline
- contact_activities: Customer interaction history

-- Indexes for Performance
- Tenant-based partitioning
- Email and search optimization
- Time-based queries
```

### API Endpoints

#### Customers
```
POST   /api/v1/customers           # Create customer
GET    /api/v1/customers           # List customers (paginated)
GET    /api/v1/customers/search    # Search customers
GET    /api/v1/customers/stats     # Customer statistics
GET    /api/v1/customers/:id       # Get customer by ID
PUT    /api/v1/customers/:id       # Update customer
DELETE /api/v1/customers/:id       # Delete customer
```

#### Leads
```
POST   /api/v1/leads               # Create lead
GET    /api/v1/leads               # List leads (paginated)
GET    /api/v1/leads/stats         # Lead statistics
GET    /api/v1/leads/status/:status # Get leads by status
GET    /api/v1/leads/:id           # Get lead by ID
PUT    /api/v1/leads/:id           # Update lead
DELETE /api/v1/leads/:id           # Delete lead
POST   /api/v1/leads/:id/convert   # Convert lead to customer
POST   /api/v1/leads/:id/score     # Update lead score
```

#### Opportunities
```
POST   /api/v1/opportunities              # Create opportunity
GET    /api/v1/opportunities              # List opportunities (paginated)
GET    /api/v1/opportunities/stats        # Opportunity statistics
GET    /api/v1/opportunities/pipeline     # Sales pipeline data
GET    /api/v1/opportunities/stage/:stage # Get opportunities by stage
GET    /api/v1/opportunities/customer/:id # Get opportunities by customer
GET    /api/v1/opportunities/:id          # Get opportunity by ID
PUT    /api/v1/opportunities/:id          # Update opportunity
DELETE /api/v1/opportunities/:id          # Delete opportunity
POST   /api/v1/opportunities/:id/close-won  # Mark as won
POST   /api/v1/opportunities/:id/close-lost # Mark as lost
```

## Configuration

### Environment Variables
```bash
PORT=8082                              # Service port
DATABASE_URL=postgres://...            # PostgreSQL connection
REDIS_URL=redis://localhost:6379       # Redis connection
JWT_SECRET=your-secret-key             # JWT secret
ENVIRONMENT=development                 # Environment
DEBUG=true                             # Debug mode
```

### Dependencies
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Language**: Go 1.24+
- **Framework**: Fiber v2
- **Authentication**: JWT

## Development

### Setup
```bash
# 1. Start infrastructure
docker-compose up -d

# 2. Run database migration
psql $DATABASE_URL -f migrations/001_create_crm_tables.sql

# 3. Start CRM service
./start-crm.sh
# OR
go run cmd/main.go
```

### Testing
```bash
# Health check
curl http://localhost:8082/health

# Create a customer
curl -X POST http://localhost:8082/api/v1/customers \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant-1" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Acme Corp"
  }'

# Get customers
curl http://localhost:8082/api/v1/customers \
  -H "X-Tenant-ID: tenant-1"
```

### Docker Deployment
```bash
# Build image
docker build -t zplus/crm-service .

# Run container
docker run -p 8082:8082 \
  -e DATABASE_URL="postgres://..." \
  -e REDIS_URL="redis://..." \
  zplus/crm-service
```

## API Documentation

### Authentication
All endpoints require:
- **X-Tenant-ID** header for tenant isolation
- **Authorization** header with Bearer token (when auth is enabled)

### Request/Response Format
```json
// Success Response
{
  "error": false,
  "message": "Operation successful",
  "data": {...},
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}

// Error Response
{
  "error": true,
  "message": "Error description"
}
```

### Sample Requests

#### Create Customer
```json
POST /api/v1/customers
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+1-555-0123",
  "company": "Acme Corp",
  "address": "123 Main St",
  "city": "New York",
  "state": "NY",
  "country": "USA",
  "zip_code": "10001",
  "source": "website",
  "tags": ["enterprise", "vip"],
  "notes": "Important client from website signup"
}
```

#### Create Lead
```json
POST /api/v1/leads
{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "phone": "+1-555-0124",
  "company": "Tech Corp",
  "title": "Marketing Director",
  "source": "referral",
  "assigned_to": 1,
  "value": 25000.00,
  "notes": "Interested in enterprise package"
}
```

#### Create Opportunity
```json
POST /api/v1/opportunities
{
  "customer_id": 1,
  "name": "Acme Corp Expansion",
  "description": "Expand services to include additional modules",
  "value": 50000.00,
  "currency": "USD",
  "source": "existing-client",
  "assigned_to": 1,
  "expected_date": "2025-08-15"
}
```

## Monitoring

### Health Check
```bash
GET /health
# Returns: {"status": "healthy", "service": "crm-service", "timestamp": 1690128000}
```

### Metrics
- Customer count
- Lead conversion rates
- Opportunity win rates
- Pipeline value
- Response times

## Security

### Multi-tenancy
- All data is isolated by tenant_id
- Database-level tenant separation
- API-level tenant validation

### Data Protection
- Input validation and sanitization
- SQL injection prevention
- XSS protection
- Rate limiting (planned)

## Roadmap

### Phase 2 Enhancements
- [ ] Contact activity tracking
- [ ] Email integration
- [ ] Calendar synchronization
- [ ] Advanced reporting
- [ ] Data import/export
- [ ] Mobile API optimization

### Phase 3 Features
- [ ] AI-powered lead scoring
- [ ] Predictive analytics
- [ ] Marketing automation
- [ ] Social media integration
- [ ] Advanced workflow automation

## Support

For technical support or questions:
- Check the main project README
- Review API documentation
- Submit issues via GitHub
- Contact the development team

---

**Service Status**: ✅ **PRODUCTION READY**  
**Last Updated**: July 23, 2025  
**Version**: 1.0.0
