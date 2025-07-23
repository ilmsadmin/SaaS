# HRM Service Documentation

## Overview

The HRM (Human Resource Management) Service is a microservice that handles all employee-related operations including employee management, department management, leave management, and performance tracking.

## Features

### 📋 Employee Management
- Complete CRUD operations for employees
- Employee search and filtering
- Employee hierarchy management (manager-subordinate relationships)
- Employee status tracking (active, inactive, terminated)
- Emergency contact management
- Employee statistics and analytics

### 🏢 Department Management
- Department creation and management
- Department budget tracking
- Manager assignment to departments
- Employee count per department
- Department location management

### 🏖️ Leave Management
- Leave request creation and management
- Multiple leave types (annual, sick, maternity, paternity, personal, emergency)
- Leave approval/rejection workflow
- Leave balance tracking
- Automatic leave days calculation (excluding weekends)
- Leave statistics

### 📊 Performance Management
- Performance review creation and management
- Multiple review types (quarterly, annual, probation)
- Rating system (1-5 scale)
- Review workflow (draft → submitted → completed)
- Performance statistics and analytics
- Department-wise performance tracking

## API Endpoints

### Employee Endpoints
```
POST   /api/v1/employees                    # Create employee
GET    /api/v1/employees/:id               # Get employee by ID
GET    /api/v1/employees                   # Get all employees (with filters)
PUT    /api/v1/employees/:id               # Update employee
DELETE /api/v1/employees/:id               # Delete employee
GET    /api/v1/employees/search            # Search employees
GET    /api/v1/employees/by-email          # Get employee by email
GET    /api/v1/employees/statistics/hrm    # Get HRM statistics
```

### Department Endpoints
```
POST   /api/v1/departments                     # Create department
GET    /api/v1/departments/:id                 # Get department by ID
GET    /api/v1/departments                     # Get all departments
PUT    /api/v1/departments/:id                 # Update department
DELETE /api/v1/departments/:id                 # Delete department
GET    /api/v1/departments/with-employee-count # Get departments with employee count
```

### Leave Endpoints
```
POST   /api/v1/leaves                          # Create leave request
GET    /api/v1/leaves/:id                      # Get leave by ID
GET    /api/v1/leaves                          # Get all leaves (with filters)
PUT    /api/v1/leaves/:id                      # Update leave
DELETE /api/v1/leaves/:id                      # Delete leave
POST   /api/v1/leaves/:id/approve              # Approve leave
POST   /api/v1/leaves/:id/reject               # Reject leave
GET    /api/v1/leaves/employee/:employee_id    # Get leaves by employee
GET    /api/v1/leaves/statistics/pending-count # Get pending leaves count
GET    /api/v1/leaves/balance                  # Get leave balance
```

### Performance Endpoints
```
POST   /api/v1/performance                           # Create performance review
GET    /api/v1/performance/:id                       # Get performance review by ID
GET    /api/v1/performance                           # Get all performance reviews (with filters)
PUT    /api/v1/performance/:id                       # Update performance review
DELETE /api/v1/performance/:id                       # Delete performance review
POST   /api/v1/performance/:id/submit                # Submit performance review
POST   /api/v1/performance/:id/complete              # Complete performance review
GET    /api/v1/performance/employee/:employee_id     # Get performance by employee
GET    /api/v1/performance/statistics/average-rating # Get average performance rating
GET    /api/v1/performance/statistics/by-department  # Get performance stats by department
```

## Data Models

### Employee
```json
{
  "id": 1,
  "tenant_id": "tenant-123",
  "employee_code": "EMP001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@company.com",
  "phone": "+1234567890",
  "department_id": 1,
  "position": "Software Engineer",
  "hire_date": "2024-01-15",
  "salary": 75000.00,
  "status": "active",
  "manager_id": 2,
  "address": "123 Main St, City, Country",
  "date_of_birth": "1990-05-15",
  "gender": "male",
  "emergency_name": "Jane Doe",
  "emergency_phone": "+1234567891",
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Department
```json
{
  "id": 1,
  "tenant_id": "tenant-123",
  "name": "Engineering",
  "description": "Software development team",
  "manager_id": 2,
  "budget": 500000.00,
  "location": "Building A, Floor 3",
  "is_active": true,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### Leave
```json
{
  "id": 1,
  "tenant_id": "tenant-123",
  "employee_id": 1,
  "leave_type": "annual",
  "start_date": "2024-12-20",
  "end_date": "2024-12-25",
  "days": 4,
  "reason": "Holiday vacation",
  "status": "pending",
  "approved_by": null,
  "approved_at": null,
  "comments": "",
  "is_active": true,
  "created_at": "2024-12-01T10:00:00Z",
  "updated_at": "2024-12-01T10:00:00Z"
}
```

### Performance Review
```json
{
  "id": 1,
  "tenant_id": "tenant-123",
  "employee_id": 1,
  "reviewer_id": 2,
  "period": "Q4-2024",
  "review_type": "quarterly",
  "overall_rating": 4.5,
  "goals": "Complete project X and improve skills",
  "achievements": "Successfully delivered project on time",
  "strengths": "Strong technical skills and teamwork",
  "areas_for_improvement": "Communication and leadership",
  "comments": "Excellent performance overall",
  "status": "completed",
  "is_active": true,
  "created_at": "2024-12-01T10:00:00Z",
  "updated_at": "2024-12-01T10:00:00Z"
}
```

## Configuration

### Environment Variables
```
HRM_SERVICE_PORT=8083
DATABASE_URL=postgresql://username:password@localhost:5432/zplus_saas
CORS_ORIGINS=http://localhost:3000
```

### Database
The service uses PostgreSQL with the following tables:
- `departments` - Department information
- `employees` - Employee records
- `leaves` - Leave requests and approvals
- `performance_reviews` - Performance evaluation records

## Multi-tenant Support

All data is isolated by `tenant_id`:
- Every database operation includes tenant filtering
- Tenant ID is extracted from the request headers via middleware
- Cross-tenant data access is prevented at the database level

## Authentication & Authorization

- Uses JWT-based authentication via API Gateway
- Tenant middleware extracts tenant information from tokens
- All endpoints require valid authentication
- Role-based access control can be implemented at the handler level

## Error Handling

Standard HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

Error response format:
```json
{
  "error": true,
  "message": "Error description"
}
```

## Health Check

```
GET /health
```

Response:
```json
{
  "status": "ok",
  "service": "hrm-service",
  "version": "1.0.0",
  "time": "2024-12-01T10:00:00Z"
}
```

## Performance Considerations

- Database indexes on frequently queried columns
- Pagination for list endpoints (default: 10 items per page, max: 100)
- Soft deletes for data integrity
- Optimized queries with JOINs where needed
- Connection pooling for database connections

## Business Rules

### Leave Management
- Leave requests cannot be made for past dates
- Weekend days are excluded from leave day calculations
- Leave balance is calculated based on entitlements and used days
- Only pending leaves can be updated or deleted
- Leave approval requires manager or HR role

### Performance Reviews
- Rating scale is 1-5 (with decimals)
- Reviews follow workflow: draft → submitted → completed
- Only draft reviews can be deleted
- Performance statistics exclude draft reviews

### Employee Management
- Employee codes must be unique per tenant
- Email addresses must be unique per tenant
- Manager must be an existing active employee
- Employee status changes are tracked with timestamps

## Development

### Running Locally
```bash
# Set environment variables
export HRM_SERVICE_PORT=8083
export DATABASE_URL=postgresql://postgres:password@localhost:5432/zplus_saas

# Run the service
go run apps/backend/hrm-service/cmd/main.go
```

### Running with Docker
```bash
# Build image
docker build -t hrm-service -f apps/backend/hrm-service/Dockerfile .

# Run container
docker run -p 8083:8083 \
  -e DATABASE_URL=postgresql://postgres:password@host.docker.internal:5432/zplus_saas \
  hrm-service
```

### Testing
```bash
# Run tests
go test ./apps/backend/hrm-service/...

# Test health check
curl http://localhost:8083/health
```

## Integration

### API Gateway Integration
Add route configuration to API Gateway:
```go
// In api-gateway routes
app.All("/hrm/*", handlers.ProxyToHRM)
```

### Frontend Integration
Use the service endpoints with proper authentication headers:
```javascript
// Example: Get employees
const response = await fetch('/api/hrm/v1/employees', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'X-Tenant-ID': tenantId
  }
});
```

## Monitoring & Logging

- Structured logging with log levels
- Request/response logging via middleware
- Performance metrics tracking
- Health check endpoint for monitoring
- Error tracking and alerting

## Security

- Input validation on all endpoints
- SQL injection prevention via parameterized queries
- XSS protection through proper response headers
- Rate limiting (implemented at API Gateway level)
- Tenant isolation at database level
