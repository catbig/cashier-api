# Cashier API - Session 2: Layered Architecture with Selective Category Display

A professional cashier system API built with Go implementing layered architecture, PostgreSQL database (Supabase), and intelligent product-category relationships with selective display.

## 🚀 Project Overview

Enhanced Cashier API implementing industry best practices with optimized data retrieval:
- **Complete Layered Architecture** (Handler-Service-Repository-Model)
- **PostgreSQL Database** with Supabase integration
- **Smart Product-Category Relationships** - category shown only where needed
- **Configuration Management** with Viper
- **Performance Optimized** data retrieval

## ✅ Task Requirements Completed

### 1. **Categories moved to layered architecture** ✅
- Complete separation of concerns for Categories
- Proper dependency injection setup
- Consistent error handling across all layers

### 2. **Challenge: Implement JOIN with selective display** ✅
- Added `category_id` foreign key to products table
- **Smart JOIN implementation**: category only in detail view
- **Optimized list view**: no unnecessary JOINs
- Added validation for category existence
- Protected category deletion when products exist

## 📋 Prerequisites

### Required Tools
- [Go 1.21+](https://golang.org/dl/)
- [Git](https://git-scm.com/)
- [Supabase Account](https://supabase.com/) (free tier)

### Environment Setup
```bash
# Check Go installation
go version

# Check Git installation
git --version
```

## 🏗️ Project Structure (Layered Architecture)

```
cashier-api/
├── database/              # Database connection and setup
├── models/                # Data structures
│   ├── product.go        # Product models (list vs detail)
│   └── category.go       # Category model
├── repositories/          # Database operations
│   ├── product_repository.go  # Smart DB operations
│   └── category_repository.go # Category DB operations
├── services/              # Business logic
│   ├── product_service.go     # Product business logic
│   └── category_service.go    # Category business logic
├── handlers/              # HTTP handlers
│   ├── product_handler.go     # Product HTTP handlers
│   └── category_handler.go    # Category HTTP handlers
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
├── .env                  # Environment variables
└── .gitignore           # Git ignore file
```

### Architecture Flow
```
HTTP Request → Handler → Service → Repository → Database
HTTP Response ← Handler ← Service ← Repository ← Database
```

## ⚡ Quick Start

### 1. Clone and Setup
```bash
# Create project directory
mkdir cashier-api
cd cashier-api

# Initialize Go module
go mod init cashier-api

# Create directory structure
mkdir -p database models repositories services handlers
```

### 2. Install Dependencies
```bash
# Install required packages
go get github.com/spf13/viper
go get github.com/lib/pq
go mod tidy
```

### 3. Supabase Database Setup

1. **Create a Supabase Project** at [supabase.com](https://supabase.com)
2. **Get Connection String**:
   - Go to Project Settings → Database
   - Find "Connection string" section
   - Use "Transaction Pooler" URL
   - Format: `postgresql://postgres:[YOUR_PASSWORD]@db.[YOUR_REF].supabase.co:5432/postgres`  
     Example: `postgresql://postgres.abcdefghi:ABCD1234@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres`

3. **Create Configuration File**:
```bash
# Create .env file
cat > .env << 'EOF'
PORT=8080
DB_CONN=postgresql://postgres:[YOUR_PASSWORD]@db.[YOUR_REF].supabase.co:5432/postgres
EOF
```

4. **Run Database Schema Setup** in Supabase SQL Editor:
```sql
-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create products table with foreign key
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

-- Insert sample categories
INSERT INTO categories (name, description) VALUES 
    ('Food', 'Food and snacks'),
    ('Beverages', 'Drinks and beverages'),
    ('Condiments', 'Sauces and seasonings'),
    ('Electronics', 'Electronic devices'),
    ('Clothing', 'Clothes and accessories')
ON CONFLICT (name) DO NOTHING;

-- Insert sample products with categories
INSERT INTO products (name, price, stock, category_id) VALUES 
    ('Indomie Godog', 3500, 10, 1),
    ('Vit 1000ml', 3000, 40, 2),
    ('Kecap', 12000, 20, 3),
    ('Smartphone', 2500000, 5, 4),
    ('T-Shirt', 150000, 50, 5)
ON CONFLICT DO NOTHING;
```

### 4. Run the Application
```bash
# Development mode
go run main.go

# Production build
go build -o cashier-api
./cashier-api  # Mac/Linux
cashier-api.exe  # Windows

# Build with optimizations
go build -ldflags="-s -w" -o cashier-api
```

## 📡 API Endpoints

### Health Check
- **GET** `/health` - Check API status

### Product Management (Smart Category Display)
| Method | Endpoint | Description | Category Display | Request Body |
|--------|----------|-------------|------------------|--------------|
| GET | `/api/products` | Get all products | ❌ **NO category** | None |
| POST | `/api/products` | Create new product | N/A | `{"name": "string", "price": int, "stock": int, "category_id": int}` |
| GET | `/api/products/{id}` | Get product by ID | ✅ **WITH category_name** | None |
| PUT | `/api/products/{id}` | Update product | N/A | `{"name": "string", "price": int, "stock": int, "category_id": int}` |
| DELETE | `/api/products/{id}` | Delete product | N/A | None |

### Category Management
| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/api/categories` | Get all categories | None |
| POST | `/api/categories` | Create new category | `{"name": "string", "description": "string"}` |
| GET | `/api/categories/{id}` | Get category by ID | None |
| PUT | `/api/categories/{id}` | Update category | `{"name": "string", "description": "string"}` |
| DELETE | `/api/categories/{id}` | Delete category (fails if products exist) | None |

## 🧪 API Testing Examples

### Products - Smart Category Display
```bash
# Get all products (NO category info - optimized)
curl http://localhost:8080/api/products
# Response: id, name, price, stock only

# Get product detail (WITH category_name)
curl http://localhost:8080/api/products/1
# Response: id, name, price, stock, category_id, category_name

# Create product (requires category_id)
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kopi Kapal Api",
    "price": 2500,
    "stock": 100,
    "category_id": 2
  }'

# Update product (requires category_id)
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 5000,
    "stock": 50,
    "category_id": 3
  }'
```

### Category Operations
```bash
# Get all categories
curl http://localhost:8080/api/categories

# Create new category
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Snacks",
    "description": "Various snacks"
  }'

# Try to delete category with products (will fail)
curl -X DELETE http://localhost:8080/api/categories/1

# Delete category without products
curl -X DELETE http://localhost:8080/api/categories/5
```

## 📊 Database Schema

### Categories Table
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table (with Foreign Key)
```sql
CREATE TABLE products (
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
```

## 🔑 Key Features

### 1. **Complete Layered Architecture**
- Clear separation of concerns
- Easy testing and maintenance
- Scalable structure for future features

### 2. **Smart Data Retrieval** 🎯
- **Product List**: Simple query, no JOIN (faster)
- **Product Detail**: JOIN with categories for complete info
- Optimized performance for common operations

### 3. **Database Relationships**
- Foreign key constraints
- Selective JOIN operations
- Data integrity enforcement

### 4. **Input Validation**
- Product price must be > 0
- Stock cannot be negative
- Category must exist for products
- Category name is required

### 5. **Error Handling**
- Proper HTTP status codes
- Clear error messages
- Protection against orphaned records

### 6. **Configuration Management**
- Environment-based configuration
- Secure credential management
- Easy deployment across environments

## 📝 API Response Examples

### GET /api/products (List - Optimized, NO category)
```json
[
  {
    "id": 1,
    "name": "Indomie Godog",
    "price": 3500,
    "stock": 10
  },
  {
    "id": 2,
    "name": "Vit 1000ml",
    "price": 3000,
    "stock": 40
  }
]
```

### GET /api/products/1 (Detail - WITH category_name)
```json
{
  "id": 1,
  "name": "Indomie Godog",
  "price": 3500,
  "stock": 10,
  "category_id": 1,
  "category_name": "Food"
}
```

### Error Responses
```json
{
  "error": "category not found"
}
// Status: 404 Not Found

{
  "error": "cannot delete category that has products"
}
// Status: 409 Conflict
```

## 🚨 Error Handling

### HTTP Status Codes
| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET/PUT requests |
| 201 | Created | Successful POST requests |
| 400 | Bad Request | Invalid input data |
| 404 | Not Found | Resource not found |
| 405 | Method Not Allowed | Invalid HTTP method |
| 409 | Conflict | Cannot delete category with products |
| 500 | Internal Server Error | Server-side errors |

## 🐛 Troubleshooting

### Common Issues
1. **Database Connection Failed**
   ```bash
   # Check connection string format
   echo $DB_CONN
   
   # Test connection manually
   psql "$DB_CONN" -c "SELECT 1"
   ```

2. **Category Validation Error**
   - Ensure category exists before creating product
   - Check category_id is valid integer

3. **Port Already in Use**
   ```bash
   # Find process using port 8080
   lsof -i :8080
   
   # Or use different port
   export PORT=8081
   ```

### Logs and Debugging
- Check console output for connection messages
- Enable debug mode in `.env`: `DEBUG=true`
- Monitor Supabase dashboard for query performance

## ☁️ Deployment

### Railway (Recommended)
1. Push code to GitHub
2. Go to [railway.app](https://railway.app)
3. New Project → Deploy from GitHub
4. Add environment variables:
   - `PORT`: 8080
   - `DB_CONN`: Your Supabase connection string
5. Deploy!

### Manual Deployment
```bash
# Build for production
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o cashier-api

# Copy to server and run
./cashier-api
```

## 📚 Learning Points

### Layered Architecture Benefits
1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Clear separation of concerns
3. **Scalability**: Easy to add new features
4. **Reusability**: Services can be reused across handlers

### Smart Database Design
1. **Selective JOINs**: Only join when necessary
2. **Performance**: List views are faster without unnecessary joins
3. **Data Integrity**: Constraints protect data quality
4. **Optimization**: Different models for different use cases

## 🎓 Session 2 Achievements

✅ **Complete Layered Architecture** for both Products and Categories  
✅ **Smart Database Relationships** - JOIN only in detail view  
✅ **Selective Data Display** - category only where needed  
✅ **Input Validation** and error handling  
✅ **Configuration Management** with Viper  
✅ **Production-ready** code structure  
✅ **Optimized Performance** for common operations  

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

This project is created for educational purposes as part of the KodingWorks Bootcamp.

---

**Built with ❤️ for KodingWorks Bootcamp Session 2**

**Happy Coding!** 🚀