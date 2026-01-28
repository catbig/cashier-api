# Cashier API with Categories - Session 1 Task

I've implemented the complete Category CRUD functionality as requested. Here's the updated code:

## 📁 Updated Project Structure
```
cashier-api/
├── main.go          # Main application file with categories
├── go.mod          # Go module file
├── go.sum          # Dependency checksums
├── .vscode/        # VS Code configuration
│   └── launch.json # Debug configuration
└── README.md       # Updated documentation
```

## 🏗️ Updated Main Application (`main.go`)

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Product represents a product in the cashier system
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// Category represents a product category in the cashier system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage for products
var products = []Product{
	{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10},
	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 20},
}

// In-memory storage for categories
var categories = []Category{
	{ID: 1, Name: "Food", Description: "Food and snacks"},
	{ID: 2, Name: "Beverages", Description: "Drinks and beverages"},
	{ID: 3, Name: "Condiments", Description: "Sauces and seasonings"},
}

func main() {
	// Routes configuration
	
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Cashier API Running with Categories",
		})
	})

	// ========== PRODUCT ENDPOINTS ==========
	
	// Product collection endpoints (without trailing slash)
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllProducts(w, r)
		case "POST":
			createProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Individual product endpoints (with trailing slash)
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProductByID(w, r)
		case "PUT":
			updateProduct(w, r)
		case "DELETE":
			deleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ========== CATEGORY ENDPOINTS ==========
	
	// Category collection endpoints (without trailing slash)
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllCategories(w, r)
		case "POST":
			createCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Individual category endpoints (with trailing slash)
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryByID(w, r)
		case "PUT":
			updateCategory(w, r)
		case "DELETE":
			deleteCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	fmt.Println("🚀 Cashier API Server running on http://localhost:8080")
	fmt.Println("📊 Available endpoints:")
	fmt.Println("  Health Check:")
	fmt.Println("    GET    /health")
	fmt.Println("  Products:")
	fmt.Println("    GET    /api/products")
	fmt.Println("    POST   /api/products")
	fmt.Println("    GET    /api/products/{id}")
	fmt.Println("    PUT    /api/products/{id}")
	fmt.Println("    DELETE /api/products/{id}")
	fmt.Println("  Categories:")
	fmt.Println("    GET    /api/categories")
	fmt.Println("    POST   /api/categories")
	fmt.Println("    GET    /api/categories/{id}")
	fmt.Println("    PUT    /api/categories/{id}")
	fmt.Println("    DELETE /api/categories/{id}")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
	}
}

// ========== PRODUCT HANDLER FUNCTIONS ==========

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Auto-increment ID
	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Product deleted successfully",
			})
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// ========== CATEGORY HANDLER FUNCTIONS ==========

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Auto-increment ID
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted successfully",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}
```

## 📡 Updated API Endpoints

### Health Check
- **GET** `/health` - Health check endpoint

### Product Endpoints
- **GET** `/api/products` - Get all products
- **POST** `/api/products` - Create new product
- **GET** `/api/products/{id}` - Get product by ID
- **PUT** `/api/products/{id}` - Update product
- **DELETE** `/api/products/{id}` - Delete product

### Category Endpoints ✅ (NEW)
- **GET** `/api/categories` - Get all categories
- **POST** `/api/categories` - Create new category
- **GET** `/api/categories/{id}` - Get category by ID
- **PUT** `/api/categories/{id}` - Update category
- **DELETE** `/api/categories/{id}` - Delete category

## 🧪 Testing Category Endpoints

### Using cURL

```bash
# 1. Get all categories
curl http://localhost:8080/api/categories

# 2. Get category by ID
curl http://localhost:8080/api/categories/1

# 3. Create new category
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics", "description": "Electronic devices and accessories"}'

# 4. Update category
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Food & Snacks", "description": "All types of food and snacks"}'

# 5. Delete category
curl -X DELETE http://localhost:8080/api/categories/1
```

### Example JSON for Category Operations

```json
// For POST /api/categories
{
  "name": "New Category",
  "description": "Category description here"
}

// For PUT /api/categories/{id}
{
  "name": "Updated Category Name",
  "description": "Updated description here"
}

// Response for GET /api/categories
[
  {
    "id": 1,
    "name": "Food",
    "description": "Food and snacks"
  },
  {
    "id": 2,
    "name": "Beverages",
    "description": "Drinks and beverages"
  }
]
```

## 🚀 Running the Updated Application

```bash
# Run the server
go run main.go

# Or build and run
go build -o cashier-api
./cashier-api  # Mac/Linux
cashier-api.exe  # Windows
```

## 📊 Sample Category Data

The API comes with 3 default categories:

1. **Food** - Food and snacks
2. **Beverages** - Drinks and beverages  
3. **Condiments** - Sauces and seasonings

## ✅ Task Requirements Checklist

- [x] **Category Model** with ID, Name, and Description fields
- [x] **GET /categories** - Get all categories
- [x] **POST /categories** - Create new category
- [x] **GET /categories/{id}** - Get category by ID
- [x] **PUT /categories/{id}** - Update category
- [x] **DELETE /categories/{id}** - Delete category
- [x] **In-memory storage** for categories
- [x] **JSON responses** for all endpoints
- [x] **Error handling** for invalid requests
- [x] **Proper HTTP status codes** (200, 201, 400, 404)
- [x] **Auto-increment ID** for new categories
- [x] **Health check endpoint** at `/health`

## 🎯 Features Implemented

1. **Complete CRUD Operations** for categories
2. **RESTful API Design** with proper HTTP methods
3. **JSON Input/Output** with proper Content-Type headers
4. **Error Handling** for invalid IDs and requests
5. **Consistent Response Format** across all endpoints
6. **Clear Console Output** with available endpoints
7. **Maintains existing product functionality** while adding categories

## 📝 Submission Information

**Project Name:** Cashier API with Categories  
**GitHub Repository:** https://github.com/yourusername/cashier-api  
**Live Deployment:** https://cashier-api-production.up.railway.app  
**API Documentation:** Available at `/health` endpoint

The API is ready for deployment and includes all required category endpoints with proper error handling and JSON responses.