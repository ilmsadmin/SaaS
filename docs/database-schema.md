# Database Schema - Zplus SaaS

## 1. Tổng quan Database Architecture

Zplus SaaS sử dụng chiến lược **Schema-per-Tenant** để đảm bảo data isolation hoàn toàn giữa các tenant.

### Database Structure

```
PostgreSQL Server
├── Schema: public (System Level)
├── Schema: tenant_demo (Tenant Level)
├── Schema: tenant_acme (Tenant Level)
└── Schema: tenant_xyz (Tenant Level)
```

### Multi-tenant Strategy Benefits

- **Data Isolation**: Hoàn toàn tách biệt dữ liệu giữa các tenant
- **Security**: Không thể access cross-tenant data
- **Backup**: Có thể backup/restore từng tenant riêng biệt
- **Scaling**: Dễ dàng move tenant sang database khác
- **Compliance**: Đáp ứng yêu cầu GDPR, SOC2

## 2. System Schema (public)

### 2.1 System Users & Authentication

```sql
-- System level users (Super Admin, System Admin)
CREATE TABLE system_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar TEXT,
    status VARCHAR(20) DEFAULT 'active',
    role VARCHAR(50) NOT NULL, -- 'super_admin', 'system_admin', 'support'
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- System sessions
CREATE TABLE system_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES system_users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 2.2 Tenant Management

```sql
-- Tenant registry
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) UNIQUE NOT NULL, -- Used in subdomain
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255), -- Custom domain (optional)
    logo TEXT,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'suspended', 'trial'
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Tenant database schemas
CREATE TABLE tenant_schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    schema_name VARCHAR(100) NOT NULL,
    database_host VARCHAR(255) DEFAULT 'localhost',
    database_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 2.3 Plans & Subscriptions

```sql
-- Service plans
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL, -- 'monthly', 'yearly'
    max_users INTEGER,
    max_storage BIGINT, -- in bytes
    features JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Tenant subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID REFERENCES plans(id),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'cancelled', 'expired'
    trial_end_at TIMESTAMP,
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payment history
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) NOT NULL, -- 'pending', 'completed', 'failed'
    payment_method VARCHAR(50),
    external_payment_id VARCHAR(255),
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 2.4 Module System

```sql
-- Available modules
CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    category VARCHAR(50), -- 'crm', 'lms', 'pos', 'hrm', 'checkin'
    icon VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    dependencies JSONB DEFAULT '[]',
    permissions JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Modules enabled for each tenant
CREATE TABLE tenant_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    is_enabled BOOLEAN DEFAULT true,
    config JSONB DEFAULT '{}',
    installed_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, module_id)
);
```

### 2.5 System Analytics

```sql
-- Usage analytics
CREATE TABLE usage_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL, -- 'api_requests', 'storage_used', 'active_users'
    value BIGINT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, metric_name, date)
);

-- Audit logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
    user_id UUID, -- Can be system_user or tenant_user
    user_type VARCHAR(20), -- 'system', 'tenant'
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 3. Tenant Schema Template

Mỗi tenant sẽ có schema riêng với cấu trúc như sau:

### 3.1 User Management

```sql
-- Tenant users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar TEXT,
    phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'suspended'
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false, -- Cannot be deleted
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(100) NOT NULL, -- 'users', 'customers', 'products'
    action VARCHAR(50) NOT NULL, -- 'create', 'read', 'update', 'delete'
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Role permissions mapping
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE(role_id, permission_id)
);

-- User roles mapping
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);
```

### 3.2 Customer Management

```sql
-- End customers (CRM)
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    company VARCHAR(255),
    job_title VARCHAR(255),
    status VARCHAR(20) DEFAULT 'lead', -- 'lead', 'prospect', 'active', 'inactive', 'churned'
    tags JSONB DEFAULT '[]',
    address JSONB DEFAULT '{}',
    custom_fields JSONB DEFAULT '{}',
    notes TEXT,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    last_contact_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Customer interactions/activities
CREATE TABLE customer_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- 'call', 'email', 'meeting', 'note'
    title VARCHAR(255) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3.3 File Management

```sql
-- Files
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    path TEXT NOT NULL,
    url TEXT,
    category VARCHAR(50), -- 'avatar', 'document', 'image'
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- File associations
CREATE TABLE file_associations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL, -- 'user', 'customer', 'product'
    entity_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3.4 Notifications

```sql
-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'info', 'warning', 'error', 'success'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}',
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 4. Module-Specific Tables

### 4.1 CRM Module

```sql
-- Leads/Opportunities
CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    value DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    probability INTEGER DEFAULT 0, -- 0-100
    stage VARCHAR(50) NOT NULL, -- 'prospecting', 'qualification', 'proposal', 'negotiation', 'closed_won', 'closed_lost'
    source VARCHAR(100), -- 'website', 'referral', 'cold_call', 'social_media'
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    expected_close_date DATE,
    closed_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Pipeline stages
CREATE TABLE pipeline_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL,
    probability INTEGER DEFAULT 0,
    is_won BOOLEAN DEFAULT false,
    is_lost BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 4.2 LMS Module

```sql
-- Courses
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail TEXT,
    price DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'draft', -- 'draft', 'published', 'archived'
    category_id UUID REFERENCES course_categories(id) ON DELETE SET NULL,
    instructor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    duration_minutes INTEGER,
    level VARCHAR(20) DEFAULT 'beginner', -- 'beginner', 'intermediate', 'advanced'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Lessons
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    video_url TEXT,
    duration_minutes INTEGER,
    position INTEGER NOT NULL,
    is_free BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Student enrollments
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    student_id UUID REFERENCES customers(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    progress INTEGER DEFAULT 0, -- 0-100
    certificate_issued BOOLEAN DEFAULT false
);
```

### 4.3 POS Module

```sql
-- Products
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    price DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2),
    category_id UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    stock_quantity INTEGER DEFAULT 0,
    low_stock_threshold INTEGER DEFAULT 10,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'discontinued'
    images JSONB DEFAULT '[]',
    attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Orders
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(100) UNIQUE NOT NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Sales person
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'cancelled'
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50),
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'partial', 'refunded'
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Order items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 4.4 HRM Module

```sql
-- Employees
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    department VARCHAR(100),
    position VARCHAR(100),
    hire_date DATE,
    salary DECIMAL(10,2),
    employment_type VARCHAR(20) DEFAULT 'full_time', -- 'full_time', 'part_time', 'contract'
    manager_id UUID REFERENCES employees(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'terminated'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Attendance
CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID REFERENCES employees(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    break_duration_minutes INTEGER DEFAULT 0,
    worked_hours DECIMAL(4,2),
    status VARCHAR(20) DEFAULT 'present', -- 'present', 'absent', 'late', 'half_day'
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(employee_id, date)
);
```

## 5. Indexes và Performance

### 5.1 System Schema Indexes

```sql
-- Tenants
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);

-- Subscriptions
CREATE INDEX idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

-- Usage stats
CREATE INDEX idx_usage_stats_tenant_date ON usage_stats(tenant_id, date);
CREATE INDEX idx_usage_stats_metric ON usage_stats(metric_name, date);

-- Audit logs
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
```

### 5.2 Tenant Schema Indexes

```sql
-- Users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- Customers
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_assigned_to ON customers(assigned_to);
CREATE INDEX idx_customers_created_at ON customers(created_at);

-- Notifications
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_read_at ON notifications(read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Files
CREATE INDEX idx_files_category ON files(category);
CREATE INDEX idx_files_uploaded_by ON files(uploaded_by);

-- File associations
CREATE INDEX idx_file_associations_entity ON file_associations(entity_type, entity_id);
```

### 5.3 Module-specific Indexes

```sql
-- CRM
CREATE INDEX idx_opportunities_customer_id ON opportunities(customer_id);
CREATE INDEX idx_opportunities_stage ON opportunities(stage);
CREATE INDEX idx_opportunities_assigned_to ON opportunities(assigned_to);

-- LMS
CREATE INDEX idx_courses_status ON courses(status);
CREATE INDEX idx_courses_instructor_id ON courses(instructor_id);
CREATE INDEX idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX idx_enrollments_course_id ON enrollments(course_id);

-- POS
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- HRM
CREATE INDEX idx_employees_employee_id ON employees(employee_id);
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_attendance_employee_date ON attendance_records(employee_id, date);
```

## 6. Triggers và Functions

### 6.1 Audit Trigger

```sql
-- Function to handle audit logging
CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit_logs (
            tenant_id, action, resource, resource_id, details, created_at
        ) VALUES (
            current_setting('app.tenant_id')::UUID,
            'CREATE',
            TG_TABLE_NAME,
            NEW.id,
            row_to_json(NEW),
            NOW()
        );
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_logs (
            tenant_id, action, resource, resource_id, details, created_at
        ) VALUES (
            current_setting('app.tenant_id')::UUID,
            'UPDATE',
            TG_TABLE_NAME,
            NEW.id,
            jsonb_build_object('old', row_to_json(OLD), 'new', row_to_json(NEW)),
            NOW()
        );
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_logs (
            tenant_id, action, resource, resource_id, details, created_at
        ) VALUES (
            current_setting('app.tenant_id')::UUID,
            'DELETE',
            TG_TABLE_NAME,
            OLD.id,
            row_to_json(OLD),
            NOW()
        );
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Apply audit trigger to important tables
CREATE TRIGGER audit_users AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();
    
CREATE TRIGGER audit_customers AFTER INSERT OR UPDATE OR DELETE ON customers
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();
```

### 6.2 Updated At Trigger

```sql
-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## 7. Database Migrations

### 7.1 Migration Structure

```
infra/db/migrations/
├── system/
│   ├── 001_create_system_users.up.sql
│   ├── 001_create_system_users.down.sql
│   ├── 002_create_tenants.up.sql
│   └── 002_create_tenants.down.sql
└── tenant/
    ├── 001_create_users.up.sql
    ├── 001_create_users.down.sql
    ├── 002_create_customers.up.sql
    └── 002_create_customers.down.sql
```

### 7.2 Migration Commands

```bash
# System migrations
migrate -path infra/db/migrations/system -database "postgres://..." up

# Tenant migrations (for specific tenant)
migrate -path infra/db/migrations/tenant -database "postgres://..." up

# Create new migration
migrate create -ext sql -dir infra/db/migrations/system add_audit_logs
```

## 8. Backup & Recovery

### 8.1 Backup Strategy

```bash
# Full system backup
pg_dump -h localhost -U postgres zplus_system > system_backup.sql

# Tenant-specific backup
pg_dump -h localhost -U postgres -n tenant_demo zplus_system > tenant_demo_backup.sql

# All tenants backup
pg_dump -h localhost -U postgres zplus_system > full_backup.sql
```

### 8.2 Point-in-Time Recovery

```bash
# Enable WAL archiving in postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'cp %p /backup/wal/%f'

# Create base backup
pg_basebackup -D /backup/base -Ft -z -P

# Restore to specific point
pg_ctl stop -D /var/lib/postgresql/data
rm -rf /var/lib/postgresql/data/*
tar -xf /backup/base/base.tar.gz -C /var/lib/postgresql/data/
# Create recovery.conf with target time
pg_ctl start -D /var/lib/postgresql/data
```

## 9. Performance Monitoring

### 9.1 Database Statistics

```sql
-- Monitor table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname LIKE 'tenant_%' 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Monitor query performance
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;
```

### 9.2 Index Usage

```sql
-- Check unused indexes
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes 
WHERE idx_scan = 0;

-- Check index efficiency
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes 
ORDER BY idx_scan DESC;
```

## 10. Security Considerations

### 10.1 Row Level Security (RLS)

```sql
-- Enable RLS on tenant tables
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

-- Create policy for tenant isolation
CREATE POLICY tenant_isolation ON customers
    FOR ALL
    TO application_role
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
```

### 10.2 Database Roles

```sql
-- Application role
CREATE ROLE app_user;
GRANT CONNECT ON DATABASE zplus_system TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;

-- Read-only role for analytics
CREATE ROLE analytics_user;
GRANT CONNECT ON DATABASE zplus_system TO analytics_user;
GRANT USAGE ON ALL SCHEMAS IN DATABASE zplus_system TO analytics_user;
GRANT SELECT ON ALL TABLES IN DATABASE zplus_system TO analytics_user;
```

### 10.3 Data Encryption

```sql
-- Encrypt sensitive fields
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Store encrypted data
INSERT INTO users (email, password_hash) 
VALUES ('user@example.com', crypt('password123', gen_salt('bf')));

-- Verify password
SELECT * FROM users 
WHERE email = 'user@example.com' 
AND password_hash = crypt('password123', password_hash);
```
