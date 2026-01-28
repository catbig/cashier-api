-- Products table
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  price INTEGER NOT NULL,
  stock INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Insert sample data
INSERT INTO products (name, price, stock) VALUES 
('Indomie Godog', 3500, 10),
('Vit 1000ml', 3000, 40),
('Kecap', 12000, 20);
