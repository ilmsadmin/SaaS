# Checkin Service

Attendance and checkin management service for the Zplus SaaS Platform. This service handles employee checkin/checkout, attendance tracking, and time management features.

## Features

### 🕒 Checkin Management
- Employee checkin/checkout tracking
- Break time management (break start/end)
- Real-time attendance recording
- Location-based checkin (GPS coordinates)
- Photo verification support
- Device and IP tracking

### 📊 Attendance Tracking  
- Daily attendance summaries
- Work hours calculation
- Break time tracking
- Overtime calculation
- Attendance statistics and analytics
- Attendance rate monitoring

### 📝 Record Management
- Checkin record CRUD operations
- Attendance approval workflow
- Record validation and verification
- Historical data access
- Multi-tenant data isolation

### 🎯 Business Logic
- Checkin sequence validation
- Work time calculations
- Late arrival detection
- Early departure tracking
- Attendance policy enforcement

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Checkin Records
- `POST /api/v1/checkin` - Create new checkin record
- `GET /api/v1/checkin` - Get checkin records (with filtering)
- `GET /api/v1/checkin/employee/:employee_id/today` - Get today's records for employee

### Request/Response Format

#### Create Checkin
```json
POST /api/v1/checkin
{
  "employee_id": 123,
  "checkin_type": "checkin", // "checkin", "checkout", "break_start", "break_end"
  "location": "Office Main Building",
  "latitude": 10.762622,
  "longitude": 106.660172,
  "photo": "base64_encoded_photo_or_url",
  "notes": "Working from office today"
}
```

#### Response
```json
{
  "error": false,
  "message": "Checkin recorded successfully",
  "data": {
    "id": 1,
    "tenant_id": "demo-tenant",
    "employee_id": 123,
    "employee_name": "Employee_123",
    "checkin_type": "checkin",
    "timestamp": "2025-07-23T09:00:00Z",
    "location": "Office Main Building",
    "latitude": 10.762622,
    "longitude": 106.660172,
    "ip_address": "192.168.1.100",
    "device_info": "Mozilla/5.0...",
    "photo": "base64_encoded_photo_or_url",
    "notes": "Working from office today",
    "status": "approved",
    "created_at": "2025-07-23T09:00:00Z",
    "updated_at": "2025-07-23T09:00:00Z"
  }
}
```

## Database Schema

### checkin_records
- `id` - Primary key
- `tenant_id` - Tenant identifier for multi-tenancy
- `employee_id` - Employee identifier
- `employee_name` - Employee display name
- `checkin_type` - Type of checkin (checkin/checkout/break_start/break_end)
- `timestamp` - Checkin/checkout timestamp
- `location` - Location description
- `latitude/longitude` - GPS coordinates
- `ip_address` - Client IP address
- `device_info` - User agent/device information
- `photo` - Photo verification (URL or base64)
- `notes` - Additional notes
- `status` - Approval status (approved/pending/rejected)
- `approved_by/approved_at` - Approval information
- `created_at/updated_at` - Timestamps

### attendance_summary
- Daily attendance summaries with work hours calculation
- Break time tracking and overtime calculation
- Attendance status (present/absent/late/early_leave/partial)

### attendance_policies  
- Configurable attendance policies per tenant
- Work schedule and time thresholds
- Photo and location requirements
- Work days configuration

## Configuration

### Environment Variables
- `CHECKIN_SERVICE_PORT` - Service port (default: 8086)
- `DATABASE_HOST` - PostgreSQL host
- `DATABASE_PORT` - PostgreSQL port
- `DATABASE_USER` - Database user
- `DATABASE_PASSWORD` - Database password
- `DATABASE_NAME` - Database name

### Default Configuration
- Port: 8086
- Database: PostgreSQL
- Auto-approval: Enabled
- Tenant support: Multi-tenant with tenant_id

## Multi-tenant Support

- All data is isolated by `tenant_id`
- Tenant identification via `X-Tenant-ID` header
- Default tenant: `demo-tenant` for development
- Database-level tenant separation

## Running the Service

### Prerequisites
- Go 1.21+
- PostgreSQL database
- Required environment variables

### Database Setup
```sql
-- Run the migration file
psql -d zplus_saas -f migrations/001_create_checkin_tables.sql
```

### Start Service
```bash
# Set environment variables
export CHECKIN_SERVICE_PORT=8086
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_NAME=zplus_saas

# Run the service
go run cmd/main.go
```

### Docker
```bash
docker build -t checkin-service .
docker run -p 8086:8086 --env-file .env checkin-service
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

### Employee Checkin Flow
1. Employee checks in: `POST /api/v1/checkin` with `checkin_type: "checkin"`
2. Employee takes break: `POST /api/v1/checkin` with `checkin_type: "break_start"`
3. Employee returns from break: `POST /api/v1/checkin` with `checkin_type: "break_end"`
4. Employee checks out: `POST /api/v1/checkin` with `checkin_type: "checkout"`

### Get Today's Records
```bash
curl -X GET "http://localhost:8086/api/v1/checkin/employee/123/today" \
  -H "X-Tenant-ID: demo-tenant"
```

### Health Check
```bash
curl http://localhost:8086/health
```

## Development Status

- ✅ Database schema design
- ✅ Basic service structure  
- ✅ Core API endpoints
- ✅ Health check endpoint
- ✅ Database integration
- ✅ Multi-tenant support
- ✅ Request validation
- ✅ Location and photo support
- 🚧 Advanced attendance policies
- 🚧 Approval workflow
- ⏳ Attendance analytics
- ⏳ Policy enforcement
- ⏳ Integration with HRM service
- ⏳ Real-time notifications
- ⏳ Mobile app support

## Integration

The Checkin service integrates with:
- **API Gateway** - Routing and authentication
- **Auth Service** - User authentication and authorization  
- **HRM Service** - Employee information and management
- **Tenant Service** - Multi-tenant configuration
- **Notification Service** - Attendance alerts and reminders (future)

## Business Logic

### Checkin Sequence Validation
- First action of the day must be checkin
- Cannot checkin twice without checkout
- Cannot checkout without checkin first
- Break start/end must follow proper sequence

### Attendance Calculation
- Work hours = checkout time - checkin time - break hours
- Break hours calculated from break_start/break_end pairs
- Late detection based on work start time + threshold
- Attendance rate = present days / total days

## Security Features

- Multi-tenant data isolation
- IP address tracking
- Device fingerprinting
- Photo verification support
- Input validation and sanitization
- SQL injection prevention

## Performance Features

- Database indexing for common queries
- Optimized queries for attendance stats
- Pagination support for large datasets
- Efficient date-based filtering

## License

Part of the Zplus SaaS Platform - Internal Development Project
