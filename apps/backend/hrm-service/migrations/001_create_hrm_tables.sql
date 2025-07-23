-- HRM Service Database Migration - Create HRM Tables
-- Version: 001
-- Created: 2025-07-23

-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    manager_id INT,
    budget DECIMAL(15, 2) DEFAULT 0.00,
    location VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    employee_code VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    department_id INT NOT NULL,
    position VARCHAR(100) NOT NULL,
    hire_date DATE NOT NULL,
    salary DECIMAL(15, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'terminated')),
    manager_id INT,
    address TEXT,
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    emergency_name VARCHAR(100),
    emergency_phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (manager_id) REFERENCES employees(id)
);

-- Create leaves table
CREATE TABLE IF NOT EXISTS leaves (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    employee_id INT NOT NULL,
    leave_type VARCHAR(20) NOT NULL CHECK (leave_type IN ('annual', 'sick', 'maternity', 'paternity', 'personal', 'emergency')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days INT NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    approved_by INT,
    approved_at TIMESTAMP,
    comments TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (approved_by) REFERENCES employees(id)
);

-- Create performance_reviews table
CREATE TABLE IF NOT EXISTS performance_reviews (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    employee_id INT NOT NULL,
    reviewer_id INT NOT NULL,
    period VARCHAR(20) NOT NULL,
    review_type VARCHAR(20) NOT NULL CHECK (review_type IN ('quarterly', 'annual', 'probation')),
    overall_rating DECIMAL(3, 2) NOT NULL CHECK (overall_rating >= 1.0 AND overall_rating <= 5.0),
    goals TEXT,
    achievements TEXT,
    strengths TEXT,
    areas_for_improvement TEXT,
    comments TEXT,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'completed')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (reviewer_id) REFERENCES employees(id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_departments_tenant_id ON departments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_departments_manager_id ON departments(manager_id);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments(is_active);

CREATE INDEX IF NOT EXISTS idx_employees_tenant_id ON employees(tenant_id);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees(department_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees(manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
CREATE INDEX IF NOT EXISTS idx_employees_employee_code ON employees(employee_code);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);
CREATE INDEX IF NOT EXISTS idx_employees_is_active ON employees(is_active);
CREATE INDEX IF NOT EXISTS idx_employees_hire_date ON employees(hire_date);

CREATE INDEX IF NOT EXISTS idx_leaves_tenant_id ON leaves(tenant_id);
CREATE INDEX IF NOT EXISTS idx_leaves_employee_id ON leaves(employee_id);
CREATE INDEX IF NOT EXISTS idx_leaves_status ON leaves(status);
CREATE INDEX IF NOT EXISTS idx_leaves_leave_type ON leaves(leave_type);
CREATE INDEX IF NOT EXISTS idx_leaves_start_date ON leaves(start_date);
CREATE INDEX IF NOT EXISTS idx_leaves_approved_by ON leaves(approved_by);
CREATE INDEX IF NOT EXISTS idx_leaves_is_active ON leaves(is_active);

CREATE INDEX IF NOT EXISTS idx_performance_tenant_id ON performance_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_performance_employee_id ON performance_reviews(employee_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviewer_id ON performance_reviews(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_performance_period ON performance_reviews(period);
CREATE INDEX IF NOT EXISTS idx_performance_review_type ON performance_reviews(review_type);
CREATE INDEX IF NOT EXISTS idx_performance_status ON performance_reviews(status);
CREATE INDEX IF NOT EXISTS idx_performance_is_active ON performance_reviews(is_active);

-- Create unique constraints
CREATE UNIQUE INDEX IF NOT EXISTS idx_employee_code_tenant ON employees(tenant_id, employee_code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_employee_email_tenant ON employees(tenant_id, email);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS update_departments_updated_at ON departments;
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
DROP TRIGGER IF EXISTS update_leaves_updated_at ON leaves;
DROP TRIGGER IF EXISTS update_performance_reviews_updated_at ON performance_reviews;

-- Create triggers
CREATE TRIGGER update_departments_updated_at
    BEFORE UPDATE ON departments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_employees_updated_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_leaves_updated_at
    BEFORE UPDATE ON leaves
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_performance_reviews_updated_at
    BEFORE UPDATE ON performance_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default departments for testing/demo purposes
INSERT INTO departments (tenant_id, name, description, budget, location) VALUES
('demo-tenant', 'Human Resources', 'Manages employee relations, recruitment, and policies', 50000.00, 'Head Office'),
('demo-tenant', 'Engineering', 'Software development and technical operations', 200000.00, 'Tech Hub'),
('demo-tenant', 'Sales', 'Customer acquisition and revenue generation', 80000.00, 'Sales Floor'),
('demo-tenant', 'Marketing', 'Brand promotion and market research', 60000.00, 'Creative Studio'),
('demo-tenant', 'Finance', 'Financial planning and accounting', 70000.00, 'Finance Wing')
ON CONFLICT DO NOTHING;
