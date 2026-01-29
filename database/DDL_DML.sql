-- 1. Update categories table (if not exists)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Update products table with category_id foreign key
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    stock INTEGER DEFAULT 0 CHECK (stock >= 0),
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT
);

-- 3. Insert sample categories (if not exists)
INSERT INTO categories (name, description) VALUES 
    ('Food', 'Food and snacks'),
    ('Beverages', 'Drinks and beverages'),
    ('Condiments', 'Sauces and seasonings'),
    ('Electronics', 'Electronic devices'),
    ('Clothing', 'Clothes and accessories')
ON CONFLICT (name) DO NOTHING;

-- 4. Insert sample products with category relationships
INSERT INTO products (name, price, stock, category_id) VALUES 
    ('Indomie Godog', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap', 12000, 20, 3),
    ('Smartphone', 2500000, 5, 4),
    ('T-Shirt', 150000, 50, 5)
ON CONFLICT DO NOTHING;