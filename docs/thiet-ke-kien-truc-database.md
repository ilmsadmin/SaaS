# Thiết kế Kiến trúc Database - Zplus SaaS

## 1. Tổng quan Database Architecture

Zplus SaaS sử dụng **Hybrid Database Strategy** với PostgreSQL cho dữ liệu quan hệ và MongoDB cho dữ liệu phi cấu trúc, kết hợp với chiến lược **Schema-per-Tenant** để đảm bảo cô lập dữ liệu giữa các tenant trong khi vẫn tối ưu hóa hiệu suất và chi phí vận hành.

## 2. Multi-Database Strategy

### 2.1 Database Technology Selection

| Database | Use Cases | Purpose |
|----------|-----------|---------|
| **PostgreSQL** | User data, Transactions, Relations | Structured data with ACID properties |
| **MongoDB** | Files, Logs, Analytics, Flexible schemas | Document storage and unstructured data |
| **Redis** | Session, Cache, Queue | High-performance caching and real-time data |

### 2.2 Schema Separation Strategy (PostgreSQL)

```sql
-- PostgreSQL Multi-tenant Architecture
-- Multiple Schemas for tenant isolation

┌─────────────────────────────────────────────────────────────┐
│                   PostgreSQL Server                         │
├─────────────────────────────────────────────────────────────┤
│  Schema: public (hoặc system)                               │
│  └── System-wide data                                       │
├─────────────────────────────────────────────────────────────┤
│  Schema: tenant_company_a                                   │
│  └── Company A specific data                                │
├─────────────────────────────────────────────────────────────┤
│  Schema: tenant_company_b                                   │
│  └── Company B specific data                                │
├─────────────────────────────────────────────────────────────┤
│  Schema: tenant_zin100                                      │
│  └── Zin100 specific data                                   │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 Collection Separation Strategy (MongoDB)

```javascript
// MongoDB Multi-tenant Architecture
// Database-per-tenant approach for MongoDB

┌─────────────────────────────────────────────────────────────┐
│                     MongoDB Server                          │
├─────────────────────────────────────────────────────────────┤
│  Database: system                                           │
│  └── System-wide logs and analytics                         │
├─────────────────────────────────────────────────────────────┤
│  Database: tenant_company_a                                 │
│  └── Company A files, logs, flexible data                   │
├─────────────────────────────────────────────────────────────┤
│  Database: tenant_company_b                                 │
│  └── Company B files, logs, flexible data                   │
├─────────────────────────────────────────────────────────────┤
│  Database: tenant_zin100                                    │
│  └── Zin100 files, logs, flexible data                      │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Lợi ích của Schema-per-Tenant và Database-per-Tenant

**Ưu điểm:**
- **Data Isolation**: Dữ liệu được tách biệt hoàn toàn
- **Security**: Không thể truy cập cross-tenant data
- **Backup**: Có thể backup riêng từng tenant
- **Scalability**: Dễ dàng migrate tenant sang server khác
- **Customization**: Mỗi tenant có thể có schema/database tùy chỉnh
- **Performance**: MongoDB cung cấp hiệu suất cao cho document queries

**Nhược điểm:**
- **Complexity**: Phức tạp hơn single-database approach
- **Maintenance**: Cần migration cho tất cả schemas/databases
- **Resource**: Có thể tốn nhiều connection pool
- **Consistency**: Cần xử lý cross-database transactions cẩn thận

## 3. System Layer Database Design

### 3.1 Schema: public/system

```sql
-- Quản lý Admin hệ thống
CREATE TABLE system_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- super_admin, admin, support
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Tenant
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL, -- cho subdomain
    description TEXT,
    custom_domain VARCHAR(255),
    schema_name VARCHAR(63) NOT NULL, -- tenant_slug
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, deleted
    settings JSONB, -- cấu hình tùy chỉnh
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Gói dịch vụ
CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    billing_cycle VARCHAR(20) NOT NULL, -- monthly, yearly
    max_users INTEGER,
    max_storage_gb INTEGER,
    features JSONB, -- danh sách tính năng
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Subscription
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    plan_id INTEGER REFERENCES plans(id),
    status VARCHAR(20) DEFAULT 'active', -- active, cancelled, expired
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    auto_renew BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Module
CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL, -- crm, lms, pos, hrm
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    dependencies JSONB, -- modules cần thiết
    default_permissions JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Module cho từng Tenant
CREATE TABLE tenant_modules (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    module_id INTEGER REFERENCES modules(id),
    is_enabled BOOLEAN DEFAULT true,
    configuration JSONB, -- cấu hình riêng của module
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, module_id)
);

-- Quản lý Billing/Payment
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER REFERENCES subscriptions(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) NOT NULL, -- pending, completed, failed
    payment_method VARCHAR(50), -- stripe, paypal, bank_transfer
    payment_id VARCHAR(255), -- external payment ID
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 Indexes for System Schema

```sql
-- Performance indexes
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_subscriptions_tenant ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_tenant_modules_tenant ON tenant_modules(tenant_id);
CREATE INDEX idx_tenant_modules_module ON tenant_modules(module_id);
CREATE INDEX idx_payments_subscription ON payments(subscription_id);
```

## 4. Tenant Layer Database Design

### 4.1 Schema: tenant_[slug]

Mỗi tenant sẽ có schema riêng với cấu trúc chuẩn:

```sql
-- Sử dụng schema riêng: tenant_zin100, tenant_acme, etc.
SET search_path TO tenant_zin100;

-- Quản lý User trong tenant
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Role
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false, -- admin, manager, user
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quản lý Permission
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(100) NOT NULL, -- users, customers, reports
    action VARCHAR(50) NOT NULL, -- create, read, update, delete
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gán Role cho User
CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- Gán Permission cho Role
CREATE TABLE role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE(role_id, permission_id)
);

-- Quản lý Customer/End User
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE, -- mã khách hàng
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    company VARCHAR(255),
    customer_type VARCHAR(50), -- individual, company
    status VARCHAR(20) DEFAULT 'active',
    notes TEXT,
    custom_fields JSONB, -- trường tùy chỉnh
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cấu hình Module cho Tenant 
CREATE TABLE modules_config (
    id SERIAL PRIMARY KEY,
    module_name VARCHAR(100) NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value JSONB,
    description TEXT,
    updated_by INTEGER REFERENCES users(id),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(module_name, config_key)
);

-- Audit Log
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id INTEGER,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4.2 Tenant Schema Indexes

```sql
-- Performance indexes for tenant schema
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
CREATE INDEX idx_customers_code ON customers(code);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);
```

## 5. Module-specific Database Design

### 5.1 CRM Module Tables

```sql
-- Trong schema tenant_xxx
CREATE TABLE crm_leads (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    value DECIMAL(12,2),
    status VARCHAR(50) DEFAULT 'new', -- new, contacted, qualified, closed
    source VARCHAR(100), -- website, referral, ads
    assigned_to INTEGER REFERENCES users(id),
    expected_close_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE crm_opportunities (
    id SERIAL PRIMARY KEY,
    lead_id INTEGER REFERENCES crm_leads(id),
    customer_id INTEGER REFERENCES customers(id),
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12,2),
    stage VARCHAR(50) NOT NULL, -- prospect, proposal, negotiation, closed
    probability INTEGER DEFAULT 0, -- 0-100%
    close_date DATE,
    owner_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE crm_activities (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    lead_id INTEGER REFERENCES crm_leads(id),
    activity_type VARCHAR(50) NOT NULL, -- call, email, meeting, note
    subject VARCHAR(255),
    description TEXT,
    duration_minutes INTEGER,
    completed BOOLEAN DEFAULT false,
    scheduled_at TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 5.2 LMS Module Tables

```sql
-- Learning Management System tables
CREATE TABLE lms_courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url VARCHAR(500),
    duration_hours INTEGER,
    price DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'draft', -- draft, published, archived
    instructor_id INTEGER REFERENCES users(id),
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lms_lessons (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES lms_courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    video_url VARCHAR(500),
    duration_minutes INTEGER,
    order_index INTEGER NOT NULL,
    is_free BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lms_students (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    student_code VARCHAR(50) UNIQUE,
    enrollment_date DATE DEFAULT CURRENT_DATE,
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, graduated
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lms_enrollments (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES lms_students(id),
    course_id INTEGER REFERENCES lms_courses(id),
    enrollment_date DATE DEFAULT CURRENT_DATE,
    completion_date DATE,
    progress_percentage INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'enrolled', -- enrolled, completed, dropped
    UNIQUE(student_id, course_id)
);

CREATE TABLE lms_lesson_progress (
    id SERIAL PRIMARY KEY,
    enrollment_id INTEGER REFERENCES lms_enrollments(id),
    lesson_id INTEGER REFERENCES lms_lessons(id),
    completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP,
    watch_time_seconds INTEGER DEFAULT 0,
    UNIQUE(enrollment_id, lesson_id)
);

CREATE TABLE lms_quizzes (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES lms_courses(id),
    lesson_id INTEGER REFERENCES lms_lessons(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    time_limit_minutes INTEGER,
    passing_score INTEGER DEFAULT 70,
    questions JSONB NOT NULL, -- câu hỏi và đáp án
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lms_quiz_attempts (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER REFERENCES lms_quizzes(id),
    student_id INTEGER REFERENCES lms_students(id),
    answers JSONB, -- câu trả lời
    score INTEGER,
    passed BOOLEAN,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

### 5.3 POS Module Tables

```sql
-- Point of Sale tables
CREATE TABLE pos_products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2),
    category VARCHAR(100),
    stock_quantity INTEGER DEFAULT 0,
    min_stock_level INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pos_orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id INTEGER REFERENCES customers(id),
    total_amount DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    payment_method VARCHAR(50), -- cash, card, transfer
    payment_status VARCHAR(20) DEFAULT 'pending', -- pending, paid, refunded
    order_status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, cancelled
    cashier_id INTEGER REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pos_order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES pos_orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES pos_products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL
);

CREATE TABLE pos_inventory_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES pos_products(id),
    movement_type VARCHAR(20) NOT NULL, -- in, out, adjustment
    quantity INTEGER NOT NULL,
    reference_type VARCHAR(50), -- order, adjustment, return
    reference_id INTEGER,
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 6. Database Security & Performance

### 6.1 Row Level Security (RLS)

```sql
-- Enable RLS for multi-tenant security
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

-- Policy để chỉ cho phép truy cập dữ liệu của tenant hiện tại
CREATE POLICY tenant_isolation_policy ON customers
FOR ALL TO application_role
USING (pg_has_role(current_user, 'tenant_' || current_setting('app.tenant_id'), 'USAGE'));
```

### 6.2 Connection Management

```go
// Connection pool per tenant
type DatabaseManager struct {
    systemDB *sql.DB
    tenantPools map[string]*sql.DB
}

func (dm *DatabaseManager) GetTenantDB(tenantID string) (*sql.DB, error) {
    if pool, exists := dm.tenantPools[tenantID]; exists {
        return pool, nil
    }
    
    // Create new connection pool for tenant
    dsn := fmt.Sprintf("%s?search_path=tenant_%s", baseDSN, tenantID)
    pool, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    dm.tenantPools[tenantID] = pool
    return pool, nil
}
```

### 6.3 Backup Strategy

```bash
# System backup
pg_dump -h localhost -U postgres -n public zplus_saas > system_backup.sql

# Tenant backup
pg_dump -h localhost -U postgres -n tenant_zin100 zplus_saas > tenant_zin100_backup.sql

# Full backup
pg_dump -h localhost -U postgres zplus_saas > full_backup.sql
```

## 7. Migration Strategy

### 7.1 Schema Migration

```go
// Migration for all tenants
func RunMigration(version string) error {
    // Migrate system schema first
    if err := migrateSystemSchema(version); err != nil {
        return err
    }
    
    // Get all active tenants
    tenants, err := getAllActiveTenants()
    if err != nil {
        return err
    }
    
    // Migrate each tenant schema
    for _, tenant := range tenants {
        if err := migrateTenantSchema(tenant.SchemaName, version); err != nil {
            log.Printf("Failed to migrate tenant %s: %v", tenant.Name, err)
            continue
        }
    }
    
    return nil
}
```

### 7.2 Data Migration

```sql
-- Example: Add column to all tenant schemas
DO $$
DECLARE
    tenant_schema TEXT;
BEGIN
    FOR tenant_schema IN 
        SELECT schema_name FROM information_schema.schemata 
        WHERE schema_name LIKE 'tenant_%'
    LOOP
        EXECUTE format('ALTER TABLE %I.customers ADD COLUMN IF NOT EXISTS tags JSONB', tenant_schema);
    END LOOP;
END $$;
```

## 8. Monitoring & Maintenance

### 8.1 Database Monitoring

```sql
-- Monitor schema sizes
SELECT 
    schemaname,
    pg_size_pretty(SUM(pg_total_relation_size(schemaname||'.'||tablename))::bigint) as size
FROM pg_tables 
WHERE schemaname LIKE 'tenant_%'
GROUP BY schemaname
ORDER BY SUM(pg_total_relation_size(schemaname||'.'||tablename)) DESC;

-- Monitor connection usage
SELECT 
    datname,
    usename,
    client_addr,
    COUNT(*) as connections
FROM pg_stat_activity 
WHERE state = 'active'
GROUP BY datname, usename, client_addr;
```

### 8.2 Performance Optimization

```sql
-- Analyze tables regularly
ANALYZE;

-- Vacuum old data
VACUUM (ANALYZE, VERBOSE);

-- Check slow queries
SELECT 
    query,
    mean_time,
    calls,
    total_time
FROM pg_stat_statements 
ORDER BY total_time DESC
LIMIT 10;
```

## 9. MongoDB Database Design

### 9.1 System Database Collections

```javascript
// Database: system
// Collection: system_logs
{
  _id: ObjectId,
  timestamp: ISODate,
  level: "info" | "warning" | "error",
  service: "gateway" | "auth" | "crm" | "hrm" | "pos",
  message: String,
  metadata: Object,
  tenant_id: String // for filtering
}

// Collection: system_analytics
{
  _id: ObjectId,
  date: ISODate,
  metric_type: "usage" | "performance" | "business",
  tenant_id: String,
  data: {
    users_active: Number,
    requests_count: Number,
    response_time_avg: Number,
    // ... other metrics
  }
}
```

### 9.2 Tenant Database Collections

```javascript
// Database: tenant_{slug}
// Collection: files
{
  _id: ObjectId,
  filename: String,
  original_name: String,
  mime_type: String,
  size: Number,
  path: String,
  uploaded_by: ObjectId, // Reference to user
  module: "crm" | "hrm" | "pos" | "lms",
  tags: [String],
  metadata: Object,
  created_at: ISODate,
  updated_at: ISODate
}

// Collection: audit_logs
{
  _id: ObjectId,
  user_id: ObjectId,
  action: String,
  resource: String,
  resource_id: String,
  changes: {
    before: Object,
    after: Object
  },
  ip_address: String,
  user_agent: String,
  timestamp: ISODate
}

// Collection: notifications
{
  _id: ObjectId,
  user_id: ObjectId,
  type: "email" | "push" | "sms",
  title: String,
  message: String,
  data: Object,
  read: Boolean,
  sent: Boolean,
  created_at: ISODate
}
```

### 9.3 Module-specific Collections

```javascript
// CRM Module
// Collection: crm_activities
{
  _id: ObjectId,
  customer_id: ObjectId, // Reference to PostgreSQL
  type: "call" | "email" | "meeting" | "note",
  subject: String,
  description: String,
  attachments: [ObjectId], // References to files collection
  user_id: ObjectId,
  timestamp: ISODate,
  metadata: Object
}

// LMS Module  
// Collection: lms_course_content
{
  _id: ObjectId,
  course_id: ObjectId, // Reference to PostgreSQL
  lesson_id: ObjectId, // Reference to PostgreSQL
  content_type: "video" | "document" | "quiz" | "assignment",
  content_data: Object,
  files: [ObjectId], // References to files collection
  created_at: ISODate,
  updated_at: ISODate
}

// HRM Module
// Collection: hrm_timesheets
{
  _id: ObjectId,
  employee_id: ObjectId, // Reference to PostgreSQL
  date: ISODate,
  entries: [{
    start_time: ISODate,
    end_time: ISODate,
    break_duration: Number,
    project_id: ObjectId,
    task_description: String,
    location: {
      type: "Point",
      coordinates: [Number, Number]
    }
  }],
  total_hours: Number,
  status: "draft" | "submitted" | "approved"
}
```

## 10. Redis Cache Strategy

### 10.1 Cache Structure

```javascript
// Session Cache
session:{session_id} = {
  user_id: Number,
  tenant_id: String,
  roles: [String],
  permissions: [String],
  last_activity: Timestamp,
  expires_at: Timestamp
}

// Tenant Cache
tenant:{tenant_slug} = {
  id: Number,
  name: String,
  schema_name: String,
  settings: Object,
  active_modules: [String],
  expires_at: Timestamp
}

// User Cache
user:{user_id}:{tenant_id} = {
  id: Number,
  email: String,
  name: String,
  roles: [String],
  permissions: [String],
  expires_at: Timestamp
}
```

### 10.2 Cache Patterns

```javascript
// Cache Aside Pattern
function getUser(userId, tenantId) {
  key = `user:${userId}:${tenantId}`
  
  // Try cache first
  user = redis.get(key)
  if (user) return user
  
  // Fallback to database
  user = database.getUser(userId, tenantId)
  
  // Cache the result
  redis.setex(key, 3600, user)
  return user
}

// Write Through Pattern
function updateUser(userId, tenantId, userData) {
  key = `user:${userId}:${tenantId}`
  
  // Update database
  database.updateUser(userId, tenantId, userData)
  
  // Update cache
  redis.setex(key, 3600, userData)
}
```

## 11. Data Integration Patterns

### 11.1 PostgreSQL ↔ MongoDB Sync

```go
// GORM Event Hooks for MongoDB sync
func (u *User) AfterCreate(tx *gorm.DB) error {
    // Sync to MongoDB for audit logging
    auditLog := AuditLog{
        UserID:    u.ID,
        Action:    "create",
        Resource:  "user",
        Changes:   map[string]interface{}{"after": u},
        Timestamp: time.Now(),
    }
    
    return mongoClient.Database(u.TenantID).
        Collection("audit_logs").
        InsertOne(context.Background(), auditLog)
}
```

### 11.2 Cross-Database Queries

```go
// Service layer handling multiple databases
type CustomerService struct {
    pgDB    *gorm.DB
    mongoDB *mongo.Database
}

func (s *CustomerService) GetCustomerWithActivities(customerID uint) (*CustomerWithActivities, error) {
    // Get customer from PostgreSQL
    var customer Customer
    err := s.pgDB.First(&customer, customerID).Error
    if err != nil {
        return nil, err
    }
    
    // Get activities from MongoDB
    var activities []Activity
    cursor, err := s.mongoDB.Collection("crm_activities").
        Find(context.Background(), bson.M{"customer_id": customerID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    cursor.All(context.Background(), &activities)
    
    return &CustomerWithActivities{
        Customer:   customer,
        Activities: activities,
    }, nil
}
```