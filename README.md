# Cashier API - Complete Tutorial Session 1

A simple cashier system API built with Go from scratch. This tutorial covers building a complete CRUD API with in-memory storage, HTTP handling, and deployment.

## 🚀 What We're Building

A Cashier API that can:
- Manage product data (full CRUD operations)
- Handle HTTP requests
- Send JSON responses
- Deploy to free cloud hosting

## 📋 Prerequisites

### Required Tools
- [Go](https://golang.org/dl/) (check with `go version`)
- [VS Code](https://code.visualstudio.com/) (recommended)
- Terminal/Command Prompt
- [Git](https://git-scm.com/) (for deployment)

## 🛠️ Project Setup

### Create Project Directory
```bash
# Create project folder
mkdir cashier-api
cd cashier-api

# Initialize Go module
go mod init cashier-api

# Create main.go file
touch main.go  # Mac/Linux
# or
type nul > main.go  # Windows
```

## 📁 Project Structure
```
cashier-api/
├── main.go          # Main application file
├── go.mod          # Go module file
├── go.sum          # Dependency checksums
├── .vscode/        # VS Code configuration
│   └── launch.json # Debug configuration
└── README.md       # This file
```

## 🖥️ VS Code Configuration

### Install Go Extension
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Go" by Microsoft
4. Install the extension

### Create Launch Configuration
Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Cashier API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "env": {},
            "args": []
        }
    ]
}
```

## 📦 Dependencies

The project uses only standard Go libraries:
- `encoding/json` - JSON encoding/decoding
- `fmt` - Formatting and printing
- `net/http` - HTTP server and handling
- `strconv` - String conversion
- `strings` - String manipulation

## 🏗️ Code Implementation

### Main Application (`main.go`)

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

// In-memory storage (temporary, will replace with database later)
var products = []Product{
	{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10},
	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 20},
}

func main() {
	// Routes configuration
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

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

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  GET    /api/products")
	fmt.Println("  POST   /api/products")
	fmt.Println("  GET    /api/products/{id}")
	fmt.Println("  PUT    /api/products/{id}")
	fmt.Println("  DELETE /api/products/{id}")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// Handler functions
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
```

## 🚀 Running the Application

### Development Mode
```bash
# Run directly
go run main.go
```

### Build and Run
```bash
# Build binary
go build -o cashier-api

# Run binary (Mac/Linux)
./cashier-api

# Run binary (Windows)
cashier-api.exe
```

### Build for Production
```bash
# Smaller binary size
go build -ldflags="-s -w" -o cashier-api

# Cross-compilation examples
GOOS=windows GOARCH=amd64 go build -o cashier-api.exe  # Windows
GOOS=linux GOARCH=amd64 go build -o cashier-api        # Linux
GOOS=darwin GOARCH=amd64 go build -o cashier-api       # macOS
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/products` | Get all products |
| POST | `/api/products` | Create new product |
| GET | `/api/products/{id}` | Get product by ID |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |

## 🧪 Testing the API

### Using cURL
```bash
# Health check
curl http://localhost:8080/health

# Get all products
curl http://localhost:8080/api/products

# Get product by ID
curl http://localhost:8080/api/products/1

# Create product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name": "New Product", "price": 5000, "stock": 100}'

# Update product
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Product", "price": 6000, "stock": 50}'

# Delete product
curl -X DELETE http://localhost:8080/api/products/1
```

### Using GUI Tools
- **Postman**: Popular API testing tool
- **Thunder Client**: VS Code extension
- **Bruno**: Open-source API client

## ☁️ Deployment

### Deploy to Railway (Free)

1. **Prepare for Git**
```bash
# Create .gitignore
echo "cashier-api\ncashier-api.exe\n*.log\n.env" > .gitignore

# Initialize Git repository
git init
git add .
git commit -m "Initial commit"
git branch -M main
```

2. **Push to GitHub**
```bash
git remote add origin https://github.com/yourusername/cashier-api.git
git push -u origin main
```

3. **Deploy on Railway**
   - Go to [railway.app](https://railway.app)
   - Login with GitHub
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway automatically detects Go and deploys

4. **Get Your URL**
   - After deployment, you'll get a URL like:
     `https://cashier-api-production.up.railway.app`

## 📝 VS Code Tips

### Useful Shortcuts
- `Ctrl+Shift+P` - Open command palette
- `Ctrl+` ` - Open integrated terminal
- `F5` - Start debugging
- `Ctrl+F5` - Run without debugging

### Recommended Extensions
1. **Go** - Official Go extension
2. **Go Test Explorer** - Test management
3. **Code Spell Checker** - Spelling checker
4. **Better Comments** - Colorful comments
5. **REST Client** - API testing from VS Code

## 🐛 Debugging

1. Set breakpoints by clicking left of line numbers
2. Press `F5` to start debugging
3. Use debug toolbar:
   - Continue (`F5`)
   - Step Over (`F10`)
   - Step Into (`F11`)
   - Step Out (`Shift+F11`)
   - Restart (`Ctrl+Shift+F5`)
   - Stop (`Shift+F5`)

## 📊 Project Status

✅ **Session 1 Complete** - Basic CRUD API with in-memory storage  
⏳ **Session 2** - Database integration (SQLite)  
⏳ **Session 3** - Authentication & authorization  
⏳ **Session 4** - Advanced features (search, pagination)  

## 📚 Learning Outcomes

✅ Go basics (package, import, struct, function)  
✅ HTTP handling (request, response, routing)  
✅ JSON encoding/decoding  
✅ CRUD operations (in-memory)  
✅ URL path parsing  
✅ Error handling  
✅ VS Code configuration  
✅ Build & deployment  

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

This project is created for educational purposes.

## 🆘 Support

For issues and questions:
1. Check the [Go Documentation](https://golang.org/doc/)
2. Search existing issues
3. Create a new issue with detailed description

---

**Happy Coding!** 🚀