package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cashier-api/database"
	"cashier-api/handlers"
	"cashier-api/repositories"
	"cashier-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Load configuration with Viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Warning: Error reading .env file: %v", err)
		}
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DBConn == "" {
		log.Fatal("DB_CONN environment variable is required")
	}

	// Initialize database connection
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Dependency Injection Setup
	// Category layer first (products need categories)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Product layer (depends on category repo for validation)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Cashier API Running with Complete Layered Architecture",
			"version": "2.1.0",
		})
	})

	// Product routes
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPost:
			productHandler.HandleProducts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPut, http.MethodDelete:
			productHandler.HandleProductByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Category routes
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPost:
			categoryHandler.HandleCategories(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPut, http.MethodDelete:
			categoryHandler.HandleCategoryByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("🚀 Cashier API Server running on http://" + addr)
	fmt.Println("📊 Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  Products (with category relationships):")
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

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
