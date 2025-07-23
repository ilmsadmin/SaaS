-- POS Service Database Migration
-- Create POS tables with multi-tenant support

-- Categories table for product categorization
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table for product catalog
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sku VARCHAR(100) UNIQUE NOT NULL,
    barcode VARCHAR(100),
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    cost_price DECIMAL(10,2) DEFAULT 0.00,
    stock_quantity INTEGER DEFAULT 0,
    min_stock_level INTEGER DEFAULT 0,
    max_stock_level INTEGER DEFAULT 0,
    unit VARCHAR(50) DEFAULT 'pcs',
    weight DECIMAL(8,3),
    dimensions VARCHAR(100),
    image_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table for sales transactions
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    order_number VARCHAR(100) UNIQUE NOT NULL,
    customer_name VARCHAR(255),
    customer_phone VARCHAR(20),
    customer_email VARCHAR(255),
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    payment_method VARCHAR(50) DEFAULT 'cash',
    payment_status VARCHAR(20) DEFAULT 'pending',
    order_status VARCHAR(20) DEFAULT 'pending',
    served_by VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items table for detailed order contents
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Inventory transactions for stock management
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    transaction_type VARCHAR(20) NOT NULL, -- 'in', 'out', 'adjustment'
    quantity INTEGER NOT NULL,
    unit_cost DECIMAL(10,2),
    total_cost DECIMAL(10,2),
    reference_type VARCHAR(50), -- 'purchase', 'sale', 'adjustment'
    reference_id INTEGER,
    notes TEXT,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Suppliers table for inventory management
CREATE TABLE IF NOT EXISTS suppliers (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    payment_terms VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase orders for inventory restocking
CREATE TABLE IF NOT EXISTS purchase_orders (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE SET NULL,
    po_number VARCHAR(100) UNIQUE NOT NULL,
    order_date DATE DEFAULT CURRENT_DATE,
    expected_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    subtotal DECIMAL(10,2) DEFAULT 0.00,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) DEFAULT 0.00,
    notes TEXT,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase order items
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    purchase_order_id INTEGER REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    unit_cost DECIMAL(10,2) NOT NULL,
    total_cost DECIMAL(10,2) NOT NULL,
    received_quantity INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_tenant_id ON products(tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_tenant_id ON order_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_inventory_transactions_tenant_id ON inventory_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_transactions_product_id ON inventory_transactions(product_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_tenant_id ON suppliers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_tenant_id ON purchase_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_tenant_id ON purchase_order_items(tenant_id);

-- Create triggers for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_suppliers_updated_at BEFORE UPDATE ON suppliers FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_purchase_orders_updated_at BEFORE UPDATE ON purchase_orders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
